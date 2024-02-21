package services

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

func ListPodNameAndStatus(client kubernetes.Interface, namespace string) (map[string]string, error) {
	fmt.Println("Get Kubernetes Pods")
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("error getting pods: %v \n", err)
		return nil, err
	}

	podNamesAndStatus := make(map[string]string)
	for _, pod := range pods.Items {
		podNamesAndStatus[pod.Name] = string(pod.Status.Phase)
	}

	return podNamesAndStatus, nil  
}