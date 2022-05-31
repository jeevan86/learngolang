package ip

import (
	"github.com/google/gopacket"
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/ip/base"
)

func IsIpPacket(packet gopacket.Packet) bool {
	return base.IsIpPacket(packet)
}

var ipPacketProcessor *base.PacketProcessor

func NewPacketProcessor() *base.PacketProcessor {
	if ipPacketProcessor == nil {
		ipPacketProcessor = base.NewPacketProcessor(processIp4Packets, processIp6Packets)
	}
	isLocal := func(ip string) bool {
		return ipPacketProcessor.IsLocalIp(ip)
	}
	base.SetCheckLocalIpFunc(isLocal)
	return ipPacketProcessor
}
func packetBatches(clz base.ProtocolClass,
	prev, curr, next base.ProtocolBatch) (base.PacketBatch, base.PacketBatch, base.PacketBatch) {
	return prev[clz], curr[clz], next[clz]
}
