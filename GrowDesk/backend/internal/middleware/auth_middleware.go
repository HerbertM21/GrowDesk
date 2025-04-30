package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/pkg/auth"
)

// AuthRequired middleware para verificar que el usuario esté autenticado
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Se requiere autenticación"})
			c.Abort()
			return
		}

		// Extraer token del formato "Bearer {token}"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		// Validar token
		claims, err := auth.ValidateToken(tokenParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido: " + err.Error()})
			c.Abort()
			return
		}

		// Verificar que no sea un token de refresh
		if claims.TokenType == "refresh" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de refresh no válido para acceso"})
			c.Abort()
			return
		}

		// Guardar información del usuario en el contexto
		c.Set("userId", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("firstName", claims.FirstName)
		c.Set("lastName", claims.LastName)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// RoleRequired middleware para verificar que el usuario tenga uno de los roles especificados
func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Este middleware debe usarse después de AuthRequired
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Se requiere autenticación"})
			c.Abort()
			return
		}

		// Verificar si el rol del usuario está en la lista de roles permitidos
		userRoleStr := userRole.(string)
		allowed := false
		for _, role := range roles {
			if userRoleStr == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos suficientes para esta acción"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth middleware para autenticación opcional
// Útil para endpoints que pueden ser accedidos con o sin autenticación
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Continuar sin autenticación
			c.Next()
			return
		}

		// Verificar formato del token (Bearer token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Token con formato incorrecto, pero continuamos sin autenticar
			c.Next()
			return
		}

		tokenString := parts[1]

		// Intentar validar token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			// Token inválido, pero continuamos sin autenticar
			c.Next()
			return
		}

		// Almacenar datos del usuario en el contexto
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("firstName", claims.FirstName)
		c.Set("lastName", claims.LastName)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// ExtractBearerToken extrae el token JWT de un header Authorization
func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("header de autorización vacío")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("formato de token inválido")
	}

	return parts[1], nil
}
