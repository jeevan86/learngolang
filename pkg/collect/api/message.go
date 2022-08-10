package api

type Ip struct {
	SourceIpAddr       string `json:"source_ip_addr"`       // 源端ip
	TargetIpAddr       string `json:"target_ip_addr"`       // 目标ip
	SendTotalByte      uint64 `json:"send_total_byte"`      // 发送字节数
	ReceiveTotalByte   uint64 `json:"receive_total_byte"`   // 接收字节数
	SendTotalPacket    uint64 `json:"send_total_packet"`    // 客户端数据包数
	ReceiveTotalPacket uint64 `json:"receive_total_packet"` // 服务器数据包数
	TotalPacket        uint64 `json:"total_packet"`         // 总数据包量
	TotalBytes         uint64 `json:"total_bytes"`          // 总数据包字节数
}

type Tcp struct {
	Ip
	SourcePort     uint16 `json:"source_port"`       // 源端口
	TargetPort     uint16 `json:"target_port"`       // 目标端口
	SynCount       uint64 `json:"syn_count"`         // 连接syn包总数
	SynAckCount    uint64 `json:"syn_ack_count"`     // 连接syn-ack包总数
	SynAckAckCount uint64 `json:"syn_ack_ack_count"` // 连接syn-ack-ack包总数
	FinCount       uint64 `json:"fin_count"`         // 断开fin包总数
	FinAckCount    uint64 `json:"fin_ack_count"`     // 断开fin-ack包总数
	AckCount       uint64 `json:"ack_count"`         // 普通ack包总数
	ResetCount     uint64 `json:"reset_count"`       // 连接reset总数
	Retransmit     uint64 `json:"retransmit"`        // 重传的总次数
	Rtt            int64  `json:"rtt"`               // tcp 套接字的平均往返时间
}

type Udp struct {
	Ip
	SourcePort int16 `json:"source_port"` // 源端口
	TargetPort int16 `json:"target_port"` // 目标端口
}

type Icmp struct {
	Ip
}

type Igmp struct {
	Ip
}

type Protocol struct {
	Tcp  []*Tcp  `json:"tcp"`
	Udp  []*Udp  `json:"udp"`
	Icmp []*Icmp `json:"icmp"`
	Igmp []*Igmp `json:"igmp"`
}

type NetStatics struct {
	Timestamp int64     `json:"timestamp"`
	Ip6       *Protocol `json:"ip6"`
	Ip4       *Protocol `json:"ip4"`
	GatherIp  string    `json:"gather_ip"`
}

type LocalIpLst struct {
	Data IpList
}

type IpList []string
