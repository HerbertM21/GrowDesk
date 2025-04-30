package middleware

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hmdev/GrowDesk/backend/models"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
)

// WidgetAuth middleware para autenticar solicitudes del widget
func WidgetAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener API key del header o query parameter
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "API key requerida",
			})
			return
		}

		// Verificar origen si está disponible
		origin := c.GetHeader("Origin")
		referer := c.GetHeader("Referer")

		db := database.GetDB()
		if db != nil {
			// Buscar configuración del widget por API key
			var config models.WidgetConfig
			if err := db.Where("api_key = ?", apiKey).First(&config).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status":  "error",
					"message": "API key inválida",
				})
				return
			}

			// Verificar dominio de origen si está presente y hay dominios permitidos configurados
			if (origin != "" || referer != "") && len(config.AllowedDomains) > 0 {
				allowed := false

				// Normalizar y verificar dominio de Origin
				if origin != "" {
					originURL, err := url.Parse(origin)
					if err == nil {
						originDomain := originURL.Hostname()
						allowed = isDomainAllowed(originDomain, config.AllowedDomains)
					}
				}

				// Si no se permitió por Origin, verificar Referer
				if !allowed && referer != "" {
					refererURL, err := url.Parse(referer)
					if err == nil {
						refererDomain := refererURL.Hostname()
						allowed = isDomainAllowed(refererDomain, config.AllowedDomains)
					}
				}

				if !allowed {
					log.Printf("Dominio no permitido: Origin=%s, Referer=%s", origin, referer)
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
						"status":  "error",
						"message": "Dominio de origen no autorizado",
					})
					return
				}
			}

			// Almacenar configuración y ID del widget en el contexto
			c.Set("widgetID", config.ID)
			c.Set("widgetConfig", config)
		} else {
			// En modo desarrollo, aceptar cualquier API key y simular una configuración
			widgetID := "WIDGET-DEV"

			// Si la API key tiene formato "WIDGET-XXX", usarla como ID
			if strings.HasPrefix(apiKey, "WIDGET-") {
				widgetID = apiKey
			}

			c.Set("widgetID", widgetID)
		}

		// Configurar headers CORS
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, Origin, Authorization")

		c.Next()
	}
}

// ValidateWidgetJWT middleware para validar JWT de widget
func ValidateWidgetJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authHeader := c.GetHeader("Authorization")
		tokenString, err := ExtractBearerToken(authHeader)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token no proporcionado o formato inválido",
			})
			return
		}

		// Validar token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verificar algoritmo de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}

			// Aquí se debería utilizar la misma clave secreta que se usó para firmar
			return []byte("tu_clave_secreta"), nil // Debe ser la misma que en auth.SetJWTSecret
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token inválido: " + err.Error(),
			})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token inválido",
			})
			return
		}

		// Extraer claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "No se pudieron extraer los claims del token",
			})
			return
		}

		// Verificar el widgetId en el token
		widgetID, exists := claims["widgetId"]
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token no contiene ID de widget",
			})
			return
		}

		// Almacenar el ID del widget en el contexto
		c.Set("widgetID", widgetID)

		c.Next()
	}
}

// isDomainAllowed verifica si un dominio está permitido
func isDomainAllowed(domain string, allowedDomains []string) bool {
	for _, allowed := range allowedDomains {
		// Permitir subdominios si el dominio permitido comienza con punto
		if strings.HasPrefix(allowed, ".") && (strings.HasSuffix(domain, allowed) || domain == allowed[1:]) {
			return true
		}

		// De lo contrario, exigir coincidencia exacta
		if domain == allowed {
			return true
		}
	}

	return false
}
