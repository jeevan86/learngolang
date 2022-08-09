package main

import (
	"github.com/jeevan86/learngolang/cmd/util"
	"github.com/jeevan86/learngolang/pkg/capture"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
	"github.com/jeevan86/learngolang/pkg/collect/client"
	"github.com/jeevan86/learngolang/pkg/config"
	"github.com/jeevan86/learngolang/pkg/server/actuator"
	"github.com/jeevan86/learngolang/pkg/server/http"
)

var packetProcessor *base.PacketProcessor
var collector *client.Collect

/*
export PATH=$GOPATH/bin:$PATH
export PRJ_DIR=`pwd`
cd ${PRJ_DIR}/pkg/collect/api/grpc/pb && go generate && cd ${PRJ_DIR}
go build --ldflags "-extldflags -static" \
 -o ${PRJ_DIR}/dist/binary/gopcap-agent  \
${PRJ_DIR}/cmd/pcap/main.go

cd ${PRJ_DIR}/dist && \
dlv --listen=:8626 --headless=true --api-version=2 --accept-multiclient exec \
./binary/gopcap-agent \
-- -config-file ${PRJ_DIR}/dist/config/agent.yml

cd ${PRJ_DIR}/dist && ./binary/gopcap-agent \
 -config-file ${PRJ_DIR}/dist/config/agent.yml
*/
func main() {
	packetProcessor = ip.NewPacketProcessor()
	collector = client.NewCollector()
	actuator.Init()
	http.Start()
	collector.Start(packetProcessor.Out())
	localIpFetcher := func() []string {
		return collector.Api.GetLocalIpList(*config.GetConfig().NodeIp)
	}
	packetProcessor.Start(localIpFetcher)
	capture.StartCapture(packetProcessor.Process)
	util.WaitForSig()
	packetProcessor.Stop()
	collector.Stop()
	http.Stop()
}
