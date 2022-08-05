package grpc

import (
	"fmt"
	"github.com/jeevan86/learngolang/pkg/collect/api/grpc/pb"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/cmdb"
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/k8s"
	"github.com/jeevan86/learngolang/pkg/config"
	"github.com/jeevan86/learngolang/pkg/flag"
	"github.com/jeevan86/learngolang/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	serv "net/http"
)

var logger = log.NewLogger()

// _ 服务器名称，这个包含三个元素：监听、控制器、服务
const _ = "CollectorGrpcServer"

var grpcListener net.Listener
var grpcController = &controller{}
var grpcServer *grpc.Server

func Start() {
	k8s.Start()
	cmdb.Start()
	grpcServer = newGrpcServer()
	grpcListener = startListen()
	startServer()
}

func startListen() net.Listener {
	lis, err := net.Listen(newLisAttr())
	if err != nil || lis == nil {
		logger.Fatal("failed to listen: %v\n", err)
	}
	return lis
}

func newGrpcServer() *grpc.Server {
	return grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
	)
}

func newLisAttr() (lisType string, lisAddr string) {
	lisType = "tcp"
	if *flag.GrpcHost != "" && *flag.GrpcPort == -1 {
		lisAddr = fmt.Sprintf("%s:%d", *flag.GrpcHost, *flag.GrpcPort)
	} else {
		collectorCfg := config.GetConfig().Collector
		if collectorCfg == nil {
			logger.Fatal("Unable to create grpc server, CollectorConfig missed.")
		}
		if collectorCfg.GrpcHost == nil || collectorCfg.GrpcPort == nil {
			lisAddr = fmt.Sprintf("%s:%d", "0.0.0.0", 50051)
			logger.Warn("Configure grpc server listen address as %s", lisAddr)
		} else {
			lisAddr = fmt.Sprintf("%s:%d", *collectorCfg.GrpcHost, *collectorCfg.GrpcPort)
		}
	}
	return
}

func startServer() {
	logger.Info("Prepared server listening at %v\n", grpcListener.Addr())
	go func() {
		pb.RegisterCollectServer(grpcServer, grpcController)
		err := serv.Serve(grpcListener, grpcServer) // should block here, because of internal looping.
		if err != nil {
			logger.Fatal("Failed to serve: %s\n", err.Error())
		}
	}()
}

func Stop() {
	_ = grpcListener.Close()
	grpcServer.Stop()
	k8s.Stop()
	cmdb.Stop()
}
