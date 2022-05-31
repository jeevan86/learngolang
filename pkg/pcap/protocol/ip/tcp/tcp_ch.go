package tcp

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/ip/base"
)

type Channel struct {
	SrcIp, DstIp     string
	SrcPort, DstPort uint16
}

func (o *Channel) ToString() string {
	return fmt.Sprintf("%s:%d->%s:%d", o.SrcIp, o.SrcPort, o.DstIp, o.DstPort)
}

func reverseChannel(ip base.LayerIp, tcp *layers.TCP) *Channel {
	return &Channel{
		SrcIp:   ip.GetDstIp(),
		DstIp:   ip.GetSrcIp(),
		SrcPort: uint16(tcp.DstPort),
		DstPort: uint16(tcp.SrcPort),
	}
}

func forwardChannel(ip base.LayerIp, tcp *layers.TCP) *Channel {
	return &Channel{
		SrcIp:   ip.GetSrcIp(),
		DstIp:   ip.GetDstIp(),
		SrcPort: uint16(tcp.SrcPort),
		DstPort: uint16(tcp.DstPort),
	}
}
