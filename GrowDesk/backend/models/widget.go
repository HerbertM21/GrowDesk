package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// WidgetConfig almacena la configuración del widget de chat
type WidgetConfig struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name" gorm:"not null"`
	ApiKey         string    `json:"apiKey" gorm:"unique;not null"`
	AllowedDomains []string  `json:"allowedDomains" gorm:"type:json"`
	WelcomeMessage string    `json:"welcomeMessage"`
	PrimaryColor   string    `json:"primaryColor" gorm:"default:#4caf50"`
	BrandName      string    `json:"brandName"`
	Position       string    `json:"position" gorm:"default:bottom-right"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CreatedBy      string    `json:"createdBy"`
	IsActive       bool      `json:"isActive" gorm:"default:true"`
	EmbedCode      string    `json:"embedCode" gorm:"-"` // Campo calculado, no se guarda en BD
}

// GenerateApiKey genera una nueva API Key segura
func GenerateApiKey() (string, error) {
	bytes := make([]byte, 16) // 32 caracteres hex
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateEmbedCode crea el código HTML para incrustar el widget
func (w *WidgetConfig) GenerateEmbedCode(baseUrl string) string {
	// El puerto 8082 es el que corresponde a la API del widget
	apiUrl := "http://localhost:8082"

	// Construye el código de inserción basado en los parámetros de configuración
	return `<script src="` + baseUrl + `/widget.js" id="growdesk-widget"
  data-widget-id="` + w.ID + `"
  data-widget-token="` + w.ApiKey + `"
  data-api-url="` + apiUrl + `"
  data-brand-name="` + w.BrandName + `"
  data-welcome-message="` + w.WelcomeMessage + `"
  data-primary-color="` + w.PrimaryColor + `"
  data-position="` + w.Position + `">
</script>`
}

// GenerateEmbedCodeWithApiUrl permite especificar tanto la URL base como la URL de API
func (w *WidgetConfig) GenerateEmbedCodeWithApiUrl(baseUrl, apiUrl string) string {
	// código de inserción basado en los parámetros de configuración
	return `<script src="` + baseUrl + `/widget.js" id="growdesk-widget"
  data-widget-id="` + w.ID + `"
  data-widget-token="` + w.ApiKey + `"
  data-api-url="` + apiUrl + `"
  data-brand-name="` + w.BrandName + `"
  data-welcome-message="` + w.WelcomeMessage + `"
  data-primary-color="` + w.PrimaryColor + `"
  data-position="` + w.Position + `">
</script>`
}
