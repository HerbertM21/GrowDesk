package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/data"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/middleware"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/utils"
)

// TicketHandler contiene manejadores para operaciones de tickets
type TicketHandler struct {
	Store data.DataStore
}

// GetAllTickets maneja la obtención de todos los tickets
func (h *TicketHandler) GetAllTickets(w http.ResponseWriter, r *http.Request) {
	// Esta función solo maneja solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Establecer CORS
	utils.SetCORS(w)

	// Obtener tickets del almacén
	tickets, err := h.Store.GetTickets()
	if err != nil {
		http.Error(w, "Error al obtener tickets", http.StatusInternalServerError)
		return
	}

	// Devolver tickets como JSON
	utils.WriteJSON(w, http.StatusOK, tickets)
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

// CreateTicket maneja la creación de un nuevo ticket
func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	// Esta función solo maneja solicitudes POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Establecer CORS
	utils.SetCORS(w)

	// Obtener ID de usuario del contexto
	userID := r.Context().Value(middleware.UserIDKey).(string)
	if userID == "" {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	// Decodificar cuerpo de la solicitud
	var ticketReq models.TicketRequest
	if err := utils.DecodeJSON(r, &ticketReq); err != nil {
		http.Error(w, "Error al leer datos del ticket", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if ticketReq.Title == "" || ticketReq.Description == "" || ticketReq.CategoryID == "" {
		http.Error(w, "Título, descripción y categoría son requeridos", http.StatusBadRequest)
		return
	}

	// Crear mensaje inicial
	initialMessage := models.Message{
		ID:        uuid.New().String(),
		Content:   ticketReq.Description,
		UserID:    userID,
		UserName:  ticketReq.UserName,
		IsClient:  ticketReq.IsClient,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}

	// Crear nuevo ticket
	newTicket := models.Ticket{
		ID:          fmt.Sprintf("TICKET-%s", time.Now().Format("20060102-150405")),
		Title:       ticketReq.Title,
		Description: ticketReq.Description,
		CategoryID:  ticketReq.CategoryID,
		Status:      "open",
		Priority:    ticketReq.Priority,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Messages:    []models.Message{initialMessage},
		Metadata:    ticketReq.Metadata,
	}

	// Agregar ticket al almacén
	if err := h.Store.CreateTicket(newTicket); err != nil {
		http.Error(w, "Error al crear ticket", http.StatusInternalServerError)
		return
	}

	// Devolver ticket creado
	utils.WriteJSON(w, http.StatusCreated, newTicket)
}

// UpdateTicket maneja la actualización de un ticket existente
func (h *TicketHandler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Establecer CORS
	utils.SetCORS(w)

	// Obtener ID del ticket de la URL
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 4 {
		http.Error(w, "URL de ticket inválida", http.StatusBadRequest)
		return
	}

	ticketID := segments[3]

	// Obtener el ticket existente
	ticket, err := h.Store.GetTicket(ticketID)
	if err != nil {
		http.Error(w, "Ticket no encontrado", http.StatusNotFound)
		return
	}

	// Decodificar cuerpo de la solicitud
	var updates models.TicketUpdateRequest
	if err := utils.DecodeJSON(r, &updates); err != nil {
		http.Error(w, "Error al leer datos de actualización", http.StatusBadRequest)
		return
	}

	// Actualizar los campos del ticket
	if updates.Status != "" {
		ticket.Status = updates.Status
	}
	if updates.Priority != "" {
		ticket.Priority = updates.Priority
	}
	if updates.AssignedTo != "" {
		ticket.AssignedTo = updates.AssignedTo
	}
	if updates.Category != "" {
		ticket.Category = updates.Category
	}
	if updates.Department != "" {
		ticket.Department = updates.Department
	}
	if updates.Subject != "" {
		ticket.Subject = updates.Subject
	}

	// Actualizar timestamp
	ticket.UpdatedAt = time.Now()

	// Guardar en el almacén
	if err := h.Store.UpdateTicket(*ticket); err != nil {
		http.Error(w, "Error al actualizar ticket", http.StatusInternalServerError)
		return
	}

	// Devolver ticket actualizado
	utils.WriteJSON(w, http.StatusOK, ticket)
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
	if err := h.Store.AddTicketMessage(ticketID, message); err != nil {
		http.Error(w, "Failed to add message: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Broadcast a los clientes WebSocket
	h.Store.BroadcastMessage(ticketID, message)

	// Devolver respuesta de éxito
	response := struct {
		Success bool           `json:"success"`
		Message string         `json:"message"`
		Data    models.Message `json:"data"`
	}{
		Success: true,
		Message: "Message added successfully",
		Data:    message,
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
		h.Store.CreateTicket(ticket)

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
	h.Store.CreateTicket(ticket)

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
