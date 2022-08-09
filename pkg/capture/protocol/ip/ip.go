package ip

import (
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
)

var ipPacketProcessor *base.PacketProcessor

// NewPacketProcessor
// @title       创建包处理器结构指针
// @description 创建包处理器结构指针，并设置判断是否本地IP的函数，供其他包使用。
// @auth        小卒    2022/08/03 10:57
// @return      r      *base.PacketProcessor  "包处理器结构指针"
func NewPacketProcessor() *base.PacketProcessor {
	if ipPacketProcessor == nil {
		ipPacketProcessor = base.NewPacketProcessor(processIp4Packets, processIp6Packets)
	}
	isLocal := func(ip string) bool {
		return ipPacketProcessor.IsLocalIp(ip)
	}
	base.SetCheckLocalIpFunc(isLocal)
	return ipPacketProcessor
}

func packetBatches(clz base.ProtocolClass,
	prev, curr, next base.ProtocolBatch) (base.PacketBatch, base.PacketBatch, base.PacketBatch) {
	return prev[clz], curr[clz], next[clz]
}
