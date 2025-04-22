package models

import (
	"time"

	"gorm.io/gorm"
)

// TicketStatus define los posibles valores de estado para un ticket
type TicketStatus string

const (
	StatusOpen       TicketStatus = "open"
	StatusAssigned   TicketStatus = "assigned"
	StatusInProgress TicketStatus = "in_progress"
	StatusResolved   TicketStatus = "resolved"
	StatusClosed     TicketStatus = "closed"
)

// TicketPriority define los posibles valores de prioridad para un ticket
type TicketPriority string

const (
	PriorityLow    TicketPriority = "low"
	PriorityMedium TicketPriority = "medium"
	PriorityHigh   TicketPriority = "high"
	PriorityUrgent TicketPriority = "urgent"
)

// Ticket representa un ticket de soporte en el sistema
type Ticket struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text;not null"`
	Status      TicketStatus   `json:"status" gorm:"not null;default:'open'"`
	Priority    TicketPriority `json:"priority" gorm:"not null;default:'medium'"`
	Category    string         `json:"category" gorm:"not null"`
	CreatedBy   string         `json:"createdBy" gorm:"not null"`
	AssignedTo  *string        `json:"assignedTo,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`

	// Relations
	User     User      `json:"-" gorm:"foreignKey:CreatedBy"`
	Agent    *User     `json:"-" gorm:"foreignKey:AssignedTo"`
	Messages []Message `json:"-" gorm:"foreignKey:TicketID"`
}

// BeforeCreate es un hook de GORM que genera un UUID para nuevos tickets
func (t *Ticket) BeforeSave(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = generateUUID()
	}
	return nil
}

type Tag struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
