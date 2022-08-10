package udp

import (
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
)

func ProcessUdp4Packets(prev, curr, next base.PacketBatch) *ChannelAggregatedValues {
	result := make(map[Channel]*AggregatedValues, 32)
	for _, item := range curr {
		p := item.Packet
		ip4 := p.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
		udp := p.Layer(layers.LayerTypeUDP).(*layers.UDP)
		aggregate(
			result,
			ip4.SrcIP.String(),
			ip4.DstIP.String(),
			ip4.Length,
			udp,
		)
	}
	return &ChannelAggregatedValues{
		Values: result,
	}
}
