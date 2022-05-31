package icmp

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/ip/base"
)

func ProcessIcmp6Packets(prev, curr, next base.PacketBatch) *ChannelAggregatedValues {
	result := make(map[Channel]*AggregatedValues, 32)
	for _, item := range curr {
		p := item.Packet
		ip6 := p.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
		aggregate(
			result,
			ip6.SrcIP.String(),
			ip6.DstIP.String(),
			ip6.Length,
		)
	}
	return &ChannelAggregatedValues{
		Values: result,
	}
}

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
