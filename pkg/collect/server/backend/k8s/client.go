package k8s

import (
	"github.com/jeevan86/learngolang/pkg/config"
	"github.com/jeevan86/learngolang/pkg/flag"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

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
	var cfg *rest.Config
	var err error
	if *flag.KubeConfigFile != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", *flag.KubeConfigFile)
	} else {
		collectorCfg := config.GetConfig().Collector
		if collectorCfg == nil {
			logger.Fatal("Unable to create k8s client, KubeConfig missed.")
		}
		if collectorCfg.KubeConfigData == nil && collectorCfg.KubeConfigFile == nil {
			logger.Fatal("Unable to create k8s client, KubeConfigData or KubeConfigFile missed.")
		}
		if collectorCfg.KubeConfigData != nil {
			data := collectorCfg.KubeConfigData
			bytes := []byte(*data)
			cfg, err = clientcmd.RESTConfigFromKubeConfig(bytes)
		} else if collectorCfg.KubeConfigFile != nil {
			cfg, err = clientcmd.BuildConfigFromFlags("", *collectorCfg.KubeConfigFile)
		}
	}
	if err != nil {
		logger.Fatal("Failed to create k8s client config, err => %s", err.Error())
	}
	// create the client
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logger.Fatal("Failed to create k8s client, err => %s", err.Error())
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
