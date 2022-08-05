package grpc

import (
	"fmt"
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/cmdb"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/k8s"
	"github.com/jeevan86/learngolang/pkg/flag"
	"github.com/jeevan86/learngolang/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

var logger = log.NewLogger()

var listener net.Listener
var serv *grpc.Server

func Start() {
	k8s.Start()
	cmdb.Start()
	listener, serv = newGrpcServer()
	startServer()
}

func newGrpcServer() (net.Listener, *grpc.Server) {
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

func startServer() {
	logger.Info("Prepared server listening at %v\n", listener.Addr())
	go func() {
		pb.RegisterCollectServer(serv, &controller{})
		err := serv.Serve(listener) // should block here, because of internal looping.
		if err != nil {
			logger.Fatal("Failed to serve: %s\n", err.Error())
		}
	}()
}

func Stop() {
	serv.Stop()
	_ = listener.Close()
	k8s.Stop()
	cmdb.Stop()
}
