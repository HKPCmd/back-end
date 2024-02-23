package main

import (
	"fmt"

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
	
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	controllers.HandleConnections(server)

	r.GET("/ws", func(c *gin.Context) {
		server.ServeHTTP(c.Writer, c.Request)
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
	})

	if err := r.Run(":8000"); err != nil {
		fmt.Printf("Error to connect 8000")
	}
}