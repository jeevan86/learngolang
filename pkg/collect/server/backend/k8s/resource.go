package k8s

type ResIpMeta struct {
	ClusterId string `json:"cluster-id,omitempty" yaml:"cluster-id"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace"`
}

func resourceKey(namespace, name string) string {
	return namespace + "/" + name
}

func labelAsKey(namespace, label, value string) string {
	return namespace + "/" + label + "=" + value
}
