package capture

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"strings"
)

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
