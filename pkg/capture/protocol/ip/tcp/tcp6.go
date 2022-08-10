package tcp

import (
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
)

func ProcessTcp6Packets(prev, curr, next base.PacketBatch) *ChannelAggregatedValues {
	return processPackets(prev, curr, next, v6IpTcpLayer)
}

var v6IpTcpLayer ipTcpLayerFunc = func(item *base.PacketItem) (base.LayerIp, *layers.TCP) {
	p := item.Packet
	ip6 := p.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
	tcp := p.Layer(layers.LayerTypeTCP).(*layers.TCP)
	ip := base.NewLayerIp6(ip6)
	return ip, tcp
}
