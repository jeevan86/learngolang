package client

import (
	"encoding/json"
	"github.com/jeevan86/learngolang/pkg/collect/api"
	"github.com/jeevan86/learngolang/pkg/collect/api/http"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
	"github.com/jeevan86/learngolang/pkg/util/arr"
	"log"
)

type httpCollector struct {
	converter *converter
	handler   *http.Handler
}

func newHttpCollector(serverAddr string) *httpCollector {
	return &httpCollector{
		converter: newConverter(),
		handler:   http.NewHandler(serverAddr),
	}
}

func (c *httpCollector) Collect(msg *base.OutputStruct) {
	body := c.converter.convert(msg).(*api.NetStatics)
	rsp, err := c.handler.Save(body)
	if err != nil {
		log.Printf("Could not save: %v", err)
	} else {
		log.Printf("Save success: %s", *rsp)
	}
}

func (c *httpCollector) GetLocalIpList(nodeIp string) []string {
	rsp, err := c.handler.LocalIp(nodeIp)
	if err != nil {
		return arr.EMPTY()
	}
	ipList := &api.LocalIpLst{}
	err = json.Unmarshal([]byte(*rsp), ipList)
	if err != nil {
		return arr.EMPTY()
	}
	if data := ipList.Data; data != nil {
		return data
	}
	return arr.EMPTY()
}
