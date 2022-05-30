package udp

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func ProcessUdp6Packet(ip *layers.IPv6, packet gopacket.Packet) bool {
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp, ok := udpLayer.(*layers.UDP)
		if !ok {
			logger.Error("Convert to TCP failed!")
		} else if udp != nil {
			printUdp6Packet(ip, udp)
		}
		return true
	}
	return false
}

func printUdp6Packet(ip *layers.IPv6, udp *layers.UDP) {
	// Checksum, SrcIP, DstIP
	format := fmt.Sprintf("UDP%d-[%s:%d -> %s:%d]-[len:%d]",
		ip.Version,
		ip.SrcIP, udp.SrcPort, ip.DstIP, udp.DstPort,
		ip.Length)
	logger.Info(format)
}
