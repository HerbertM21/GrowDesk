package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User representa un usuario en el sistema
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	FirstName string    `json:"firstName" gorm:"not null"`
	LastName  string    `json:"lastName" gorm:"not null"`
	Password  string    `json:"-" gorm:"not null"`                       // Password is never returned in JSON
	Role      string    `json:"role" gorm:"not null;default:'customer'"` // customer, agent, admin
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BeforeSave es un hook de GORM que hashea la contraseña antes de guardar
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = generateUUID()
	}

	// Si la contraseña está siendo cambiada y es menor a 60 caracteres, la hashea
	if u.Password != "" && len(u.Password) < 60 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}

	return nil
}

// CheckPassword verifica la contraseña contra la contraseña hasheada
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
