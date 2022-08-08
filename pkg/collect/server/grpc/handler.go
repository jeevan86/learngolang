package grpc

import (
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/types"
	"github.com/jeevan86/learngolang/pkg/collect/server/output"
)

var consumer = output.NewConsumer()

func ip4(ipData *pb.Protocol) {
	ip(types.IpFamilyIp4, ipData)
}

func ip6(ipData *pb.Protocol) {
	ip(types.IpFamilyIp6, ipData)
}

func ip(ipFamily types.IpFamily, ipData *pb.Protocol) {
	if ipData != nil {
		// TODO: ignore icmp/igmp now
		tcp(ipFamily, ipData.Tcp)
		udp(ipFamily, ipData.Udp)
	}
}

func tcp(ipFamily types.IpFamily, tcpData []*pb.Tcp) {
	v := tcpVer(ipFamily)
	for _, t := range tcpData {
		each(
			v,
			t.SourceIpAddr,
			t.TargetIpAddr,
			t.SourcePort,
			t.TargetPort,
		)
	}
}

func tcpVer(ipFamily types.IpFamily) types.IpProtocol {
	var v types.IpProtocol
	if ipFamily == types.IpFamilyIp4 {
		v = types.IpProtocolTcp4
	} else if ipFamily == types.IpFamilyIp6 {
		v = types.IpProtocolTcp6
	}
	return v
}

func udp(ipFamily types.IpFamily, udpData []*pb.Udp) {
	v := udpVer(ipFamily)
	for _, u := range udpData {
		each(
			v,
			u.SourceIpAddr,
			u.TargetIpAddr,
			u.SourcePort,
			u.TargetPort,
		)
	}
}

func udpVer(ipFamily types.IpFamily) types.IpProtocol {
	var v types.IpProtocol
	if ipFamily == types.IpFamilyIp4 {
		v = types.IpProtocolUdp4
	} else if ipFamily == types.IpFamilyIp6 {
		v = types.IpProtocolUdp6
	}
	return v
}

func each(ver types.IpProtocol, srcIp, dstIp string, srcPort, dstPort int32) {
	consumer.Apply(
		backend.ChIpPortMeta(ver, srcIp, dstIp, srcPort, dstPort),
	)
}
