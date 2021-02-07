package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func demo(c *gin.Context) {
	c.String(http.StatusOK, "DEMO")
}
