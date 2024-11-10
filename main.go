package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/mdkhanga/kvstore/grpcserver"
	m "github.com/mdkhanga/kvstore/models"

	// client "github.com/mdkhanga/kvstore/tcpclient"
	client "github.com/mdkhanga/kvstore/grpcclient"

	"github.com/gin-gonic/gin"
)

var kvMap map[string]string

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	fmt.Println("Welcome to key value store")

	argsWithProg := os.Args

	fmt.Println(argsWithProg)
	for i := 0; i < len(argsWithProg); i++ {
		fmt.Println(argsWithProg[i])
	}

	portPtr := flag.String("p", "8081", "tcp port to listenon")
	seed := flag.String("seed", "", "ip of server to connect to")
	httpPort := flag.String("h", "8080", "http port to listenon")

	flag.Parse()
	fmt.Println("Going to listen on port ", *portPtr)
	fmt.Println("Seed to connect to ", *seed)
	fmt.Println("Going to listen on http port ", *httpPort)

	kvMap = make(map[string]string)
	kvMap["hello"] = "world"

	router := gin.Default()
	router.GET("/kvstore", getInfo)
	router.GET("/kvstore/:key", getValue)
	router.POST("/kvstore", setValue)

	// go server.Listen(*portPtr)
	/* lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *portPtr))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}


	s := grpc.NewServer()
	pb.RegisterKVSeviceServer(s, &grpcserver.Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	} *

	*/

	go grpcserver.StartGrpcServer(portPtr)

	if *seed != "" {
		fmt.Println("mjjjjjj Seed is not nill", *seed)
		// go client.CallServer(*seed)
		go client.CallGrpcServer(*seed)
	}

	router.Run(":" + *httpPort)

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
