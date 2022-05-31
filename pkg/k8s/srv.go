package k8s

import (
	"github.com/jeevan86/learngolang/pkg/flag"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

/*
metadata:
  name: metrics-datatier
  namespace: gops-bk-apps
  labels:
    app.kubernetes.io/managed-by: '11'
apiVersion: v1
kind: Service
spec:
  sessionAffinity: None
  externalTrafficPolicy: Cluster
  selector:
    app.kubernetes.io/name: metrics-datatier
  ports:
  - protocol: TCP
    port: 8080
    name: portc396
    nodePort: 17630
    targetPort: 8080
  type: NodePort
  clusterIP: 10.68.215.125
*/

const SrvSelectorKeyAppName = "app.kubernetes.io/name"

type SrvIpMeta struct {
	ResIpMeta
	Name    string `json:"name,omitempty" yaml:"name"`
	Ip      string `json:"ip,omitempty" yaml:"ip"`
	AppKind string `json:"app-kind,omitempty" yaml:"app-kind"`
	AppName string `json:"app-name,omitempty" yaml:"app-name"`
}

var srvIpMetaMap = make(map[string]*SrvIpMeta)

func newSrvEventHandler() cache.ResourceEventHandler {
	return &cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			updateSrv(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			updateSrv(newObj)
		},
		DeleteFunc: func(obj interface{}) {
			deleteSrv(obj)
		},
	}
}

func deleteSrv(obj interface{}) {
	srv := obj.(*coreV1.Service)
	clstrIp := srv.Spec.ClusterIP
	delete(srvIpMetaMap, clstrIp)
}

func updateSrv(obj interface{}) {
	srv := obj.(*coreV1.Service)
	namespace := srv.Namespace
	srvName := srv.Name
	clstrIp := srv.Spec.ClusterIP
	appName := workload(namespace, srv.Spec.Selector)
	srvMeta := newSrvIpMeta(namespace, srvName, clstrIp, appName)
	srvIpMetaMap[clstrIp] = srvMeta
}

func workload(namespace string, selector map[string]string) string {
	var podIpMeta *PodIpMeta
	for k, v := range selector {
		key := labelAsKey(namespace, k, v)
		subMap, exists := lblIpMetaMap[key]
		if exists && subMap != nil {
			for _, meta := range subMap {
				podIpMeta = meta
				break
			}
		}
		if podIpMeta != nil {
			break
		}
	}
	if podIpMeta != nil {
		return podIpMeta.AppName
	}
	return ""
}

func newSrvIpMeta(namespace, name, ip, appName string) *SrvIpMeta {
	meta := &SrvIpMeta{
		ResIpMeta: ResIpMeta{
			ClusterId: *flag.ClusterId,
			Namespace: namespace,
		},
		Name:    name,
		Ip:      ip,
		AppName: appName,
	}
	return meta
}
