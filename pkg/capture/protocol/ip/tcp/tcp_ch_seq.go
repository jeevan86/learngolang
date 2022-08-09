package tcp

import (
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
)

type chanSeq struct {
	tcpCh Channel
	seq   uint32
}

type chanSeqMap map[chanSeq]uint32

func registerChSeq(m chanSeqMap, ip base.LayerIp, tcp *layers.TCP) {
	k := chSeqKey(ip, tcp)
	seqCnt, ok := m[chSeqKey(ip, tcp)]
	if !ok {
		m[k] = 1
	} else {
		m[k] = seqCnt + 1
	}
}

func exceptedChSeqKey(ip base.LayerIp, tcp *layers.TCP) chanSeq {
	tcpCh := Channel{
		SrcIp:   ip.GetDstIp(),       // <- 这里要换一下
		DstIp:   ip.GetSrcIp(),       // <- 这里要换一下
		SrcPort: uint16(tcp.DstPort), // <- 这里要换一下
		DstPort: uint16(tcp.SrcPort), // <- 这里要换一下
	}
	chAck := chanSeq{
		tcpCh: tcpCh,
		seq:   tcp.Seq + uint32(ip.GetPktSz()),
	}
	return chAck
}

func chSeqKey(ip base.LayerIp, tcp *layers.TCP) chanSeq {
	srcIp := ip.GetSrcIp()
	dstIp := ip.GetDstIp()
	srcPort := uint16(tcp.SrcPort)
	dstPort := uint16(tcp.DstPort)
	tcpCh := Channel{
		SrcIp:   srcIp,
		DstIp:   dstIp,
		SrcPort: srcPort,
		DstPort: dstPort,
	}
	chSeq := chanSeq{
		tcpCh: tcpCh,
		seq:   tcp.Seq,
	}
	return chSeq
}
