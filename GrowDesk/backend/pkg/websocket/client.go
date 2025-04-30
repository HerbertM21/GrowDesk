package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Tiempo máximo de espera para escritura
	writeWait = 10 * time.Second

	// Tiempo máximo de espera para leer el siguiente mensaje pong
	pongWait = 60 * time.Second

	// Enviar pings al cliente con esta periodicidad (debe ser menor que pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Tamaño máximo del mensaje
	maxMessageSize = 1024 * 16
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// En producción aquí se verificaría el origen
		return true
	},
}

// ServeWs maneja las conexiones WebSocket
func ServeWs(c *gin.Context) {
	// Verificar que hay un hub de websocket
	if hub == nil {
		NewHub() // Crear el hub si no existe
	}

	// Actualizar la conexión HTTP a WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket error:", err)
		return
	}

	// Obtener datos de autenticación
	userID := c.GetString("userID")
	ticketID := c.Query("ticketId") // Opcional: para suscribirse a un ticket específico
	role := c.GetString("role")

	// Crear nuevo cliente
	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan interface{}, 256),
		userID:   userID,
		ticketID: ticketID,
		role:     role,
	}

	// Registrar cliente
	client.hub.Register(client)

	// Iniciar rutinas para leer/escribir mensajes
	go client.readPump()
	go client.writePump()

	// Enviar mensaje de bienvenida
	client.send <- Message{
		Type: "welcome",
		Data: gin.H{
			"message": "Conexión WebSocket establecida",
			"userId":  userID,
			"role":    role,
		},
	}
}

// readPump bombea mensajes desde la conexión WebSocket al hub
func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Procesar mensaje recibido
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))

		// Aquí se puede analizar el mensaje si se espera input del cliente
		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err == nil {
			// Ejemplo: Manejar ping del cliente
			if action, ok := data["action"].(string); ok && action == "ping" {
				c.send <- gin.H{"type": "pong", "time": time.Now().Unix()}
			}
		}
	}
}

// writePump bombea mensajes desde el hub a la conexión WebSocket
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// El hub ha cerrado el canal
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Serializar y enviar mensaje
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("WebSocket error al serializar: %v", err)
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(data)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			// Enviar ping periódico para mantener la conexión viva
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
