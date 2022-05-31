package k8s

import (
	cfg "github.com/jeevan86/learngolang/pkg/flag"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type ResIpMeta struct {
	ClusterId string `json:"cluster-id,omitempty" yaml:"cluster-id"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace"`
}

type interestInformers struct {
	stopCh             chan struct{}
	replicaSetInformer cache.SharedIndexInformer
	serviceInformer    cache.SharedIndexInformer
	podInformer        cache.SharedIndexInformer
	nodeInformer       cache.SharedIndexInformer
}

func (i *interestInformers) run() {
	i.nodeInformer.AddEventHandler(newNodeEventHandler())
	i.podInformer.AddEventHandler(newPodEventHandler())
	i.serviceInformer.AddEventHandler(newSrvEventHandler())
	go i.replicaSetInformer.Run(i.stopCh)
	go i.serviceInformer.Run(i.stopCh)
	go i.podInformer.Run(i.stopCh)
	go i.nodeInformer.Run(i.stopCh)
}

func (i *interestInformers) stop() {
	close(i.stopCh)
}

var interesting = new(interestInformers)

func newClient() kubernetes.Interface {
	// use the current context in kube config
	config, err := clientcmd.BuildConfigFromFlags("", *cfg.KubeConfig)
	if err != nil {
		panic(err.Error())
	}
	// create the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return client
}

func newInformerFactory(client kubernetes.Interface) informers.SharedInformerFactory {
	return informers.NewSharedInformerFactory(client, time.Minute*10)
}

func initInteresting(factory informers.SharedInformerFactory) {
	interesting.stopCh = make(chan struct{})
	interesting.replicaSetInformer = factory.Apps().V1().ReplicaSets().Informer()
	interesting.serviceInformer = factory.Core().V1().Services().Informer()
	interesting.podInformer = factory.Core().V1().Pods().Informer()
	interesting.nodeInformer = factory.Core().V1().Nodes().Informer()
}

func listen() {
	//pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	//watching, err := clientset.CoreV1().Pods("").Watch(metav1.ListOptions{})
	//if err != nil {
	//	panic(err.Error())
	//}
	//select {
	//case event := <-watching.ResultChan():
	//
	//}
}

func resourceKey(namespace, name string) string {
	return namespace + "/" + name
}

func labelAsKey(namespace, label, value string) string {
	return namespace + "/" + label + "=" + value
}
