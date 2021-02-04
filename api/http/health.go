package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

//TODO: add services healthchecks db, etc
func health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func readiness(c *gin.Context) {
	c.String(http.StatusOK, "READY %s", os.Getenv("POD_NAME"))
}
