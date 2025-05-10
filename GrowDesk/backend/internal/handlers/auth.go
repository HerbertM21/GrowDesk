package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/data"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/middleware"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/utils"
)

// AuthHandler contiene manejadores para autenticación
type AuthHandler struct {
	Store data.DataStore
}

// Login maneja solicitudes de inicio de sesión de usuarios
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear el cuerpo de la solicitud
	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "El cuerpo de la solicitud es inválido", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if loginReq.Email == "" || loginReq.Password == "" {
		http.Error(w, "Email y contraseña son requeridos", http.StatusBadRequest)
		return
	}

	// En una implementación real, validaríamos las credenciales contra la base de datos
	// Por ahora, usaremos un token fijo para cualquier inicio de sesión válido

	// Generar token
	token := utils.GenerateMockToken()

	// Preparar respuesta
	resp := models.AuthResponse{
		Token: token,
		User: models.User{
			ID:        "admin-123",
			Email:     loginReq.Email,
			FirstName: "Admin",
			LastName:  "User",
			Role:      "admin",
		},
	}

	// Devolver token y información de usuario
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Register maneja el registro de usuarios
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear el cuerpo de la solicitud
	var registerReq models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		http.Error(w, "El cuerpo de la solicitud es inválido", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if registerReq.Email == "" || registerReq.Password == "" ||
		registerReq.FirstName == "" || registerReq.LastName == "" {
		http.Error(w, "Todos los campos son requeridos", http.StatusBadRequest)
		return
	}

	// En una implementación real, almacenaríamos al usuario en la base de datos
	// Por ahora, devolveremos una respuesta de éxito con un token fijo

	// Generar token
	token := utils.GenerateMockToken()

	// Preparar respuesta
	resp := models.AuthResponse{
		Token: token,
		User: models.User{
			ID:        "user-" + utils.GenerateTimestamp(),
			Email:     registerReq.Email,
			FirstName: registerReq.FirstName,
			LastName:  registerReq.LastName,
			Role:      "customer",
		},
	}

	// Devolver token y información de usuario
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Me devuelve la información del usuario actual
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener información del usuario desde el contexto
	// Esto será establecido por el middleware de autenticación
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	email, _ := r.Context().Value(middleware.EmailKey).(string)
	role, _ := r.Context().Value(middleware.RoleKey).(string)

	// Preparar respuesta
	user := models.User{
		ID:        userID,
		Email:     email,
		FirstName: "Admin",
		LastName:  "User",
		Role:      role,
	}

	// Devolver información de usuario
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
