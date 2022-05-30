package ip

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"gopackettest/protocol/ip/dup"
	"gopackettest/protocol/ip/icmp"
	"gopackettest/protocol/ip/igmp"
	"gopackettest/protocol/ip/tcp"
)

func processIp6Packet(packet gopacket.Packet) {
	iPv6, _ := packet.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
	if iPv6.DstIP.IsLoopback() || iPv6.SrcIP.IsLoopback() {
		return
	}
	if igmp.ProcessIgmp6Packet(iPv6, packet) {
		return
	}
	if icmp.ProcessIcmp6Packet(iPv6, packet) {
		return
	}
	if tcp.ProcessTcp6Packet(iPv6, packet) {
		return
	}
	if udp.ProcessUdp6Packet(iPv6, packet) {
		return
	}
}
