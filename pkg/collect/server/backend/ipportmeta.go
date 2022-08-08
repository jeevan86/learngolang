package backend

import (
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/cmdb"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/k8s"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/types"
)

// ChIpPortMeta
// @title       获取远端和目标两端的IP端口的信息
// @description 获取远端和目标两端的IP端口的信息，先从K8S获取，没有的话在从CMDB获取
// @auth        小卒     2022/08/03 10:57
// @param       ver             backend.IpProtocol  "协议类型（tcp4、udp6...等）"
// @param       srcIp,dstIp     string              "源和目标的IP"
// @param       srcPort,dstPort int32               "源和目标的端口"
// @return      meta            *backend.ChannelPacketMeta "远端和目标两端的IP端口的信息"
func ChIpPortMeta(ver types.IpProtocol, srcIp, dstIp string, srcPort, dstPort int32) *types.ChannelPacketMeta {
	srcMeta := GetIpPortMeta(&types.IpPort{Ip: srcIp, Port: srcPort})
	dstMeta := GetIpPortMeta(&types.IpPort{Ip: dstIp, Port: dstPort})
	return &types.ChannelPacketMeta{
		Protocol:  ver,
		Type:      JudgePacketType(srcMeta, dstMeta),
		SrcIpPort: srcMeta,
		DstIpPort: dstMeta,
	}
}

// GetIpPortMeta
// @title       获取IP端口的信息
// @description 先从K8S获取，没有的话在从CMDB获取
// @auth        小卒     2022/08/03 10:57
// @param       ipPort *backend.IpPort     "IP端口"
// @return      ipPort *backend.IpPortMeta "IP端口的信息"
func GetIpPortMeta(ipPort *types.IpPort) (meta *types.IpPortMeta) {
	meta, _ = k8s.IsPodOrSrv(ipPort)
	if meta != nil {
		return
	}
	meta, _ = cmdb.IsComponent(ipPort)
	if meta != nil {
		return
	}
	meta, _ = cmdb.IsHost(ipPort)
	if meta != nil {
		return
	}
	meta = nil
	return
}

// JudgePacketType
// @title       为了判断是请求流量还是响应流量
// @description 为了判断是请求流量还是响应流量
// @auth        小卒     2022/08/03 10:57
// @param       srcMeta *backend.IpPortMeta "源IP端口的信息"
// @param       dstMeta *backend.IpPortMeta "目的IP端口的信息"
// @return      backend.PacketType          "流量类型"
func JudgePacketType(srcMeta, dstMeta *types.IpPortMeta) types.PacketType {
	pType := types.PacketTypeUnknown
	if srcMeta != nil && IsServer(srcMeta.Type) {
		// 服务端发出的（源端是服务端），则认为是响应包
		pType = types.PacketTypeRes
	} else if dstMeta != nil && IsServer(dstMeta.Type) {
		// 发向服务端的（目标是服务端），则认为是请求包
		pType = types.PacketTypeReq
	}
	return pType
}

// IsServer
// @title       这个IP端口是否是已注册的服务端
// @description 根据类型判断这个IP端口是否是已注册的服务端
// @auth        小卒    2022/08/03 10:57
// @param       tp backend.IpPortType "IP端口的类型"
// @return      b  bool               "是否服务端"
func IsServer(tp types.IpPortType) bool {
	return tp == types.IpPortTypeK8sSrv || tp == types.IpPortTypeComponent
}
