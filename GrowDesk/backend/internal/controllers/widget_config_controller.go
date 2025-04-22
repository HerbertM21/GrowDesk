package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hmdev/GrowDesk/backend/models"
)

// Estructura para solicitudes de creación/actualización
type WidgetConfigRequest struct {
	Name           string   `json:"name" binding:"required"`
	AllowedDomains []string `json:"allowedDomains"`
	WelcomeMessage string   `json:"welcomeMessage"`
	PrimaryColor   string   `json:"primaryColor"`
	BrandName      string   `json:"brandName" binding:"required"`
	Position       string   `json:"position"`
}

// GetWidgetConfigs obtiene todas las configuraciones de widgets
func GetWidgetConfigs(c *gin.Context) {
	// implementar despues
	// var configs []models.WidgetConfig
	// database.DB.Find(&configs)

	// Para demo, devolver una configuración de ejemplo
	configs := []models.WidgetConfig{
		{
			ID:             "default-widget",
			Name:           "Widget Principal",
			ApiKey:         "demo-token",
			AllowedDomains: []string{"localhost", "mitienda.com"},
			WelcomeMessage: "¿En qué podemos ayudarte hoy?",
			PrimaryColor:   "#4caf50",
			BrandName:      "MiTienda",
			Position:       "bottom-right",
			CreatedAt:      time.Now().Add(-24 * time.Hour),
			UpdatedAt:      time.Now(),
			CreatedBy:      "admin",
			IsActive:       true,
		},
	}

	// Generar el código de inserción para cada configuración
	baseUrl := "http://localhost:3030"
	for i := range configs {
		configs[i].EmbedCode = generateEmbedCode(configs[i].ID, configs[i].ApiKey, configs[i].BrandName, configs[i].WelcomeMessage, configs[i].PrimaryColor, configs[i].Position, baseUrl, "http://localhost:8082")
	}

	c.JSON(http.StatusOK, gin.H{
		"configs": configs,
	})
}

// GetWidgetConfig obtiene una configuración específica
func GetWidgetConfig(c *gin.Context) {
	id := c.Param("id")

	// implementar despues
	// var config models.WidgetConfig
	// if err := database.DB.Where("id = ?", id).First(&config).Error; err != nil {
	//     c.JSON(http.StatusNotFound, gin.H{"error": "Configuración no encontrada"})
	//     return
	// }

	// Para demo, devolver una configuración de ejemplo
	config := models.WidgetConfig{
		ID:             id,
		Name:           "Widget Principal",
		ApiKey:         "demo-token",
		AllowedDomains: []string{"localhost", "growdesk.com"},
		WelcomeMessage: "¿En qué podemos ayudarte hoy?",
		PrimaryColor:   "#4caf50",
		BrandName:      "vidrieria",
		Position:       "bottom-right",
		CreatedAt:      time.Now().Add(-24 * time.Hour),
		UpdatedAt:      time.Now(),
		CreatedBy:      "admin",
		IsActive:       true,
	}

	// Generar el código de inserción
	baseUrl := "http://localhost:3030"
	widgetApiUrl := "http://localhost:8082"
	config.EmbedCode = generateEmbedCode(config.ID, config.ApiKey, config.BrandName, config.WelcomeMessage, config.PrimaryColor, config.Position, baseUrl, widgetApiUrl)

	c.JSON(http.StatusOK, config)
}

// CreateWidgetConfig crea una nueva configuración de widget
func CreateWidgetConfig(c *gin.Context) {
	var req WidgetConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generar API Key
	apiKey, err := models.GenerateApiKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar API Key"})
		return
	}

	// Obtener ID del usuario actual
	userID, _ := c.Get("userID")
	if userID == nil {
		userID = "admin" // Valor por defecto para demo
	}

	// Crear nueva configuración
	config := models.WidgetConfig{
		ID:             uuid.New().String(),
		Name:           req.Name,
		ApiKey:         apiKey,
		AllowedDomains: req.AllowedDomains,
		WelcomeMessage: req.WelcomeMessage,
		PrimaryColor:   req.PrimaryColor,
		BrandName:      req.BrandName,
		Position:       req.Position,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      userID.(string),
		IsActive:       true,
	}

	// implementar dsp
	// if err := database.DB.Create(&config).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar configuración"})
	// 	return
	// }

	// Generar el código de inserción
	baseUrl := "http://localhost:3030"
	widgetApiUrl := "http://localhost:8082"
	config.EmbedCode = generateEmbedCode(config.ID, config.ApiKey, config.BrandName, config.WelcomeMessage, config.PrimaryColor, config.Position, baseUrl, widgetApiUrl)

	c.JSON(http.StatusCreated, config)
}

