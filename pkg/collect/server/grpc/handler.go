package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/flag"
	"github.com/jeevan86/learngolang/pkg/k8s"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

type server struct {
	pb.UnimplementedCollectServer
}

func (s *server) Save(ctx context.Context, in *pb.NetStaticsReq) (*pb.NetStaticsRsp, error) {
	// ignore icmp/igmp now
	ip4(in.Ip4, apply)
	ip6(in.Ip6, apply)
	return &pb.NetStaticsRsp{
		Message: fmt.Sprintf("from %s. done.", in.GatherIp),
	}, nil
}

func apply(meta *IpPortMeta) {
	metaBytes, _ := json.Marshal(meta)
	metaJson := string(metaBytes)
	logger.Info("Received: %s", metaJson)
}

func ip4(ip *pb.Protocol, f func(*IpPortMeta)) {
	if ip != nil {
		for _, tcp := range ip.Tcp {
			ipPortMeta := tcpMeta(tcp)
			ipPortMeta.Protocol = IpProtocolTcp4
			f(ipPortMeta)
		}
		for _, udp := range ip.Udp {
			ipPortMeta := udpMeta(udp)
			ipPortMeta.Protocol = IpProtocolUdp4
			f(ipPortMeta)
		}
	}
}

func ip6(ip *pb.Protocol, f func(*IpPortMeta)) {
	if ip != nil {
		for _, tcp := range ip.Tcp {
			ipPortMeta := tcpMeta(tcp)
			ipPortMeta.Protocol = IpProtocolTcp6
			f(ipPortMeta)
		}
		for _, udp := range ip.Udp {
			ipPortMeta := udpMeta(udp)
			ipPortMeta.Protocol = IpProtocolUdp6
			f(ipPortMeta)
		}
	}
}

func (s *server) LocalIp(ctx context.Context, in *pb.LocalIpReq) (*pb.LocalIpRsp, error) {
	return &pb.LocalIpRsp{
		Data: &pb.LocalIpLst{
			IpList: k8s.GetNodePodIpList(in.NodeIp),
		},
	}, nil
}

func newServer() (net.Listener, *grpc.Server) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *flag.GrpcHost, *flag.GrpcPort))
	if err != nil {
		logger.Fatal("failed to listen: %v\n", err)
	}
	s := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
	)
	return lis, s
}
