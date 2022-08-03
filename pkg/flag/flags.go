package flag

import (
	"flag"
)

var (
	ConfigFile     = flag.String("config-file", "./config/default.yml", "The config file path.")
	GrpcHost       = flag.String("grpc-host", "0.0.0.0", "the grpc server host")
	GrpcPort       = flag.Int("grpc-port", 50051, "the grpc server port")
	KubeClusterId  = flag.String("kube-cluster-id", "", "the id of this k8s cluster")
	KubeConfigFile = flag.String("kube-config-file", "./config/kube-conf.yml", "absolute path to the kube config file")
	CmdbConfigFile = flag.String("cmdb-config-file", "./config/cmdb-conf.yml", "absolute path to the cmdb config file")
)

func init() {
	flag.Parse()
}
