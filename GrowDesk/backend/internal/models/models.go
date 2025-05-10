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
	Subject     string    `json:"subject,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Priority    string    `json:"priority,omitempty"`
	Category    string    `json:"category,omitempty"`
	CategoryID  string    `json:"categoryId,omitempty"`
	AssignedTo  string    `json:"assignedTo,omitempty"`
	CreatedBy   string    `json:"createdBy,omitempty"`
	UserID      string    `json:"userId,omitempty"`
	Description string    `json:"description,omitempty"`
	Customer    Customer  `json:"customer"`
	Messages    []Message `json:"messages,omitempty"`
	Source      string    `json:"source,omitempty"`
	WidgetID    string    `json:"widgetId,omitempty"`
	Department  string    `json:"department,omitempty"`
	Metadata    *Metadata `json:"metadata,omitempty"`
}

// Customer representa a un cliente de un ticket
type Customer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Message representa un mensaje en un ticket de soporte
type Message struct {
	ID         string    `json:"id"`
	Content    string    `json:"content"`
	IsClient   bool      `json:"isClient"`
	IsInternal bool      `json:"isInternal,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	UserID     string    `json:"userId,omitempty"`
	UserName   string    `json:"userName,omitempty"`
	UserEmail  string    `json:"userEmail,omitempty"`
}

// NewMessageRequest representa una solicitud para agregar un nuevo mensaje
type NewMessageRequest struct {
	Content    string `json:"content"`
	IsClient   bool   `json:"isClient"`
	IsInternal bool   `json:"isInternal,omitempty"`
	UserID     string `json:"userId,omitempty"`
	UserName   string `json:"userName,omitempty"`
	UserEmail  string `json:"userEmail,omitempty"`
}

// TicketUpdateRequest representa una solicitud para actualizar un ticket
type TicketUpdateRequest struct {
	Status     string `json:"status,omitempty"`
	Priority   string `json:"priority,omitempty"`
	AssignedTo string `json:"assignedTo,omitempty"`
	Category   string `json:"category,omitempty"`
	Department string `json:"department,omitempty"`
	Subject    string `json:"subject,omitempty"`
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

// Activity representa una actividad o evento en el sistema
type Activity struct {
	ID          string         `json:"id"`
	UserID      string         `json:"userId"`
	Type        string         `json:"type"`
	TargetID    string         `json:"targetId,omitempty"`
	Description string         `json:"description"`
	Timestamp   time.Time      `json:"timestamp"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// Notification representa una notificación para un usuario
type Notification struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Message     string    `json:"message"`
	Type        string    `json:"type"`
	Read        bool      `json:"read"`
	CreatedAt   time.Time `json:"createdAt"`
	RelatedID   string    `json:"relatedId,omitempty"`
	RelatedType string    `json:"relatedType,omitempty"`
}

// Attachment representa un archivo adjunto a un mensaje
type Attachment struct {
	ID        string    `json:"id"`
	MessageID string    `json:"messageId"`
	FileName  string    `json:"fileName"`
	FileType  string    `json:"fileType"`
	FileSize  int       `json:"fileSize"`
	FileURL   string    `json:"fileUrl"`
	CreatedAt time.Time `json:"createdAt"`
}

// WidgetSetting representa la configuración de un widget
type WidgetSetting struct {
	ID             string    `json:"id"`
	WidgetID       string    `json:"widgetId"`
	WidgetToken    string    `json:"widgetToken"`
	BrandName      string    `json:"brandName"`
	WelcomeMessage string    `json:"welcomeMessage"`
	PrimaryColor   string    `json:"primaryColor"`
	Position       string    `json:"position"`
	LogoURL        string    `json:"logoUrl,omitempty"`
	AllowedDomains []string  `json:"allowedDomains,omitempty"`
	IsActive       bool      `json:"isActive"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// WidgetTicket representa un ticket creado desde el widget
type WidgetTicket struct {
	ID          string    `json:"id"`
	TicketID    string    `json:"ticketId"`
	Title       string    `json:"title"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ClientName  string    `json:"clientName"`
	ClientEmail string    `json:"clientEmail"`
	WidgetID    string    `json:"widgetId"`
	Department  string    `json:"department"`
	Synced      bool      `json:"synced"`
}

// WidgetMessage representa un mensaje en un ticket de widget
type WidgetMessage struct {
	ID             string    `json:"id"`
	WidgetTicketID string    `json:"widgetTicketId"`
	Content        string    `json:"content"`
	IsClient       bool      `json:"isClient"`
	CreatedAt      time.Time `json:"createdAt"`
	UserName       string    `json:"userName"`
	UserEmail      string    `json:"userEmail"`
	Synced         bool      `json:"synced"`
}

// WidgetSession representa una sesión de usuario en el widget
type WidgetSession struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	TicketID   string    `json:"ticketId"`
	WidgetID   string    `json:"widgetId"`
	CreatedAt  time.Time `json:"createdAt"`
	ExpiresAt  time.Time `json:"expiresAt"`
	LastActive time.Time `json:"lastActive"`
}

// WidgetMetadata representa información adicional sobre una sesión de widget
type WidgetMetadata struct {
	ID         string    `json:"id"`
	SessionID  string    `json:"sessionId"`
	URL        string    `json:"url"`
	Referrer   string    `json:"referrer"`
	UserAgent  string    `json:"userAgent"`
	ScreenSize string    `json:"screenSize"`
	Browser    string    `json:"browser"`
	OS         string    `json:"os"`
	IPAddress  string    `json:"ipAddress"`
	CreatedAt  time.Time `json:"createdAt"`
}

// TicketRequest representa una solicitud para crear un nuevo ticket
type TicketRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CategoryID  string    `json:"categoryId"`
	Priority    string    `json:"priority"`
	UserName    string    `json:"userName,omitempty"`
	UserEmail   string    `json:"userEmail,omitempty"`
	IsClient    bool      `json:"isClient"`
	Metadata    *Metadata `json:"metadata,omitempty"`
}
