package service

import (
	"fmt"
	"context"

	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListNamespaces(client kubernetes.Interface) ([]string, error) {
	fmt.Println("Get Kubernetes Namespaces")
	namespaces, err := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("error getting namespaces: %v\n", err)
		return nil, err
	}

	namespaceList := make([]string, len(namespaces.Items))
	for i, ns := range namespaces.Items {
		namespaceList[i] = ns.Name
	}

	return namespaceList, nil
} 

