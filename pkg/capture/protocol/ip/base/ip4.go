package base

import "github.com/google/gopacket/layers"

// 已经没用
//  1000 -- minimize delay #最小延迟
//  0100 -- maximize throughput #最大吞吐量
//  0010 -- maximize reliability #最高可靠性
//  0001 -- minimize monetary cost #最小费用
//  0000 -- normal service #一般服务
type tos int8

// 包头长度
type ihl int8

// TTL
//  每经过一层路由减1，默认64?
type ttl int8

// Flag（标志位）： 标志字段在IP报头中占3位。
//  第1位作为保留；
//  第2位，分段，是否允许分片;（如果不允许分片，包超过了数据连路支持的最大长度，则丢弃该包，返回发送者一个 ICMP 错误）
//  第3位，更多分段。表示是否最后一个分片。
// 当目的主机接收到一个IP数据报时，会首先查看该数据报的标识符，并且检查标志位的第3位是置0或置1，以确定是否还有更多的分段。
// 如果还有后续报文，接收主机则将接收到的报文放在缓存直到接收完所有具有相同标识符的数据报，然后再进行重组。
type flag int8

// flushIp4
// @title       flushIp4
// @description 数据交给ip4Processor函数处理、清理、返回打包后的批次数据
// @auth        小卒  2022/08/03 10:57
// @param       keyNext int64 "下一个时间窗的时间戳，作为缓存的KEY"
// @param       keyCurr int64 "当前的时间窗的时间戳，作为缓存的KEY"
// @param       keyPrev int64 "上一个时间窗的时间戳，作为缓存的KEY"
// @return      r       map[ProtocolClass]interface{} "打包后的批次"
func (p *PacketProcessor) flushIp4(keyNext, keyCurr, keyPrev int64) map[ProtocolClass]interface{} {
	prev := p.ip4PacketCache.GetBatch(keyPrev)
	curr := p.ip4PacketCache.GetBatch(keyCurr)
	next := p.ip4PacketCache.GetBatch(keyNext)
	ip4 := p.ip4Processor(prev, curr, next)
	p.ip4PacketCache.DelBatch(keyPrev)
	return ip4
}

// LayerIp4 使用gopacket的layers.Ipv4结构
type LayerIp4 struct {
	layer *layers.IPv4
}

func NewLayerIp4(l *layers.IPv4) *LayerIp4 {
	return &LayerIp4{
		layer: l,
	}
}

func (l *LayerIp4) GetSrcIp() string {
	return l.layer.SrcIP.String()
}

func (l *LayerIp4) GetDstIp() string {
	return l.layer.DstIP.String()
}

func (l *LayerIp4) GetPktSz() uint16 {
	return l.layer.Length
}
