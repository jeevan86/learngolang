package icmp

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func ProcessIcmp6Packet(ip *layers.IPv6, packet gopacket.Packet) bool {
	icmpLayer := packet.Layer(layers.LayerTypeICMPv6)
	if icmpLayer != nil {
		icmp, ok := icmpLayer.(*layers.ICMPv6)
		if !ok {
			logger.Error("Convert to ICMPv4 failed!")
		} else if icmp != nil {
			printIcmp6Packet(ip, icmp)
		}
		return true
	}
	return false
}

func printIcmp6Packet(ip *layers.IPv6, icmp *layers.ICMPv6) {
	format := fmt.Sprintf("ICMP%d-[%s -> %s][len:%d]",
		ip.Version, ip.SrcIP, ip.DstIP, ip.Length)
	logger.Info(format)
}
