package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Attachment representa un archivo adjunto a un mensaje
type Attachment struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	MessageID string    `json:"messageId" gorm:"not null;index"`
	Filename  string    `json:"filename" gorm:"not null;type:varchar(255)"`
	FilePath  string    `json:"filePath" gorm:"not null;type:varchar(255)"`
	FileSize  int64     `json:"fileSize" gorm:"not null"`
	MimeType  string    `json:"mimeType" gorm:"type:varchar(100)"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BeforeSave genera un ID único para el adjunto si no tiene uno
func (a *Attachment) BeforeSave(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = "ATTACH-" + uuid.New().String()
	}
	return nil
}
