package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/data"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/utils"
)

// TicketHandler contiene manejadores para operaciones de tickets
type TicketHandler struct {
	Store *data.Store
}

// GetAllTickets devuelve todos los tickets
func (h *TicketHandler) GetAllTickets(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Devolver todos los tickets
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.Store.Tickets)
}

// GetTicket devuelve un ticket específico por ID
func (h *TicketHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer el ID del ticket desde la URL
	// Formato de URL: /api/tickets/:id
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "ID de ticket inválido", http.StatusBadRequest)
		return
	}

	ticketID := parts[len(parts)-1]

	// Obtener el ticket
	ticket, err := h.Store.GetTicket(ticketID)
	if err != nil {
		http.Error(w, "Ticket no encontrado", http.StatusNotFound)
		return
	}

	// Devolver el ticket
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
}

// CreateTicket crea un nuevo ticket
func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear el cuerpo de la solicitud
	var ticketReq struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Priority    string `json:"priority"`
		Customer    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"customer"`
	}

	if err := json.NewDecoder(r.Body).Decode(&ticketReq); err != nil {
		http.Error(w, "El cuerpo de la solicitud es inválido", http.StatusBadRequest)
		return
	}

	// Crear nuevo ticket
	ticket := models.Ticket{
		ID:          utils.GenerateTicketID(),
		Title:       ticketReq.Title,
		Description: ticketReq.Description,
		Status:      "open",
		Priority:    ticketReq.Priority,
		Category:    ticketReq.Category,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Customer: models.Customer{
			Name:  ticketReq.Customer.Name,
			Email: ticketReq.Customer.Email,
		},
		Messages: []models.Message{},
	}

	// Agregar mensaje inicial si se proporciona una descripción
	if ticketReq.Description != "" {
		ticket.Messages = append(ticket.Messages, models.Message{
			ID:        utils.GenerateMessageID(),
			Content:   ticketReq.Description,
			IsClient:  true,
			Timestamp: time.Now(),
			UserName:  ticketReq.Customer.Name,
		})
	}

	// Agregar ticket al almacén
	h.Store.AddTicket(ticket)

	// Devolver el ticket creado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ticket)
}

// UpdateTicket updates an existing ticket
func (h *TicketHandler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer el ID del ticket desde la URL
	// Formato de URL: /api/tickets/:id
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "ID de ticket inválido", http.StatusBadRequest)
		return
	}

	ticketID := parts[len(parts)-1]

	// Parsear el cuerpo de la solicitud
	var updateReq models.TicketUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "El cuerpo de la solicitud es inválido", http.StatusBadRequest)
		return
	}

	// Actualizar el ticket
	ticket, err := h.Store.UpdateTicket(ticketID, updateReq)
	if err != nil {
		http.Error(w, "Ticket no encontrado", http.StatusNotFound)
		return
	}

	// Devolver el ticket actualizado
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
}

// GetTicketMessages devuelve mensajes para un ticket específico
func (h *TicketHandler) GetTicketMessages(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer el ID del ticket desde la URL
	// Formato de URL: /api/tickets/:id/messages
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "ID de ticket inválido", http.StatusBadRequest)
		return
	}

	// Obtener el ID desde la URL (asumiendo formato /tickets/ID/messages)
	ticketID := parts[len(parts)-2]

	// Obtener el ticket
	ticket, err := h.Store.GetTicket(ticketID)
	if err != nil {
		http.Error(w, "Ticket no encontrado", http.StatusNotFound)
		return
	}

	// Devolver los mensajes
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket.Messages)
}

// AddTicketMessage agrega un mensaje a un ticket
func (h *TicketHandler) AddTicketMessage(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer el ID del ticket desde la URL
	// Formato de URL: /api/tickets/:id/messages
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "ID de ticket inválido", http.StatusBadRequest)
		return
	}

	// Obtener el ID desde la URL (asumiendo formato /tickets/ID/messages)
	ticketID := parts[len(parts)-2]

	// Parsear el cuerpo de la solicitud
	var messageReq models.NewMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&messageReq); err != nil {
		http.Error(w, "El cuerpo de la solicitud es inválido", http.StatusBadRequest)
		return
	}

	// Validar contenido
	if messageReq.Content == "" {
		http.Error(w, "El contenido del mensaje es requerido", http.StatusBadRequest)
		return
	}

	// Crear nuevo mensaje
	message := models.Message{
		ID:        utils.GenerateMessageID(),
		Content:   messageReq.Content,
		IsClient:  messageReq.IsClient,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
		UserName:  messageReq.UserName,
		UserEmail: messageReq.UserEmail,
	}

	// Agregar mensaje al ticket
	addedMessage, err := h.Store.AddMessageToTicket(ticketID, message)
	if err != nil {
		http.Error(w, "Failed to add message: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Broadcast a los clientes WebSocket
	h.Store.BroadcastMessage(ticketID, *addedMessage)

	// Devolver respuesta de éxito
	response := struct {
		Success bool           `json:"success"`
		Message string         `json:"message"`
		Data    models.Message `json:"data"`
	}{
		Success: true,
		Message: "Message added successfully",
		Data:    *addedMessage,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// CreateWidgetTicket crea un nuevo ticket desde el widget
func (h *TicketHandler) CreateWidgetTicket(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Log para depuración
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	fmt.Printf("Recibido del widget: %s\n", string(body))

	// Estructura para recibir el formato específico del widget
	var widgetRequest struct {
		Subject    string                 `json:"subject"`    // El widget envía "subject" en lugar de "title"
		Message    string                 `json:"message"`    // Mensaje inicial
		Name       string                 `json:"name"`       // Nombre del cliente
		Email      string                 `json:"email"`      // Email del cliente
		Priority   string                 `json:"priority"`   // Opcional
		Department string                 `json:"department"` // Opcional
		WidgetId   string                 `json:"widgetId"`   // ID del widget
		Metadata   map[string]interface{} `json:"metadata"`   // Metadatos
	}

	// Intentar decodificar primero con el formato del widget
	if err := json.Unmarshal(body, &widgetRequest); err != nil {
		// Si falla, intentar con el formato original
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		var reqBody models.WidgetTicketRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Continuar con el formato original
		name := reqBody.Name
		if name == "" {
			name = "Anónimo"
		}

		email := reqBody.Email
		if email == "" {
			email = "anonymous@example.com"
		}

		// Crear nuevo ticket
		ticketID := utils.GenerateTicketID()
		now := time.Now()

		ticket := models.Ticket{
			ID:          ticketID,
			Title:       fmt.Sprintf("Soporte Web - %s", name),
			Status:      "open",
			CreatedAt:   now,
			UpdatedAt:   now,
			Description: reqBody.Message,
			Priority:    "medium",
			Category:    "soporte",
			CreatedBy:   email,
			Customer: models.Customer{
				Name:  name,
				Email: email,
			},
			Messages: []models.Message{
				{
					ID:        utils.GenerateMessageID(),
					Content:   reqBody.Message,
					IsClient:  true,
					Timestamp: now,
					UserName:  name,
					UserEmail: email,
				},
			},
		}

		// Add ticket to store
		h.Store.AddTicket(ticket)

		// Devolver respuesta de éxito
		response := models.TicketResponse{
			Success:           true,
			TicketID:          ticketID,
			LiveChatAvailable: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Proceso con el formato del widget
	name := widgetRequest.Name
	if name == "" {
		name = "Anónimo"
	}

	email := widgetRequest.Email
	if email == "" {
		email = "anonymous@example.com"
	}

	// Crear nuevo ticket
	ticketID := utils.GenerateTicketID()
	now := time.Now()

	priority := widgetRequest.Priority
	if priority == "" {
		priority = "medium"
	}

	department := widgetRequest.Department
	if department == "" {
		department = "soporte"
	}

	ticket := models.Ticket{
		ID:          ticketID,
		Title:       widgetRequest.Subject,
		Status:      "open",
		CreatedAt:   now,
		UpdatedAt:   now,
		Description: widgetRequest.Message,
		Priority:    priority,
		Category:    department,
		CreatedBy:   email,
		Customer: models.Customer{
			Name:  name,
			Email: email,
		},
		Messages: []models.Message{
			{
				ID:        utils.GenerateMessageID(),
				Content:   widgetRequest.Message,
				IsClient:  true,
				Timestamp: now,
				UserName:  name,
				UserEmail: email,
			},
		},
	}

	// Agregar ticket al almacén
	h.Store.AddTicket(ticket)

	// Devolver respuesta de éxito con el formato que el widget espera
	response := map[string]interface{}{
		"success":           true,
		"ticketId":          ticketID,
		"liveChatAvailable": true,
		"id":                ticketID, // Agregar id porque el widget lo espera
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
