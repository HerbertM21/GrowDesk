package models

import (
	"time"

	"gorm.io/gorm"
)

// Attachment representa un archivo adjunto a un mensaje
type Attachment struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	MessageID string    `json:"messageId" gorm:"not null"`
	FileName  string    `json:"fileName" gorm:"not null"`
	FileType  string    `json:"fileType" gorm:"not null"`
	FileSize  int64     `json:"fileSize" gorm:"not null"`
	FileURL   string    `json:"fileURL" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Message Message `json:"-" gorm:"foreignKey:MessageID"`
}

// BeforeSave es un hook de GORM que genera un UUID para nuevos adjuntos
func (a *Attachment) BeforeSave(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = generateUUID()
	}
	return nil
}
