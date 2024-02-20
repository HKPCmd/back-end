package controllers

import (
	"net/http"

	svc "main/services"

	"github.com/gin-gonic/gin"
)

func GetNamespaces(c *gin.Context) {
	client := getKubeConfig()

	namespaceList, err := svc.ListNamespaces(client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, namespaceList)
}

func GetPods(c *gin.Context) {
	namespace := c.Query("namespace")
	client := getKubeConfig()

	podsList, err := svc.ListPodNameAndStatus(client, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, podsList)
}