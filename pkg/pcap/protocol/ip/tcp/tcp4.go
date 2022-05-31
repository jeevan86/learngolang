// Package tcp
//   header->source：16位源端口号
//   header->dest：16位目的端口号
//   header->seq：表示此次发送的数据在整个报文段中的起始字节数。序号是32 bit的无符号数。为了安全起见，它的初始值是一个随机生成的数，它到达32位最大值后，又从零开始。
//   header->ack_seq：指定的是对方所期望接收的字节。
//   header->doff：TCP头长度，指明了在TCP头部包含多少个32位的字（即单位为4字节）。此信息是必须的，因为options域的长度是可变的，所以整个TCP头部的长度也是变化的。从技术上讲，这个域实际上指明了数据部分在段内部的其起始地址(以32位字作为单位进行计量)，因为这个数值正好是按字为单位的TCP头部的长度，所以，二者的效果是等同的。
//   header->res1：为保留位。
//   header->window：是16位滑动窗口的大小，单位为字节，起始于确认序列号字段指明的值，这个值是接收端正期望接收的字节数，其最大值是63353字节。TCP中的流量控制是通过一个可变大小的滑动窗口来完成的。window域指定了从被确认的字节算起可以接收的多少个字节。window = 0也是合法的，这相当于说，到现在为止多达ack_seq-1个字节已经接收到了，但是接收方现在状态不佳，需要休息一下，等一会儿再继续接收更多的数据，谢谢。以后，接收方可以通过发送一个同样ack_seq但是window不为0的数据段，告诉发送方继续发送数据段。
//   header->check：检验和，覆盖了整个的TCP报文段，这是一个强制性的字段，一定是由发送端计算和存储，并由接收端进行验证。
//   header->urg_ptr：这个域被用来指示紧急数据在当前数据段中的位置，它是一个相对于当前序列号的字节偏移值。这个设施可以代替中断信息。
//
//  fin, syn, rst, psh, ack, urg为6个标志位这6个位域已经保留了超过四分之一个世纪的时间而仍然原封未动，
//  这样的事实正好也说明了TCP的设计者们考虑的是多么的周到。它们的含义如下：
//   header->fin：被用于释放一个连接。它表示发送方已经没有数据要传输了。
//   header->syn：同步序号，用来发起一个连接。syn位被用于建立连接的过程。在连接请求中，syn=1; ack=0表示该数据段没有使用捎带的确认域。连接应答捎带了一个确认，所以有syn=1; ack=1。本质上，syn位被用来表示connection request和connection accepted，然而进一步用ack位来区分这两种情况。
//   header->rst：该为用于重置一个已经混乱的连接，之所以会混乱，可能是由于主机崩溃，或者其他的原因。该位也可以被用来拒绝一个无效的数据段，或者拒绝一个连接请求。一般而言，如果你得到的数据段设置了rst位，那说明你这一端有了问题。
//   header->psh：接收方在收到数据后应立即请求将数据递交给应用程序，而不是将它缓冲起来直到整个缓冲区接收满为止(这样做的目的可能是为了效率的原因)。
//   header->ack：ack位被设置为1表示header->ack_seq是有效的，需要check ack_seq字段是否为所需要。如果ack为0，则该数据段不包含确认信息，所以，header->ack_seq域应该被忽略。
//   header->urg：紧急指针有效。
//   header->ece：用途暂时不明。
//   header->cwr：用途暂时不明。
// 内核源代码在函数tcp_transmit_skb()中建立tcp首部。
package tcp

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/ip/base"
)

func ProcessTcp4Packets(prev, curr, next base.PacketBatch) *ChannelAggregatedValues {
	return processPackets(prev, curr, next, v4IpTcpLayer)
}

func v4IpTcpLayer(item *base.PacketItem) (base.LayerIp, *layers.TCP) {
	p := item.Packet
	ip4 := p.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	tcp := p.Layer(layers.LayerTypeTCP).(*layers.TCP)
	ip := base.NewLayerIp4(ip4)
	return ip, tcp
}

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
