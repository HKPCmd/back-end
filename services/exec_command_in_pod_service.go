package services

import (
	"log"
	"fmt"
	"bytes"
	"context"

	"main/models"

	"k8s.io/kubectl/pkg/scheme"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"

	corev1 "k8s.io/api/core/v1"
	restclient "k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getContainerNameInPod(client kubernetes.Interface, podName string, namespace string) (string, error) {
	pod, err := client.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if len(pod.Spec.Containers) > 0 {
		return pod.Spec.Containers[0].Name, nil
	}

	return "", fmt.Errorf("no containers in pod %s", podName)
}

func ExecCommandInPod(client kubernetes.Interface, config *restclient.Config, msg models.Message) (string, string, error) {
	command, podName, namespace := msg.Command, msg.PodName, msg.Namespace

	containerName, err := getContainerNameInPod(client, podName, namespace)
	if err != nil {
		return "", "", err
	}

	req := client.CoreV1().RESTClient().Post().
	Resource("pods").
	Name(podName).
	Namespace(namespace).
	SubResource("exec").
	Param("container", containerName)

	req.VersionedParams(&corev1.PodExecOptions{
		Container: containerName,
		Command: []string{"/bin/sh", "-c", command},
		Stdin: false,
		Stdout: true,
		Stderr: true,
		TTY: false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		log.Printf("error remotecommand.NewSPDYExecutor: %v", err)
		return "", "", err
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		log.Printf("error exec.Stream: %v", err)
		return "", stderr.String(), err
	}

	return stdout.String(), stderr.String(), nil
}