package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/internal/controllers"
	"github.com/hmdev/GrowDesk/backend/internal/middleware"
)

// SetupRoutes configura todas las rutas de la API
func SetupRoutes(router *gin.Engine) {
	// Ruta de salud para verificar que el servidor esté en funcionamiento
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Servidor en funcionamiento",
		})
	})

	// Grupo de rutas de autenticación
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)

		// Rutas protegidas que requieren autenticación
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/me", controllers.GetMe)
		}
	}

	// Grupo de rutas para tickets
	tickets := router.Group("/api/tickets")
	tickets.Use(middleware.AuthMiddleware())
	{
		tickets.GET("", func(c *gin.Context) {
			// Implementación provisional
			c.JSON(http.StatusOK, gin.H{"message": "Lista de tickets"})
		})
		tickets.POST("", func(c *gin.Context) {
			// Implementación provisional
			c.JSON(http.StatusCreated, gin.H{"message": "Ticket creado"})
		})
		tickets.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{"id": id, "message": "Detalle de ticket"})
		})
		tickets.PUT("/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{"id": id, "message": "Ticket actualizado"})
		})
		tickets.DELETE("/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{"id": id, "message": "Ticket eliminado"})
		})

		// Rutas para mensajes en tickets
		tickets.GET("/:id/messages", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{"ticketId": id, "messages": []string{}})
		})
		tickets.POST("/:id/messages", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusCreated, gin.H{"ticketId": id, "message": "Mensaje creado"})
		})
	}

	// Rutas para widgets sin autenticación JWT
	widgetApi := router.Group("/api/widget")
	{
		// Endpoint para crear tickets desde el widget
		widgetApi.POST("/tickets", func(c *gin.Context) {
			var ticket struct {
				Title       string                 `json:"title"`
				Description string                 `json:"description"`
				Email       string                 `json:"email"`
				Name        string                 `json:"name"`
				Source      string                 `json:"source"`
				Metadata    map[string]interface{} `json:"metadata"`
			}

			if err := c.ShouldBindJSON(&ticket); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
				return
			}

			// En una implementación real, guardaríamos en la base de datos
			c.JSON(http.StatusCreated, gin.H{
				"id":      "W-" + time.Now().Format("20060102150405"),
				"message": "Ticket creado desde el widget",
				"status":  "open",
			})
		})

		// Endpoint para enviar mensajes desde el widget
		widgetApi.POST("/chat/messages", func(c *gin.Context) {
			var msg struct {
				TicketID string `json:"ticketId"`
				Content  string `json:"content"`
				IsClient bool   `json:"isClient"`
			}

			if err := c.ShouldBindJSON(&msg); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
				return
			}

			// En una implementación real, guardaríamos en la base de datos
			c.JSON(http.StatusCreated, gin.H{
				"id":       "MSG-" + time.Now().Format("20060102150405"),
				"ticketId": msg.TicketID,
				"status":   "sent",
			})
		})

		// Endpoint para obtener mensajes de un ticket desde el widget
		widgetApi.GET("/chat/messages/:id", func(c *gin.Context) {
			ticketID := c.Param("id")

			// En una implementación real, obtendríamos los mensajes de la base de datos
			c.JSON(http.StatusOK, gin.H{
				"messages": []gin.H{
					{
						"id":        "server-msg-1",
						"content":   "Hola, ¿en qué podemos ayudarte?",
						"isClient":  false,
						"timestamp": time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
					},
					{
						"id":        "client-msg-1",
						"content":   "Necesito ayuda con mi cuenta",
						"isClient":  true,
						"timestamp": time.Now().Add(-3 * time.Minute).Format(time.RFC3339),
					},
					{
						"id":        "server-msg-2",
						"content":   "Claro, dime qué problema tienes",
						"isClient":  false,
						"timestamp": time.Now().Add(-1 * time.Minute).Format(time.RFC3339),
					},
				},
				"ticketId": ticketID,
			})
		})
	}

	// Grupo de rutas para usuarios (solo para administradores)
	users := router.Group("/api/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Lista de usuarios"})
		})
		users.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{"id": id, "message": "Detalle de usuario"})
		})
		users.PUT("/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{"id": id, "message": "Usuario actualizado"})
		})
		users.DELETE("/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{"id": id, "message": "Usuario eliminado"})
		})
	}

	// Chatbox websocket
	router.GET("/api/ws/chat/:ticketId", controllers.WebSocketHandler)
}
