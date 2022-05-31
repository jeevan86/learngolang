package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/jeevan86/learngolang/cmd"
	"github.com/jeevan86/learngolang/pkg/collect/client"
	"github.com/jeevan86/learngolang/pkg/config"
	"github.com/jeevan86/learngolang/pkg/log"
	"github.com/jeevan86/learngolang/pkg/pcap"
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/ip"
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/ip/base"
	"github.com/jeevan86/learngolang/pkg/server/http"
	"github.com/jeevan86/learngolang/pkg/util/parallel"
	"github.com/jeevan86/learngolang/pkg/util/reactor"
	"os"
)

var logger = log.NewLogger()

var packetProcessor *base.PacketProcessor
var collector *client.Collect

/*
export PATH=$GOPATH/bin:$PATH
export PRJ_DIR=`pwd`
cd ${PRJ_DIR}/pkg/collect/api/grpc/pb && go generate && cd ${PRJ_DIR}
go build --ldflags "-extldflags -static" \
 -o ${PRJ_DIR}/dist/binary/gopcap-agent  \
${PRJ_DIR}/cmd/pcap/main.go

cd ${PRJ_DIR}/dist && \
dlv --listen=:8626 --headless=true --api-version=2 --accept-multiclient exec \
./binary/gopcap-agent \
-- -config-file ${PRJ_DIR}/dist/config/agent.yml

cd ${PRJ_DIR}/dist && ./binary/gopcap-agent \
 -config-file ${PRJ_DIR}/dist/config/agent.yml
*/
func main() {
	packetProcessor = ip.NewPacketProcessor()
	collector = client.NewCollector()
	http.Start()
	collector.Start(packetProcessor.Out())
	packetProcessor.Start(func() []string {
		return collector.Api.GetLocalIpList(config.GetConfig().NodeIp)
	})
	startCapture(captureHandeFunc())
	fmt.Println("packet capture started.")
	cmd.WaitForSig()
	packetProcessor.Stop()
	collector.Stop()
	http.Stop()
}

func startCapture(captured func(uint64, gopacket.Packet)) {
	for _, device := range config.GetConfig().Pcap.Devices {
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
	switch config.GetConfig().Pcap.ParType {
	case config.ParTypeRoutine:
		return parRoutinesFunc() // with Routines
	case config.ParTypeReactor:
		return reactorFluxFunc() // with RxGo
	default:
		return parRoutinesFunc() // with Routines
	}
}

func parRoutinesFunc() func(uint64, gopacket.Packet) {
	parallelism := config.GetConfig().Pcap.Routine.Parallelism
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
	bufferSz := config.GetConfig().Pcap.Reactor.BufferSz
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
