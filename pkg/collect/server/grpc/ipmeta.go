package grpc

import (
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/cmdb"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/k8s"
)

type IpPortMetaConsumer interface {
	Apply(*backend.ChannelPacketMeta)
}

func ip4(ipData *pb.Protocol, consumer IpPortMetaConsumer) {
	ip(backend.IpFamilyIp4, ipData, consumer)
}

func ip6(ipData *pb.Protocol, consumer IpPortMetaConsumer) {
	ip(backend.IpFamilyIp6, ipData, consumer)
}

func ip(ipFamily backend.IpFamily, ipData *pb.Protocol, consumer IpPortMetaConsumer) {
	if ipData != nil {
		// TODO: ignore icmp/igmp now
		tcp(ipFamily, ipData.Tcp, consumer)
		udp(ipFamily, ipData.Udp, consumer)
	}
}

func tcp(ipFamily backend.IpFamily, tcpData []*pb.Tcp, consumer IpPortMetaConsumer) {
	var tcpVer backend.IpProtocol
	if ipFamily == backend.IpFamilyIp4 {
		tcpVer = backend.IpProtocolTcp4
	} else if ipFamily == backend.IpFamilyIp6 {
		tcpVer = backend.IpProtocolTcp6
	}
	for _, t := range tcpData {
		ipPortMeta := tcpMeta(t)
		ipPortMeta.Protocol = tcpVer
		consumer.Apply(ipPortMeta)
	}
}

func udp(ipFamily backend.IpFamily, udpData []*pb.Udp, consumer IpPortMetaConsumer) {
	var udpVer backend.IpProtocol
	if ipFamily == backend.IpFamilyIp4 {
		udpVer = backend.IpProtocolUdp4
	} else if ipFamily == backend.IpFamilyIp6 {
		udpVer = backend.IpProtocolUdp6
	}
	for _, u := range udpData {
		ipPortMeta := udpMeta(u)
		ipPortMeta.Protocol = udpVer
		consumer.Apply(ipPortMeta)
	}
}

func tcpMeta(tcp *pb.Tcp) *backend.ChannelPacketMeta {
	srcIp := tcp.SourceIpAddr
	srcPort := tcp.SourcePort
	dstIp := tcp.TargetIpAddr
	dstPort := tcp.TargetPort
	return ipMeta(srcIp, dstIp, srcPort, dstPort)
}

func udpMeta(udp *pb.Udp) *backend.ChannelPacketMeta {
	srcIp := udp.SourceIpAddr
	srcPort := udp.SourcePort
	dstIp := udp.TargetIpAddr
	dstPort := udp.TargetPort
	return ipMeta(srcIp, dstIp, srcPort, dstPort)
}

func ipMeta(srcIp, dstIp string, srcPort, dstPort int32) *backend.ChannelPacketMeta {
	srcMeta := getIpPortMeta(srcIp, srcPort)
	dstMeta := getIpPortMeta(dstIp, dstPort)

	pType := backend.PacketTypeUnknown
	if srcMeta != nil && isServer(srcMeta.Type) {
		// 服务端发出的（源端是服务端），则认为是响应包
		pType = backend.PacketTypeRes
	} else if dstMeta != nil && isServer(dstMeta.Type) {
		// 发向服务端的（目标是服务端），则认为是请求包
		pType = backend.PacketTypeReq
	}

	return &backend.ChannelPacketMeta{
		Type:    pType,
		SrcIp:   srcIp,
		SrcPort: srcPort,
		SrcMeta: srcMeta,
		DstIp:   dstIp,
		DstPort: dstPort,
		DstMeta: dstMeta,
	}
}

func isServer(tp backend.IpPortType) bool {
	return tp == backend.IpPortTypeK8sSrv || tp == backend.IpPortTypeComponent
}

func getIpPortMeta(ip string, port int32) (meta *backend.IpPortMeta) {
	if meta, _ = ipIsPodOrSrv(ip); meta == nil {
		meta, _ = isIpPortComp(ip, port)
	}
	return
}

func ipIsPodOrSrv(ip string) (*backend.IpPortMeta, bool) {
	if meta, exists := k8s.GetPodIpMeta(ip); exists {
		tags := make(backend.Tags, 8)
		tags["clusterId"] = meta.ClusterId
		tags["namespace"] = meta.Namespace
		tags["appName"] = meta.AppName
		tags["appKind"] = meta.AppKind
		tags["hostIp"] = meta.HostIp
		tags["podName"] = meta.PodName
		return newIpPortMeta(backend.IpPortTypeK8sPod, tags), true
	}
	if meta, exists := k8s.GetSrvIpMeta(ip); exists {
		tags := make(backend.Tags, 8)
		tags["clusterId"] = meta.ClusterId
		tags["namespace"] = meta.Namespace
		tags["appName"] = meta.AppName
		return newIpPortMeta(backend.IpPortTypeK8sSrv, tags), true
	}
	return nil, false
}

func isIpPortComp(ip string, port int32) (*backend.IpPortMeta, bool) {
	if meta, exists := cmdb.GetIpPortMeta(ip, port); exists {
		tags := make(backend.Tags, 8)
		tags["compId"] = meta.CompId
		tags["compName"] = meta.CompName
		tags["compTypeId"] = meta.CompTypeId
		tags["compTypeName"] = meta.CompTypeName
		return newIpPortMeta(backend.IpPortTypeComponent, tags), true
	}
	return nil, false
}

func newIpPortMeta(tp backend.IpPortType, tags backend.Tags) *backend.IpPortMeta {
	return &backend.IpPortMeta{Type: tp, Tags: tags}
}
