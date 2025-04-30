package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WidgetConfig representa la configuración de un widget para integración en sitios web
type WidgetConfig struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name" gorm:"not null;type:varchar(100)"`
	BrandName      string    `json:"brandName" gorm:"not null;type:varchar(100)"`
	ApiKey         string    `json:"apiKey" gorm:"unique;not null;type:varchar(100)"`
	AllowedDomains []string  `json:"allowedDomains" gorm:"type:json"`
	WelcomeMessage string    `json:"welcomeMessage" gorm:"type:text"`
	PrimaryColor   string    `json:"primaryColor" gorm:"default:'#2196F3';type:varchar(20)"`
	Position       string    `json:"position" gorm:"default:'right';type:varchar(20)"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// BeforeSave genera un ID único y una API key para la configuración de widget si no tiene
func (w *WidgetConfig) BeforeSave(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = "WIDGET-" + uuid.New().String()
	}

	if w.ApiKey == "" {
		w.ApiKey = uuid.New().String()
	}

	return nil
}
