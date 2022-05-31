package grpc

import (
	"fmt"
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/k8s"
	"github.com/jeevan86/learngolang/pkg/log"
	"google.golang.org/grpc"
	"net"
)

var logger = log.NewLogger()

var listener net.Listener
var serv *grpc.Server

func Start() {
	listener, serv = newServer()
	fmt.Printf("prepared server listening at %v\n", listener.Addr())
	startServer(listener, serv)
}

func startServer(lis net.Listener, s *grpc.Server) {
	k8s.Start()
	go func() {
		pb.RegisterCollectServer(s, &server{})
		if err := s.Serve(lis); err != nil {
			fmt.Printf("failed to serve: %v\n", err)
		}
	}()
}

func Stop() {
	serv.Stop()
	_ = listener.Close()
	k8s.Stop()
}
