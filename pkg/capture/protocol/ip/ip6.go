package ip

import (
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/icmp"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/igmp"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/tcp"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/udp"
)

var processIp6Packets = func(prev, curr, next base.ProtocolBatch) map[base.ProtocolClass]interface{} {
	result := make(map[base.ProtocolClass]interface{}, 4)
	result[base.ProtocolTcp] = tcp.ProcessTcp6Packets(packetBatches(base.ProtocolTcp, prev, curr, next))
	result[base.ProtocolUdp] = udp.ProcessUdp6Packets(packetBatches(base.ProtocolUdp, prev, curr, next))
	result[base.ProtocolIcmp] = icmp.ProcessIcmp6Packets(packetBatches(base.ProtocolIcmp, prev, curr, next))
	result[base.ProtocolIgmp] = igmp.ProcessIgmp6Packets(packetBatches(base.ProtocolIgmp, prev, curr, next))
	return result
}
