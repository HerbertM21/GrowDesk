package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Permitir cualquier origen en desarrollo
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Singleton hub para toda la aplicación
var GlobalHub *Hub

// Inicializar el hub global al importar el paquete
func init() {
	GlobalHub = NewHub()
	// Ejecutar el hub en una goroutine
	go GlobalHub.Run()
}

// ServeWs maneja las solicitudes WebSocket
func ServeWs(c *gin.Context) {
	// Extraer ID del ticket y usuario de los parámetros o header
	ticketID := c.Param("ticketId")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ticket no proporcionado"})
		return
	}

	// En producción, obtener el userID del token de autenticación
	userID := c.GetString("userId")
	if userID == "" {
		// Para desarrollo, permitir un userId en el query param
		userID = c.Query("userId")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			return
		}
	}

	// Actualizar la conexión HTTP a WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Crear un nuevo cliente
	client := NewClient(GlobalHub, conn, ticketID, userID)

	// Registrar el cliente con el hub
	client.Hub.Register <- client

	// Iniciar goroutines para bombear mensajes
	go client.WritePump()
	go client.ReadPump()
} 