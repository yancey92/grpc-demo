package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"demo.test/grpc-demo/internal/client/api"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	gin.SetMode("release")
	router := api.SetupRouter(gin.Default())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
