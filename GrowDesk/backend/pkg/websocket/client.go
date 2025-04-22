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
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Procesar el mensaje recibido
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error al deserializar mensaje: %v", err)
			continue
		}

		// Asegurarse de que el ticketID coincida con la sala
		if msg.TicketID != c.TicketID {
			log.Printf("error: intento de enviar mensaje a sala incorrecta")
			continue
		}

		// Enviar el mensaje al hub para distribución
		c.Hub.Broadcast <- msg
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