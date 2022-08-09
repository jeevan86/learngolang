package udp

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
	"github.com/jeevan86/learngolang/pkg/log"
)

var logger = log.NewLogger()

type Channel struct {
	SrcIp, DstIp     string
	SrcPort, DstPort int16
}

func (o *Channel) ToString() string {
	return fmt.Sprintf("%s:%d->%s:%d", o.SrcIp, o.SrcPort, o.DstIp, o.DstPort)
}

type AggregatedValues struct {
	base.CommonAggregatedValues
}

type ChannelAggregatedValues struct {
	Values map[Channel]*AggregatedValues
}

func aggregate(result map[Channel]*AggregatedValues, srcIp, dstIp string, pktSz uint16, udp *layers.UDP) {
	sPort := int16(udp.SrcPort)
	dPort := int16(udp.DstPort)
	channel := Channel{
		SrcIp:   srcIp,
		DstIp:   dstIp,
		SrcPort: sPort,
		DstPort: dPort,
	}
	values := getOrInit(channel, result)
	// 字节数、数量
	fillCommonValues(values, pktSz, srcIp, dstIp)
}

func getOrInit(channel Channel, result map[Channel]*AggregatedValues) *AggregatedValues {
	values, ok := result[channel]
	if !ok || values == nil {
		result[channel] = &AggregatedValues{
			CommonAggregatedValues: base.CommonAggregatedValues{
				SendCount: 0,
				SendBytes: 0,
				RecvCount: 0,
				RecvBytes: 0,
				Count:     0,
				Bytes:     0,
			},
		}
		values = result[channel]
	}
	return values
}

func fillCommonValues(values *AggregatedValues, pktSz uint16, srcIp, dstIp string) {
	// 发送的包总字节数、数量
	sBytes := values.SendBytes
	sCount := values.SendCount
	// 接收的包总字节数、数量
	rBytes := values.RecvBytes
	rCount := values.RecvCount
	// 包总字节数、数量
	bytes := values.Bytes
	count := values.Count
	isSrcLocal := base.IsLocalIp(srcIp)
	isDstLocal := base.IsLocalIp(dstIp)
	if isSrcLocal && !isDstLocal {
		sBytes += uint64(pktSz)
		sCount += 1
	} else if !isSrcLocal && isDstLocal {
		rBytes += uint64(pktSz)
		rCount += 1
	}
	//ProcessTcp4Packet(ip, p)
	bytes = sBytes + rBytes
	count = sCount + rCount

	values.SendBytes = sBytes
	values.SendCount = sCount
	values.RecvBytes = rBytes
	values.RecvCount = rCount
	values.Bytes = bytes
	values.Count = count
}
