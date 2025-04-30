package models

import (
	"time"
)

// FAQ representa una pregunta frecuente
type FAQ struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Question    string    `json:"question" gorm:"not null;type:text"`
	Answer      string    `json:"answer" gorm:"not null;type:text"`
	Category    string    `json:"category" gorm:"not null;type:varchar(100)"`
	IsPublished bool      `json:"isPublished" gorm:"default:false"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
