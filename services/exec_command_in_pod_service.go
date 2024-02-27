package services

import (
	"io"
	"log"
	"fmt"
	"context"

	"main/models"
	"github.com/gorilla/websocket"

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

type websocketWriter struct {
	ws *websocket.Conn
}

func (w *websocketWriter) Write(p []byte) (n int, err error) {
	err = w.ws.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func ExecCommandInPod(client kubernetes.Interface, config *restclient.Config, msg models.Message, ws *websocket.Conn) error {
	command, podName, namespace := msg.Command, msg.PodName, msg.Namespace

	containerName, err := getContainerNameInPod(client, podName, namespace)
	if err != nil {
		return err
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
		Stdin: true,
		Stdout: true,
		Stderr: true,
		TTY: true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		log.Printf("error remotecommand.NewSPDYExecutor: %v", err)
		return err
	}

	r, w := io.Pipe()

	go func() {
		defer w.Close()
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			if _, err := w.Write(message); err != nil {
				log.Println("write:", err)
				return
			}
		}
	}()

	stdout := &websocketWriter{ws: ws}
	stderr := &websocketWriter{ws: ws}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin: r,
		Stdout: stdout,
		Stderr: stderr,
	})
	log.Printf(stderr)
	if err != nil {
		log.Printf("error exec.Stream: %v", err)
		return err
	}

	return nil
}