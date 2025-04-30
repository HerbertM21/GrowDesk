package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS middleware para habilitar Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener origen de la solicitud
		origin := c.Request.Header.Get("Origin")

		// Configurar los orígenes permitidos
		allowOrigins := getAllowedOrigins()

		// Si no hay origen o no se han configurado orígenes permitidos
		// permitir cualquier origen en desarrollo
		if origin == "" || len(allowOrigins) == 0 {
			// En entorno de desarrollo, permitir cualquier origen
			if os.Getenv("APP_ENV") != "production" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}
		} else {
			// Verificar si el origen está permitido
			for _, allowedOrigin := range allowOrigins {
				if origin == allowedOrigin || allowedOrigin == "*" {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		// Resto de cabeceras CORS
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// Manejar la solicitud OPTIONS (preflight)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// getAllowedOrigins obtiene la lista de orígenes permitidos
func getAllowedOrigins() []string {
	// Obtener orígenes de variable de entorno
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")

	// Si no está definido, usar valores por defecto
	if corsOrigins == "" {
		// En producción se deben especificar orígenes concretos
		if os.Getenv("APP_ENV") == "production" {
			frontendURL := os.Getenv("FRONTEND_URL")
			if frontendURL != "" {
				return []string{frontendURL}
			}
			// En producción, sin configuración explícita, no permitir CORS
			return []string{}
		} else {
			// En desarrollo, permitir orígenes comunes para desarrollo local
			return []string{
				"http://localhost:3000",
				"http://localhost:5173",
				"http://localhost:5174",
				"http://localhost:8080",
				"http://127.0.0.1:3000",
				"http://127.0.0.1:5173",
				"http://127.0.0.1:5174",
				"http://127.0.0.1:8080",
			}
		}
	}

	// Dividir la cadena en un slice
	origins := strings.Split(corsOrigins, ",")

	// Eliminar espacios en blanco
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	return origins
}
