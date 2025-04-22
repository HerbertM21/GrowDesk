package websocket

import (
	"log"
	"sync"
)

// Message representa un mensaje que se enviará a través de WebSocket
type Message struct {
	Type     string      `json:"type"`
	TicketID string      `json:"ticketId"`
	Data     interface{} `json:"data"`
}

// Hub mantiene el registro de todas las conexiones activas
type Hub struct {
	// Mapa de conexiones activas por ID de ticket
	Rooms      map[string]map[*Client]bool
	RoomsMutex sync.Mutex
	
	// Canal para enviar mensajes a los clientes
	Broadcast chan Message
	
	// Registrar nuevos clientes
	Register chan *Client
	
	// Dar de baja clientes desconectados
	Unregister chan *Client
}

// NewHub crea un nuevo hub
func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		RoomsMutex: sync.Mutex{},
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run inicia el hub en una goroutine
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.RoomsMutex.Lock()
			// Crear la sala si no existe
			if _, ok := h.Rooms[client.TicketID]; !ok {
				h.Rooms[client.TicketID] = make(map[*Client]bool)
			}
			// Agregar cliente a la sala
			h.Rooms[client.TicketID][client] = true
			h.RoomsMutex.Unlock()
			log.Printf("Cliente registrado en sala: %s, total clientes: %d", client.TicketID, len(h.Rooms[client.TicketID]))
			
		case client := <-h.Unregister:
			h.RoomsMutex.Lock()
			// Verificar si la sala y el cliente existen
			if room, ok := h.Rooms[client.TicketID]; ok {
				if _, ok := room[client]; ok {
					delete(room, client)
					// Cerrar conexión del cliente
					close(client.Send)
					
					// Si no quedan clientes en la sala, eliminarla
					if len(room) == 0 {
						delete(h.Rooms, client.TicketID)
					}
				}
			}
			h.RoomsMutex.Unlock()
			log.Printf("Cliente eliminado de sala: %s", client.TicketID)
			
		case message := <-h.Broadcast:
			h.RoomsMutex.Lock()
			// Enviar mensaje a todos los clientes en la sala
			if room, ok := h.Rooms[message.TicketID]; ok {
				for client := range room {
					select {
					case client.Send <- message:
						// Mensaje enviado exitosamente
					default:
						// Si el canal está bloqueado, eliminar cliente
						close(client.Send)
						delete(room, client)
						
						// Si no quedan clientes en la sala, eliminarla
						if len(room) == 0 {
							delete(h.Rooms, message.TicketID)
						}
					}
				}
			}
			h.RoomsMutex.Unlock()
		}
	}
}

// BroadcastToTicket envía un mensaje a todos los clientes en una sala específica
func (h *Hub) BroadcastToTicket(ticketID string, messageType string, data interface{}) {
	message := Message{
		Type:     messageType,
		TicketID: ticketID,
		Data:     data,
	}
	h.Broadcast <- message
} 