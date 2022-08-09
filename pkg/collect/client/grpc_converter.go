package client

import (
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/config"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/icmp"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/tcp"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/udp"
)

type grpcConverter string

func (c *grpcConverter) convert(m *base.OutputStruct) interface{} {
	gatherIp := *config.GetConfig().NodeIp
	converted := &pb.NetStaticsReq{
		Timestamp: m.Bucket,
		GatherIp:  gatherIp,
	}
	ip4, ip4Exists := m.Values[base.Ipv4]
	if ip4Exists {
		converted.Ip4 = &pb.Protocol{
			Tcp:  c.tcpMsg(ip4),
			Udp:  c.udpMsg(ip4),
			Icmp: c.icmpMsg(ip4),
			Igmp: c.igmpMsg(ip4),
		}
	}
	ip6, ip6Exists := m.Values[base.Ipv6]
	if ip6Exists {
		converted.Ip6 = &pb.Protocol{
			Tcp:  c.tcpMsg(ip6),
			Udp:  c.udpMsg(ip6),
			Icmp: c.icmpMsg(ip6),
			Igmp: c.igmpMsg(ip6),
		}
	}
	return converted
}

func (c *grpcConverter) tcpMsg(ip map[base.ProtocolClass]interface{}) []*pb.Tcp {
	if ip == nil {
		return nil
	}
	ipTcp, exists := ip[base.ProtocolTcp].(*tcp.ChannelAggregatedValues)
	if !exists {
		return nil
	}
	tcpSlice := make([]*pb.Tcp, 0)
	for k, v := range ipTcp.Values {
		tcpSlice = append(tcpSlice,
			&pb.Tcp{
				SourceIpAddr:       k.SrcIp,
				TargetIpAddr:       k.DstIp,
				SendTotalByte:      v.SendBytes,
				ReceiveTotalByte:   v.RecvBytes,
				SendTotalPacket:    v.SendCount,
				ReceiveTotalPacket: v.RecvCount,
				TotalPacket:        v.Count,
				TotalBytes:         v.Bytes,
				SourcePort:         int32(k.SrcPort),
				TargetPort:         int32(k.DstPort),
				SynCount:           v.Syn,
				SynAckCount:        v.SynAck,
				SynAckAckCount:     v.SynAckAck,
				FinCount:           v.Fin,
				FinAckCount:        v.FinAck,
				AckCount:           v.Ack,
				ResetCount:         v.Rst,
				Retransmit:         v.Retransmit,
				Rtt:                v.Rtt,
			},
		)
	}
	return tcpSlice
}

func (c *grpcConverter) udpMsg(ip map[base.ProtocolClass]interface{}) []*pb.Udp {
	if ip == nil {
		return nil
	}
	ipUdp, exists := ip[base.ProtocolUdp].(*udp.ChannelAggregatedValues)
	if !exists {
		return nil
	}
	udpSlice := make([]*pb.Udp, 0)
	for k, v := range ipUdp.Values {
		udpSlice = append(udpSlice,
			&pb.Udp{
				SourceIpAddr:       k.SrcIp,
				TargetIpAddr:       k.DstIp,
				SendTotalByte:      v.SendBytes,
				ReceiveTotalByte:   v.RecvBytes,
				SendTotalPacket:    v.SendCount,
				ReceiveTotalPacket: v.RecvCount,
				TotalPacket:        v.Count,
				SourcePort:         int32(k.SrcPort),
				TargetPort:         int32(k.DstPort),
			},
		)
	}
	return udpSlice
}

func (c *grpcConverter) icmpMsg(ip map[base.ProtocolClass]interface{}) []*pb.Icmp {
	if ip == nil {
		return nil
	}
	ipIcmp, exists := ip[base.ProtocolIcmp].(*icmp.ChannelAggregatedValues)
	if !exists {
		return nil
	}
	icmpSlice := make([]*pb.Icmp, 0)
	for k, v := range ipIcmp.Values {
		icmpSlice = append(icmpSlice,
			&pb.Icmp{
				SourceIpAddr:       k.SrcIp,
				TargetIpAddr:       k.DstIp,
				SendTotalByte:      v.SendBytes,
				ReceiveTotalByte:   v.RecvBytes,
				SendTotalPacket:    v.SendCount,
				ReceiveTotalPacket: v.RecvCount,
				TotalPacket:        v.Count,
			},
		)
	}
	return icmpSlice
}

func (c *grpcConverter) igmpMsg(ip map[base.ProtocolClass]interface{}) []*pb.Igmp {
	if ip == nil {
		return nil
	}
	ipIgmp, exists := ip[base.ProtocolUdp].(*udp.ChannelAggregatedValues)
	if !exists {
		return nil
	}
	igmpSlice := make([]*pb.Igmp, 0)
	for k, v := range ipIgmp.Values {
		igmpSlice = append(igmpSlice,
			&pb.Igmp{
				SourceIpAddr:       k.SrcIp,
				TargetIpAddr:       k.DstIp,
				SendTotalByte:      v.SendBytes,
				ReceiveTotalByte:   v.RecvBytes,
				SendTotalPacket:    v.SendCount,
				ReceiveTotalPacket: v.RecvCount,
				TotalPacket:        v.Count,
			},
		)
	}
	return igmpSlice
}
