package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/internal/controllers"
	"github.com/hmdev/GrowDesk/backend/internal/middleware"
	"github.com/hmdev/GrowDesk/backend/pkg/websocket"
)

// SetupRoutes configura todas las rutas de la API
func SetupRoutes(router *gin.Engine) {
	// Middleware global de CORS
	router.Use(middleware.CORS())

	// Rutas públicas
	publicRoutes := router.Group("/api")
	{
		// Autenticación
		publicRoutes.POST("/auth/login", controllers.Login)
		publicRoutes.POST("/auth/refresh", controllers.RefreshToken)

		// Estado del sistema
		publicRoutes.GET("/health", controllers.HealthCheck)
	}

	// Rutas para el widget (requieren API key)
	widgetRoutes := router.Group("/api/widget")
	widgetRoutes.Use(middleware.WidgetAuth())
	{
		// Configuración del widget
		widgetRoutes.GET("/config", controllers.GetWidgetConfig)

		// Tickets del widget
		widgetRoutes.POST("/tickets", controllers.CreateWidgetTicket)
		widgetRoutes.GET("/tickets/:ticketId", controllers.GetWidgetTicket)
		widgetRoutes.POST("/tickets/:ticketId/messages", controllers.SendWidgetMessage)

		// FAQs
		widgetRoutes.GET("/faqs", controllers.GetWidgetFaqs)
	}

	// Rutas autenticadas (requieren token JWT)
	apiRoutes := router.Group("/api")
	apiRoutes.Use(middleware.AuthRequired())
	{
		// Usuarios
		apiRoutes.GET("/users/me", controllers.GetCurrentUser)

		// Tickets
		tickets := apiRoutes.Group("/tickets")
		{
			tickets.GET("", controllers.GetTickets)
			tickets.POST("", controllers.CreateTicket)
			tickets.GET("/:id", controllers.GetTicket)
			tickets.PUT("/:id", controllers.UpdateTicket)
			tickets.DELETE("/:id", middleware.RoleRequired("admin"), controllers.DeleteTicket)

			// Mensajes
			tickets.GET("/:id/messages", controllers.GetMessages)
			tickets.POST("/:id/messages", controllers.CreateMessage)

			// Asignación de tickets
			tickets.POST("/:id/assign", middleware.RoleRequired("admin", "agent"), controllers.AssignTicket)
			tickets.POST("/:id/unassign", middleware.RoleRequired("admin", "agent"), controllers.UnassignTicket)

			// Cambiar estado/prioridad
			tickets.POST("/:id/status", controllers.UpdateTicketStatus)
			tickets.POST("/:id/priority", controllers.UpdateTicketPriority)
		}

		// Usuarios (admin)
		users := apiRoutes.Group("/users")
		users.Use(middleware.RoleRequired("admin"))
		{
			users.GET("", controllers.GetUsers)
			users.POST("", controllers.CreateUser)
			users.GET("/:id", controllers.GetUser)
			users.PUT("/:id", controllers.UpdateUser)
			users.DELETE("/:id", controllers.DeleteUser)
		}

		// Categorías
		categories := apiRoutes.Group("/categories")
		{
			categories.GET("", controllers.GetCategories)
			categories.POST("", middleware.RoleRequired("admin"), controllers.CreateCategory)
			categories.PUT("/:id", middleware.RoleRequired("admin"), controllers.UpdateCategory)
			categories.DELETE("/:id", middleware.RoleRequired("admin"), controllers.DeleteCategory)
		}

		// FAQs
		faqs := apiRoutes.Group("/faqs")
		{
			faqs.GET("", controllers.GetAllFaqs)
			faqs.GET("/:id", controllers.GetFaq)
			faqs.POST("", middleware.RoleRequired("admin", "agent"), controllers.CreateFaq)
			faqs.PUT("/:id", middleware.RoleRequired("admin", "agent"), controllers.UpdateFaq)
			faqs.DELETE("/:id", middleware.RoleRequired("admin"), controllers.DeleteFaq)
			faqs.POST("/:id/publish", middleware.RoleRequired("admin", "agent"), controllers.ToggleFaqPublish)
		}

		// Configuración de widgets
		widgets := apiRoutes.Group("/widgets")
		widgets.Use(middleware.RoleRequired("admin"))
		{
			widgets.GET("", controllers.GetWidgets)
			widgets.POST("", controllers.CreateWidget)
			widgets.GET("/:id", controllers.GetWidget)
			widgets.PUT("/:id", controllers.UpdateWidget)
			widgets.DELETE("/:id", controllers.DeleteWidget)
			widgets.POST("/:id/generate-key", controllers.GenerateWidgetApiKey)
		}

		// Estadísticas
		apiRoutes.GET("/stats/overview", middleware.RoleRequired("admin", "agent"), controllers.GetStatsOverview)
		apiRoutes.GET("/stats/tickets", middleware.RoleRequired("admin", "agent"), controllers.GetTicketsStats)
	}

	// WebSocket para notificaciones en tiempo real
	router.GET("/ws", middleware.AuthRequired(), websocket.ServeWs)
}
