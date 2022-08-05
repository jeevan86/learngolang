package k8s

import "github.com/jeevan86/learngolang/pkg/log"

var logger = log.NewLogger()

func Start() {
	client := newClient()
	factory := newInformerFactory(client)
	initInteresting(factory)
	interesting.run()
}

func Stop() {
	interesting.stop()
}

func GetNodeIpMeta(ip string) (meta *NodeIpMeta, exists bool) {
	meta, exists = nodeIpMetaMap[ip]
	return
}

func GetPodIpMeta(ip string) (meta *PodIpMeta, exists bool) {
	meta, exists = podIpMetaMap[ip]
	return
}

func GetSrvIpMeta(ip string) (meta *SrvIpMeta, exists bool) {
	meta, exists = srvIpMetaMap[ip]
	return
}

func GetNodePodIpList(ip string) []string {
	list := make([]string, 256)
	for k := range nodePodIpMap[ip] {
		list = append(list, k)
	}
	return list
}
