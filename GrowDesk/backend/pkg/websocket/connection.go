package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Tiempo permitido para escribir un mensaje al peer
	writeWait = 10 * time.Second

	// Tiempo permitido para leer el siguiente mensaje del peer
	pongWait = 60 * time.Second

	// Enviar pings al peer con esta periodicidad
	pingPeriod = (pongWait * 9) / 10

	// Tamaño máximo del mensaje permitido desde el peer
	maxMessageSize = 512 * 1024 // 512KB
)

// ConnectionType define el tipo de conexión WebSocket
type ConnectionType string

const (
	// ConnectionTypeUser para conexiones de usuarios del sistema
	ConnectionTypeUser ConnectionType = "user"

	// ConnectionTypeWidget para conexiones desde widgets externos
	ConnectionTypeWidget ConnectionType = "widget"
)

// Mensaje representa la estructura de los mensajes WebSocket
type Message struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// Connection encapsula una conexión WebSocket
type Connection struct {
	ID             string
	Conn           *websocket.Conn
	Type           ConnectionType
	OrganizationID string
	UserID         string
	WidgetID       string
	Send           chan []byte
	Hub            *Hub
	mu             sync.Mutex
	Metadata       map[string]interface{}
}

// NewConnection crea una nueva conexión WebSocket
func NewConnection(conn *websocket.Conn, connType ConnectionType, hub *Hub) *Connection {
	return &Connection{
		ID:       uuid.New().String(),
		Conn:     conn,
		Type:     connType,
		Send:     make(chan []byte, 256),
		Hub:      hub,
		Metadata: make(map[string]interface{}),
	}
}

// ReadPump bombea mensajes desde la conexión WebSocket al hub
func (c *Connection) ReadPump() {
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

		// Procesar el mensaje
		c.processIncomingMessage(message)
	}
}

// WritePump bombea mensajes desde el hub a la conexión WebSocket
func (c *Connection) WritePump() {
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

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Agregar mensajes en cola al actual mensaje WebSocket
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
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

// SendJSON envía un mensaje codificado en JSON a través de la conexión
func (c *Connection) SendJSON(msgType string, data interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	msg := Message{
		Type:      msgType,
		Data:      data,
		Timestamp: time.Now(),
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error al serializar mensaje: %v", err)
	}

	select {
	case c.Send <- jsonData:
		return nil
	default:
		return fmt.Errorf("buffer de conexión lleno")
	}
}

// processIncomingMessage procesa los mensajes entrantes
func (c *Connection) processIncomingMessage(data []byte) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("Error al deserializar mensaje: %v", err)
		return
	}

	// Procesar el mensaje según su tipo
	switch msg.Type {
	case "ping":
		c.SendJSON("pong", nil)
	case "ticket_update":
		// Reenviar al hub para distribución
		c.Hub.Broadcast <- &BroadcastMessage{
			OrganizationID: c.OrganizationID,
			Message:        data,
			Sender:         c,
		}
	case "typing":
		// Notificar que el usuario está escribiendo
		c.Hub.Broadcast <- &BroadcastMessage{
			OrganizationID: c.OrganizationID,
			Message:        data,
			Sender:         c,
		}
	default:
		log.Printf("Tipo de mensaje desconocido: %s", msg.Type)
	}
}

// SetMetadata establece metadatos para la conexión
func (c *Connection) SetMetadata(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Metadata[key] = value
}

// GetMetadata obtiene un valor de metadatos de la conexión
func (c *Connection) GetMetadata(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.Metadata[key]
	return val, ok
}
