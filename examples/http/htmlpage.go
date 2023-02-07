package main

import (
	"context"
	"fmt"
	"html/template"
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
		tmpl := template.Must(template.ParseFiles("examples/http/template.html"))

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
			fmt.Fprintln(w, "Error getting namespaces:", err)
			return
		}

		// Create a slice to hold the names of the namespaces
		var names []string
		for _, namespace := range namespaces.Items {
			names = append(names, namespace.Name)
		}

		// Execute the template with the names of the namespaces
		err = tmpl.Execute(w, names)
		if err != nil {
			fmt.Fprintln(w, "Error executing template:", err)
			return
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting HTTP server:", err)
		return
	}
}
