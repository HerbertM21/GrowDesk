package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/models"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
)

// CreateTicketRequest contiene datos para crear un ticket
type CreateTicketRequest struct {
	Title       string                `json:"title" binding:"required"`
	Description string                `json:"description" binding:"required"`
	Priority    models.TicketPriority `json:"priority" binding:"required"`
	Category    string                `json:"category" binding:"required"`
}

// UpdateTicketRequest contiene campos que se pueden actualizar
type UpdateTicketRequest struct {
	Title       *string                `json:"title"`
	Description *string                `json:"description"`
	Status      *models.TicketStatus   `json:"status"`
	Priority    *models.TicketPriority `json:"priority"`
	Category    *string                `json:"category"`
	AssignedTo  *string                `json:"assignedTo"`
}

// GetAllTickets devuelve todos los tickets (filtrados por rol de usuario)
func GetAllTickets(c *gin.Context) {
	db := database.GetDB()
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")

	var tickets []models.Ticket
	query := db

	// Filtrar tickets basado en el rol de usuario
	if userRole == "customer" {
		query = query.Where("created_by = ?", userID)
	} else if userRole == "agent" {
		// Los agentes pueden ver tickets asignados a ellos o tickets sin asignar
		query = query.Where("assigned_to = ? OR assigned_to IS NULL", userID)
	}
	// Los administradores pueden ver todos los tickets (no filtrar)

	result := query.Order("created_at DESC").Find(&tickets)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// GetTicket devuelve un ticket único por ID
func GetTicket(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	var ticket models.Ticket
	db := database.GetDB()
	result := db.Where("id = ?", id).First(&ticket)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Verificar permisos
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")
	if userRole == "customer" && ticket.CreatedBy != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this ticket"})
		return
	}

	if userRole == "agent" && ticket.AssignedTo != nil && *ticket.AssignedTo != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "This ticket is assigned to another agent"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// CreateTicket crea un nuevo ticket
func CreateTicket(c *gin.Context) {
	var request CreateTicketRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	ticket := models.Ticket{
		Title:       request.Title,
		Description: request.Description,
		Priority:    request.Priority,
		Category:    request.Category,
		Status:      models.StatusOpen,
		CreatedBy:   userID.(string),
	}

	db := database.GetDB()
	result := db.Create(&ticket)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// UpdateTicket actualiza un ticket existente
func UpdateTicket(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	var request UpdateTicketRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	var ticket models.Ticket
	result := db.Where("id = ?", id).First(&ticket)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Verificar permisos
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")

	// Los clientes solo pueden actualizar sus propios tickets
	if userRole == "customer" && ticket.CreatedBy != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this ticket"})
		return
	}

	// Los agentes solo pueden actualizar tickets asignados a ellos
	if userRole == "agent" && ticket.AssignedTo != nil && *ticket.AssignedTo != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "This ticket is assigned to another agent"})
		return
	}

	// Actualizar los campos que se proporcionaron
	updates := make(map[string]interface{})

	if request.Title != nil {
		updates["title"] = *request.Title
	}

	if request.Description != nil {
		updates["description"] = *request.Description
	}

	if request.Status != nil {
		updates["status"] = *request.Status
	}

	if request.Priority != nil {
		updates["priority"] = *request.Priority
	}

	if request.Category != nil {
		updates["category"] = *request.Category
	}

	// Solo los agentes y los administradores pueden asignar tickets
	if request.AssignedTo != nil && (userRole == "agent" || userRole == "admin") {
		updates["assigned_to"] = *request.AssignedTo

		// Si un ticket está siendo asignado, actualizar el estado
		if *request.AssignedTo != "" && ticket.Status == models.StatusOpen {
			updates["status"] = models.StatusAssigned
		}
	}

	// Aplicar actualizaciones
	result = db.Model(&ticket).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket"})
		return
	}

	// Obtener el ticket actualizado
	db.Where("id = ?", id).First(&ticket)
	c.JSON(http.StatusOK, ticket)
}

// DeleteTicket elimina un ticket
func DeleteTicket(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	db := database.GetDB()
	var ticket models.Ticket
	result := db.Where("id = ?", id).First(&ticket)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Solo los administradores y el creador del ticket pueden eliminar tickets
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")
	if userRole != "admin" && ticket.CreatedBy != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this ticket"})
		return
	}

	result = db.Delete(&ticket)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}
