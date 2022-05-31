package client

import (
	"context"
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc"
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/ip/base"
	"github.com/jeevan86/learngolang/pkg/util/arr"
	"log"
	"time"
)

type grpcCollector struct {
	converter *converter
	client    pb.CollectClient
}

func newGrpcCollector(serverAddr string) *grpcCollector {
	return &grpcCollector{
		converter: newConverter(),
		client:    grpc.NewClient(serverAddr),
	}
}

func (c *grpcCollector) Collect(msg *base.OutputStruct) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	grpcMsg := c.converter.convert(msg).(*pb.NetStaticsReq)
	// TODO: Could not save: rpc error: code = Unavailable desc = connection closed before server preface received
	r, err := c.client.Save(ctx, grpcMsg)
	if err != nil {
		log.Printf("Could not save: %v", err)
	} else {
		log.Printf("Save success: %s", r.GetMessage())
	}
}

func (c *grpcCollector) GetLocalIpList(nodeIp string) []string {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	req := &pb.LocalIpReq{
		NodeIp: nodeIp,
	}
	rsp, err := c.client.LocalIp(ctx, req)
	if err != nil {
		return arr.EMPTY()
	}
	if data := rsp.GetData(); data != nil {
		lst := data.GetIpList()
		if lst != nil {
			return lst
		}
	}
	return arr.EMPTY()
}
