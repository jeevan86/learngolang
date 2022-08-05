package grpc

import (
	"github.com/jeevan86/learngolang/pkg/collect/server/backend"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/cmdb"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/k8s"
)

// chIpPortMeta
// @title       获取远端和目标两端的IP端口的信息
// @description 获取远端和目标两端的IP端口的信息，先从K8S获取，没有的话在从CMDB获取
// @auth        小卒     2022/08/03 10:57
// @param       ver             backend.IpProtocol  "协议类型（tcp4、udp6...等）"
// @param       srcIp,dstIp     string              "源和目标的IP"
// @param       srcPort,dstPort int32               "源和目标的端口"
// @return      meta            *backend.ChannelPacketMeta "远端和目标两端的IP端口的信息"
func chIpPortMeta(ver backend.IpProtocol, srcIp, dstIp string, srcPort, dstPort int32) *backend.ChannelPacketMeta {
	srcMeta := getIpPortMeta(&backend.IpPort{Ip: srcIp, Port: srcPort})
	dstMeta := getIpPortMeta(&backend.IpPort{Ip: dstIp, Port: dstPort})
	return &backend.ChannelPacketMeta{
		Protocol:  ver,
		Type:      judgePacketType(srcMeta, dstMeta),
		SrcIpPort: srcMeta,
		DstIpPort: dstMeta,
	}
}

// getIpPortMeta
// @title       获取IP端口的信息
// @description 先从K8S获取，没有的话在从CMDB获取
// @auth        小卒     2022/08/03 10:57
// @param       ipPort *backend.IpPort     "IP端口"
// @return      ipPort *backend.IpPortMeta "IP端口的信息"
func getIpPortMeta(ipPort *backend.IpPort) (meta *backend.IpPortMeta) {
	if meta, _ = isPodOrSrv(ipPort); meta == nil {
		if meta, _ = isComponent(ipPort); meta == nil {
			meta, _ = isHost(ipPort)
		}
	}
	return
}

// judgePacketType
// @title       为了判断是请求流量还是响应流量
// @description 为了判断是请求流量还是响应流量
// @auth        小卒     2022/08/03 10:57
// @param       srcMeta *backend.IpPortMeta "源IP端口的信息"
// @param       dstMeta *backend.IpPortMeta "目的IP端口的信息"
// @return      backend.PacketType          "流量类型"
func judgePacketType(srcMeta, dstMeta *backend.IpPortMeta) backend.PacketType {
	pType := backend.PacketTypeUnknown
	if srcMeta != nil && isServer(srcMeta.Type) {
		// 服务端发出的（源端是服务端），则认为是响应包
		pType = backend.PacketTypeRes
	} else if dstMeta != nil && isServer(dstMeta.Type) {
		// 发向服务端的（目标是服务端），则认为是请求包
		pType = backend.PacketTypeReq
	}
	return pType
}

// isServer
// @title       这个IP端口是否是已注册的服务端
// @description 根据类型判断这个IP端口是否是已注册的服务端
// @auth        小卒    2022/08/03 10:57
// @param       tp backend.IpPortType "IP端口的类型"
// @return      b  bool               "是否服务端"
func isServer(tp backend.IpPortType) bool {
	return tp == backend.IpPortTypeK8sSrv || tp == backend.IpPortTypeComponent
}

// isPodOrSrv
// @title       是否是K8S中的PodIp、ServiceIp、NodeIp
// @description 是否是K8S中的PodIp、ServiceIp、NodeIp
// @auth        小卒    2022/08/03 10:57
// @param       ipPort *backend.IpPort     "IP端口"
// @return      ipPort *backend.IpPortMeta "IP端口的信息"
// @return      bool                       "是否有查到数据"
func isPodOrSrv(ipPort *backend.IpPort) (*backend.IpPortMeta, bool) {
	if meta, exists := k8s.GetPodIpMeta(ipPort.Ip); exists {
		tags := make(backend.Tags, 8)
		tags["clusterId"] = meta.ClusterId
		tags["namespace"] = meta.Namespace
		tags["appName"] = meta.AppName
		tags["appKind"] = meta.AppKind
		tags["hostIp"] = meta.HostIp
		tags["podName"] = meta.PodName
		return &backend.IpPortMeta{
			IpPort: *ipPort,
			Type:   backend.IpPortTypeK8sPod,
			Tags:   tags,
		}, true
	}
	if meta, exists := k8s.GetSrvIpMeta(ipPort.Ip); exists {
		tags := make(backend.Tags, 8)
		tags["clusterId"] = meta.ClusterId
		tags["namespace"] = meta.Namespace
		tags["appName"] = meta.AppName
		return &backend.IpPortMeta{
			IpPort: *ipPort,
			Type:   backend.IpPortTypeK8sSrv,
			Tags:   tags,
		}, true
	}
	if meta, exists := k8s.GetNodeIpMeta(ipPort.Ip); exists {
		tags := make(backend.Tags, 4)
		tags["clusterId"] = meta.ClusterId
		tags["nodeName"] = meta.Name
		return &backend.IpPortMeta{
			IpPort: *ipPort,
			Type:   backend.IpPortTypeK8sNode,
			Tags:   tags,
		}, true
	}
	return nil, false
}

// isComponent
// @title       是否是CMDB中的组件实例的IP端口
// @description 是否是CMDB中的组件实例的IP端口
// @auth        小卒    2022/08/03 10:57
// @param       ipPort *backend.IpPort     "IP端口"
// @return      ipPort *backend.IpPortMeta "IP端口的信息"
// @return      bool                       "是否有查到数据"
func isComponent(ipPort *backend.IpPort) (*backend.IpPortMeta, bool) {
	if meta, exists := cmdb.GetIpPortMeta(ipPort); exists {
		tags := make(backend.Tags, 8)
		tags["compId"] = meta.CiId
		tags["compName"] = meta.CiName
		tags["compTypeId"] = meta.CiTypeId
		tags["compTypeName"] = meta.CiTypeName
		return &backend.IpPortMeta{
			IpPort: *ipPort,
			Type:   backend.IpPortTypeComponent,
			Tags:   tags,
		}, true
	}
	return nil, false
}

// isHost
// @title       是否是CMDB中的宿主机的IP和端口
// @description 是否是CMDB中的宿主机的IP和端口
// @auth        小卒    2022/08/03 10:57
// @param       ipPort *backend.IpPort     "IP端口"
// @return      ipPort *backend.IpPortMeta "IP端口的信息"
// @return      bool                       "是否有查到数据"
func isHost(ipPort *backend.IpPort) (*backend.IpPortMeta, bool) {
	if meta, exists := cmdb.GetIpMeta(ipPort); exists {
		tags := make(backend.Tags, 8)
		tags["hostId"] = meta.CiId
		tags["hostName"] = meta.CiName
		return &backend.IpPortMeta{
			IpPort: *ipPort,
			Type:   backend.IpPortTypeHost,
			Tags:   tags,
		}, true
	}
	return nil, false
}
