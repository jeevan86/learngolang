package grpc

import (
	"context"
	"flag"
	"fmt"
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/util/str"
	"google.golang.org/grpc"
	"net"
)

var (
	host = flag.String("host", str.EMPTY(), "The server host")
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedCollectServer
}

func (s *server) Save(ctx context.Context, in *pb.NetStaticsReq) (*pb.NetStaticsRsp, error) {
	fmt.Printf("Received: %v\n", in.String())
	return &pb.NetStaticsRsp{Message: "Hello " + in.GatherIp}, nil
}

func (s *server) LocalIp(ctx context.Context, in *pb.LocalIpReq) (*pb.LocalIpRsp, error) {
	fmt.Printf("Received: %v\n", in.String())
	return &pb.LocalIpRsp{
		Data: &pb.LocalIpLst{
			IpList: []string{
				"172.10.231.101", "172.10.231.102", "172.10.231.103",
			},
		},
	}, nil
}

var listener net.Listener
var serv *grpc.Server

func PrepareTestServer() {
	listener, serv = newTestServer()
	fmt.Printf("prepared server listening at %v\n", listener.Addr())
	startTestServer(listener, serv)
}

func StopTestServer() {
	serv.Stop()
	_ = listener.Close()
}

func newTestServer() (net.Listener, *grpc.Server) {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	return lis, s
}

func startTestServer(lis net.Listener, s *grpc.Server) {
	go func() {
		pb.RegisterCollectServer(s, &server{})
		if err := s.Serve(lis); err != nil {
			fmt.Printf("failed to serve: %v\n", err)
		}
	}()
}
