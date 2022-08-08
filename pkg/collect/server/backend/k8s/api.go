package k8s

import (
	"github.com/jeevan86/learngolang/pkg/collect/server/backend/types"
	"github.com/jeevan86/learngolang/pkg/log"
)

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

// IsPodOrSrv
// @title       是否是K8S中的PodIp、ServiceIp、NodeIp
// @description 是否是K8S中的PodIp、ServiceIp、NodeIp
// @auth        小卒    2022/08/03 10:57
// @param       ipPort *backend.IpPort     "IP端口"
// @return      ipPort *backend.IpPortMeta "IP端口的信息"
// @return      bool                       "是否有查到数据"
func IsPodOrSrv(ipPort *types.IpPort) (*types.IpPortMeta, bool) {
	if meta, exists := GetPodIpMeta(ipPort.Ip); exists {
		tags := make(types.Tags, 8)
		tags["clusterId"] = meta.ClusterId
		tags["namespace"] = meta.Namespace
		tags["appName"] = meta.AppName
		tags["appKind"] = meta.AppKind
		tags["hostIp"] = meta.HostIp
		tags["podName"] = meta.PodName
		return &types.IpPortMeta{
			IpPort: *ipPort,
			Type:   types.IpPortTypeK8sPod,
			Tags:   tags,
		}, true
	}
	if meta, exists := GetSrvIpMeta(ipPort.Ip); exists {
		tags := make(types.Tags, 8)
		tags["clusterId"] = meta.ClusterId
		tags["namespace"] = meta.Namespace
		tags["appName"] = meta.AppName
		return &types.IpPortMeta{
			IpPort: *ipPort,
			Type:   types.IpPortTypeK8sSrv,
			Tags:   tags,
		}, true
	}
	if meta, exists := GetNodeIpMeta(ipPort.Ip); exists {
		tags := make(types.Tags, 4)
		tags["clusterId"] = meta.ClusterId
		tags["nodeName"] = meta.Name
		return &types.IpPortMeta{
			IpPort: *ipPort,
			Type:   types.IpPortTypeK8sNode,
			Tags:   tags,
		}, true
	}
	return nil, false
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
