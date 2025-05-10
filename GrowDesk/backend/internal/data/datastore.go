package data

import (
	"github.com/gorilla/websocket"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
)

// DataStore define la interfaz para el almacenamiento de datos
type DataStore interface {
	// Métodos para usuarios
	GetUsers() ([]models.User, error)
	GetUser(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user models.User) error
	UpdateUser(user models.User) error
	DeleteUser(id string) error

	// Métodos para tickets
	GetTickets() ([]models.Ticket, error)
	GetTicket(id string) (*models.Ticket, error)
	CreateTicket(ticket models.Ticket) error
	UpdateTicket(ticket models.Ticket) error
	DeleteTicket(id string) error
	AddTicketMessage(ticketID string, message models.Message) error

	// Métodos para categorías
	GetCategories() ([]models.Category, error)
	GetCategory(id string) (*models.Category, error)
	CreateCategory(category models.Category) error
	UpdateCategory(category models.Category) error
	DeleteCategory(id string) error

	// Métodos para FAQs
	GetFAQs() ([]models.FAQ, error)
	GetFAQsByStatus(published bool) ([]models.FAQ, error)
	GetFAQ(id int) (*models.FAQ, error)
	CreateFAQ(faq models.FAQ) error
	UpdateFAQ(faq models.FAQ) error
	DeleteFAQ(id int) error
	ToggleFAQPublish(id int) error

	// Métodos para WebSocket
	AddWSConnection(ticketID string, conn *websocket.Conn) string
	RemoveWSConnection(ticketID, connectionID string)
	BroadcastMessage(ticketID string, message models.Message)
}
