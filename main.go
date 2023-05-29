package main

import (
	"flag"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	snapshotName string
)

type snapshotMethods struct {
	dynamic dynamic.Interface
}

func main() {
	config := flag.String("kubeconfig", "/home/abhijeet/.kube/config", "this is where the kubeconfig file is stored")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", *config)
	if err != nil {
		fmt.Printf("Creating kubeconfig file: %s", err.Error())

		kubeconfig, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("Creating in cluster config:%s", err.Error())
		}
	}

	//dynamic clientset
	dynclient, err := dynamic.NewForConfig(kubeconfig)

	//typed client set
	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		fmt.Printf("creating client set error : %s", err.Error())
	}

	resource := schema.GroupVersionResource{
		Group:    "mabhi.dev",
		Version:  "v1",
		Resource: "snapinputs",
	}

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynclient, time.Minute, corev1.NamespaceAll, nil)
	informer := factory.ForResource(resource).Informer()
	// fmt.Println("informer: ", informer)

	controller := newController(dynclient, informer, clientset)
	fmt.Println(controller)

	ch := make(chan struct{})
	informer.Run(ch)
	controller.run(ch)
}