// UpdateWidgetConfig actualiza una configuración existente
func UpdateWidgetConfig(c *gin.Context) {
	id := c.Param("id")
	var req WidgetConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// implementar dsp
	// var config models.WidgetConfig
	// if err := database.DB.Where("id = ?", id).First(&config).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Configuración no encontrada"})
	// 	return
	// }

	// Para demo, simular la actualización
	config := models.WidgetConfig{
		ID:             id,
		Name:           req.Name,
		ApiKey:         "demo-token", // Mantener la API key existente
		AllowedDomains: req.AllowedDomains,
		WelcomeMessage: req.WelcomeMessage,
		PrimaryColor:   req.PrimaryColor,
		BrandName:      req.BrandName,
		Position:       req.Position,
		CreatedAt:      time.Now().Add(-24 * time.Hour),
		UpdatedAt:      time.Now(),
		CreatedBy:      "admin",
		IsActive:       true,
	}

	// En una implementación real
	// database.DB.Save(&config)

	// Generar el código de inserción actualizado
	baseUrl := "http://localhost:3030"
	widgetApiUrl := "http://localhost:8082"
	config.EmbedCode = generateEmbedCode(config.ID, config.ApiKey, config.BrandName, config.WelcomeMessage, config.PrimaryColor, config.Position, baseUrl, widgetApiUrl)

	c.JSON(http.StatusOK, config)
}

// RegenerateApiKey genera una nueva API key para un widget
func RegenerateApiKey(c *gin.Context) {
	id := c.Param("id")

	// Generar nueva API Key
	apiKey, err := models.GenerateApiKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar API Key"})
		return
	}

	//
	config := models.WidgetConfig{
		ID:             id,
		Name:           "Widget Principal",
		ApiKey:         apiKey, // Nueva API key
		AllowedDomains: []string{"localhost", "growdesk.com"},
		WelcomeMessage: "¿En qué podemos ayudarte hoy?",
		PrimaryColor:   "#4caf50",
		BrandName:      "vidrieria",
		Position:       "bottom-right",
		CreatedAt:      time.Now().Add(-24 * time.Hour),
		UpdatedAt:      time.Now(),
		CreatedBy:      "admin",
		IsActive:       true,
	}

	// implementar dsp
	// config.ApiKey = apiKey
	// config.UpdatedAt = time.Now()
	// database.DB.Save(&config)

	// Generar el código de inserción con la nueva API key
	baseUrl := "http://localhost:3030"
	widgetApiUrl := "http://localhost:8082"
	config.EmbedCode = generateEmbedCode(config.ID, config.ApiKey, config.BrandName, config.WelcomeMessage, config.PrimaryColor, config.Position, baseUrl, widgetApiUrl)

	c.JSON(http.StatusOK, config)
}

// DeleteWidgetConfig desactiva una configuración de widget
func DeleteWidgetConfig(c *gin.Context) {
	id := c.Param("id")

	// implementar dsp
	// var config models.WidgetConfig
	// if err := database.DB.Where("id = ?", id).First(&config).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Configuración no encontrada"})
	// 	return
	// }
	// config.IsActive = false
	// database.DB.Save(&config)

	c.JSON(http.StatusOK, gin.H{
		"message": "Configuración con ID " + id + " desactivada correctamente",
	})
}

// generateEmbedCode crea el código HTML para incrustar el widget
func generateEmbedCode(widgetId, widgetToken, brandName, welcomeMessage, primaryColor, position, baseUrl, apiUrl string) string {
	return `<script src="` + baseUrl + `/widget.js" id="growdesk-widget"
  data-widget-id="` + widgetId + `"
  data-widget-token="` + widgetToken + `"
  data-api-url="` + apiUrl + `"
  data-brand-name="` + brandName + `"
  data-welcome-message="` + welcomeMessage + `"
  data-primary-color="` + primaryColor + `"
  data-position="` + position + `">
</script>`
}
