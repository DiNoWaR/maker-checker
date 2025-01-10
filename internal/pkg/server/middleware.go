package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

var defaultAdminToken = os.Getenv("DEFAULT_MODERATOR_TOKEN")

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		if authHeader == "" {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			context.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			context.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != defaultAdminToken {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
			context.Abort()
			return
		}

		context.Next()
	}
}
