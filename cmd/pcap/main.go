package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/jeevan86/learngolang/cmd"
	"github.com/jeevan86/learngolang/pkg/collect"
	"github.com/jeevan86/learngolang/pkg/config"
	"github.com/jeevan86/learngolang/pkg/log"
	"github.com/jeevan86/learngolang/pkg/pcap"
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/ip"
	"github.com/jeevan86/learngolang/pkg/server/http"
	"github.com/jeevan86/learngolang/pkg/util/parallel"
	"github.com/jeevan86/learngolang/pkg/util/reactor"
	"os"
)

var logger = log.NewLogger()

var packetProcessor = ip.NewPacketProcessor()
var collector = collect.NewCollector()
var localIpGetFunc = func() []string {
	return collector.Api.GetLocalIpList(config.Config.NodeIp)
}

func main() {
	http.Start()
	collector.Start(packetProcessor.Out())
	packetProcessor.Start(localIpGetFunc)
	startCapture(captureHandeFunc())
	fmt.Println("packet capture started.")
	cmd.WaitForSig()
	packetProcessor.Stop()
	collector.Stop()
	http.Stop()
}

func startCapture(captured func(uint64, gopacket.Packet)) {
	for _, device := range config.Config.Pcap.Devices {
		handle, err := pcap.CreateHandle(device.Duration, device.Prefix, device.Snaplen, device.Promisc)
		if err != nil {
			// permission issue
			logger.Fatal("Unable to create handler, cause: %s", err.Error())
			os.Exit(-1)
		} else {
			go func() { pcap.StartCapture(handle, captured) }()
		}
	}
}

func captureHandeFunc() func(uint64, gopacket.Packet) {
	switch config.Config.Pcap.ParType {
	case config.ParTypeRoutine:
		return parRoutinesFunc() // with Routines
	case config.ParTypeReactor:
		return reactorFluxFunc() // with RxGo
	default:
		return parRoutinesFunc() // with Routines
	}
}

func parRoutinesFunc() func(uint64, gopacket.Packet) {
	parallelism := config.Config.Pcap.Routine.Parallelism
	// with Routines
	routines := parallel.NewParRoutines(
		parallelism,
		2048,
		false,
		func(p interface{}) {
			processPacket(p.(gopacket.Packet))
		},
	)
	return func(i uint64, p gopacket.Packet) {
		routines.Dispatch(i, p)
	}
}

func reactorFluxFunc() func(uint64, gopacket.Packet) {
	bufferSz := config.Config.Pcap.Reactor.BufferSz
	fluxSink := reactor.NewFluxSink(bufferSz)
	fluxSink.Map(func(e interface{}) interface{} {
		return e
		//}).FlatMap(func(e interface{}) *util.FluxSink {
		//	packet := e.(gopacket.Packet)
		//	newFluxSink := util.NewFluxSink(2048)
		//	newFluxSink.Next(packet)
		//	return newFluxSink
	}).DoOnNext(func(p interface{}) {
		processPacket(p.(gopacket.Packet))
	}).Subscribe(func(e interface{}) {
		// do nothing
	})
	return func(i uint64, p gopacket.Packet) {
		fluxSink.Next(p)
	}
}

// processPacket 只处理IP包
func processPacket(packet gopacket.Packet) {
	if ip.IsIpPacket(packet) {
		packetProcessor.Process(packet)
	} else {
		// Check for errors
		if err := packet.ErrorLayer(); err != nil {
			logger.Warn("Error decoding some part of the packet: %s", err.Error())
		}
	}
}
