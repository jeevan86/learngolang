package grpc

import (
	"github.com/jeevan86/learngolang/pkg/collect/server/backend"
	"github.com/jeevan86/learngolang/pkg/util/jsonutl"
)

type loggingConsumer int8

func (c *loggingConsumer) Apply(meta *backend.ChannelPacketMeta) {
	metaJson := jsonutl.ToJsonStr(meta)
	logger.Info("Received: %s", metaJson)
}

type channelConsumer struct {
	ch chan *backend.ChannelPacketMeta
}

func (c *channelConsumer) Apply(meta *backend.ChannelPacketMeta) {
	c.ch <- meta
}

func newConsumer() IpPortMetaConsumer {
	//c := channelConsumer{
	//	ch: make(chan *ChannelPacketMeta, 1024),
	//}
	var c loggingConsumer = 0
	return &c
}
