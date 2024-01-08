package main

import (
	"fmt"
	"net/http"

	m "github.com/mdkhanga/kvstore/models"
	server "github.com/mdkhanga/kvstore/tcpserver"

	"github.com/gin-gonic/gin"
)

var kvMap map[string]string

func main() {
	fmt.Println("Welcome to key value store")

	kvMap = make(map[string]string)
	kvMap["hello"] = "world"

	router := gin.Default()
	router.GET("/kvstore", getInfo)
	router.GET("/kvstore/:key", getValue)
	router.POST("/kvstore", setValue)

	go server.Listen()
	router.Run()

}

func getInfo(c *gin.Context) {
	c.JSON(http.StatusOK, "Welcome to keystore")
}

func getValue(c *gin.Context) {
	key := c.Param("key")
	value := kvMap[key]
	jsonString := fmt.Sprintf("{\"%s\":\"%s\"}", key, value)
	c.JSON(http.StatusOK, jsonString)
}

func setValue(c *gin.Context) {
	var input m.KeyValue
	c.BindJSON(&input)
	kvMap[input.Key] = input.Value
	c.JSON(http.StatusOK, "Welcome to keystore")
}
