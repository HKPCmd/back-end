package main

import (
	"log"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"main/controllers"
)

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://172.25.219.226:8080"}
	r.Use(cors.New(config))

	r.GET("/context", controllers.GetCurrentContext)
	r.GET("/namespaces", controllers.GetNamespaces)
	r.GET("/pods", controllers.GetPods)
	
	server := socketio.NewServer(nil)
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/cmd", "message", controllers.HandleMessage)
	server.OnError("/", contorllers.HandleError)
	server.OnDisconnect("/", controllers.HandleDisconnect)

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	r.GET("/socket.io/*any", gin.WrapH(server))
	r.POST("/socket.io/*any", gin.WrapH(server))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
	})

	if err := r.Run(":8000"); err != nil {
		fmt.Printf("Error to connect 8000")
	}
}