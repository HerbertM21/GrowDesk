package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Category representa una categoría para los tickets
type Category struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;type:varchar(100)"`
	Description string    `json:"description" gorm:"type:text"`
	Color       string    `json:"color" gorm:"default:'#2196F3';type:varchar(20)"`
	Icon        string    `json:"icon" gorm:"default:'category';type:varchar(50)"`
	Active      bool      `json:"active" gorm:"default:true"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// BeforeSave genera un ID único para la categoría si no tiene uno
func (c *Category) BeforeSave(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = "CAT-" + uuid.New().String()
	}
	return nil
}
