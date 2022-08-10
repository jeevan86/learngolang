package tcp

import (
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
)

// ProcessTcp6Packets
// @title       ProcessTcp6Packets
// @description 处理ip包批次数据，返回聚合结果
// @auth        小卒    2022/08/03 10:57
// @param       prev   base.PacketBatch        "上一个时间窗的ip包数据"
// @param       curr   base.PacketBatch        "要处理的ip包数据"
// @param       next   base.PacketBatch        "下一个时间窗的ip包数据"
// @return      r      ChannelAggregatedValues "按Channel聚合处理的结果"
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
