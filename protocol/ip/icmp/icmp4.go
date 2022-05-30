package icmp

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func ProcessIcmp4Packet(ip *layers.IPv4, packet gopacket.Packet) bool {
	icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
	if icmpLayer != nil {
		icmp, ok := icmpLayer.(*layers.ICMPv4)
		if !ok {
			logger.Error("Convert to ICMPv4 failed!")
		} else if icmp != nil {
			printIcmp4Packet(ip, icmp)
		}
		return true
	}
	return false
}
func printIcmp4Packet(ip *layers.IPv4, icmp *layers.ICMPv4) {
	format := fmt.Sprintf("ICMP%d-[%s -> %s][seq:%d]-[ttl:%d][ihl:%d][tos:%d][flg:%s][len:%d]",
		ip.Version, ip.SrcIP, ip.DstIP, icmp.Seq,
		ip.TTL, ip.IHL, ip.TOS, ip.Flags.String(),
		ip.Length)
	logger.Info(format)
}
