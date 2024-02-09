package metrics

import (
	"strings"

	"demo.test/grpc-demo/internal/client/logic"
	"github.com/gin-gonic/gin"
)

func Network(c *gin.Context) {
	pretty := strings.ToLower(c.Query("pretty")) == "true"
	metrics, err := logic.MockRPC()
	if pretty {
		if err != nil {
			c.IndentedJSON(500, err.Error())
		} else {
			c.IndentedJSON(200, metrics)
		}
	} else {
		if err != nil {
			c.JSON(500, err.Error())
		} else {
			c.JSON(200, metrics)
		}
	}
}
