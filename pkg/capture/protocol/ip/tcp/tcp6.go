package tcp

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
)

func ProcessTcp6Packets(prev, curr, next base.PacketBatch) *ChannelAggregatedValues {
	return processPackets(prev, curr, next, v6IpTcpLayer)
}

func v6IpTcpLayer(item *base.PacketItem) (base.LayerIp, *layers.TCP) {
	p := item.Packet
	ip6 := p.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
	tcp := p.Layer(layers.LayerTypeTCP).(*layers.TCP)
	ip := base.NewLayerIp6(ip6)
	return ip, tcp
}

func printTcp6Packet(ip *layers.IPv6, tcp *layers.TCP, payload bool) {
	// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
	// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
	format := fmt.Sprintf("TCP%d-[%s:%d -> %s:%d]-%s%s%s%s%s%s%s%s%s[seq:%d][len:%d]",
		ip.Version,
		ip.SrcIP, tcp.SrcPort, ip.DstIP, tcp.DstPort,
		flag(tcp.FIN, "FIN"),
		flag(tcp.SYN, "SYN"),
		flag(tcp.RST, "RST"),
		flag(tcp.PSH, "PSH"),
		flag(tcp.ACK, "ACK"),
		flag(tcp.URG, "URG"),
		flag(tcp.ECE, "ECE"),
		flag(tcp.CWR, "CWR"),
		flag(tcp.NS, "NS"),
		tcp.Seq, ip.Length)
	logger.Info(format)
}
