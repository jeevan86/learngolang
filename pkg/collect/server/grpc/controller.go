package grpc

import (
	"context"
	"fmt"
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/k8s"
	"github.com/jeevan86/learngolang/pkg/flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

type server struct {
	pb.UnimplementedCollectServer
}

var consumer = newConsumer()

func (s *server) Save(ctx context.Context, in *pb.NetStaticsReq) (*pb.NetStaticsRsp, error) {
	ip4(in.Ip4, consumer)
	ip6(in.Ip6, consumer)
	return &pb.NetStaticsRsp{
		Message: fmt.Sprintf("From %s. done.", in.GatherIp),
	}, nil
}

func (s *server) LocalIp(ctx context.Context, in *pb.LocalIpReq) (*pb.LocalIpRsp, error) {
	return &pb.LocalIpRsp{
		Data: &pb.LocalIpLst{
			IpList: k8s.GetNodePodIpList(in.NodeIp),
		},
	}, nil
}

func newServer() (net.Listener, *grpc.Server) {
	lisType := "tcp"
	lisAddr := fmt.Sprintf("%s:%d", *flag.GrpcHost, *flag.GrpcPort)
	lis, err := net.Listen(lisType, lisAddr)
	if err != nil || lis == nil {
		logger.Fatal("failed to listen: %v\n", err)
		return nil, nil
	}
	s := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
	)
	return lis, s
}
