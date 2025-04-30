package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Variables globales para la configuración de JWT
var (
	jwtSecret            string
	accessTokenDuration  = 24 * time.Hour  // 24 horas por defecto
	refreshTokenDuration = 720 * time.Hour // 30 días por defecto
)

// Claims estructura personalizada para JWT claims
type Claims struct {
	UserID    string `json:"userId"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
	TokenType string `json:"tokenType"` // access o refresh
	jwt.RegisteredClaims
}

// SetJWTSecret establece la clave secreta para firmar tokens JWT
func SetJWTSecret(secret string) {
	jwtSecret = secret
}

// SetTokenDurations configura la duración de los tokens
func SetTokenDurations(accessDuration, refreshDuration time.Duration) {
	if accessDuration > 0 {
		accessTokenDuration = accessDuration
	}
	if refreshDuration > 0 {
		refreshTokenDuration = refreshDuration
	}
}

// GenerateAccessToken genera un token JWT de acceso con duración limitada
func GenerateAccessToken(userID, email, firstName, lastName, role string) (string, time.Time, error) {
	return generateToken(userID, email, firstName, lastName, role, "access", accessTokenDuration)
}

// GenerateRefreshToken genera un token JWT de renovación con mayor duración
func GenerateRefreshToken(userID, email string) (string, time.Time, error) {
	return generateToken(userID, email, "", "", "", "refresh", refreshTokenDuration)
}

// generateToken genera un token JWT con la información proporcionada
func generateToken(userID, email, firstName, lastName, role, tokenType string, duration time.Duration) (string, time.Time, error) {
	// Verificar que se haya establecido la clave secreta
	if jwtSecret == "" {
		return "", time.Time{}, errors.New("la clave secreta JWT no está configurada")
	}

	// Calcular tiempo de expiración
	expirationTime := time.Now().Add(duration)

	// Crear claims con datos del usuario
	claims := &Claims{
		UserID:    userID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "growdesk-api",
			Subject:   userID,
		},
	}

	// Crear token con los claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar token con la clave secreta
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// ValidateToken valida un token JWT y retorna los claims si es válido
func ValidateToken(tokenString string) (*Claims, error) {
	// Verificar que se haya establecido la clave secreta
	if jwtSecret == "" {
		return nil, errors.New("la clave secreta JWT no está configurada")
	}

	// Crear objeto de claims personalizado
	claims := &Claims{}

	// Parsear y validar token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verificar algoritmo de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inesperado")
		}
		return []byte(jwtSecret), nil
	})

	// Verificar errores de parsing
	if err != nil {
		return nil, err
	}

	// Verificar validez del token
	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}

// ValidateRefreshToken valida que un token sea de tipo refresh
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Verificar que sea un token de renovación
	if claims.TokenType != "refresh" {
		return nil, errors.New("el token no es de tipo refresh")
	}

	return claims, nil
}

// ExtractUserFromToken extrae la información del usuario desde un token JWT
func ExtractUserFromToken(tokenString string) (string, string, string, string, string, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", "", "", "", "", err
	}

	return claims.UserID, claims.Email, claims.FirstName, claims.LastName, claims.Role, nil
}
