package udp

import (
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
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
