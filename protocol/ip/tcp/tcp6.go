package tcp

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func ProcessTcp6Packet(ip *layers.IPv6, packet gopacket.Packet) bool {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, ok := tcpLayer.(*layers.TCP)
		if !ok {
			logger.Error("Convert to TCP failed!")
		} else if tcp != nil {
			printTcp6Packet(ip, tcp, true)
		}
		return true
	}
	return false
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
