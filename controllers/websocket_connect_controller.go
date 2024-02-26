package controllers

import (
	"log"

	"main/models"
	svc "main/services"

	socketio "github.com/googollee/go-socket.io"
)

func HandleMessage(s socketio.Conn, msg models.Message) {
	client, config := getKubeConfig()
	stdout, stderr, err := svc.ExecCommandInPod(client, config, msg)
	if err != nil {
		log.Printf("error executing command in pod: %v", err)
		s.Emit("message", models.Response{
			Command: msg.Command,
			Stdout: "",
			Stderr: stderr,
		})
		return 
	}
	s.Emit("message", models.Response{
		Command: msg.Command,
		Stdout: stdout,
		Stderr: stderr,
	})
}

func HandleError(s socketio.Conn, e error) {
	log.Println("error: ", e)
}

func HandleDisconnect(s socketio.Conn, msg string) {
	log.Println("closed", msg)
}