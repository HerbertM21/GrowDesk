package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Message representa un mensaje en un ticket
type Message struct {
	ID          string       `json:"id" gorm:"primaryKey"`
	TicketID    string       `json:"ticketId" gorm:"not null;index"`
	Content     string       `json:"content" gorm:"not null;type:text"`
	IsClient    bool         `json:"isClient" gorm:"default:false"`
	IsInternal  bool         `json:"isInternal" gorm:"default:false"`
	UserID      *string      `json:"userId" gorm:"type:varchar(100)"`
	UserName    string       `json:"userName" gorm:"type:varchar(100)"`
	UserEmail   string       `json:"userEmail" gorm:"type:varchar(100)"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	Attachments []Attachment `json:"attachments,omitempty" gorm:"foreignKey:MessageID"`
}

// BeforeSave genera un ID único para el mensaje si no tiene uno
func (m *Message) BeforeSave(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = "MSG-" + uuid.New().String()
	}
	return nil
}
