package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

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
	r.GET("/ws", controllers.HandleConnections)

	if err := r.Run(":8000"); err != nil {
		fmt.Printf("Error to connect 8000")
	}
}