package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/models"
)

// WidgetAuthRequired es un middleware para verificar la autenticación de solicitudes del widget
func WidgetAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el API key del widget desde el header
		widgetKey := c.GetHeader("X-Widget-API-Key")
		if widgetKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key del widget no proporcionada"})
			c.Abort()
			return
		}

		// Obtener la configuración del widget usando el API key
		widgetConfig, err := models.GetWidgetConfigByAPIKey(widgetKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API Key del widget inválida"})
			c.Abort()
			return
		}

		// Verificar el dominio permitido si existe referer
		referer := c.Request.Header.Get("Referer")
		if referer != "" && !isAllowedDomain(referer, widgetConfig.AllowedDomains) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Dominio no autorizado"})
			c.Abort()
			return
		}

		// Verificar la firma si está presente (para operaciones críticas)
		signature := c.GetHeader("X-Widget-Signature")
		timestamp := c.GetHeader("X-Widget-Timestamp")

		if signature != "" && timestamp != "" {
			// Construir el string a verificar (normalmente incluiría el body en una implementación completa)
			stringToSign := timestamp + "." + widgetKey

			// Calcular la firma esperada
			expectedSignature := calculateHMAC(stringToSign, widgetConfig.SecretKey)

			// Verificar la firma
			if signature != expectedSignature {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Firma inválida"})
				c.Abort()
				return
			}

			// Verificar que el timestamp no sea muy antiguo (15 minutos max)
			ts, err := time.Parse(time.RFC3339, timestamp)
			if err != nil || time.Since(ts) > 15*time.Minute {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Timestamp expirado o inválido"})
				c.Abort()
				return
			}
		}

		// Almacenar información del widget en el contexto
		c.Set("widgetId", widgetConfig.ID)
		c.Set("organizationId", widgetConfig.OrganizationID)
		c.Set("widgetConfig", widgetConfig)

		c.Next()
	}
}

// VerificaOrigen comprueba si el referer está permitido para los orígenes CORS
func VerificaOrigen(origen string, origenesPermitidos []string) bool {
	return isAllowedDomain(origen, origenesPermitidos)
}

// isAllowedDomain verifica si un referer proviene de un dominio permitido
func isAllowedDomain(referer string, allowedDomains []string) bool {
	// Si no hay dominios permitidos configurados, denegar todo
	if len(allowedDomains) == 0 {
		return false
	}

	// Si hay un wildcard en los dominios permitidos, permitir todos
	for _, domain := range allowedDomains {
		if domain == "*" {
			return true
		}
	}

	// Extraer el dominio del referer
	for _, domain := range allowedDomains {
		// La comprobación simple verifica si el referer contiene el dominio permitido
		// En una implementación real, se debería utilizar análisis de URL más robusto
		if strings.Contains(referer, domain) {
			return true
		}
	}

	return false
}

// calculateHMAC calcula la firma HMAC-SHA256 para autenticación del widget
func calculateHMAC(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
