package base

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type IpVersion int8

const (
	NotIp IpVersion = -1
	Ipv4  IpVersion = 1
	Ipv6  IpVersion = 2
)

type ProtocolClass int8

const (
	ProtocolUnknown ProtocolClass = -1
	ProtocolIgmp    ProtocolClass = 1
	ProtocolIcmp    ProtocolClass = 2
	ProtocolTcp     ProtocolClass = 3
	ProtocolUdp     ProtocolClass = 4
)

func IsIpPacket(packet gopacket.Packet) bool {
	return Version(packet) > 0
}

func Version(packet gopacket.Packet) IpVersion {
	if nil != packet.Layer(layers.LayerTypeIPv4) {
		return Ipv4
	}
	if nil != packet.Layer(layers.LayerTypeIPv6) {
		return Ipv6
	}
	return NotIp
}

type LayerIp interface {
	GetSrcIp() string
	GetDstIp() string
	GetPktSz() uint16
}
