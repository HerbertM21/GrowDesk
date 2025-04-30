package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/models"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
)

// CheckWidgetAuth verifica que el widget tenga credenciales válidas
func CheckWidgetAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		widgetID := c.GetHeader("X-Widget-ID")
		apiKey := c.GetHeader("X-Widget-Token")

		if widgetID == "" || apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Se requiere ID y token del widget"})
			c.Abort()
			return
		}

		// En modo de desarrollo sin BD, permitir cualquier credencial
		db := database.GetDB()
		if db == nil {
			c.Set("widgetID", widgetID)
			c.Set("widgetName", "Widget de desarrollo")
			c.Next()
			return
		}

		// Verificar en la base de datos
		var widget models.WidgetConfig

		result := db.Where("id = ? AND api_key = ?", widgetID, apiKey).First(&widget)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales de widget inválidas"})
			c.Abort()
			return
		}

		// Verificar dominio origen si está configurado
		if len(widget.AllowedDomains) > 0 {
			origin := c.GetHeader("Origin")
			if origin != "" {
				allowed := false
				for _, domain := range widget.AllowedDomains {
					if strings.Contains(origin, domain) {
						allowed = true
						break
					}
				}

				if !allowed {
					c.JSON(http.StatusForbidden, gin.H{"error": "Dominio no permitido"})
					c.Abort()
					return
				}
			}
		}

		// Agregar info del widget al contexto
		c.Set("widgetID", widget.ID)
		c.Set("widgetName", widget.Name)
		c.Set("widgetConfig", widget)

		c.Next()
	}
}
