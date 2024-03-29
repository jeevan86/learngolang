package flag

import (
	"flag"
)

var (
	ConfigFile     = flag.String("config-file", "./config/default.yml", "The config file path.")
	GrpcHost       = flag.String("grpc-host", "", "the grpc server host")
	GrpcPort       = flag.Int("grpc-port", -1, "the grpc server port")
	KubeClusterId  = flag.String("kube-cluster-id", "", "the id of this k8s cluster")
	KubeConfigFile = flag.String("kube-config-file", "", "absolute path to the kube config file")
	CmdbConfigFile = flag.String("cmdb-config-file", "", "absolute path to the cmdb config file")
)

func init() {
	flag.Parse()
}
