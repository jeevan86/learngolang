package tcp

import (
	"fmt"
	"github.com/google/gopacket/layers"
)

func printTcp4Packet(ip *layers.IPv4, tcp *layers.TCP, payload bool) {
	// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
	// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
	format := fmt.Sprintf("TCP%d-[%s:%d -> %s:%d]-[ttl:%d][ihl:%d][tos:%d][flg:%s]%s%s%s%s%s%s%s%s%s[seq:%d][ack:%d][win:%d][chk:%d][urg:%d][len:%d]",
		ip.Version,
		ip.SrcIP, tcp.SrcPort, ip.DstIP, tcp.DstPort,
		ip.TTL, ip.IHL, ip.TOS, ip.Flags.String(),
		flag(tcp.FIN, "FIN"),
		flag(tcp.SYN, "SYN"),
		flag(tcp.RST, "RST"),
		flag(tcp.PSH, "PSH"),
		flag(tcp.ACK, "ACK"),
		flag(tcp.URG, "URG"),
		flag(tcp.ECE, "ECE"),
		flag(tcp.CWR, "CWR"),
		flag(tcp.NS, "NS"),
		tcp.Seq, tcp.Ack, tcp.Window, tcp.Checksum, tcp.Urgent, ip.Length)
	if logger.IsDebugEnabled() {
		logger.Debug(format)
	} else {
		logger.Info(format)
	}
}

// 表示此次发送的数据在整个报文段中的起始字节数。
// 序号是32 bit的无符号数。为了安全起见，它的初始值是一个随机生成的数，它到达32位最大值后，又从零开始。
type seq uint32

// ack_seq 指定的是对方所期望接收的字节。
type ack uint32

// window
//  是16位滑动窗口的大小，单位为字节，起始于确认序列号字段指明的值，这个值是接收端正期望接收的字节数，其最大值是63353字节。
//  TCP中的流量控制是通过一个可变大小的滑动窗口来完成的。window域指定了从被确认的字节算起可以接收的多少个字节。
//  window = 0也是合法的，这相当于说，到现在为止多达ack_seq-1个字节已经接收到了，但是接收方现在状态不佳，需要休息一下，
//  等一会儿再继续接收更多的数据，谢谢。以后，接收方可以通过发送一个同样ack_seq但是window不为0的数据段，告诉发送方继续发送数据段。
type win uint8

func printTcp6Packet(ip *layers.IPv6, tcp *layers.TCP, payload bool) {
	// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
	// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
	format := fmt.Sprintf("TCP%d-[%s:%d -> %s:%d]-%s%s%s%s%s%s%s%s%s[seq:%d][len:%d]",
		ip.Version,
		ip.SrcIP, tcp.SrcPort, ip.DstIP, tcp.DstPort,
		flag(tcp.FIN, "FIN"),
		flag(tcp.SYN, "SYN"),
		flag(tcp.RST, "RST"),
		flag(tcp.PSH, "PSH"),
		flag(tcp.ACK, "ACK"),
		flag(tcp.URG, "URG"),
		flag(tcp.ECE, "ECE"),
		flag(tcp.CWR, "CWR"),
		flag(tcp.NS, "NS"),
		tcp.Seq, ip.Length)
	logger.Info(format)
}
