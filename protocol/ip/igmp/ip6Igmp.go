package igmp

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func ProcessIgmp6Packet(ip *layers.IPv6, packet gopacket.Packet) bool {
	igmpLayer := packet.Layer(layers.LayerTypeIGMP)
	if igmpLayer != nil {
		igmp1or2, ok := igmpLayer.(*layers.IGMPv1or2) // IGMPv1or2
		if !ok {
			igmp3, ok := igmpLayer.(*layers.IGMP) // IGMPv3
			if !ok {
				logger.Error("Convert to IGMPv3、IGMPv1or2 failed!")
			}
			printIp6Igmp3Packet(ip, igmp3)
		} else if igmp1or2 != nil {
			printIp6Igmp1or2Packet(ip, igmp1or2)
		}
		return true
	}
	return false
}

func printIp6Igmp3Packet(ip *layers.IPv6, igmp *layers.IGMP) {
	format := fmt.Sprintf("IGMP%d-[%s -> %s]-[len:%d]",
		igmp.Version, ip.SrcIP, ip.DstIP,
		ip.Length)
	logger.Info(format)
}

func printIp6Igmp1or2Packet(ip *layers.IPv6, igmp *layers.IGMPv1or2) {
	format := fmt.Sprintf("IGMP%d-[%s -> %s]-[len:%d]",
		igmp.Version, ip.SrcIP, ip.DstIP, ip.Length)
	logger.Info(format)

}
