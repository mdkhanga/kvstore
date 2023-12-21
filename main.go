package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Welcome to key value store")

	router := gin.Default()
	router.GET("/kvstore", getInfo)

	router.Run("localhost:8080")
}

func getInfo(c *gin.Context) {
	c.JSON(http.StatusOK, "Welcome to keystore")
}
