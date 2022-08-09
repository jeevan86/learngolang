package base

import "github.com/google/gopacket/layers"

// HopLimit
type hop uint8

func (p *PacketProcessor) flushIp6(keyNext, keyCurr, keyPrev int64) map[ProtocolClass]interface{} {
	prev := p.ip6PacketCache.GetBatch(keyPrev)
	curr := p.ip6PacketCache.GetBatch(keyCurr)
	next := p.ip6PacketCache.GetBatch(keyNext)
	ip6 := p.ip6Processor(prev, curr, next)
	p.ip6PacketCache.DelBatch(keyPrev)
	return ip6
}

type LayerIp6 struct {
	layer *layers.IPv6
}

func NewLayerIp6(l *layers.IPv6) *LayerIp6 {
	return &LayerIp6{
		layer: l,
	}
}

func (l *LayerIp6) GetSrcIp() string {
	return l.layer.SrcIP.String()
}

func (l *LayerIp6) GetDstIp() string {
	return l.layer.SrcIP.String()
}

func (l *LayerIp6) GetPktSz() uint16 {
	return l.layer.Length
}
