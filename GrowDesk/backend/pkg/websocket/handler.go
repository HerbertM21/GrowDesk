package websocket

import (
	"log"
	"net/http"
	"time"

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
			// Si no hay userID, usar uno genérico para desarrollo
			userID = "anonymous"
			log.Printf("ADVERTENCIA: Usuario no autenticado, usando ID genérico: %s", userID)
		}
	}

	log.Printf("Estableciendo conexión WebSocket para ticket: %s, usuario: %s", ticketID, userID)
	log.Printf("Headers de la solicitud WebSocket: %v", c.Request.Header)

	// Actualizar la conexión HTTP a WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error al establecer conexión WebSocket: %v", err)
		return
	}

	// Crear un nuevo cliente
	client := NewClient(GlobalHub, conn, ticketID, userID)

	// Registrar el cliente con el hub
	client.Hub.Register <- client

	// Número de conexiones actuales para este ticket
	numConnections := 0
	client.Hub.RoomsMutex.Lock()
	if room, ok := client.Hub.Rooms[ticketID]; ok {
		numConnections = len(room)
	}
	client.Hub.RoomsMutex.Unlock()
	log.Printf("Nueva conexión WebSocket para ticket: %s", ticketID)
	log.Printf("Total de conexiones para este ticket: %d", numConnections)

	// Enviar mensaje de conexión exitosa
	connectionMessage := Message{
		Type:     "connection_established",
		TicketID: ticketID,
		Data: map[string]interface{}{
			"id":        "system-" + time.Now().UnixMilli(),
			"content":   "Conexión establecida",
			"isClient":  false,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}
	client.Send <- connectionMessage

	// Obtener mensajes históricos del ticket y enviarlos al cliente
	messages := GetTicketMessages(ticketID)
	if len(messages) > 0 {
		log.Printf("Enviando %d mensajes existentes al cliente que acaba de conectarse", len(messages))

		historyMessage := Message{
			Type:     "message_history",
			TicketID: ticketID,
			Data: map[string]interface{}{
				"messages": messages, // Formato correcto: un objeto con campo "messages"
				"count":    len(messages),
			},
		}
		client.Send <- historyMessage
	}

	// Iniciar goroutines para bombear mensajes
	go client.WritePump()
	go client.ReadPump()
}

// GetTicketMessages obtiene los mensajes históricos de un ticket
// En una implementación real, esto buscaría en la base de datos
func GetTicketMessages(ticketID string) []interface{} {
	// Ejemplo de implementación simple con mensajes de demostración
	// En un sistema real, esto buscaría mensajes en la base de datos
	return []interface{}{
		map[string]interface{}{
			"id":        "MSG-" + time.Now().Add(-5*time.Minute).UnixMilli(),
			"content":   "¿Hola, cómo puedo ayudarte?",
			"isClient":  false,
			"timestamp": time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
		},
	}
}
