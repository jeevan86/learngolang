package pcap

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/jeevan86/learngolang/pkg/config"
	"github.com/jeevan86/learngolang/pkg/log"
	"os"
	"strings"
	"time"
)

var logger = log.NewLogger()

func StartCapture(process func(packet gopacket.Packet)) {
	for _, device := range config.GetConfig().Agent.Pcap.Devices {
		handle, err := createHandle(device.Duration, device.Prefix, device.Snaplen, device.Promisc)
		if err != nil {
			// permission issue
			logger.Fatal("Unable to create handler, cause: %s", err.Error())
			os.Exit(-1)
		} else {
			go func() { startCapture(handle, process) }()
		}
	}
	logger.Info("Packet capture started.")
}

func createHandle(duration, device string, snaplen int32, promisc bool) (*pcap.Handle, error) {
	// 设置1秒？
	secs, _ := time.ParseDuration(duration)
	//handle, err := pcap.OpenLive("en0", 40, true, secs)
	h, err := pcap.OpenLive(device, snaplen, promisc, secs)
	if err != nil {
		return nil, err
	}
	return h, nil
}

const maxIdxPerCycle = ^uint64(0) - 99999999

func startCapture(handle *pcap.Handle, process func(packet gopacket.Packet)) {
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
	captured := captureHandeFunc(process)
	for packet := range packetSource.Packets() {
		captured(idx, packet)
		postCaptured()
	}
}

func PrintCaptureInfo() {
	captureInfo := gopacket.CaptureInfo{}
	layer := layers.CiscoDiscovery{}
	fmt.Println(captureInfo, layer)
}

func PrintAppLayer(packet gopacket.Packet) {
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

func PrintAllLayer(packet gopacket.Packet) {
	// Iterate over all layers, printing out each layer type
	fmt.Println("All packet layers:")
	for _, layer := range packet.Layers() {
		fmt.Println("- ", layer.LayerType())
	}
}
