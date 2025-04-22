package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/models"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
)

type CreateMessageRequest struct {
	TicketID   string `json:"ticketId" binding:"required"`
	Content    string `json:"content" binding:"required"`
	IsInternal bool   `json:"isInternal"`
}

func GetMessages(c *gin.Context) {
	ticketID := c.Param("ticketId")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	db := database.GetDB()
	var ticket models.Ticket
	result := db.Where("id = ?", ticketID).First(&ticket)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")
	if userRole == "customer" && ticket.CreatedBy != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view messages for this ticket"})
		return
	}

	if userRole == "agent" && ticket.AssignedTo != nil && *ticket.AssignedTo != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "This ticket is assigned to another agent"})
		return
	}

	var messages []models.Message
	query := db.Where("ticket_id = ?", ticketID)

	if userRole == "customer" {
		query = query.Where("is_internal = ?", false)
	}

	result = query.Order("created_at ASC").Find(&messages)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	for i := range messages {
		db.Model(&messages[i]).Association("Attachments").Find(&messages[i].Attachments)
	}

	c.JSON(http.StatusOK, messages)
}

func CreateMessage(c *gin.Context) {
	var request CreateMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the user has access to this ticket
	db := database.GetDB()
	var ticket models.Ticket
	result := db.Where("id = ?", request.TicketID).First(&ticket)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")

	if userRole == "customer" && ticket.CreatedBy != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to add messages to this ticket"})
		return
	}

	if userRole == "agent" && ticket.AssignedTo != nil && *ticket.AssignedTo != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "This ticket is assigned to another agent"})
		return
	}

	if request.IsInternal && userRole == "customer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Customers cannot create internal notes"})
		return
	}

	message := models.Message{
		TicketID:   request.TicketID,
		UserID:     userID.(string),
		Content:    request.Content,
		IsInternal: request.IsInternal,
	}

	result = db.Create(&message)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	// Se actualiza el estado del ticket si el usuario es cliente y el ticket está resuelto
	if userRole == "customer" && ticket.Status == models.StatusResolved {
		db.Model(&ticket).Update("status", models.StatusOpen)
	}

	// If agent adds message to ticket, set to in_progress if not already resolved
	if (userRole == "agent" || userRole == "admin") && ticket.Status != models.StatusResolved && ticket.Status != models.StatusClosed {
		db.Model(&ticket).Update("status", models.StatusInProgress)
	}

	c.JSON(http.StatusCreated, message)
}

// UploadAttachment  sube archivos adjuntos a un mensaje
func UploadAttachment(c *gin.Context) {
	ticketID := c.PostForm("ticketId")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	// Verificar si el usuario tiene acceso a este ticket
	db := database.GetDB()
	var ticket models.Ticket
	result := db.Where("id = ?", ticketID).First(&ticket)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Verificar permisos
	userID, _ := c.Get("userID")
	userRole, _ := c.Get("role")

	if userRole == "customer" && ticket.CreatedBy != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to add attachments to this ticket"})
		return
	}

	if userRole == "agent" && ticket.AssignedTo != nil && *ticket.AssignedTo != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "This ticket is assigned to another agent"})
		return
	}

	// Obtener el archivo del formulario
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Generar un nombre de archivo único
	filename := generateFilename(file.Filename)

	// Guardar el archivo en el disco (en una aplicación real, podrías usar almacenamiento en la nube)
	if err := c.SaveUploadedFile(file, "./uploads/"+filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Crear mensaje si hay contenido
	messageID := ""
	content := c.PostForm("content")
	if content != "" {
		message := models.Message{
			TicketID:   ticketID,
			UserID:     userID.(string),
			Content:    content,
			IsInternal: false,
		}

		db.Create(&message)
		messageID = message.ID
	}

	// Crear registro de archivo adjunto
	attachment := models.Attachment{
		MessageID: messageID,
		FileName:  file.Filename,
		FileType:  file.Header.Get("Content-Type"),
		FileSize:  file.Size,
		FileURL:   "/uploads/" + filename,
	}

	result = db.Create(&attachment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attachment record"})
		return
	}

	c.JSON(http.StatusCreated, attachment)
}

// Función auxiliar para generar un nombre de archivo único
func generateFilename(originalName string) string {
	return models.GenerateUUID() + "_" + originalName
}
