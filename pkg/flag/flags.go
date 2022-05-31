package flag

import (
	"flag"
)

var (
	ConfFile   = flag.String("config-file", "", "The config file path.")
	ClusterId  = flag.String("cluster-id", "", "the id of this k8s cluster")
	KubeConfig = flag.String("kube-config", "", "absolute path to the kube config file")
	GrpcHost   = flag.String("grpc-host", "", "the grpc server host")
	GrpcPort   = flag.Int("grpc-port", 50051, "the grpc server port")
)

func init() {
	flag.Parse()
}
