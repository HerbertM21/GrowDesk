package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User representa un usuario en el sistema
type User struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	Email      string    `json:"email" gorm:"unique;not null"`
	FirstName  string    `json:"firstName" gorm:"not null"`
	LastName   string    `json:"lastName" gorm:"not null"`
	Password   string    `json:"-" gorm:"not null"`                       // Password is never returned in JSON
	Role       string    `json:"role" gorm:"not null;default:'customer'"` // customer, agent, admin
	Department string    `json:"department"`
	Active     bool      `json:"active" gorm:"default:true"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// BeforeSave es un hook de GORM que hashea la contraseña antes de guardar
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	// Si la contraseña no está hasheada (longitud típica de hash bcrypt > 50)
	if len(u.Password) < 50 {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPass)
	}
	return nil
}

// CheckPassword verifica la contraseña contra la contraseña hasheada
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ExcludePassword retorna una copia del usuario sin el campo contraseña
func (u *User) ExcludePassword() User {
	userCopy := *u
	userCopy.Password = ""
	return userCopy
}
