package base

type TagsV4 struct {
	Tags
}

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
