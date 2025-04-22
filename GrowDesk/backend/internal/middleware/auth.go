package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/pkg/auth"
)

// AuthMiddleware revisa si el usuario está autenticado
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		tokenDetails, err := auth.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", tokenDetails.UserID)
		c.Set("email", tokenDetails.Email)
		c.Set("role", tokenDetails.Role)

		c.Next()
	}
}

// RoleMiddleware revisa si el usuario tiene el rol requerido
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role format"})
			c.Abort()
			return
		}

		for _, role := range roles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
		c.Abort()
	}
}

// AdminMiddleware
func AdminMiddleware() gin.HandlerFunc {
	return RoleMiddleware("admin")
}

// AgentMiddleware revisa si el usuario es un agente o administrador
func AgentMiddleware() gin.HandlerFunc {
	return RoleMiddleware("agent", "admin")
}
