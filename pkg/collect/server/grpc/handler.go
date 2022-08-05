package grpc

import (
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend"
)

func ip4(ipData *pb.Protocol) {
	ip(backend.IpFamilyIp4, ipData)
}

func ip6(ipData *pb.Protocol) {
	ip(backend.IpFamilyIp6, ipData)
}

func ip(ipFamily backend.IpFamily, ipData *pb.Protocol) {
	if ipData != nil {
		// TODO: ignore icmp/igmp now
		tcp(ipFamily, ipData.Tcp)
		udp(ipFamily, ipData.Udp)
	}
}

func tcp(ipFamily backend.IpFamily, tcpData []*pb.Tcp) {
	var tcpVer backend.IpProtocol
	if ipFamily == backend.IpFamilyIp4 {
		tcpVer = backend.IpProtocolTcp4
	} else if ipFamily == backend.IpFamilyIp6 {
		tcpVer = backend.IpProtocolTcp6
	}
	for _, t := range tcpData {
		ipPortMeta := chIpPortMeta(
			tcpVer,
			t.SourceIpAddr,
			t.TargetIpAddr,
			t.SourcePort,
			t.TargetPort,
		)
		consumer.Apply(ipPortMeta)
	}
}

func udp(ipFamily backend.IpFamily, udpData []*pb.Udp) {
	var udpVer backend.IpProtocol
	if ipFamily == backend.IpFamilyIp4 {
		udpVer = backend.IpProtocolUdp4
	} else if ipFamily == backend.IpFamilyIp6 {
		udpVer = backend.IpProtocolUdp6
	}
	for _, u := range udpData {
		ipPortMeta := chIpPortMeta(
			udpVer,
			u.SourceIpAddr,
			u.TargetIpAddr,
			u.SourcePort,
			u.TargetPort,
		)
		consumer.Apply(ipPortMeta)
	}
}
