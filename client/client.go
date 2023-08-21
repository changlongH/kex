package client

import (
	"context"
	"errors"
	"log"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubeCli struct {
	RestConfig *rest.Config
	ClientSet  *kubernetes.Clientset
}

func NewClient() *KubeCli {
	home := homedir.HomeDir()
	if home == "" {
		err := errors.New("not found cluster config ~/.kube/config file")
		panic(err.Error())
	}

	config := ".kube/config"
	kubeconfig := filepath.Join(home, config)
	// use the current context in kubeconfig
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return &KubeCli{
		RestConfig: restConfig,
	}
}

func (cli *KubeCli) InitClientSet() {
	// create the clientset
	clientset, err := kubernetes.NewForConfig(cli.RestConfig)
	if err != nil {
		panic(err.Error())
	}
	cli.ClientSet = clientset
}

func (cli *KubeCli) InitRestCLient()      {}
func (cli *KubeCli) InitDynamicClient()   {}
func (cli *KubeCli) InitDiscoveryClient() {}

// get pods in all the namespaces by omitting namespace
// Or specify namespace to get pods in particular namespace
func (cli *KubeCli) GetPods(namespace string) *v1.PodList {
	pods, err := cli.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	log.Printf("There are %d pods in the cluster\n", len(pods.Items))
	// https://pkg.go.dev/k8s.io/api/core/v1#Pod
	return pods
}
