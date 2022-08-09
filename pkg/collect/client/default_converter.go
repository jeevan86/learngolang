package client

import (
	"github.com/jeevan86/learngolang/pkg/collect/api"
	"github.com/jeevan86/learngolang/pkg/config"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/icmp"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/tcp"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/udp"
)

type defaultConverter string

func (c *defaultConverter) convert(m *base.OutputStruct) interface{} {
	gatherIp := *config.GetConfig().NodeIp
	converted := &api.NetStatics{
		Timestamp: m.Bucket,
		GatherIp:  gatherIp,
	}
	ip4, ip4Exists := m.Values[base.Ipv4]
	if ip4Exists {
		converted.Ip4 = &api.Protocol{
			Tcp:  c.tcpMsg(ip4),
			Udp:  c.udpMsg(ip4),
			Icmp: c.icmpMsg(ip4),
			Igmp: c.igmpMsg(ip4),
		}
	}
	ip6, ip6Exists := m.Values[base.Ipv6]
	if ip6Exists {
		converted.Ip6 = &api.Protocol{
			Tcp:  c.tcpMsg(ip6),
			Udp:  c.udpMsg(ip6),
			Icmp: c.icmpMsg(ip6),
			Igmp: c.igmpMsg(ip6),
		}
	}
	return converted
}

func (c *defaultConverter) tcpMsg(ip map[base.ProtocolClass]interface{}) []*api.Tcp {
	if ip == nil {
		return nil
	}
	ipTcp, exists := ip[base.ProtocolTcp].(*tcp.ChannelAggregatedValues)
	if !exists {
		return nil
	}
	tcpSlice := make([]*api.Tcp, 0)
	for k, v := range ipTcp.Values {
		tcpSlice = append(tcpSlice,
			&api.Tcp{
				Ip: api.Ip{
					SourceIpAddr:       k.SrcIp,
					TargetIpAddr:       k.DstIp,
					SendTotalByte:      v.SendBytes,
					ReceiveTotalByte:   v.RecvBytes,
					SendTotalPacket:    v.SendCount,
					ReceiveTotalPacket: v.RecvCount,
					TotalPacket:        v.Count,
					TotalBytes:         v.Bytes,
				},
				SourcePort:     k.SrcPort,
				TargetPort:     k.DstPort,
				SynCount:       v.Syn,
				SynAckCount:    v.SynAck,
				SynAckAckCount: v.SynAckAck,
				FinCount:       v.Fin,
				FinAckCount:    v.FinAck,
				AckCount:       v.Ack,
				ResetCount:     v.Rst,
				Retransmit:     v.Retransmit,
				Rtt:            v.Rtt,
			},
		)
	}
	return tcpSlice
}

func (c *defaultConverter) udpMsg(ip map[base.ProtocolClass]interface{}) []*api.Udp {
	if ip == nil {
		return nil
	}
	ipUdp, exists := ip[base.ProtocolUdp].(*udp.ChannelAggregatedValues)
	if !exists {
		return nil
	}
	udpSlice := make([]*api.Udp, 0)
	for k, v := range ipUdp.Values {
		udpSlice = append(udpSlice,
			&api.Udp{
				Ip: api.Ip{
					SourceIpAddr:       k.SrcIp,
					TargetIpAddr:       k.DstIp,
					SendTotalByte:      v.SendBytes,
					ReceiveTotalByte:   v.RecvBytes,
					SendTotalPacket:    v.SendCount,
					ReceiveTotalPacket: v.RecvCount,
					TotalPacket:        v.Count,
					TotalBytes:         v.Bytes,
				},
				SourcePort: k.SrcPort,
				TargetPort: k.DstPort,
			},
		)
	}
	return udpSlice
}

func (c *defaultConverter) icmpMsg(ip map[base.ProtocolClass]interface{}) []*api.Icmp {
	if ip == nil {
		return nil
	}
	ipIcmp, exists := ip[base.ProtocolIcmp].(*icmp.ChannelAggregatedValues)
	if !exists {
		return nil
	}
	icmpSlice := make([]*api.Icmp, 0)
	for k, v := range ipIcmp.Values {
		icmpSlice = append(icmpSlice,
			&api.Icmp{
				Ip: api.Ip{
					SourceIpAddr:       k.SrcIp,
					TargetIpAddr:       k.DstIp,
					SendTotalByte:      v.SendBytes,
					ReceiveTotalByte:   v.RecvBytes,
					SendTotalPacket:    v.SendCount,
					ReceiveTotalPacket: v.RecvCount,
					TotalPacket:        v.Count,
					TotalBytes:         v.Bytes,
				},
			},
		)
	}
	return icmpSlice
}

func (c *defaultConverter) igmpMsg(ip map[base.ProtocolClass]interface{}) []*api.Igmp {
	if ip == nil {
		return nil
	}
	ipIgmp, exists := ip[base.ProtocolUdp].(*udp.ChannelAggregatedValues)
	if !exists {
		return nil
	}
	igmpSlice := make([]*api.Igmp, 0)
	for k, v := range ipIgmp.Values {
		igmpSlice = append(igmpSlice,
			&api.Igmp{
				Ip: api.Ip{
					SourceIpAddr:       k.SrcIp,
					TargetIpAddr:       k.DstIp,
					SendTotalByte:      v.SendBytes,
					ReceiveTotalByte:   v.RecvBytes,
					SendTotalPacket:    v.SendCount,
					ReceiveTotalPacket: v.RecvCount,
					TotalPacket:        v.Count,
					TotalBytes:         v.Bytes,
				},
			},
		)
	}
	return igmpSlice
}
