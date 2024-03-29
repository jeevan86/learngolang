package base

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Ip6PacketCache 实现了通用缓存接口PacketCache
type Ip6PacketCache struct {
	delegate *DefaultPacketCache
}

func (c *Ip6PacketCache) PutPacket(bucket, millis int64, p gopacket.Packet) {
	c.delegate.putPacket(bucket, millis, c.distinct(p), c.protocol(p), p)
}

func (c *Ip6PacketCache) GetBatch(bucket int64) ProtocolBatch {
	return c.delegate.getBatch(bucket)
}

func (c *Ip6PacketCache) DelBatch(bucket int64) {
	c.delegate.delBatch(bucket)
}

func (c *Ip6PacketCache) distinct(p gopacket.Packet) DistinctPacketId {
	iPv6, _ := p.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
	return DistinctPacketId{
		DstIp: iPv6.DstIP.String(),
		PktId: iPv6.FlowLabel,
	}
}

func (c *Ip6PacketCache) protocol(p gopacket.Packet) ProtocolClass {
	if nil != p.Layer(layers.LayerTypeTCP) {
		return ProtocolTcp
	}
	if nil != p.Layer(layers.LayerTypeUDP) {
		return ProtocolUdp
	}
	if nil != p.Layer(layers.LayerTypeICMPv6) {
		return ProtocolIcmp
	}
	if nil != p.Layer(layers.LayerTypeIGMP) {
		return ProtocolIgmp
	}
	return ProtocolUnknown
}

// init
// @title       包初始化执行的函数
// @description 包初始化执行的函数，将Ipv6对应的缓存创建函数注册上
// @auth        小卒     2022/08/03 10:57
func init() {
	cacheCreator[Ipv6] = func(delegate *DefaultPacketCache) PacketCache {
		return &Ip6PacketCache{delegate: delegate}
	}
}
