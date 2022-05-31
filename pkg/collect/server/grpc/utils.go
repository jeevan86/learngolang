package grpc

import (
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/k8s"
)

type IpType int16

const (
	IpTypeNode IpType = 0
	IpTypePod  IpType = 1
	IpTypeSrv  IpType = 2
)

type IpProtocol int16

const (
	IpProtocolTcp4 IpProtocol = 11
	IpProtocolUdp4 IpProtocol = 12
	IpProtocolTcp6 IpProtocol = 21
	IpProtocolUdp6 IpProtocol = 22
)

type K8sIpMeta struct {
	ClusterId string            `json:"cluster-id,omitempty"`
	IpType    IpType            `json:"ip-type,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
}

type IpPortMeta struct {
	Protocol  IpProtocol `json:"protocol,omitempty"`
	SrcIp     string     `json:"src-ip,omitempty"`
	SrcIpMeta *K8sIpMeta `json:"src-ip-meta,omitempty"`
	SrcPort   int32      `json:"src-port,omitempty"`
	DstIp     string     `json:"dst-ip,omitempty"`
	DstIpMeta *K8sIpMeta `json:"dst-ip-meta,omitempty"`
	DstPort   int32      `json:"dst-port,omitempty"`
}

func tcpMeta(tcp *pb.Tcp) *IpPortMeta {
	srcIp := tcp.SourceIpAddr
	srcPort := tcp.SourcePort
	dstIp := tcp.TargetIpAddr
	dstPort := tcp.TargetPort
	return ipMeta(srcIp, dstIp, srcPort, dstPort)
}

func udpMeta(udp *pb.Udp) *IpPortMeta {
	srcIp := udp.SourceIpAddr
	srcPort := udp.SourcePort
	dstIp := udp.TargetIpAddr
	dstPort := udp.TargetPort
	return ipMeta(srcIp, dstIp, srcPort, dstPort)
}

func ipMeta(srcIp, dstIp string, srcPort, dstPort int32) *IpPortMeta {
	srcIpMeta := ipIsPodOrSrv(srcIp)
	dstIpMeta := ipIsPodOrSrv(dstIp)
	return &IpPortMeta{
		SrcIp:     srcIp,
		SrcIpMeta: srcIpMeta,
		SrcPort:   srcPort,
		DstIp:     dstIp,
		DstIpMeta: dstIpMeta,
		DstPort:   dstPort,
	}
}

func ipIsPodOrSrv(ip string) *K8sIpMeta {
	if meta, exists := k8s.GetPodIpMeta(ip); exists {
		podIpMeta := &K8sIpMeta{
			ClusterId: meta.ClusterId,
			IpType:    IpTypePod,
			Tags:      make(map[string]string, 8),
		}
		podIpMeta.Tags["namespace"] = meta.Namespace
		podIpMeta.Tags["hostIp"] = meta.HostIp
		podIpMeta.Tags["appKind"] = meta.AppKind
		podIpMeta.Tags["appName"] = meta.AppName
		podIpMeta.Tags["podName"] = meta.PodName
		return podIpMeta
	} else if meta, exists := k8s.GetSrvIpMeta(ip); exists {
		srvIpMeta := &K8sIpMeta{
			ClusterId: meta.ClusterId,
			IpType:    IpTypeSrv,
			Tags:      make(map[string]string, 4),
		}
		srvIpMeta.Tags["namespace"] = meta.Namespace
		srvIpMeta.Tags["appName"] = meta.AppName
		return srvIpMeta
	}
	return nil
}
