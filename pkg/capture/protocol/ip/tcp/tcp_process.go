package tcp

import (
	"github.com/google/gopacket/layers"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
)

// ipTcpLayerFunc
type ipTcpLayerFunc func(item *base.PacketItem) (base.LayerIp, *layers.TCP)

func processPackets(
	prev, curr, next base.PacketBatch, f ipTcpLayerFunc) *ChannelAggregatedValues {
	seqMap := make(chanSeqMap, len(curr))
	ackMap := make(chanAckMap, (len(prev)+len(curr)+len(next))/2)
	prepareAll(seqMap, ackMap, f, prev, curr, next)
	return &ChannelAggregatedValues{
		Values: aggregateCurr(curr, seqMap, ackMap, f),
	}
}

func prepareAll(seqMap chanSeqMap, ackMap chanAckMap, f ipTcpLayerFunc, all ...base.PacketBatch) {
	for _, batch := range all {
		for _, item := range batch {
			prepare(seqMap, ackMap, item, f)
		}
	}
}

func prepare(seqMap chanSeqMap, ackMap chanAckMap, item *base.PacketItem, f ipTcpLayerFunc) {
	ip, tcp := f(item)
	if tcp.ACK {
		// 包括 Syn+Ack、Fin+Ack
		registerChAck(ackMap, ip, tcp, item.Millis)
	}
	registerChSeq(seqMap, ip, tcp)
}

func aggregateCurr(curr base.PacketBatch,
	seqMap chanSeqMap, ackMap chanAckMap, f ipTcpLayerFunc) map[Channel]*AggregatedValues {
	result := make(map[Channel]*AggregatedValues, 32)
	for _, item := range curr {
		ip, tcp := f(item)
		chSeq := chSeqKey(ip, tcp)
		aggregate(
			result,
			seqMap,
			ackMap,
			&chSeq,
			item,
			ip,
			tcp,
		)
	}
	return result
}

// aggregate
// @title       aggregate
// @description 使用cacheCreator创建缓存
// @auth        小卒     2022/08/03 10:57
// @param       result map[Channel]*AggregatedValues "按Channel聚合处理的结果"
// @param       seqMap chanSeqMap                    "chanSeqMap"
// @param       ackMap chanAckMap                    "chanAckMap"
// @param       chSeq  *chanSeq                      "*chanSeq"
// @param       item   *base.PacketItem              "*base.PacketItem"
// @param       ip     base.LayerIp                  "base.LayerIp"
// @param       tcp    *layers.TCP                   "tcp数据结构"
func aggregate(result map[Channel]*AggregatedValues, seqMap chanSeqMap, ackMap chanAckMap,
	chSeq *chanSeq, item *base.PacketItem, ip base.LayerIp, tcp *layers.TCP) {
	values := getOrInit(&chSeq.tcpCh, result)
	// 字节数、数量
	fillCommonValues(values, ip.GetPktSz(), &chSeq.tcpCh)
	// 连接
	fillConnectValues(values, seqMap, ip, tcp)
	// 重置
	fillResetValues(values, tcp)
	// 重传
	fillRetransmitValues(values, chSeq, seqMap)
	// rtt
	fillRttValues(values, item, ip, tcp, ackMap)
	// 关闭
	fillCloseValues(values, tcp)
}

// getOrInit
// @title       getOrInit
// @description 根据channel从result中获得或者初始化一个AggregatedValues结构指针
// @auth        小卒    2022/08/03 10:57
// @param       channel Channel                       "tcp通道"
// @param       result  map[Channel]*AggregatedValues "聚合处理的结果"
// @return      r       *AggregatedValues             "聚合结构的指针"
func getOrInit(channel *Channel, result map[Channel]*AggregatedValues) *AggregatedValues {
	ch := *channel
	values, ok := result[ch]
	if !ok || values == nil {
		result[ch] = &AggregatedValues{
			CommonAggregatedValues: base.CommonAggregatedValues{
				SendCount: 0,
				SendBytes: 0,
				RecvCount: 0,
				RecvBytes: 0,
				Count:     0,
				Bytes:     0,
			},
			Syn:        0,
			SynAck:     0,
			SynAckAck:  0,
			Ack:        0,
			Rst:        0,
			Retransmit: 0,
			Rtt:        0,
		}
		values = result[ch]
	}
	return values
}

// fillCommonValues
// @title       fillCommonValues
// @description 计算并修改values值
// @auth        小卒     2022/08/03 10:57
// @param       values *AggregatedValues "聚合结构的指针"
// @param       pktSz  uint16            "包大小"
// @param       ch     *Channel          "tcp通道"
func fillCommonValues(values *AggregatedValues, pktSz uint16, ch *Channel) {
	// 发送的包总字节数、数量
	sBytes := values.SendBytes
	sCount := values.SendCount
	// 接收的包总字节数、数量
	rBytes := values.RecvBytes
	rCount := values.RecvCount
	// 包总字节数、数量
	bytes := values.Bytes
	count := values.Count
	isSrcLocal := base.IsLocalIp(ch.SrcIp)
	isDstLocal := base.IsLocalIp(ch.DstIp)
	if isSrcLocal && !isDstLocal {
		sBytes += uint64(pktSz)
		sCount += 1
	} else if !isSrcLocal && isDstLocal {
		rBytes += uint64(pktSz)
		rCount += 1
	}

	bytes = sBytes + rBytes
	count = sCount + rCount

	values.SendBytes = sBytes
	values.SendCount = sCount
	values.RecvBytes = rBytes
	values.RecvCount = rCount
	values.Bytes = bytes
	values.Count = count
}

