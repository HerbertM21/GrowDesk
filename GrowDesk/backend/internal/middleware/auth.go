package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/utils"
)

// ContextKey representa una clave para los valores del contexto
type ContextKey string

// Constantes para las claves del contexto
const (
	UserIDKey ContextKey = "userID"
	EmailKey  ContextKey = "email"
	RoleKey   ContextKey = "role"
)

// ExtractToken extrae el token JWT del encabezado de autorización
func ExtractToken(r *http.Request) string {
	// Obtener el encabezado de autorización
	authHeader := r.Header.Get("Authorization")

	// Comprobar si el encabezado está presente
	if authHeader == "" {
		return ""
	}

	// Comprobar si el encabezado sigue el patrón "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// Middleware de autenticación para validar tokens JWT
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extraer token de la solicitud
		tokenString := ExtractToken(r)

		// Comprobar si el token existe
		if tokenString == "" {
			http.Error(w, "No autorizado: No se proporcionó un token", http.StatusUnauthorized)
			return
		}

		// Validar token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "No autorizado: Token inválido", http.StatusUnauthorized)
			return
		}

		// Agregar reclamaciones al contexto
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)

		// Llamar al siguiente controlador con el contexto actualizado
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole middleware para verificar si el usuario tiene un rol específico
func RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener el rol del usuario desde el contexto
		userRole, ok := r.Context().Value(RoleKey).(string)
		if !ok {
			http.Error(w, "No autorizado: No se proporcionó información de rol", http.StatusUnauthorized)
			return
		}

		// Comprobar si el usuario tiene el rol requerido
		if userRole != role {
			http.Error(w, "Prohibido: Permisos insuficientes", http.StatusForbidden)
			return
		}

		// Llamar al siguiente controlador
		next.ServeHTTP(w, r)
	})
}

// MockAuth middleware de autenticación para desarrollo
// Este siempre tiene éxito y establece el rol de administrador
func MockAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Agregar reclamaciones de prueba al contexto
		ctx := context.WithValue(r.Context(), UserIDKey, "admin-123")
		ctx = context.WithValue(ctx, EmailKey, "admin@growdesk.com")
		ctx = context.WithValue(ctx, RoleKey, "admin")

		// Llamar al siguiente controlador con el contexto actualizado
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
