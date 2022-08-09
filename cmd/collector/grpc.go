package main

import (
	"github.com/jeevan86/learngolang/cmd/util"
	"github.com/jeevan86/learngolang/pkg/collect/server/grpc"
	"github.com/jeevan86/learngolang/pkg/server/actuator"
	"github.com/jeevan86/learngolang/pkg/server/http"
)

/*
export PATH=$GOPATH/bin:$PATH
export PRJ_DIR=`pwd`
cd ${PRJ_DIR}/pkg/collect/api/grpc/pb && go generate && cd ${PRJ_DIR}
go build --ldflags "-extldflags -static"         \
 -o ${PRJ_DIR}/dist/binary/gopcap-collector-grpc \
${PRJ_DIR}/cmd/collector/grpc.go

export GRPC_GO_LOG_VERBOSITY_LEVEL=99
export GRPC_GO_LOG_SEVERITY_LEVEL=debug

cd ${PRJ_DIR}/dist && \
dlv --listen=:8616 --headless=true --api-version=2 --accept-multiclient exec \
 ./binary/gopcap-collector-grpc       \
-- -cluster-id kubernetes              \
   -grpc-host 0.0.0.0                  \
   -grpc-port 50051                    \
   -config-file ./config/collector.yml \
   -kube-config ./config/kube-conf.yml

cd ${PRJ_DIR}/dist && \
 ./binary/gopcap-collector-grpc       \
   -cluster-id kubernetes              \
   -grpc-host 0.0.0.0                  \
   -grpc-port 50051                    \
   -config-file ./config/collector.yml \
   -kube-config ./config/kube-conf.yml
*/
func main() {
	grpc.Start()
	actuator.Init()
	http.Start()
	util.WaitForSig()
	http.Stop()
	grpc.Stop()
}
