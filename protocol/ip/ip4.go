package ip

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"gopackettest/protocol/ip/dup"
	"gopackettest/protocol/ip/icmp"
	"gopackettest/protocol/ip/igmp"
	"gopackettest/protocol/ip/tcp"
)

func processIp4Packet(packet gopacket.Packet) {
	iPv4, _ := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	if iPv4.DstIP.IsLoopback() || iPv4.SrcIP.IsLoopback() {
		return
	}
	if iPv4.Protocol == layers.IPProtocolIGMP {
		igmp.ProcessIgmp4Packet(iPv4, packet)
	} else if iPv4.Protocol == layers.IPProtocolICMPv4 {
		icmp.ProcessIcmp4Packet(iPv4, packet)
	} else if iPv4.Protocol == layers.IPProtocolTCP {
		tcp.ProcessTcp4Packet(iPv4, packet)
	} else if iPv4.Protocol == layers.IPProtocolUDP {
		udp.ProcessUdp4Packet(iPv4, packet)
	}
}
