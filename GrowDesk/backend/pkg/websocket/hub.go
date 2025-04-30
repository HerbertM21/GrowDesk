package websocket

import (
	"fmt"
	"log"
	"sync"
)

// BroadcastMessage representa un mensaje para ser distribuido a múltiples conexiones
type BroadcastMessage struct {
	OrganizationID string
	Message        []byte
	Sender         *Connection
	TargetUsers    []string               // IDs de usuarios específicos a los que enviar
	Filters        map[string]interface{} // Filtros adicionales
}

// Hub mantiene un conjunto de conexiones activas y transmite mensajes
type Hub struct {
	// Conexiones registradas, mapeadas por organizaciónID
	connections map[string]map[*Connection]bool

	// Conexiones de widgets, mapeadas por widgetID
	widgetConnections map[string]*Connection

	// Canal para registrar nuevas conexiones
	Register chan *Connection

	// Canal para eliminar conexiones
	Unregister chan *Connection

	// Canal para mensajes entrantes que serán transmitidos
	Broadcast chan *BroadcastMessage

	// Mutex para operaciones seguras en el mapa de conexiones
	mu sync.RWMutex
}

// NewHub crea un nuevo hub de conexiones
func NewHub() *Hub {
	return &Hub{
		connections:       make(map[string]map[*Connection]bool),
		widgetConnections: make(map[string]*Connection),
		Register:          make(chan *Connection),
		Unregister:        make(chan *Connection),
		Broadcast:         make(chan *BroadcastMessage),
	}
}

// Run inicia el hub para manejar registros, desregistros y broadcasts
func (h *Hub) Run() {
	for {
		select {
		case connection := <-h.Register:
			h.registerConnection(connection)

		case connection := <-h.Unregister:
			h.unregisterConnection(connection)

		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerConnection registra una nueva conexión en el hub
func (h *Hub) registerConnection(c *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Para conexiones de widgets, las almacenamos por su ID único
	if c.Type == ConnectionTypeWidget && c.WidgetID != "" {
		h.widgetConnections[c.WidgetID] = c
		log.Printf("Widget %s conectado", c.WidgetID)
		return
	}

	// Para conexiones de usuarios, las organizamos por organización
	if c.OrganizationID != "" {
		if _, ok := h.connections[c.OrganizationID]; !ok {
			h.connections[c.OrganizationID] = make(map[*Connection]bool)
		}
		h.connections[c.OrganizationID][c] = true
		log.Printf("Usuario %s conectado a la organización %s", c.UserID, c.OrganizationID)
	} else {
		log.Printf("Conexión registrada sin OrganizationID: %s", c.ID)
	}
}

// unregisterConnection elimina una conexión del hub
func (h *Hub) unregisterConnection(c *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Para conexiones de widgets
	if c.Type == ConnectionTypeWidget && c.WidgetID != "" {
		delete(h.widgetConnections, c.WidgetID)
		log.Printf("Widget %s desconectado", c.WidgetID)
		close(c.Send)
		return
	}

	// Para conexiones de usuarios
	if c.OrganizationID != "" {
		if _, ok := h.connections[c.OrganizationID]; ok {
			if _, ok := h.connections[c.OrganizationID][c]; ok {
				delete(h.connections[c.OrganizationID], c)
				close(c.Send)

				// Si no hay más conexiones para esta organización, limpiamos el mapa
				if len(h.connections[c.OrganizationID]) == 0 {
					delete(h.connections, c.OrganizationID)
				}

				log.Printf("Usuario %s desconectado de la organización %s", c.UserID, c.OrganizationID)
			}
		}
	}
}

// broadcastMessage envía un mensaje a las conexiones adecuadas
func (h *Hub) broadcastMessage(bm *BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Si hay usuarios específicos como objetivo
	if len(bm.TargetUsers) > 0 {
		h.sendToSpecificUsers(bm)
		return
	}

	// Si tiene una organización específica, enviar a todos en esa organización
	if orgConns, ok := h.connections[bm.OrganizationID]; ok {
		for conn := range orgConns {
			// No enviar al remitente original
			if conn != bm.Sender {
				select {
				case conn.Send <- bm.Message:
				default:
					h.mu.RUnlock()
					h.Unregister <- conn
					h.mu.RLock()
				}
			}
		}
	}

	// Manejar posibles envíos a conexiones de widgets
	h.sendToWidgets(bm)
}

// sendToSpecificUsers envía un mensaje a usuarios específicos
func (h *Hub) sendToSpecificUsers(bm *BroadcastMessage) {
	if orgConns, ok := h.connections[bm.OrganizationID]; ok {
		for conn := range orgConns {
			for _, userID := range bm.TargetUsers {
				if conn.UserID == userID {
					select {
					case conn.Send <- bm.Message:
					default:
						go func(c *Connection) {
							h.Unregister <- c
						}(conn)
					}
					break
				}
			}
		}
	}
}

// sendToWidgets envía mensajes a conexiones de widgets cuando sea necesario
func (h *Hub) sendToWidgets(bm *BroadcastMessage) {
	// Si hay alguna lógica específica para determinar qué widgets deben recibir el mensaje
	// Por ejemplo, si el mensaje está relacionado con un ticket específico

	// Ejemplo: si el mensaje es para un widget específico (extraer widgetID de los filtros)
	if widgetID, ok := bm.Filters["widgetID"].(string); ok {
		if conn, exists := h.widgetConnections[widgetID]; exists {
			select {
			case conn.Send <- bm.Message:
			default:
				go func(c *Connection) {
					h.Unregister <- c
				}(conn)
			}
		}
	}
}

// GetConnectionCount devuelve el número total de conexiones en el hub
func (h *Hub) GetConnectionCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0
	for _, connections := range h.connections {
		count += len(connections)
	}
	count += len(h.widgetConnections)

	return count
}

// GetConnectionCountByOrg devuelve el número de conexiones por organización
func (h *Hub) GetConnectionCountByOrg(orgID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if connections, ok := h.connections[orgID]; ok {
		return len(connections)
	}
	return 0
}

// GetWidgetConnection obtiene una conexión de widget específica por ID
func (h *Hub) GetWidgetConnection(widgetID string) (*Connection, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if conn, ok := h.widgetConnections[widgetID]; ok {
		return conn, nil
	}
	return nil, fmt.Errorf("conexión de widget no encontrada para ID: %s", widgetID)
}

// SendToUser envía un mensaje a un usuario específico
func (h *Hub) SendToUser(orgID string, userID string, message []byte) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if orgConnections, ok := h.connections[orgID]; ok {
		for conn := range orgConnections {
			if conn.UserID == userID {
				select {
				case conn.Send <- message:
					return true
				default:
					return false
				}
			}
		}
	}
	return false
}
