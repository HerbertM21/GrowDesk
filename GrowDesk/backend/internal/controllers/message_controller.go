package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/models"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
	"github.com/hmdev/GrowDesk/backend/pkg/websocket"
)

// MessageRequest estructura para enviar un mensaje
type MessageRequest struct {
	Content    string `json:"content" binding:"required"`
	IsInternal bool   `json:"isInternal"` // Solo visible para agentes/admin
}

// GetMessages obtiene todos los mensajes de un ticket
func GetMessages(c *gin.Context) {
	ticketID := c.Param("ticketId")
	
	// Verificar permisos
	userRole, _ := c.Get("role")
	userID, _ := c.Get("userID")
	
	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, devolver mensajes de ejemplo
		mockMessages := []models.Message{
			{
				ID:        "MSG-20230001-1",
				TicketID:  ticketID,
				Content:   "No puedo iniciar sesión en el sistema desde ayer",
				IsClient:  true,
				UserName:  "Cliente Ejemplo",
				CreatedAt: time.Now().Add(-24 * time.Hour),
			},
			{
				ID:        "MSG-20230001-2",
				TicketID:  ticketID,
				Content:   "Estamos revisando su caso. ¿Podría indicarnos qué mensaje de error recibe?",
				IsClient:  false,
				UserName:  "Agente Soporte",
				CreatedAt: time.Now().Add(-20 * time.Hour),
			},
		}
		
		// Si es cliente, verificar que sea su ticket antes de devolver mensajes
		if userRole == "customer" {
			// Simular verificación
			if ticketID != "TICKET-123" && ticketID != "TICKET-20230001" {
				c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permiso para ver estos mensajes"})
				return
			}
		}
		
		c.JSON(http.StatusOK, mockMessages)
		return
	}

	// Verificar que el ticket existe y que el usuario tiene permisos
	var ticket models.Ticket
	if err := db.Where("id = ?", ticketID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
		return
	}
	
	// Si es cliente, verificar que sea su ticket
	if userRole == "customer" && ticket.CreatedBy != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permiso para ver estos mensajes"})
		return
	}
	
	// Obtener mensajes
	var messages []models.Message
	query := db.Where("ticket_id = ?", ticketID)
	
	// Si es cliente, no mostrar mensajes internos
	if userRole == "customer" {
		query = query.Where("is_internal = ?", false)
	}
	
	if err := query.Order("created_at asc").Find(&messages).Error; err != nil {
		log.Printf("GetMessages error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener mensajes"})
		return
	}
	
	c.JSON(http.StatusOK, messages)
}

// SendMessage envía un nuevo mensaje a un ticket
func SendMessage(c *gin.Context) {
	ticketID := c.Param("ticketId")
	var request MessageRequest
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}
	
	// Obtener información del usuario
	userID, _ := c.Get("userID")
	userEmail, _ := c.Get("email")
	userRole, _ := c.Get("role")
	isClient := userRole == "customer"
	
	// Si es mensaje interno, verificar que no sea cliente
	if request.IsInternal && isClient {
		c.JSON(http.StatusForbidden, gin.H{"error": "Los clientes no pueden enviar mensajes internos"})
		return
	}
	
	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular envío de mensaje
		newMessage := models.Message{
			ID:         "MSG-" + time.Now().Format("20060102150405"),
			TicketID:   ticketID,
			Content:    request.Content,
			IsClient:   isClient,
			IsInternal: request.IsInternal,
			UserID:     userID.(string),
			UserName:   isClient ? "Cliente" : "Agente",
			UserEmail:  userEmail.(string),
			CreatedAt:  time.Now(),
		}
		
		// Enviar notificación por WebSocket si hay un hub disponible
		if websocket.GetHub() != nil {
			websocket.GetHub().BroadcastToTicket(ticketID, gin.H{
				"type":     "new_message",
				"ticketId": ticketID,
				"data":     newMessage,
			})
		}
		
		c.JSON(http.StatusCreated, newMessage)
		return
	}

	// Verificar que el ticket existe
	var ticket models.Ticket
	if err := db.Where("id = ?", ticketID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
		return
	}
	
	// Si es cliente, verificar que sea su ticket
	if isClient && ticket.CreatedBy != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permiso para enviar mensajes a este ticket"})
		return
	}
	
	// Crear nuevo mensaje
	message := models.Message{
		TicketID:   ticketID,
		Content:    request.Content,
		IsClient:   isClient,
		IsInternal: request.IsInternal,
		UserID:     userID.(string),
		UserName:   getUserName(c),
		UserEmail:  userEmail.(string),
		CreatedAt:  time.Now(),
	}
	
	if err := db.Create(&message).Error; err != nil {
		log.Printf("SendMessage error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar mensaje"})
		return
	}
	
	// Actualizar estado del ticket si está cerrado o resuelto
	if ticket.Status == models.StatusClosed || ticket.Status == models.StatusResolved {
		updates := map[string]interface{}{
			"status":     models.StatusInProcess,
			"updated_at": time.Now(),
		}
		db.Model(&ticket).Updates(updates)
	} else {
		// Solo actualizar timestamp
		db.Model(&ticket).Update("updated_at", time.Now())
	}
	
	// Enviar notificación por WebSocket si hay un hub disponible
	if websocket.GetHub() != nil {
		websocket.GetHub().BroadcastToTicket(ticketID, gin.H{
			"type":     "new_message",
			"ticketId": ticketID,
			"data":     message,
		})
	}
	
	c.JSON(http.StatusCreated, message)
}

// DeleteMessage elimina un mensaje (solo para administradores)
func DeleteMessage(c *gin.Context) {
	messageID := c.Param("messageId")
	
	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular eliminación
		c.JSON(http.StatusOK, gin.H{
			"message": "Mensaje eliminado correctamente (simulado)",
			"id":      messageID,
		})
		return
	}

	// Eliminar el mensaje
	result := db.Delete(&models.Message{}, "id = ?", messageID)
	if result.Error != nil {
		log.Printf("DeleteMessage error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar mensaje"})
		return
	}
	
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mensaje no encontrado"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Mensaje eliminado correctamente"})
}

// Helper para obtener nombre completo del usuario
func getUserName(c *gin.Context) string {
	firstName, _ := c.Get("firstName")
	lastName, _ := c.Get("lastName")
	
	if firstName == nil || lastName == nil {
		return c.GetString("email")
	}
	
	return firstName.(string) + " " + lastName.(string)
}
