package controllers

import (
	"log"
	"time"
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

	go func() {
		var retryInterval = time.Second
		var maxRetryInterval = time.Minute * 5
		var retries = 0

		for {
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("websocket connection lost: %v", err)

				for {
					ws, _, err = websocket.DefaultDialer.Dial("ws://172.25.219.226:8000/ws", nil)
					if err != nil {
						log.Printf("error retrying connection: %v", err)
						retries++
						if retries > 10 {
							log.Printf("max retries exceeded")
							return 
						}
						time.Sleep(retryInterval)
						if retryInterval < maxRetryInterval {
							retryInterval *= 2
						}
					} else {
						retryInterval = time.Second
						retries = 0
						break
					}
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()

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