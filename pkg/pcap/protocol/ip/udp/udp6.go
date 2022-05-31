package udp

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/ip/base"
)

func ProcessUdp6Packets(prev, curr, next base.PacketBatch) *ChannelAggregatedValues {
	result := make(map[Channel]*AggregatedValues, 32)
	for _, item := range curr {
		p := item.Packet
		ip6 := p.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
		udp := p.Layer(layers.LayerTypeUDP).(*layers.UDP)
		aggregate(
			result,
			ip6.SrcIP.String(),
			ip6.DstIP.String(),
			ip6.Length,
			udp,
		)
	}
	return &ChannelAggregatedValues{
		Values: result,
	}
}

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
