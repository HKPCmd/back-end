package controllers

import (
	"log"
	"net/http"

	"main/models"
	svc "main/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnections(c *gin.Context) {
	client, config := getKubeConfig()

	ws, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("error: %v", err)
	}
	defer ws.Close()

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		log.Printf(msg.Command)
		log.Printf(msg.PodName)

		err = svc.ExecCommandInPod(client, config, msg, ws)
		if err != nil {
			log.Printf("error executing command in pod: %v", err)
			continue
		}
	}
}