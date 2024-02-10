package main

import (
	"strconv"

	"demo.test/grpc-demo/internal/client"
	"demo.test/grpc-demo/internal/client/api"
	"github.com/gin-gonic/gin"
)

func main() {

	client.InitMyFlagSet()
	client.LogSet()

	// http
	if client.HttpPort != 0 {
		api.SetupRouter(gin.Default()).Run(":" + strconv.Itoa(client.HttpPort))
	}

}
