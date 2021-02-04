package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func (s *restApiAdapter) health(c *gin.Context) {
	h, err := s.ts.IsServiceHealthy(c.Request.Context())
	if !h {
		c.String(http.StatusInternalServerError, "Not healthy cause of %s", err)
	}
	c.String(http.StatusOK, "OK")
}

func (s *restApiAdapter) readiness(c *gin.Context) {
	h, err := s.ts.IsServiceHealthy(c.Request.Context())
	if !h {
		c.String(http.StatusInternalServerError, "Not ready yet %s", err)
	}
	c.String(http.StatusOK, "READY %s", os.Getenv("POD_NAME"))
}
