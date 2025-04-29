package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Tiempo permitido para escribir un mensaje al par
	writeWait = 10 * time.Second

	// Tiempo permitido para leer el siguiente mensaje del par
	pongWait = 60 * time.Second

	// Enviar pings al par con este período. Debe ser menor que pongWait
	pingPeriod = (pongWait * 9) / 10

	// Tamaño máximo de mensaje permitido desde el par
	maxMessageSize = 512
)

// Client representa una conexión WebSocket
type Client struct {
	// Hub es el hub al que pertenece este cliente
	Hub *Hub

	// Conn es la conexión WebSocket
	Conn *websocket.Conn

	// TicketID es el ID del ticket al que pertenece la sala
	TicketID string

	// UserID es el ID del usuario conectado
	UserID string

	// Buffer de mensajes enviados
	Send chan Message
}

// NewClient crea un nuevo cliente WebSocket
func NewClient(hub *Hub, conn *websocket.Conn, ticketID string, userID string) *Client {
	return &Client{
		Hub:      hub,
		Conn:     conn,
		TicketID: ticketID,
		UserID:   userID,
		Send:     make(chan Message, 256),
	}
}

// ReadPump bombea mensajes desde la conexión WebSocket al hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
		log.Printf("Cliente WebSocket desconectado: %s", c.TicketID)
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Esto soluciona el problema de "Formato de mensaje no reconocido"
	// enviando un mensaje explícito de confirmación
	// para los mensajes de identificación inicial
	identifySuccessMsg := Message{
		Type:     "identify_success",
		TicketID: c.TicketID,
		Data: map[string]interface{}{
			"message": "Identificación exitosa",
			"userId":  c.UserID,
			"status":  "connected",
		},
	}
	// Enviar confirmación inmediata
	if err := c.Conn.WriteJSON(identifySuccessMsg); err != nil {
		log.Printf("Error al enviar confirmación de identificación: %v", err)
	}

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error en conexión WebSocket: %v", err)
			}
			break
		}

		// Depuración del mensaje recibido
		log.Printf("Mensaje WebSocket recibido para ticket %s: %s", c.TicketID, string(message))

		// Analizar el mensaje como JSON genérico
		var msgMap map[string]interface{}
		if err := json.Unmarshal(message, &msgMap); err != nil {
			log.Printf("Error al deserializar mensaje JSON: %v", err)

			// Enviar mensaje de error al cliente
			c.Send <- Message{
				Type:     "error",
				TicketID: c.TicketID,
				Data:     "Formato JSON inválido",
			}
			continue
		}

		// Extraer información relevante del mensaje
		msgType, _ := msgMap["type"].(string)
		msgTicketID, _ := msgMap["ticketId"].(string)

		// Manejar el mensaje según su tipo
		switch msgType {
		case "identify":
			// Mensaje de identificación del cliente
			userID, _ := msgMap["userId"].(string)
			if userID != "" {
				log.Printf("Actualizado ID de usuario de %s a %s", c.UserID, userID)
				c.UserID = userID
			}

			// Responder con confirmación de identificación
			c.Send <- Message{
				Type:     "identify_success",
				TicketID: c.TicketID,
				Data: map[string]interface{}{
					"message": "Identificación exitosa",
					"userId":  c.UserID,
				},
			}

		case "new_message":
			// Mensaje normal de chat
			var messageData interface{}

			// Verificar dónde están los datos del mensaje (data o message)
			if data, ok := msgMap["data"].(map[string]interface{}); ok {
				messageData = data
			} else if msg, ok := msgMap["message"].(map[string]interface{}); ok {
				messageData = msg
			} else {
				// Si no contiene datos estructurados, usar el mensaje completo
				delete(msgMap, "type")
				delete(msgMap, "ticketId")
				messageData = msgMap
			}

			// Reenviar el mensaje al hub en el formato estándar
			c.Hub.BroadcastToTicket(c.TicketID, "new_message", messageData)

		default:
			// Para cualquier otro tipo de mensaje
			log.Printf("Recibido mensaje de tipo %s", msgType)

			// Si el mensaje no tiene un tipo conocido, intentar reenviar como está
			// pero adaptándolo al formato estándar del hub
			if msgTicketID == c.TicketID || msgTicketID == "" {
				var messageData interface{}

				// Extraer los datos relevantes
				if data, ok := msgMap["data"].(map[string]interface{}); ok {
					messageData = data
				} else {
					// Usar todo el mensaje menos el tipo y ticketId
					cleanMsg := make(map[string]interface{})
					for k, v := range msgMap {
						if k != "type" && k != "ticketId" {
							cleanMsg[k] = v
						}
					}
					messageData = cleanMsg
				}

				// Reenviar al hub
				c.Hub.BroadcastToTicket(c.TicketID, msgType, messageData)
			} else {
				log.Printf("Error: mensaje para ticket incorrecto. Esperado %s, recibido %s",
					c.TicketID, msgTicketID)

				c.Send <- Message{
					Type:     "error",
					TicketID: c.TicketID,
					Data:     "ID de ticket no coincide",
				}
			}
		}
	}
}

// WritePump bombea mensajes desde el hub a la conexión WebSocket
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// El hub cerró el canal
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Enviar el mensaje como JSON
			err := c.Conn.WriteJSON(message)
			if err != nil {
				log.Printf("error al enviar mensaje: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
