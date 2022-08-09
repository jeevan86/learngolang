package client

import (
	"context"
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc"
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
	"github.com/jeevan86/learngolang/pkg/util/arr"
	"log"
	"time"
)

type grpcCollector struct {
	converter *converter
	client    pb.CollectClient
}

// newGrpcCollector
// @title       创建grpcCollector客户端
// @description 创建grpcCollector客户端
// @auth        小卒    2022/08/03 10:57
// @param       serverAddr string         "服务端地址"
// @return      r          *grpcCollector "grpcCollector客户端"
func newGrpcCollector(serverAddr string) *grpcCollector {
	return &grpcCollector{
		converter: newConverter(),
		client:    grpc.NewClient(serverAddr),
	}
}

// Collect
// @title       执行收集
// @description 执行收集：执行GRPC请求，发送数据
// @auth        小卒  2022/08/03 10:57
// @param       msg  *base.OutputStruct "输出的消息"
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

// GetLocalIpList
// @title       查询本地IP列表
// @description 查询本地IP列表
// @auth        小卒    2022/08/03 10:57
// @param       nodeIp  string    "节点的IP"
// @return      l       []string  "本地IP列表"
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
