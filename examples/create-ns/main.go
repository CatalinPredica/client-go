package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {

	// kubeconfig - "/path/to/kubeconfig"
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	//fmt.Println("/path/to/kubeconfig", *kubeconfig, "\n")

	// Load the kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Input manual

	// Define the namespace object
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example-namespace1",
		},
	}

	// Create the namespace
	result, err := clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Namespace created: %q\n", result.GetObjectMeta().GetName())
}
