package models

import (
	"time"
)

// User representa un usuario en el sistema
type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Role       string    `json:"role"`
	Department string    `json:"department,omitempty"`
	Active     bool      `json:"active"`
	Password   string    `json:"password,omitempty"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty"`
	Position   string    `json:"position,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	Language   string    `json:"language,omitempty"`
}

// LoginRequest representa los datos de la solicitud de inicio de sesión
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest representa los datos de la solicitud de registro
type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// AuthResponse representa la respuesta enviada después de una autenticación exitosa
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Ticket representa un ticket de soporte
type Ticket struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Priority    string    `json:"priority,omitempty"`
	Category    string    `json:"category,omitempty"`
	AssignedTo  string    `json:"assignedTo,omitempty"`
	CreatedBy   string    `json:"createdBy,omitempty"`
	Description string    `json:"description,omitempty"`
	Customer    Customer  `json:"customer"`
	Messages    []Message `json:"messages,omitempty"`
}

// Customer representa a un cliente de un ticket
type Customer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Message representa un mensaje en un ticket de soporte
type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	IsClient  bool      `json:"isClient"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UserName  string    `json:"userName,omitempty"`
	UserEmail string    `json:"userEmail,omitempty"`
}

// NewMessageRequest representa una solicitud para agregar un nuevo mensaje
type NewMessageRequest struct {
	Content   string `json:"content"`
	IsClient  bool   `json:"isClient"`
	UserName  string `json:"userName,omitempty"`
	UserEmail string `json:"userEmail,omitempty"`
}

// TicketUpdateRequest representa una solicitud para actualizar un ticket
type TicketUpdateRequest struct {
	Status     string `json:"status,omitempty"`
	Priority   string `json:"priority,omitempty"`
	AssignedTo string `json:"assignedTo,omitempty"`
	Category   string `json:"category,omitempty"`
}

// Category representa una categoría de ticket
type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Color       string    `json:"color,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// FAQ representa una pregunta frecuente
type FAQ struct {
	ID          int       `json:"id"`
	Question    string    `json:"question"`
	Answer      string    `json:"answer"`
	Category    string    `json:"category"`
	IsPublished bool      `json:"isPublished"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// WidgetTicketRequest representa una solicitud de creación de ticket desde el widget
type WidgetTicketRequest struct {
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Message  string    `json:"message"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

// Metadata contiene información adicional para un ticket de widget
type Metadata struct {
	URL        string `json:"url,omitempty"`
	Referrer   string `json:"referrer,omitempty"`
	UserAgent  string `json:"userAgent,omitempty"`
	ScreenSize string `json:"screenSize,omitempty"`
	ExternalID string `json:"externalId,omitempty"`
}

// TicketResponse es la respuesta después de crear un ticket desde el widget
type TicketResponse struct {
	Success           bool   `json:"success"`
	TicketID          string `json:"ticketId"`
	LiveChatAvailable bool   `json:"liveChatAvailable"`
}

// WebSocketMessage representa un mensaje enviado a través de WebSocket
type WebSocketMessage struct {
	Type     string      `json:"type"`
	TicketID string      `json:"ticketId,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Messages []Message   `json:"messages,omitempty"`
}

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Error string `json:"error"`
}
