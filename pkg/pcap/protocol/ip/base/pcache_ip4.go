package base

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Ip4PacketCache struct {
	delegate *DefaultPacketCache
}

func (c *Ip4PacketCache) PutPacket(bucket, millis int64, p gopacket.Packet) {
	c.delegate.putPacket(bucket, millis, c.distinct(p), c.protocol(p), p)
}

func (c *Ip4PacketCache) GetBatch(bucket int64) ProtocolBatch {
	return c.delegate.getBatch(bucket)
}

func (c *Ip4PacketCache) DelBatch(bucket int64) {
	c.delegate.delBatch(bucket)
}

func (c *Ip4PacketCache) distinct(p gopacket.Packet) DistinctPacketId {
	iPv4, _ := p.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	return DistinctPacketId{
		DstIp: iPv4.DstIP.String(),
		PktId: uint32(iPv4.Id),
	}
}

func (c *Ip4PacketCache) protocol(packet gopacket.Packet) ProtocolClass {
	iPv4, _ := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	switch iPv4.Protocol {
	case layers.IPProtocolIGMP:
		return ProtocolIgmp
	case layers.IPProtocolICMPv4:
		return ProtocolIcmp
	case layers.IPProtocolTCP:
		return ProtocolTcp
	case layers.IPProtocolUDP:
		return ProtocolUdp
	}
	return ProtocolUnknown
}

func init() {
	cacheCreator[Ipv4] = func(delegate *DefaultPacketCache) PacketCache {
		return &Ip4PacketCache{delegate: delegate}
	}
}
