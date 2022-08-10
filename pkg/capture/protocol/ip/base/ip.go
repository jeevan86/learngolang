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

// IsIpPacket
// @title       是否是IP包
// @description 是否是IP包
// @auth        小卒  2022/08/03 10:57
// @param       packet gopacket.Packet "网络包对象"
// @return      r      bool            "是否是IP包"
func IsIpPacket(packet gopacket.Packet) bool {
	return Version(packet) > 0
}

// Version
// @title       版本
// @description IP包对应的IP版本，ip4\ip6
// @auth        小卒  2022/08/03 10:57
// @param       packet gopacket.Packet "网络包对象"
// @return      r      IpVersion       "版本常量值"
func Version(packet gopacket.Packet) IpVersion {
	if nil != packet.Layer(layers.LayerTypeIPv4) {
		return Ipv4
	}
	if nil != packet.Layer(layers.LayerTypeIPv6) {
		return Ipv6
	}
	return NotIp
}

// LayerIp Ip层接口
type LayerIp interface {
	GetSrcIp() string
	GetDstIp() string
	GetPktSz() uint16
}
