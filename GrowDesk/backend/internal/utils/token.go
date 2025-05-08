package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Constantes para la generación de tokens
const (
	// Este sería típicamente proporcionado por variables de entorno
	jwtSecret = "your-super-secret-key-change-in-production"
	// El token expira en 24 horas
	tokenExpiration = 24 * time.Hour
)

// Claims define las reclamaciones personalizadas para JWT
type Claims struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken crea un nuevo token JWT para un usuario
func GenerateToken(userID, email, role string) (string, error) {
	// Establecer el tiempo de expiración
	expirationTime := time.Now().Add(tokenExpiration)

	// Crear las reclamaciones de JWT
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Crear el token con las reclamaciones
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("error firmando el token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken valida un token JWT y devuelve las reclamaciones
func ValidateToken(tokenString string) (*Claims, error) {
	// Comprobar si el token está vacío
	if tokenString == "" {
		return nil, fmt.Errorf("el token está vacío")
	}

	// Parsear el token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validar el método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error al parsear el token: %w", err)
	}

	// Validar el token
	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	// Extraer las reclamaciones
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("estructura de reclamaciones inválida")
	}

	// Comprobar si el token ha expirado
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("el token ha expirado")
	}

	return claims, nil
}

// GenerateMockToken genera un token fijo para pruebas
func GenerateMockToken() string {
	// Este coincide con el token utilizado en los servidores de prueba
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiJhZG1pbi0xMjMiLCJlbWFpbCI6ImFkbWluQGdyb3dkZXNrLmNvbSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTcyNDA4ODQwMH0.8J5ayPvA4B-1vF5NaqRXycW1pmIl9qjKP6Ddj4Ot_Cw"
}