// fillConnectValues
// @title       fillConnectValues
// @description 计算并填充tcp连接相关的聚合结果
// @auth        小卒     2022/08/03 10:57
// @param       values *AggregatedValues "聚合结构的指针"
// @param       seqMap chanSeqMap                    "chanSeqMap"
// @param       ip     base.LayerIp                  "ip层信息"
// @param       tcp    *layers.TCP                   "tcp数据结构"
func fillConnectValues(values *AggregatedValues, seqMap chanSeqMap, ip base.LayerIp, tcp *layers.TCP) {
	// 第一次握手
	synCount := values.Syn
	// 第二次握手
	synAckCount := values.SynAck
	// 第三次握手
	synAckAckCount := values.SynAckAck
	// 普通的ACK
	ackCount := values.Ack
	if tcp.SYN && !tcp.ACK {
		// 第一次握手：主机A发送位码为syn＝1，随机产生seq的数据包到服务器
		//           主机B由SYN=1知道，A要求建立联机；
		synCount++
	} else if tcp.SYN && tcp.ACK {
		// 第二次握手：主机B收到请求后要确认联机信息，向A发送syn=1，ack=1的数据包
		//           并且随机产生seq、ack_seq=主机A的seq+1
		synAckCount++
	} else if !tcp.SYN && tcp.ACK && !tcp.FIN {
		// 第三次握手：主机A收到后检查ack_seq是否正确，即第一次发送的seq number+1，以及位码ack是否为1
		//           若正确，主机A会再发送ack number=(主机B的seq+1)，ack=1，主机B收到后确认seq值与ack=1则连接建立成功。
		chSeq := chanSeq{
			tcpCh: Channel{
				// 这里的源、目标要互换一下
				SrcIp:   ip.GetDstIp(),
				DstIp:   ip.GetSrcIp(),
				SrcPort: uint16(tcp.DstPort),
				DstPort: uint16(tcp.SrcPort),
			},
			seq: tcp.Seq - uint32(1),
		}
		cnt := seqMap[chSeq]
		if cnt > 0 {
			synAckAckCount++
		} else {
			// 普通-ACK
			ackCount++
		}
	}
	values.Syn = synCount
	values.SynAck = synAckCount
	values.SynAckAck = synAckAckCount
	values.Ack = ackCount
}

// fillResetValues
// @title       fillResetValues
// @description 计算并填充tcp的重置相关的聚合结果
// @auth        小卒     2022/08/03 10:57
// @param       values *AggregatedValues "聚合结构的指针"
// @param       tcp    *layers.TCP                   "tcp数据结构"
func fillResetValues(values *AggregatedValues, tcp *layers.TCP) {
	rstCount := values.Rst
	if tcp.RST {
		rstCount++
	}
	values.Rst = rstCount
}

// fillRetransmitValues
// @title       fillRetransmitValues
// @description 计算并填充tcp的重传相关的聚合结果
// @auth        小卒     2022/08/03 10:57
// @param       values *AggregatedValues "聚合结构的指针"
// @param       chSeq  *chanSeq          "*chanSeq"
// @param       seqMap chanSeqMap        "chanSeqMap"
func fillRetransmitValues(values *AggregatedValues, chSeq *chanSeq, seqMap chanSeqMap) {
	if isRetransmit(chSeq, seqMap) {
		values.Retransmit = values.Retransmit + 1
	}
}

// isRetransmit
// @title       isRetransmit
// @description 判断是否是tcp重传
// @auth        小卒    2022/08/03 10:57
// @param       chSeq  *chanSeq          "*chanSeq"
// @param       seqMap chanSeqMap        "chanSeqMap"
func isRetransmit(chSeq *chanSeq, seqMap chanSeqMap) bool {
	key := *chSeq
	retransmit := false
	seqCnt, ok := seqMap[key]
	if !ok {
		seqMap[key] = 1
	} else {
		seqMap[key] = seqCnt + 1
		retransmit = true
	}
	return retransmit
}

// fillRttValues
// @title       fillRttValues
// @description 计算并填充tcp的Rtt相关的聚合结果
// @auth        小卒     2022/08/03 10:57
// @param       values *AggregatedValues "聚合结构的指针"
// @param       item   *base.PacketItem              "Ip包信息"
// @param       ip     base.LayerIp                  "Ip层结构"
// @param       tcp    *layers.TCP                   "tcp数据结构"
// @param       ackMap chanAckMap                    "chanAckMap"
func fillRttValues(values *AggregatedValues,
	item *base.PacketItem, ip base.LayerIp, tcp *layers.TCP, ackMap chanAckMap) {
	// 从非ACK的包开始找，仅处理普通ACK，不处理连接和关闭的
	if !tcp.ACK && !tcp.SYN && !tcp.FIN {
		chAck := exceptedChAckKey(ip, tcp)
		pktTs := item.Millis
		ackTs := ackMap[chAck]
		values.Rtt = values.Rtt + (ackTs - pktTs)
	}
}

// fillCloseValues
// @title       fillCloseValues
// @description 计算并填充tcp的挥手相关的聚合结果，对于确定源、目标的tcp包，挥手只是Fin、FinAck
// @auth        小卒     2022/08/03 10:57
// @param       values  *AggregatedValues "聚合结构的指针"
// @param       tcp     *layers.TCP       "tcp数据结构"
func fillCloseValues(values *AggregatedValues, tcp *layers.TCP) {
	// Fin
	finCount := values.Fin
	// FinAck
	finAckCount := values.FinAck
	// 普通的ACK
	ackCount := values.Ack
	if tcp.FIN && !tcp.ACK {
		finCount++
	} else if tcp.ACK {
		if tcp.FIN {
			finAckCount++
		} else {
			ackCount++
		}
	}
	values.Fin = finCount
	values.FinAck = finAckCount
	values.Ack = ackCount
}
