package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TicketStatus es el estado de un ticket
type TicketStatus string

const (
	StatusOpen      TicketStatus = "open"
	StatusAssigned  TicketStatus = "assigned"
	StatusInProcess TicketStatus = "in_process"
	StatusResolved  TicketStatus = "resolved"
	StatusClosed    TicketStatus = "closed"
)

// TicketPriority es la prioridad de un ticket
type TicketPriority string

const (
	PriorityLow      TicketPriority = "low"
	PriorityMedium   TicketPriority = "medium"
	PriorityHigh     TicketPriority = "high"
	PriorityCritical TicketPriority = "critical"
)

// Customer representa información del cliente
type Customer struct {
	Name  string `json:"name" gorm:"type:varchar(100)"`
	Email string `json:"email" gorm:"type:varchar(100)"`
}

// Ticket representa un ticket de soporte
type Ticket struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null;type:varchar(150)"`
	Description string         `json:"description" gorm:"type:text"`
	Status      TicketStatus   `json:"status" gorm:"default:'open';type:varchar(20)"`
	Priority    TicketPriority `json:"priority" gorm:"default:'medium';type:varchar(20)"`
	Category    string         `json:"category" gorm:"type:varchar(50)"`
	Customer    Customer       `json:"customer" gorm:"embedded"`
	CreatedBy   string         `json:"createdBy" gorm:"type:varchar(100)"`
	AssignedTo  *string        `json:"assignedTo" gorm:"type:varchar(100)"`
	Source      string         `json:"source" gorm:"default:'backend';type:varchar(20)"` // backend, widget, api
	Metadata    string         `json:"metadata" gorm:"type:json"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Messages    []Message      `json:"messages,omitempty" gorm:"foreignKey:TicketID"`
}

// BeforeSave genera un ID único para el ticket si no tiene uno
func (t *Ticket) BeforeSave(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = "TICKET-" + uuid.New().String()
	}
	return nil
}
