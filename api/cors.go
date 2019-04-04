package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS ...
func CORS(allowed []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		switch {
		case len(allowed) <= 0:
			setCorsHeaders(c.Writer, "*")
		case isAllowedOrigin(allowed, origin):
			setCorsHeaders(c.Writer, origin)
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func setCorsHeaders(writer http.ResponseWriter, origin string) {
	writer.Header().Set("Access-Control-Allow-Origin", origin)
	writer.Header().Set("Access-Control-Allow-Methods", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, *")
}

func isAllowedOrigin(allowed []string, current string) bool {
	for _, origin := range allowed {
		if strings.Contains(current, origin) {
			return true
		}
	}

	return false
}
