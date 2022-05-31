package tcp

import (
	"github.com/jeevan86/learngolang/pkg/log"
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/ip/base"
	"github.com/jeevan86/learngolang/pkg/util/str"
)

var logger = log.NewLogger()

type AggregatedValues struct {
	base.CommonAggregatedValues
	Syn        uint64
	SynAck     uint64
	SynAckAck  uint64
	Fin        uint64
	FinAck     uint64
	Ack        uint64
	Rst        uint64
	Retransmit uint64
	Rtt        int64
}

type ChannelAggregatedValues struct {
	Values map[Channel]*AggregatedValues
}

// flag 是否启用助记符
func flag(enabled bool, name string) string {
	if enabled {
		return "[" + name + "]"
	}
	return str.EMPTY()
}
