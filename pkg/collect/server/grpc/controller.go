package grpc

import (
	"context"
	"fmt"
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/k8s"
)

type controller struct {
	pb.UnimplementedCollectServer
}

func (s *controller) Save(ctx context.Context, in *pb.NetStaticsReq) (*pb.NetStaticsRsp, error) {
	ip4(in.Ip4)
	ip6(in.Ip6)
	return &pb.NetStaticsRsp{
		Message: fmt.Sprintf("From %s. done.", in.GatherIp),
	}, nil
}

func (s *controller) LocalIp(ctx context.Context, in *pb.LocalIpReq) (*pb.LocalIpRsp, error) {
	return &pb.LocalIpRsp{
		Data: &pb.LocalIpLst{
			IpList: k8s.GetNodePodIpList(in.NodeIp),
		},
	}, nil
}
