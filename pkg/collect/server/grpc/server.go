package grpc

import (
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/k8s"
	"github.com/jeevan86/learngolang/pkg/log"
	"google.golang.org/grpc"
	"net"
)

var logger = log.NewLogger()

var listener net.Listener
var serv *grpc.Server

func Start() {
	listener, serv = newServer()
	logger.Info("Prepared server listening at %v\n", listener.Addr())
	startServer(listener, serv)
}

func startServer(lis net.Listener, s *grpc.Server) {
	k8s.Start()
	go func() {
		pb.RegisterCollectServer(s, &server{})
		if err := s.Serve(lis); err != nil {
			logger.Fatal("Failed to serve: %s\n", err.Error())
		}
	}()
}

func Stop() {
	serv.Stop()
	_ = listener.Close()
	k8s.Stop()
}
