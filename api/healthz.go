package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//HealthzHandler will be called to get application status
type HealthzHandler func() interface{}

var (
	//DefautlHealthz just check if the app isnt locked.
	DefautlHealthz = func() interface{} {
		return true
	}
)

func (server *Server) healtz() gin.HandlerFunc {
	return func(c *gin.Context) {
		var healthz interface{}
		if server.Healthz != nil {
			healthz = server.Healthz()
		}
		c.JSON(http.StatusOK, healthz)
	}
}
