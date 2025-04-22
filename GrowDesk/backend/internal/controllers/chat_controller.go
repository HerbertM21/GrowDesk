package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hmdev/GrowDesk/backend/models"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
	"github.com/hmdev/GrowDesk/backend/pkg/websocket"
)

type MessageRequest struct {
	TicketID   string `json:"ticketId" binding:"required"`
	Content    string `json:"content" binding:"required"`
	IsInternal bool   `json:"isInternal"`
}

// GetRealTimeMessages obtiene todos los mensajes para un ticket (versión websocket)
func GetRealTimeMessages(c *gin.Context) {
	ticketID := c.Param("ticketId")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ticket no proporcionado"})
		return
	}

	var messages []models.Message
	result := database.DB.Where("ticket_id = ?", ticketID).Order("created_at").Find(&messages)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// CreateRealTimeMessage crea un nuevo mensaje para un ticket (versión websocket)
func CreateRealTimeMessage(c *gin.Context) {
	var req MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar que el ticket existe
	var ticket models.Ticket
	if err := database.DB.First(&ticket, "id = ?", req.TicketID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
		return
	}

	// Obtener el ID del usuario actual del token (en producción)
	userID := c.GetString("userId")
	if userID == "" {
		userID = "testuser123" // ID de ejemplo para desarrollo
	}

	// Crear nuevo mensaje
	message := models.Message{
		ID:         uuid.New().String(),
		TicketID:   req.TicketID,
		UserID:     userID,
		Content:    req.Content,
		IsInternal: req.IsInternal,
	}

	// Guardar en la base de datos
	if err := database.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	websocket.GlobalHub.BroadcastToTicket(req.TicketID, "new_message", message)

	c.JSON(http.StatusCreated, message)
}

func WebSocketHandler(c *gin.Context) {
	websocket.ServeWs(c)
}
