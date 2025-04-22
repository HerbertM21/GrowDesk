package auth

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// TokenDetails contiene los metadatos de un token
type TokenDetails struct {
	UserID string
	Email  string
	Role   string
	Exp    int64
}

var secretKey string

// InitAuth inicializa el sistema de autenticación
func init() {
	// Intentar cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	secretKey = getSecretKey()
	log.Println("JWT authentication system initialized")
}

// InitJWT inicializa el sistema de autenticación con una clave secreta específica
func InitJWT(secret string) {
	if secret != "" {
		secretKey = secret
		log.Println("JWT authentication system initialized with custom secret")
	} else {
		secretKey = getSecretKey()
		log.Println("JWT authentication system initialized with default secret")
	}
}

// GenerateToken crea un nuevo token JWT
func GenerateToken(userID, email, role string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24 * 3) // Token expira en 3 días

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(), // Tiempo de emisión
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", err
	}

	return tokenString, nil
}

// VerifyToken valida un token JWT
func VerifyToken(tokenString string) (*TokenDetails, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("reclamaciones de token inválidas")
	}

	// Verificar si el token ha expirado
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, errors.New("el token ha expirado")
		}
	} else {
		return nil, errors.New("fecha de expiración no válida en el token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("ID de usuario inválido en el token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("email inválido en el token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("rol inválido en el token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("tiempo de expiración inválido en el token")
	}

	return &TokenDetails{
		UserID: userID,
		Email:  email,
		Role:   role,
		Exp:    int64(exp),
	}, nil
}

// RefreshToken genera un nuevo token basado en el actual
func RefreshToken(td *TokenDetails) (string, error) {
	// Verificar que el token no esté a punto de expirar
	if time.Now().Unix() > td.Exp-300 { // 5 minutos antes de la expiración
		return GenerateToken(td.UserID, td.Email, td.Role)
	}
	return "", errors.New("el token actual aún es válido")
}

// getSecretKey obtiene la clave secreta JWT de variables de entorno
func getSecretKey() string {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		// en producción, esto debería fallar, pero para desarrollo permitimos un valor predeterminado
		defaultKey := "growdesk-development-secret-key-2024"
		log.Printf("ADVERTENCIA: JWT_SECRET no está configurado. Usando clave de desarrollo por defecto. NO USAR EN PRODUCCIÓN.")
		return defaultKey
	}
	return key
}
