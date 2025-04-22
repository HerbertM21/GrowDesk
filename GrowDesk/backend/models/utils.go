package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DB es un alias para gorm.DB para simplificar las referencias
type DB = gorm.DB

// GenerateUUID crea una nueva cadena UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// generateUUID es un helper interno para los modelos
func generateUUID() string {
	return GenerateUUID()
}

// generateSecureToken genera un token seguro aleatorio de la longitud especificada
func generateSecureToken(length int) string {
	// El token será el doble de largo ya que cada byte se convierte en 2 caracteres hexadecimales
	b := make([]byte, length/2)
	if _, err := rand.Read(b); err != nil {
		// En caso de error, usar un método alternativo menos seguro pero funcional
		return generateFallbackToken(length)
	}

	return hex.EncodeToString(b)
}

// generateFallbackToken genera un token usando un método menos seguro pero funcional
func generateFallbackToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Crear un slice para almacenar el token
	b := make([]byte, length)

	seed := time.Now().UnixNano()

	for i := range b {
		// Usar la semilla para generar un índice
		seed = (seed*1103515245 + 12345) & 0x7fffffff
		idx := int(seed) % len(charset)
		b[i] = charset[idx]
	}

	return string(b)
}
