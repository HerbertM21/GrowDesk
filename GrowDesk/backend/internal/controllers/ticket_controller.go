package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/models"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
)

// GetAllTickets obtiene todos los tickets con filtro opcional
func GetAllTickets(c *gin.Context) {
	// Parámetros opcionales de filtro
	status := c.Query("status")
	assignedTo := c.Query("assignedTo")
	priority := c.Query("priority")
	category := c.Query("category")

	// Obtener rol y ID del usuario actual
	userRole, _ := c.Get("role")
	userID, _ := c.Get("userID")

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, devolver tickets de ejemplo
		mockTickets := []models.Ticket{
			{
				ID:          "TICKET-20230001",
				Title:       "Problema con inicio de sesión",
				Description: "No puedo iniciar sesión en el sistema",
				Status:      models.StatusOpen,
				Priority:    models.PriorityMedium,
				Category:    "Soporte Técnico",
				Customer: models.Customer{
					Name:  "Cliente Ejemplo",
					Email: "cliente@ejemplo.com",
				},
				CreatedBy:  "cliente@ejemplo.com",
				AssignedTo: nil,
				CreatedAt:  time.Now().Add(-24 * time.Hour),
				UpdatedAt:  time.Now().Add(-24 * time.Hour),
			},
			{
				ID:          "TICKET-20230002",
				Title:       "Solicitud de nueva característica",
				Description: "Me gustaría sugerir una nueva funcionalidad",
				Status:      models.StatusAssigned,
				Priority:    models.PriorityLow,
				Category:    "Características",
				Customer: models.Customer{
					Name:  "Otro Cliente",
					Email: "otro@ejemplo.com",
				},
				CreatedBy:  "otro@ejemplo.com",
				AssignedTo: strPtr("agent-456"),
				CreatedAt:  time.Now().Add(-48 * time.Hour),
				UpdatedAt:  time.Now().Add(-12 * time.Hour),
			},
		}

		// Aplicar filtros básicos en modo sin DB
		var filtered []models.Ticket
		for _, ticket := range mockTickets {
			include := true

			if status != "" && string(ticket.Status) != status {
				include = false
			}

			if assignedTo != "" {
				if ticket.AssignedTo == nil || *ticket.AssignedTo != assignedTo {
					include = false
				}
			}

			if priority != "" && string(ticket.Priority) != priority {
				include = false
			}

			if category != "" && ticket.Category != category {
				include = false
			}

			// Si es cliente, solo mostrar sus propios tickets
			if userRole == "customer" && ticket.CreatedBy != userID.(string) {
				include = false
			}

			if include {
				filtered = append(filtered, ticket)
			}
		}

		c.JSON(http.StatusOK, filtered)
		return
	}

	// Construir consulta con filtros
	query := db.Model(&models.Ticket{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if assignedTo != "" {
		query = query.Where("assigned_to = ?", assignedTo)
	}

	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Si es cliente, solo mostrar sus propios tickets
	if userRole == "customer" {
		query = query.Where("created_by = ?", userID)
	}

	var tickets []models.Ticket
	result := query.Order("updated_at desc").Find(&tickets)
	if result.Error != nil {
		log.Printf("GetAllTickets error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// GetTicket obtiene un ticket específico por su ID
func GetTicket(c *gin.Context) {
	id := c.Param("id")

	// Obtener rol y ID del usuario actual
	userRole, _ := c.Get("role")
	userID, _ := c.Get("userID")

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, devolver ticket de ejemplo
		if id == "TICKET-20230001" {
			ticket := models.Ticket{
				ID:          "TICKET-20230001",
				Title:       "Problema con inicio de sesión",
				Description: "No puedo iniciar sesión en el sistema",
				Status:      models.StatusOpen,
				Priority:    models.PriorityMedium,
				Category:    "Soporte Técnico",
				Customer: models.Customer{
					Name:  "Cliente Ejemplo",
					Email: "cliente@ejemplo.com",
				},
				CreatedBy:  "cliente@ejemplo.com",
				AssignedTo: nil,
				CreatedAt:  time.Now().Add(-24 * time.Hour),
				UpdatedAt:  time.Now().Add(-24 * time.Hour),
				Messages: []models.Message{
					{
						ID:        "MSG-20230001-1",
						TicketID:  "TICKET-20230001",
						Content:   "No puedo iniciar sesión en el sistema desde ayer",
						IsClient:  true,
						UserName:  "Cliente Ejemplo",
						CreatedAt: time.Now().Add(-24 * time.Hour),
					},
				},
			}

			// Si es cliente, verificar que sea su ticket
			if userRole == "customer" && ticket.CreatedBy != userID.(string) {
				c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permiso para ver este ticket"})
				return
			}

			c.JSON(http.StatusOK, ticket)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
		return
	}

	var ticket models.Ticket
	result := db.Where("id = ?", id).Preload("Messages").First(&ticket)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
		return
	}

	// Si es cliente, verificar que sea su ticket
	if userRole == "customer" && ticket.CreatedBy != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permiso para ver este ticket"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// CreateTicket crea un nuevo ticket
func CreateTicket(c *gin.Context) {
	var request struct {
		Title         string `json:"title" binding:"required"`
		Description   string `json:"description" binding:"required"`
		Priority      string `json:"priority"`
		Category      string `json:"category"`
		CustomerName  string `json:"customerName"`
		CustomerEmail string `json:"customerEmail"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	// Obtener información del usuario creador
	userID, _ := c.Get("userID")
	userEmail, _ := c.Get("email")
	userRole, _ := c.Get("role")

	// Establecer prioridad predeterminada si no se proporciona
	priority := models.TicketPriority(request.Priority)
	if priority == "" {
		priority = models.PriorityMedium
	}

	// Crear estructura del ticket
	ticket := models.Ticket{
		Title:       request.Title,
		Description: request.Description,
		Status:      models.StatusOpen,
		Priority:    priority,
		Category:    request.Category,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Configurar información del cliente según el rol del creador
	if userRole == "customer" {
		ticket.CreatedBy = userID.(string)
		ticket.Customer = models.Customer{
			Name:  c.GetString("firstName") + " " + c.GetString("lastName"),
			Email: userEmail.(string),
		}
	} else { // agente o admin creando ticket para un cliente
		ticket.CreatedBy = userID.(string)

		// Si se proporcionan datos del cliente, usarlos
		if request.CustomerEmail != "" {
			ticket.Customer = models.Customer{
				Name:  request.CustomerName,
				Email: request.CustomerEmail,
			}
		} else {
			ticket.Customer = models.Customer{
				Name:  "Cliente sin especificar",
				Email: "sin@email.com",
			}
		}
	}

	// Mensaje inicial basado en la descripción
	initialMessage := models.Message{
		Content:   request.Description,
		IsClient:  userRole == "customer",
		UserName:  ticket.Customer.Name,
		UserEmail: ticket.Customer.Email,
		CreatedAt: time.Now(),
	}

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular creación
		ticket.ID = "TICKET-" + time.Now().Format("20060102150405")
		initialMessage.ID = "MSG-" + time.Now().Format("20060102150405")
		initialMessage.TicketID = ticket.ID

		ticket.Messages = append(ticket.Messages, initialMessage)

		c.JSON(http.StatusCreated, ticket)
		return
	}

	// Crear ticket y mensaje inicial en transacción
	tx := db.Begin()

	if err := tx.Create(&ticket).Error; err != nil {
		tx.Rollback()
		log.Printf("CreateTicket error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear ticket"})
		return
	}

	initialMessage.TicketID = ticket.ID
	if err := tx.Create(&initialMessage).Error; err != nil {
		tx.Rollback()
		log.Printf("CreateTicket error (message): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear mensaje inicial"})
		return
	}

	tx.Commit()

	// Cargar el ticket completo con mensaje
	db.Where("id = ?", ticket.ID).Preload("Messages").First(&ticket)

	c.JSON(http.StatusCreated, ticket)
}

// UpdateTicket actualiza un ticket existente
func UpdateTicket(c *gin.Context) {
	id := c.Param("id")

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// No permitir modificar campos sensibles
	delete(updateData, "id")
	delete(updateData, "created_at")
	delete(updateData, "created_by")

	// Obtener rol del usuario
	userRole, _ := c.Get("role")
	userID, _ := c.Get("userID")

	// Actualizar timestamp
	updateData["updated_at"] = time.Now()

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular actualización
		c.JSON(http.StatusOK, gin.H{
			"id":      id,
			"message": "Ticket actualizado correctamente (simulado)",
			"updated": updateData,
		})
		return
	}

	var ticket models.Ticket
	if err := db.Where("id = ?", id).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
		return
	}

	// Si es cliente, verificar que sea su ticket
	if userRole == "customer" && ticket.CreatedBy != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permiso para modificar este ticket"})
		return
	}

	// Actualizar ticket
	if err := db.Model(&ticket).Updates(updateData).Error; err != nil {
		log.Printf("UpdateTicket error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar ticket"})
		return
	}

	// Obtener ticket actualizado con mensajes
	if err := db.Where("id = ?", id).Preload("Messages").First(&ticket).Error; err != nil {
		log.Printf("UpdateTicket error (reload): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar ticket actualizado"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// DeleteTicket elimina un ticket (solo para administradores)
func DeleteTicket(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular eliminación
		c.JSON(http.StatusOK, gin.H{
			"message": "Ticket eliminado correctamente (simulado)",
			"id":      id,
		})
		return
	}

	// Eliminar mensajes asociados primero
	if err := db.Where("ticket_id = ?", id).Delete(&models.Message{}).Error; err != nil {
		log.Printf("DeleteTicket error (messages): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar mensajes del ticket"})
		return
	}

	// Eliminar ticket
	result := db.Delete(&models.Ticket{}, "id = ?", id)
	if result.Error != nil {
		log.Printf("DeleteTicket error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar ticket"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket eliminado correctamente"})
}

// AssignTicket asigna un ticket a un agente
func AssignTicket(c *gin.Context) {
	id := c.Param("id")

	var request struct {
		AgentID string `json:"agentId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de agente no especificado"})
		return
	}

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular asignación
		c.JSON(http.StatusOK, gin.H{
			"id":      id,
			"message": "Ticket asignado correctamente (simulado)",
			"agentId": request.AgentID,
			"status":  "assigned",
		})
		return
	}

	var ticket models.Ticket
	if err := db.Where("id = ?", id).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
		return
	}

	// Verificar que el agente existe
	var agent models.User
	if err := db.Where("id = ? AND (role = 'agent' OR role = 'admin')", request.AgentID).First(&agent).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Agente no encontrado o no tiene permisos de agente"})
		return
	}

	// Actualizar ticket
	updates := map[string]interface{}{
		"assigned_to": request.AgentID,
		"status":      models.StatusAssigned,
		"updated_at":  time.Now(),
	}

	if err := db.Model(&ticket).Updates(updates).Error; err != nil {
		log.Printf("AssignTicket error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al asignar ticket"})
		return
	}

	// Obtener ticket actualizado
	if err := db.Where("id = ?", id).First(&ticket).Error; err != nil {
		log.Printf("AssignTicket error (reload): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar ticket actualizado"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// Helper para crear puntero a string
func strPtr(s string) *string {
	return &s
}
