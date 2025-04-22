package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/models"
)

// Estructura para recibir datos desde el widget
type WidgetTicketData struct {
	Name     string                 `json:"name" binding:"required"`
	Email    string                 `json:"email" binding:"required,email"`
	Message  string                 `json:"message" binding:"required"`
	Metadata map[string]interface{} `json:"metadata"`
}

// Estructura para recibir mensajes desde el widget
type WidgetMessageData struct {
	TicketID string `json:"ticketId" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

// CreateTicketFromWidget crea un ticket a partir de datos enviados por el widget
func CreateTicketFromWidget(c *gin.Context) {
	var widgetData WidgetTicketData
	if err := c.ShouldBindJSON(&widgetData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de ticket inválidos: " + err.Error(),
		})
		return
	}

	// Obtener ID del widget si está disponible
	widgetID, exists := c.Get("widgetID")
	widgetSource := "widget"
	if exists {
		widgetSource = "widget-" + widgetID.(string)
	}

	// Crear el ticket en el sistema
	ticket := models.Ticket{
		Title:       "Soporte Web - " + widgetData.Name,
		Description: widgetData.Message,
		Status:      models.StatusOpen,
		Priority:    models.PriorityMedium,
		Category:    widgetSource,     // Usando Category en lugar de Source
		CreatedBy:   widgetData.Email, // Usando email como creador
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// implementar dsp
	// db.Create(&ticket)

	// Para demo, asignar un ID
	ticket.ID = strconv.FormatInt(time.Now().Unix(), 10)

	// implementar dsp
	// Crear también un usuario asociado o usar uno existente
	//  buscar o crear usuario

	// Responder con éxito y el ID del ticket
	c.JSON(http.StatusCreated, gin.H{
		"success":  true,
		"ticketId": ticket.ID,
		"message":  "Ticket creado correctamente",
	})
}

// CreateMessageFromWidget guarda un mensaje enviado desde el widget
func CreateMessageFromWidget(c *gin.Context) {
	var messageData WidgetMessageData
	if err := c.ShouldBindJSON(&messageData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de mensaje inválidos: " + err.Error(),
		})
		return
	}

	// implementar dsp
	// y guardar el mensaje en la base de datos

	// Para demo, crear un ID de mensaje
	messageID := time.Now().UnixNano()

	// Responder con éxito
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"messageId": messageID,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// GetMessagesForWidget obtiene los mensajes de un ticket para mostrarlos en el widget
func GetMessagesForWidget(c *gin.Context) {
	ticketID := c.Param("id")

	// implementar dsp
	// messages := db.Where("ticket_id = ?", ticketID).Find(&models.Message{})

	// Para demo, devolver mensajes de ejemplo
	messages := []gin.H{
		{
			"id":        "msg-1",
			"ticketId":  ticketID,
			"content":   "¡Hola! ¿En qué podemos ayudarte?",
			"isClient":  false,
			"createdAt": time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
		},
		{
			"id":        "msg-2",
			"ticketId":  ticketID,
			"content":   "Tengo problemas para acceder a mi cuenta",
			"isClient":  true,
			"createdAt": time.Now().Add(-25 * time.Minute).Format(time.RFC3339),
		},
		{
			"id":        "msg-3",
			"ticketId":  ticketID,
			"content":   "¿Podrías proporcionar más detalles sobre el problema?",
			"isClient":  false,
			"createdAt": time.Now().Add(-20 * time.Minute).Format(time.RFC3339),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}
