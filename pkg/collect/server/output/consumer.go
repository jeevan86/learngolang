package output

import (
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/types"
	"github.com/jeevan86/learngolang/pkg/log"
	"github.com/jeevan86/learngolang/pkg/util/jsonutl"
)

var logger = log.NewLogger()

type FuncConsumer interface {
	Apply(interface{})
}

type IpPortMetaConsumer struct {
	apply func(*types.ChannelPacketMeta)
}

func (c *IpPortMetaConsumer) Apply(o interface{}) {
	c.apply(o.(*types.ChannelPacketMeta))
}

var loggingConsumer = IpPortMetaConsumer{
	apply: func(meta *types.ChannelPacketMeta) {
		metaJson := jsonutl.ToJsonStr(meta)
		logger.Info("Received: %s", metaJson)
	},
}

func NewConsumer() FuncConsumer {
	return &loggingConsumer
}
