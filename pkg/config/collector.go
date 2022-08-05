package config

type collectorConfig struct {
	KubeClusterId  *string `json:"KubeClusterId,omitempty" yaml:"kube-cluster-id"`
	KubeConfigFile *string `json:"kubeConfigFile,omitempty" yaml:"kube-config-file"`
	KubeConfigData *string `json:"kubeConfigData,omitempty" yaml:"kube-config-data"`
	GrpcHost       *string `json:"grpcHost,omitempty" yaml:"grpc-host"`
	GrpcPort       *int32  `json:"grpcPort,omitempty" yaml:"grpc-port"`
}
