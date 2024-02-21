package models

type Message struct {
	PodName   string `json:"podName"`
	Namespace string `json:"namespace"`
	Command   string `json:"command"`
}

type Response struct {
	Command string `json:"command"`
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
}