package base

import "github.com/google/gopacket/layers"

// HopLimit
type hop uint8

// flushIp6
// @title       flushIp6
// @description 数据交给ip6Processor函数处理、清理、返回打包后的批次数据
// @auth        小卒  2022/08/03 10:57
// @param       keyNext int64 "下一个时间窗的时间戳，作为缓存的KEY"
// @param       keyCurr int64 "当前的时间窗的时间戳，作为缓存的KEY"
// @param       keyPrev int64 "上一个时间窗的时间戳，作为缓存的KEY"
// @return      r       map[ProtocolClass]interface{} "打包后的批次"
func (p *PacketProcessor) flushIp6(keyNext, keyCurr, keyPrev int64) map[ProtocolClass]interface{} {
	prev := p.ip6PacketCache.GetBatch(keyPrev)
	curr := p.ip6PacketCache.GetBatch(keyCurr)
	next := p.ip6PacketCache.GetBatch(keyNext)
	ip6 := p.ip6Processor(prev, curr, next)
	p.ip6PacketCache.DelBatch(keyPrev)
	return ip6
}

// LayerIp6 使用gopacket的layers.Ipv6结构
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
	return l.layer.DstIP.String()
}

func (l *LayerIp6) GetPktSz() uint16 {
	return l.layer.Length
}
