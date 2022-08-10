package k8s

type ResIpMeta struct {
	ClusterId string `json:"clusterId,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

func resourceKey(namespace, name string) string {
	return namespace + "/" + name
}

func labelAsKey(namespace, label, value string) string {
	return namespace + "/" + label + "=" + value
}
