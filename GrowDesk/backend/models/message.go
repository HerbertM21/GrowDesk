package models

import (
	"time"

	"gorm.io/gorm"
)

// Message representa un mensaje en una conversación de ticket
type Message struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	TicketID   string    `json:"ticketId" gorm:"not null"`
	UserID     string    `json:"userId" gorm:"not null"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	IsInternal bool      `json:"isInternal" gorm:"not null;default:false"` // Internal notes visible only to agents
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	// Relations
	Ticket      Ticket       `json:"-" gorm:"foreignKey:TicketID"`
	User        User         `json:"-" gorm:"foreignKey:UserID"`
	Attachments []Attachment `json:"attachments,omitempty" gorm:"foreignKey:MessageID"`
}

// BeforeCreate es un hook de GORM que genera un UUID para nuevos mensajes
func (m *Message) BeforeSave(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = generateUUID()
	}
	return nil
}
