package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//HealthzHandler will be called to get application status
type HealthzHandler func() interface{}

//HealthChecks joins liveness and readiness health checks.
type HealthChecks struct {
	Liveness  HealthzHandler
	Readiness HealthzHandler
}

//DefaultHealthz just check if the app isnt locked.
func DefaultHealthz() HealthzHandler {
	return func() interface{} {
		return true
	}
}

//DefaultHealthChecks returns default liveness and readiness checks.
func DefaultHealthChecks() *HealthChecks {
	return &HealthChecks{
		Liveness:  DefaultHealthz(),
		Readiness: DefaultHealthz(),
	}
}

func (server *Server) liveness() gin.HandlerFunc {
	return func(c *gin.Context) {
		var healthz interface{}
		if server.Healthz.Liveness != nil {
			healthz = server.Healthz.Liveness()
		}
		c.JSON(http.StatusOK, healthz)
	}
}

func (server *Server) readiness() gin.HandlerFunc {
	return func(c *gin.Context) {
		var healthz interface{}
		if server.Healthz.Readiness != nil {
			healthz = server.Healthz.Readiness()
		}
		c.JSON(http.StatusOK, healthz)
	}
}
