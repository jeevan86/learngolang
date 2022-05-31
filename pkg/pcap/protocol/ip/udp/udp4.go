package udp

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/ip/base"
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

func ProcessUdp4Packet(ip *layers.IPv4, packet gopacket.Packet) bool {
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp, ok := udpLayer.(*layers.UDP)
		if !ok {
			logger.Error("Convert to TCP failed!")
		} else if udp != nil {
			printUdp4Packet(ip, udp)
		}
		return true
	}
	return false
}

func printUdp4Packet(ip *layers.IPv4, udp *layers.UDP) {
	// Checksum, SrcIP, DstIP
	format := fmt.Sprintf("UDP%d-[%s:%d -> %s:%d]-[ttl:%d][ihl:%d][tos:%d][flg:%s][len:%d]",
		ip.Version,
		ip.SrcIP, udp.SrcPort, ip.DstIP, udp.DstPort,
		ip.TTL, ip.IHL, ip.TOS, ip.Flags.String(),
		ip.Length)
	if logger.IsDebugEnabled() {
		logger.Debug(format)
	} else {
		logger.Info(format)
	}
}
