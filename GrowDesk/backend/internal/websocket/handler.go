package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/data"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
)

// Upgrader para conexiones WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Permitir todas las origenes para desarrollo
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ChatHandler maneja las conexiones WebSocket para el chat de tickets
func ChatHandler(store *data.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extraer el ID del ticket desde la URL
		// Formato de URL: /api/ws/chat/:ticketID
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, "URL de chat inválida", http.StatusBadRequest)
			return
		}

		ticketID := pathParts[len(pathParts)-1]
		fmt.Printf("Solicitud de conexión WebSocket para el ticket: %s\n", ticketID)

		// Verificar si el ticket existe
		ticket, err := store.GetTicket(ticketID)
		if err != nil {
			http.Error(w, "Ticket no encontrado", http.StatusNotFound)
			return
		}

		// Actualizar la conexión a WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("Error al actualizar a WebSocket: %v\n", err)
			return
		}

		// Agregar la conexión al almacén
		connectionID := store.AddWSConnection(ticketID, conn)

		// Enviar mensaje de bienvenida
		welcomeMsg := models.WebSocketMessage{
			Type: "connection_established",
			Data: map[string]interface{}{
				"id":        fmt.Sprintf("system-%d", time.Now().Unix()),
				"content":   "Conexión establecida",
				"isClient":  false,
				"timestamp": time.Now(),
			},
		}

		if err := conn.WriteJSON(welcomeMsg); err != nil {
			fmt.Printf("Error al enviar el mensaje de bienvenida: %v\n", err)
		}

		// Enviar historial de mensajes
		if ticket != nil && len(ticket.Messages) > 0 {
			historyMsg := models.WebSocketMessage{
				Type:     "message_history",
				TicketID: ticketID,
				Messages: ticket.Messages,
			}

			if err := conn.WriteJSON(historyMsg); err != nil {
				fmt.Printf("Error al enviar el historial de mensajes: %v\n", err)
			}
		}

		// Manejar mensajes entrantes en una goroutine
		go handleMessages(conn, store, ticketID, connectionID)
	}
}

// handleMessages procesa mensajes entrantes de WebSocket
func handleMessages(conn *websocket.Conn, store *data.Store, ticketID, connectionID string) {
	defer func() {
		conn.Close()
		store.RemoveWSConnection(ticketID, connectionID)
		fmt.Printf("Conexión WebSocket cerrada para el ticket: %s\n", ticketID)
	}()

	// Configurar el manejador de ping
	conn.SetPingHandler(func(appData string) error {
		err := conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(10*time.Second))
		if err != nil {
			fmt.Printf("Error al enviar pong: %v\n", err)
			return err
		}
		return nil
	})

	// Mantener la conexión viva con pings
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				fmt.Printf("Error al enviar ping: %v\n", err)
				return
			}
		}
	}()

	// Procesar mensajes entrantes
	for {
		// Leer mensaje
		_, data, err := conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Error al leer mensaje: %v\n", err)
			}
			break
		}

		// Parsear mensaje
		var message struct {
			Type     string                 `json:"type"`
			TicketID string                 `json:"ticketId,omitempty"`
			UserID   string                 `json:"userId,omitempty"`
			Data     map[string]interface{} `json:"data,omitempty"`
			Content  string                 `json:"content,omitempty"`
			IsClient bool                   `json:"isClient,omitempty"`
		}

		if err := json.Unmarshal(data, &message); err != nil {
			fmt.Printf("Error al analizar mensaje: %v\n", err)
			continue
		}

		fmt.Printf("Mensaje recibido: %s\n", string(data))

		// Manejar mensaje basado en el tipo
		switch message.Type {
		case "identify":
			// El cliente está identificándose
			resp := models.WebSocketMessage{
				Type:     "identify_success",
				TicketID: ticketID,
				Data: map[string]interface{}{
					"message": "Identificación exitosa",
					"userId":  message.UserID,
				},
			}

			if err := conn.WriteJSON(resp); err != nil {
				fmt.Printf("Error al enviar la respuesta de identificación: %v\n", err)
			}

		case "new_message":
			// Extraer contenido del mensaje
			var content string
			var isClient bool
			var userName string

			if message.Data != nil {
				// Intentar extraer del campo de datos
				if c, ok := message.Data["content"].(string); ok {
					content = c
				}
				if ic, ok := message.Data["isClient"].(bool); ok {
					isClient = ic
				}
				if un, ok := message.Data["userName"].(string); ok {
					userName = un
				}
			} else {
				// Campos directos
				content = message.Content
				isClient = message.IsClient
			}

			// Validar contenido
			if content == "" {
				continue
			}

			// Crear objeto de mensaje
			newMessage := models.Message{
				ID:        fmt.Sprintf("MSG-%s", time.Now().Format("20060102150405.000")),
				Content:   content,
				IsClient:  isClient,
				Timestamp: time.Now(),
				CreatedAt: time.Now(),
				UserName:  userName,
			}

			// Agregar mensaje al ticket
			_, err := store.AddMessageToTicket(ticketID, newMessage)
			if err != nil {
				fmt.Printf("Error agregando mensaje al ticket: %v\n", err)
				continue
			}

			// Enviar confirmación
			resp := models.WebSocketMessage{
				Type:     "message_received",
				TicketID: ticketID,
				Data:     newMessage,
			}

			if err := conn.WriteJSON(resp); err != nil {
				fmt.Printf("Error al enviar la confirmación de mensaje: %v\n", err)
			}

			// Broadcast a todos los clientes
			store.BroadcastMessage(ticketID, newMessage)
		}
	}
}
