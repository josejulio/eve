package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthGetAPI(c *gin.Context) {
	c.Status(http.StatusOK)
}
