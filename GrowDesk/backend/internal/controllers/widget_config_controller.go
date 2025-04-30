package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hmdev/GrowDesk/backend/models"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
)

// WidgetConfigRequest estructura para solicitudes de creación/actualización
type WidgetConfigRequest struct {
	Name           string   `json:"name" binding:"required"`
	BrandName      string   `json:"brandName" binding:"required"`
	AllowedDomains []string `json:"allowedDomains"`
	WelcomeMessage string   `json:"welcomeMessage"`
	PrimaryColor   string   `json:"primaryColor"`
	Position       string   `json:"position"`
}

// GetWidgetConfigs obtiene todas las configuraciones de widgets
func GetWidgetConfigs(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, devolver configuraciones de ejemplo
		mockConfigs := []models.WidgetConfig{
			{
				ID:             "WIDGET-123456",
				Name:           "Widget Principal",
				BrandName:      "GrowDesk",
				ApiKey:         "api-key-1234567890",
				AllowedDomains: []string{"example.com", "mysite.com"},
				WelcomeMessage: "¡Bienvenido! ¿En qué podemos ayudarte?",
				PrimaryColor:   "#2196F3",
				Position:       "right",
				CreatedAt:      time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt:      time.Now().Add(-5 * 24 * time.Hour),
			},
			{
				ID:             "WIDGET-789012",
				Name:           "Widget de Pruebas",
				BrandName:      "GrowDesk Test",
				ApiKey:         "api-key-0987654321",
				AllowedDomains: []string{"test.com"},
				WelcomeMessage: "Estamos en modo de prueba",
				PrimaryColor:   "#4CAF50",
				Position:       "left",
				CreatedAt:      time.Now().Add(-15 * 24 * time.Hour),
				UpdatedAt:      time.Now().Add(-2 * 24 * time.Hour),
			},
		}

		c.JSON(http.StatusOK, mockConfigs)
		return
	}

	var configs []models.WidgetConfig
	if err := db.Order("created_at desc").Find(&configs).Error; err != nil {
		log.Printf("GetWidgetConfigs error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener configuraciones"})
		return
	}

	c.JSON(http.StatusOK, configs)
}

// GetWidgetConfig obtiene una configuración de widget específica
func GetWidgetConfig(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, devolver configuración de ejemplo
		if id == "WIDGET-123456" {
			mockConfig := models.WidgetConfig{
				ID:             "WIDGET-123456",
				Name:           "Widget Principal",
				BrandName:      "GrowDesk",
				ApiKey:         "api-key-1234567890",
				AllowedDomains: []string{"example.com", "mysite.com"},
				WelcomeMessage: "¡Bienvenido! ¿En qué podemos ayudarte?",
				PrimaryColor:   "#2196F3",
				Position:       "right",
				CreatedAt:      time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt:      time.Now().Add(-5 * 24 * time.Hour),
			}

			c.JSON(http.StatusOK, mockConfig)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Configuración no encontrada"})
		return
	}

	var config models.WidgetConfig
	if err := db.Where("id = ?", id).First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuración no encontrada"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// CreateWidgetConfig crea una nueva configuración de widget
func CreateWidgetConfig(c *gin.Context) {
	var request WidgetConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	// Crear una nueva configuración
	config := models.WidgetConfig{
		Name:           request.Name,
		BrandName:      request.BrandName,
		ApiKey:         uuid.New().String(),
		AllowedDomains: request.AllowedDomains,
		WelcomeMessage: request.WelcomeMessage,
		PrimaryColor:   request.PrimaryColor,
		Position:       request.Position,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Valores predeterminados
	if config.PrimaryColor == "" {
		config.PrimaryColor = "#2196F3"
	}

	if config.Position == "" {
		config.Position = "right"
	}

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular creación
		config.ID = "WIDGET-" + time.Now().Format("20060102150405")

		c.JSON(http.StatusCreated, config)
		return
	}

	if err := db.Create(&config).Error; err != nil {
		log.Printf("CreateWidgetConfig error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear configuración"})
		return
	}

	c.JSON(http.StatusCreated, config)
}

// UpdateWidgetConfig actualiza una configuración de widget existente
func UpdateWidgetConfig(c *gin.Context) {
	id := c.Param("id")

	var request WidgetConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular actualización
		mockConfig := models.WidgetConfig{
			ID:             id,
			Name:           request.Name,
			BrandName:      request.BrandName,
			ApiKey:         "api-key-1234567890",
			AllowedDomains: request.AllowedDomains,
			WelcomeMessage: request.WelcomeMessage,
			PrimaryColor:   request.PrimaryColor,
			Position:       request.Position,
			UpdatedAt:      time.Now(),
		}

		c.JSON(http.StatusOK, mockConfig)
		return
	}

	var config models.WidgetConfig
	if err := db.Where("id = ?", id).First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuración no encontrada"})
		return
	}

	// Actualizar campos
	config.Name = request.Name
	config.BrandName = request.BrandName
	config.AllowedDomains = request.AllowedDomains
	config.WelcomeMessage = request.WelcomeMessage
	config.PrimaryColor = request.PrimaryColor
	config.Position = request.Position
	config.UpdatedAt = time.Now()

	if err := db.Save(&config).Error; err != nil {
		log.Printf("UpdateWidgetConfig error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar configuración"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// DeleteWidgetConfig elimina una configuración de widget
func DeleteWidgetConfig(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular eliminación
		c.JSON(http.StatusOK, gin.H{
			"message": "Configuración eliminada correctamente (simulado)",
			"id":      id,
		})
		return
	}

	result := db.Delete(&models.WidgetConfig{}, "id = ?", id)
	if result.Error != nil {
		log.Printf("DeleteWidgetConfig error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar configuración"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuración no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuración eliminada correctamente"})
}

// RegenerateApiKey regenera la API key para un widget
func RegenerateApiKey(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular regeneración
		newApiKey := uuid.New().String()
		c.JSON(http.StatusOK, gin.H{
			"message": "API key regenerada correctamente (simulado)",
			"id":      id,
			"apiKey":  newApiKey,
		})
		return
	}

	var config models.WidgetConfig
	if err := db.Where("id = ?", id).First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuración no encontrada"})
		return
	}

	// Generar nueva API key
	config.ApiKey = uuid.New().String()
	config.UpdatedAt = time.Now()

	if err := db.Save(&config).Error; err != nil {
		log.Printf("RegenerateApiKey error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al regenerar API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API key regenerada correctamente",
		"id":      id,
		"apiKey":  config.ApiKey,
	})
}
