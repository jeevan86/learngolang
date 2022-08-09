package client

import (
	"encoding/json"
	"github.com/jeevan86/learngolang/pkg/collect/api"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
)

type logCollector struct {
	converter *converter
}

func newLogCollector(name string) *logCollector {
	return &logCollector{
		converter: newConverter(),
	}
}

func (c *logCollector) Collect(msg *base.OutputStruct) {
	stuMsg := c.converter.convert(msg).(*api.NetStatics)
	bytes, _ := json.Marshal(stuMsg)
	logger.Info("Collected: %s", string(bytes))
}

func (c *logCollector) GetLocalIpList(nodeIp string) []string {
	return []string{nodeIp}
}
