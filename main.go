package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"main/controllers"
)

func main() {
	r := gin.Default()

	r.GET("/namespaces", controllers.GetNamespaces)
	r.GET("/pods", controllers.GetPods)
	r.GET("/ws", controllers.HandleConnections)

	if err := r.Run(":8000"); err != nil {
		fmt.Printf("Error to connect 8000")
	}
}