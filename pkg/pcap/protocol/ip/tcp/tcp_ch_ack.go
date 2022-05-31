package tcp

import (
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/ip/base"
)

type chanAck struct {
	tcpCh Channel
	ackSq uint32
}

type chanAckMap map[chanAck]int64

func registerChAck(m chanAckMap, ip base.LayerIp, tcp *layers.TCP, millis int64) {
	m[chAckKey(ip, tcp)] = millis
}

// exceptedChAckKey 预期要收到的ACK包信息
func exceptedChAckKey(ip base.LayerIp, tcp *layers.TCP) chanAck {
	chSeq := chanAck{
		tcpCh: *reverseChannel(ip, tcp),
		ackSq: tcp.Seq + uint32(ip.GetPktSz()),
	}
	return chSeq
}

func chAckKey(ip base.LayerIp, tcp *layers.TCP) chanAck {
	chAck := chanAck{
		tcpCh: Channel{
			SrcIp:   ip.GetSrcIp(),
			DstIp:   ip.GetDstIp(),
			SrcPort: uint16(tcp.SrcPort),
			DstPort: uint16(tcp.DstPort),
		},
		ackSq: tcp.Ack,
	}
	return chAck
}
