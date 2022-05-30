package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"gopackettest/actuator"
	logging "gopackettest/logger"
	"gopackettest/protocol/device"
	"gopackettest/protocol/ip"
	"gopackettest/util"
	"strings"
	"time"
)

var logger = logging.LoggerFactory.NewLogger([]string{"stdout"}, []string{"stderr"})

type Routine struct {
	ch          chan interface{}
	disposables []*util.DisposableGoRoutine
}

type Routines struct {
	parallelism int
	routines    []Routine
}

func NewRoutines(parallelism int, process func(interface{})) *Routines {
	internal := make([]Routine, parallelism)
	routines := Routines{
		parallelism: parallelism,
		routines:    internal,
	}
	// 初始化4个协程来处理数据
	for i := 0; i < parallelism; i++ {
		routine := Routine{
			ch: make(chan interface{}, 2048),
		}
		routine.disposables = []*util.DisposableGoRoutine{
			util.StartDisposable(func() {
				p := <-routine.ch
				process(p)
			})}
		routines.routines[i] = routine
	}
	return &routines
}

func (r *Routines) dispatch(i uint64, p gopacket.Packet) {
	routines := r.routines
	parallelism := uint64(r.parallelism)
	routines[i%parallelism].ch <- p
}

func (r *Routines) close() {
	for _, routine := range r.routines {
		for _, disposable := range routine.disposables {
			disposable.Dispose()
		}
	}
}

var handle *pcap.Handle

func main() {
	actuator.StartActuator()
	// 定时更新本地的IP地址信息
	startLocalIpRefresher()
	// handle
	err := createHandle()
	if err != nil {
		// permission issue
		fmt.Println(err.Error())
		return
	}
	// with Routines
	routines := NewRoutines(
		4,
		func(packet interface{}) {
			processPacket(packet.(gopacket.Packet))
		},
	)
	nextItem := func(i uint64, p gopacket.Packet) {
		routines.dispatch(i, p)
	}
	startCapture(handle, nextItem)

	// with rxgo
	//fluxSink := util.NewFluxSink(2048)
	//fluxSink.Map(func(e interface{}) interface{} {
	//	return e
	//	//}).FlatMap(func(e interface{}) *util.FluxSink {
	//	//	packet := e.(gopacket.Packet)
	//	//	newFluxSink := util.NewFluxSink(2048)
	//	//	newFluxSink.Next(packet)
	//	//	return newFluxSink
	//}).DoOnNext(func(p interface{}) {
	//	processPacket(p.(gopacket.Packet))
	//}).Subscribe(func(e interface{}) {
	//	// do nothing
	//})
	//nextItem := func(i uint64, p gopacket.Packet) {
	//	fluxSink.Next(p)
	//}
	//startCapture(handle, nextItem)
}

func startLocalIpRefresher() {
	go func() {
		ticker := time.NewTicker(time.Second * 30)
		for t := range ticker.C {
			t.Second()
			// 更新本地的IP地址
			for _, inf := range device.AllDevices() {
				addresses := make([]pcap.InterfaceAddress, 128)
				for _, address := range inf.Addresses {
					addresses = append(addresses, address)
				}
				ip.UpdateAddresses(addresses)
			}
		}
	}()
}

func createHandle() error {
	// 设置1秒？
	secs, _ := time.ParseDuration("2s")
	//handle, err := pcap.OpenLive("en0", 40, true, secs)
	h, err := pcap.OpenLive("any", 60, true, secs)
	if err != nil {
		return err
	}
	handle = h
	return nil
}

const maxIdxPerCycle = ^uint64(0) - 99999999

func startCapture(handle *pcap.Handle, captured func(i uint64, p gopacket.Packet)) {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	cycle := uint64(0)
	idx := uint64(0)
	postCaptured := func() {
		if idx >= maxIdxPerCycle {
			idx = 0
			cycle++
		}
		idx++
	}
	for packet := range packetSource.Packets() {
		captured(idx, packet)
		postCaptured()
	}
}

func printCaptureInfo() {
	captureInfo := gopacket.CaptureInfo{}
	layer := layers.CiscoDiscovery{}
	fmt.Println(captureInfo, layer)
}

func processPacket(packet gopacket.Packet) {
	if ip.IsIpPacket(packet) {
		ip.ProcessIpPacket(packet)
	}
	//// Check for errors
	//if err := packet.ErrorLayer(); err != nil {
	//	logger.Warn(fmt.Sprintf("Error decoding some part of the packet: %s", err.Error()))
	//}
}

func printAppLayer(packet gopacket.Packet) {
	// When iterating through packet.Layers() above,
	// if it lists Payload layer then that is the same as
	// this applicationLayer. applicationLayer contains the payload
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		fmt.Println("Application layer/Payload found.")
		fmt.Printf("%s\n", applicationLayer.Payload())
		// Search for a string inside the payload
		payload := string(applicationLayer.Payload())
		if strings.Contains(payload, "HTTP") {
			fmt.Println("HTTP found!")
		}
	}
}

func printAllLayer(packet gopacket.Packet) {
	// Iterate over all layers, printing out each layer type
	fmt.Println("All packet layers:")
	for _, layer := range packet.Layers() {
		fmt.Println("- ", layer.LayerType())
	}
}
