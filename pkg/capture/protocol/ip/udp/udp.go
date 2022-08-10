package udp

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
	"github.com/jeevan86/learngolang/pkg/log"
)

var logger = log.NewLogger()

// Channel 一个逻辑上的通道，包含源目标的ip和端口
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

// aggregate
// @title       aggregate
// @description 使用cacheCreator创建缓存
// @auth        小卒     2022/08/03 10:57
// @param       result map[Channel]*AggregatedValues "按Channel聚合处理的结果"
// @param       srcIp  string                        "源Ip"
// @param       dstIp  string                        "目标IP"
// @param       pktSz  uint16                        "包大小"
// @param       udp    *layers.UDP                   "udp数据结构"
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

// getOrInit
// @title       getOrInit
// @description 根据channel从result中获得或者初始化一个AggregatedValues结构指针
// @auth        小卒    2022/08/03 10:57
// @param       channel Channel                       "UDP逻辑上的通道"
// @param       result  map[Channel]*AggregatedValues "聚合处理的结果"
// @return      r       *AggregatedValues             "聚合结构的指针"
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

// fillCommonValues
// @title       fillCommonValues
// @description 计算并修改values值
// @auth        小卒     2022/08/03 10:57
// @param       values *AggregatedValues "聚合结构的指针"
// @param       pktSz  uint16            "包大小"
// @param       srcIp  string            "源Ip"
// @param       dstIp  string            "目标IP"
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
