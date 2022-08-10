package k8s

import (
	"github.com/jeevan86/learngolang/pkg/flag"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

/*
apiVersion: v1
kind: Node
metadata:
  annotations:
    kubeadm.alpha.kubernetes.io/cri-socket: /var/run/dockershim.sock
    node.alpha.kubernetes.io/ttl: "0"
    volumes.kubernetes.io/controller-managed-attach-detach: "true"
  creationTimestamp: "2021-07-20T10:01:46Z"
  labels:
    beta.kubernetes.io/arch: amd64
    beta.kubernetes.io/os: linux
    kubernetes.io/arch: amd64
    kubernetes.io/hostname: biyi-04
    kubernetes.io/os: linux
    storagenode: glusterfs
  managedFields:
    ...
  name: biyi-04
  resourceVersion: "72043585"
  uid: 68e4a42a-222e-4476-8c78-7c6567c04b9e
spec:
  podCIDR: 172.168.3.0/24
  podCIDRs:
  - 172.168.3.0/24
status:
  addresses:
  - address: 192.168.7.164
    type: InternalIP
  - address: biyi-04
    type: Hostname
*/

type NodeIpMeta struct {
	ClusterId string `json:"clusterId,omitempty"`
	Name      string `json:"name,omitempty"`
	Ip        string `json:"ip,omitempty"`
}

var nodeIpMetaMap = make(map[string]*NodeIpMeta)

func newNodeEventHandler() cache.ResourceEventHandler {
	return &cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			updateNode(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			updateNode(newObj)
		},
		DeleteFunc: func(obj interface{}) {
			deleteNode(obj)
		},
	}
}

func deleteNode(obj interface{}) {
	node := obj.(*coreV1.Node)
	for _, addr := range node.Status.Addresses {
		if addr.Type == coreV1.NodeInternalIP {
			nodeIp := addr.Address
			delete(nodeIpMetaMap, nodeIp)
			break
		}
	}
}

func updateNode(obj interface{}) {
	node := obj.(*coreV1.Node)
	for _, addr := range node.Status.Addresses {
		if addr.Type == coreV1.NodeInternalIP {
			nodeIp := addr.Address
			nodeIpMetaMap[nodeIp] = newNodeIpMeta(node.Name, nodeIp)
			break
		}
	}
}

func newNodeIpMeta(name, ip string) *NodeIpMeta {
	meta := &NodeIpMeta{
		ClusterId: *flag.KubeClusterId,
		Name:      name,
		Ip:        ip,
	}
	return meta
}
