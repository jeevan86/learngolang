package k8s

import (
	"github.com/jeevan86/learngolang/pkg/flag"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

type PodIpMeta struct {
	ResIpMeta
	PodName string `json:"podName,omitempty"`
	PodIp   string `json:"podIp,omitempty"`
	HostIp  string `json:"hostIp,omitempty"`
	AppKind string `json:"appKind,omitempty"`
	AppName string `json:"appName,omitempty"`
}

var podIpMetaMap = make(map[string]*PodIpMeta)
var nodePodIpMap = make(map[string]map[string]*PodIpMeta)
var lblIpMetaMap = make(map[string]map[string]*PodIpMeta)

func newPodEventHandler() cache.ResourceEventHandler {
	return &cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			updatePod(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			updatePod(newObj)
		},
		DeleteFunc: func(obj interface{}) {
			deletePod(obj)
		},
	}
}

func deletePod(obj interface{}) {
	pod := obj.(*coreV1.Pod)
	podIp := pod.Status.PodIP
	hostIp := pod.Status.HostIP
	delete(podIpMetaMap, podIp)
	deleteNodeMap(hostIp, podIp)
	deleteLabelMap(pod)
}

func updatePod(obj interface{}) {
	pod := obj.(*coreV1.Pod)
	namespace := pod.Namespace
	podName := pod.Name
	podIp := pod.Status.PodIP
	hostIp := pod.Status.HostIP
	meta := newIpMeta(namespace, podName, podIp, hostIp)
	for _, owner := range pod.OwnerReferences {
		if KindReplicaSet == owner.Kind {
			item, exists, err := interesting.replicaSetInformer.
				GetIndexer().
				GetByKey(resourceKey(namespace, owner.Name))
			if exists && err == nil {
				replicaSet := item.(*appsV1.ReplicaSet)
				for _, dp := range replicaSet.OwnerReferences {
					updateMeta(meta, KindDeployment, dp.Name)
					break
				}
			}
		} else {
			updateMeta(meta, owner.Kind, owner.Name)
		}
		break
	}
	podIpMetaMap[podIp] = meta
	updateNodeMap(meta)
	updateLabelMap(pod, meta)
	logger.Info("Add pod %s", meta)
}

func newIpMeta(namespace, podName, podIp, hostIp string) *PodIpMeta {
	meta := &PodIpMeta{
		ResIpMeta: ResIpMeta{
			ClusterId: *flag.KubeClusterId,
			Namespace: namespace,
		},
		PodName: podName,
		PodIp:   podIp,
		HostIp:  hostIp,
	}
	return meta
}

func updateMeta(meta *PodIpMeta, appKind, appName string) {
	meta.AppKind = appKind
	meta.AppName = appName
}

func updateLabelMap(pod *coreV1.Pod, meta *PodIpMeta) {
	podLabels := pod.ObjectMeta.Labels
	for k, v := range podLabels {
		labelKey := labelAsKey(meta.Namespace, k, v)
		putLabelMap(labelKey, meta)
	}
}

func putLabelMap(labelKey string, meta *PodIpMeta) {
	podIpMap, exists := lblIpMetaMap[labelKey]
	if !exists || podIpMap == nil {
		podIpMap = make(map[string]*PodIpMeta)
	}
	podIpMap[labelKey] = meta
}

func deleteLabelMap(pod *coreV1.Pod) {
	podLabels := pod.ObjectMeta.Labels
	for k, v := range podLabels {
		labelKey := labelAsKey(pod.Namespace, k, v)
		delLabelMap(labelKey, pod.Status.PodIP)
	}
}

func delLabelMap(labelKey, podIp string) {
	podIpMap, exists := lblIpMetaMap[labelKey]
	if exists && podIpMap != nil {
		delete(podIpMap, podIp)
	}
}

func updateNodeMap(meta *PodIpMeta) {
	podIpMap, exists := nodePodIpMap[meta.HostIp]
	if !exists || podIpMap == nil {
		podIpMap = make(map[string]*PodIpMeta)
	}
	podIpMap[meta.PodIp] = meta
}

func deleteNodeMap(hostIp, podIp string) {
	podIpMap, exists := nodePodIpMap[hostIp]
	if exists && podIpMap != nil {
		delete(podIpMap, podIp)
	}
}
