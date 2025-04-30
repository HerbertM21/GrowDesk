package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/models"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
	"gorm.io/gorm"
)

// WidgetResponse estructura para respuestas al widget
type WidgetResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// WidgetTicketData estructura para recibir datos desde el widget
type WidgetTicketData struct {
	Name     string                 `json:"name" binding:"required"`
	Email    string                 `json:"email" binding:"required,email"`
	Message  string                 `json:"message" binding:"required"`
	Metadata map[string]interface{} `json:"metadata"`
}

// WidgetMessageData estructura para recibir mensajes desde el widget
type WidgetMessageData struct {
	TicketID string `json:"ticketId" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

// GetWidgetConfig retorna la configuración pública del widget
func GetWidgetConfig(c *gin.Context) {
	widgetConfig, exists := c.Get("widgetConfig")

	// Si la configuración no está en el contexto, usar ID
	if !exists {
		widgetID, exists := c.Get("widgetID")
		if !exists {
			c.JSON(http.StatusBadRequest, WidgetResponse{
				Status:  "error",
				Message: "ID de widget no encontrado",
			})
			return
		}

		db := database.GetDB()
		if db != nil {
			var config models.WidgetConfig
			if err := db.Where("id = ?", widgetID).First(&config).Error; err != nil {
				c.JSON(http.StatusNotFound, WidgetResponse{
					Status:  "error",
					Message: "Configuración de widget no encontrada",
				})
				return
			}

			widgetConfig = config
		} else {
			// En modo desarrollo sin DB, crear configuración ficticia
			widgetConfig = models.WidgetConfig{
				ID:             widgetID.(string),
				Name:           "Widget Demo",
				BrandName:      "GrowDesk Demo",
				WelcomeMessage: "¡Bienvenido al soporte! ¿En qué podemos ayudarte?",
				PrimaryColor:   "#2196F3",
				Position:       "right",
			}
		}
	}

	// Extraer solo los campos públicos necesarios para el widget
	config := widgetConfig.(models.WidgetConfig)
	publicConfig := map[string]interface{}{
		"id":             config.ID,
		"brandName":      config.BrandName,
		"welcomeMessage": config.WelcomeMessage,
		"primaryColor":   config.PrimaryColor,
		"position":       config.Position,
	}

	c.JSON(http.StatusOK, WidgetResponse{
		Status: "success",
		Data:   publicConfig,
	})
}

// CreateWidgetTicket crea un ticket desde el widget
func CreateWidgetTicket(c *gin.Context) {
	widgetID, exists := c.Get("widgetID")
	if !exists {
		c.JSON(http.StatusBadRequest, WidgetResponse{
			Status:  "error",
			Message: "ID de widget no encontrado",
		})
		return
	}

	var request struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, WidgetResponse{
			Status:  "error",
			Message: "Datos inválidos: " + err.Error(),
		})
		return
	}

	// Crear ticket
	ticket := models.Ticket{
		Title:       request.Title,
		Description: request.Description,
		Status:      models.StatusOpen,
		Priority:    models.PriorityMedium,
		Customer: models.Customer{
			Name:  request.Name,
			Email: request.Email,
		},
		CreatedBy: request.Email,
		Source:    "widget",
		Metadata:  `{"widgetId":"` + widgetID.(string) + `"}`,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular creación
		ticket.ID = "TICKET-W-" + time.Now().Format("20060102150405")

		// Mensaje inicial
		message := models.Message{
			ID:        "MSG-W-" + time.Now().Format("20060102150405"),
			TicketID:  ticket.ID,
			Content:   request.Description,
			IsClient:  true,
			UserName:  request.Name,
			UserEmail: request.Email,
			CreatedAt: time.Now(),
		}

		ticket.Messages = append(ticket.Messages, message)

		c.JSON(http.StatusCreated, WidgetResponse{
			Status: "success",
			Data:   ticket,
		})
		return
	}

	// Crear ticket y mensaje inicial en transacción
	tx := db.Begin()

	if err := tx.Create(&ticket).Error; err != nil {
		tx.Rollback()
		log.Printf("CreateWidgetTicket error: %v", err)
		c.JSON(http.StatusInternalServerError, WidgetResponse{
			Status:  "error",
			Message: "Error al crear el ticket",
		})
		return
	}

	// Crear mensaje inicial
	message := models.Message{
		TicketID:  ticket.ID,
		Content:   request.Description,
		IsClient:  true,
		UserName:  request.Name,
		UserEmail: request.Email,
		CreatedAt: time.Now(),
	}

	if err := tx.Create(&message).Error; err != nil {
		tx.Rollback()
		log.Printf("CreateWidgetTicket message error: %v", err)
		c.JSON(http.StatusInternalServerError, WidgetResponse{
			Status:  "error",
			Message: "Error al crear el mensaje inicial",
		})
		return
	}

	tx.Commit()

	// Cargar el ticket completo con mensaje
	db.Where("id = ?", ticket.ID).Preload("Messages").First(&ticket)

	c.JSON(http.StatusCreated, WidgetResponse{
		Status: "success",
		Data:   ticket,
	})
}

// GetWidgetTicket obtiene un ticket desde el widget
func GetWidgetTicket(c *gin.Context) {
	ticketID := c.Param("ticketId")
	email := c.Query("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, WidgetResponse{
			Status:  "error",
			Message: "Se requiere el correo electrónico para verificar acceso",
		})
		return
	}

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular ticket
		if ticketID == "TICKET-W-123456" {
			mockTicket := models.Ticket{
				ID:          ticketID,
				Title:       "Problema con el pedido #123",
				Description: "No he recibido mi pedido",
				Status:      models.StatusOpen,
				Priority:    models.PriorityMedium,
				Customer: models.Customer{
					Name:  "Cliente Widget",
					Email: email,
				},
				CreatedBy: email,
				Source:    "widget",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
				Messages: []models.Message{
					{
						ID:        "MSG-W-123456-1",
						TicketID:  ticketID,
						Content:   "No he recibido mi pedido #123",
						IsClient:  true,
						UserName:  "Cliente Widget",
						UserEmail: email,
						CreatedAt: time.Now().Add(-24 * time.Hour),
					},
					{
						ID:        "MSG-W-123456-2",
						TicketID:  ticketID,
						Content:   "Estamos verificando el estado de su pedido, nos comunicaremos pronto",
						IsClient:  false,
						UserName:  "Soporte",
						CreatedAt: time.Now().Add(-20 * time.Hour),
					},
				},
			}

			// Verificar correo
			if mockTicket.Customer.Email == email {
				c.JSON(http.StatusOK, WidgetResponse{
					Status: "success",
					Data:   mockTicket,
				})
				return
			}
		}

		c.JSON(http.StatusNotFound, WidgetResponse{
			Status:  "error",
			Message: "Ticket no encontrado o no tiene acceso",
		})
		return
	}

	var ticket models.Ticket
	result := db.Where("id = ?", ticketID).Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("messages.created_at ASC")
	}).First(&ticket)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, WidgetResponse{
			Status:  "error",
			Message: "Ticket no encontrado",
		})
		return
	}

	// Verificar que el correo electrónico coincida con el del ticket
	if ticket.Customer.Email != email {
		c.JSON(http.StatusForbidden, WidgetResponse{
			Status:  "error",
			Message: "No tiene acceso a este ticket",
		})
		return
	}

	c.JSON(http.StatusOK, WidgetResponse{
		Status: "success",
		Data:   ticket,
	})
}

// SendWidgetMessage envía un mensaje a un ticket desde el widget
func SendWidgetMessage(c *gin.Context) {
	ticketID := c.Param("ticketId")

	var request struct {
		Content string `json:"content" binding:"required"`
		Email   string `json:"email" binding:"required,email"`
		Name    string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, WidgetResponse{
			Status:  "error",
			Message: "Datos inválidos: " + err.Error(),
		})
		return
	}

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular envío
		mockMessage := models.Message{
			ID:        "MSG-W-" + time.Now().Format("20060102150405"),
			TicketID:  ticketID,
			Content:   request.Content,
			IsClient:  true,
			UserName:  request.Name,
			UserEmail: request.Email,
			CreatedAt: time.Now(),
		}

		c.JSON(http.StatusCreated, WidgetResponse{
			Status: "success",
			Data:   mockMessage,
		})
		return
	}

	// Verificar que el ticket existe
	var ticket models.Ticket
	if err := db.Where("id = ?", ticketID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, WidgetResponse{
			Status:  "error",
			Message: "Ticket no encontrado",
		})
		return
	}

	// Verificar que el correo electrónico coincida con el del ticket
	if ticket.Customer.Email != request.Email {
		c.JSON(http.StatusForbidden, WidgetResponse{
			Status:  "error",
			Message: "No tiene acceso a este ticket",
		})
		return
	}

	// Crear nuevo mensaje
	message := models.Message{
		TicketID:  ticketID,
		Content:   request.Content,
		IsClient:  true,
		UserName:  request.Name,
		UserEmail: request.Email,
		CreatedAt: time.Now(),
	}

	if err := db.Create(&message).Error; err != nil {
		log.Printf("SendWidgetMessage error: %v", err)
		c.JSON(http.StatusInternalServerError, WidgetResponse{
			Status:  "error",
			Message: "Error al enviar mensaje",
		})
		return
	}

	// Actualizar timestamp del ticket
	db.Model(&ticket).Update("updated_at", time.Now())

	c.JSON(http.StatusCreated, WidgetResponse{
		Status: "success",
		Data:   message,
	})
}

// GetWidgetFaqs obtiene las FAQs disponibles para el widget
func GetWidgetFaqs(c *gin.Context) {
	category := c.Query("category")

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, devolver FAQs de ejemplo
		mockFaqs := []models.FAQ{
			{
				ID:          1,
				Question:    "¿Cómo puedo rastrear mi pedido?",
				Answer:      "Puede rastrear su pedido en la sección 'Mis Pedidos' con su número de pedido.",
				Category:    "Pedidos",
				IsPublished: true,
			},
			{
				ID:          2,
				Question:    "¿Cómo puedo cambiar mi contraseña?",
				Answer:      "Puede cambiar su contraseña en la sección 'Mi Cuenta'.",
				Category:    "Cuenta",
				IsPublished: true,
			},
			{
				ID:          3,
				Question:    "¿Cuáles son los métodos de pago aceptados?",
				Answer:      "Aceptamos tarjetas de crédito, PayPal y transferencia bancaria.",
				Category:    "Pagos",
				IsPublished: true,
			},
		}

		// Filtrar por categoría si se especifica
		if category != "" {
			var filtered []models.FAQ
			for _, faq := range mockFaqs {
				if faq.Category == category {
					filtered = append(filtered, faq)
				}
			}
			c.JSON(http.StatusOK, WidgetResponse{
				Status: "success",
				Data:   filtered,
			})
			return
		}

		c.JSON(http.StatusOK, WidgetResponse{
			Status: "success",
			Data:   mockFaqs,
		})
		return
	}

	// Construir consulta
	query := db.Model(&models.FAQ{}).Where("is_published = ?", true)

	// Filtrar por categoría si se especifica
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var faqs []models.FAQ
	if err := query.Order("category, id").Find(&faqs).Error; err != nil {
		log.Printf("GetWidgetFaqs error: %v", err)
		c.JSON(http.StatusInternalServerError, WidgetResponse{
			Status:  "error",
			Message: "Error al obtener preguntas frecuentes",
		})
		return
	}

	c.JSON(http.StatusOK, WidgetResponse{
		Status: "success",
		Data:   faqs,
	})
}
