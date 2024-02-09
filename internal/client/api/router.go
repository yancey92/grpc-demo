package api

import (
	"demo.test/grpc-demo/internal/client/api/metrics"
	"github.com/gin-gonic/gin"
)

func SetupRouter(engine *gin.Engine) *gin.Engine {
	engine.GET("/", welcome)
	engine.GET("/ping", handlerPing)

	// api v1 group
	routerGroupV1 := engine.Group("/api/v1")
	{
		routerGroupV1.GET("/metrics/grpc/network", metrics.Network)
	}
	return engine
}

func welcome(c *gin.Context) {
	c.String(200, "Welcome ...")
}

func handlerPing(c *gin.Context) {
	c.String(200, "pong")
}
