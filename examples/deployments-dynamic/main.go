package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = "config"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "examples/deployments-dynamic/index.html")
	})
	http.HandleFunc("/namespaces", func(w http.ResponseWriter, r *http.Request) {
		// Load the kubeconfig file
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Fprintln(w, "Error loading kubeconfig:", err)
			return
		}

		// Create a new Kubernetes client
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Fprintln(w, "Error creating Kubernetes client:", err)
			return
		}

		// Get a list of all namespaces in the cluster
		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Fprintln(w, "Error fetching namespaces:", err)
			return
		}

		// Convert the namespaces to a list of strings
		var namespaceList []string
		for _, namespace := range namespaces.Items {
			namespaceList = append(namespaceList, namespace.Name)
		}

		// Return the list of namespaces as a JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(namespaceList)
		fmt.Println("Inside handler")
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}
}
