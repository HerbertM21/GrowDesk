package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CheckWidgetAuth verifica la autenticación para las solicitudes de widget
func CheckWidgetAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		widgetID := c.GetHeader("X-Widget-ID")
		widgetToken := c.GetHeader("X-Widget-Token")

		if widgetID == "" || widgetToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Faltan credenciales del widget"})
			c.Abort()
			return
		}

		// verificar en la base de datos
		// Aquí, permitimos demo-widget con demo-token
		if widgetID == "demo-widget" && widgetToken == "demo-token" {
			c.Set("widgetID", widgetID)
			c.Next()
			return
		}

		// Si no es la demo y no hay verificación de base de datos,
		// podemos rechazar para esta implementación simplificada
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales de widget inválidas"})
		c.Abort()
	}
}
