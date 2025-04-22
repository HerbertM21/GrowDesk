package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/internal/controllers"
	"github.com/hmdev/GrowDesk/backend/internal/middleware"
)

// SetupWidgetRoutes configura solo las rutas relacionadas con el widget
func SetupWidgetRoutes(r *gin.Engine) {
	// Grupo principal de la API
	api := r.Group("/api")

	// Rutas para integración del widget
	widget := api.Group("/widget")
	widget.Use(middleware.CheckWidgetAuth()) // Middleware específico para autenticar widgets
	{
		// Endpoint para crear tickets desde el widget
		widget.POST("/tickets", controllers.CreateTicketFromWidget)

		// Endpoint para enviar mensajes desde el widget
		widget.POST("/messages", controllers.CreateMessageFromWidget)

		// Endpoint para obtener mensajes de un ticket desde el widget
		widget.GET("/tickets/:id/messages", controllers.GetMessagesForWidget)
	}

	// Rutas para administración de widgets (solo para administradores)
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware())  // Middleware para autenticar usuarios
	admin.Use(middleware.AdminMiddleware()) // Middleware para verificar rol de administrador
	{
		widgetConfig := admin.Group("/widget-config")
		{
			widgetConfig.GET("", controllers.GetWidgetConfigs)
			widgetConfig.GET("/:id", controllers.GetWidgetConfig)
			widgetConfig.POST("", controllers.CreateWidgetConfig)
			widgetConfig.PUT("/:id", controllers.UpdateWidgetConfig)
			widgetConfig.POST("/:id/regenerate-key", controllers.RegenerateApiKey)
			widgetConfig.DELETE("/:id", controllers.DeleteWidgetConfig)
		}
	}
}
