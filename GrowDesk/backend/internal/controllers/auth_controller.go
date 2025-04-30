package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/models"
	"github.com/hmdev/GrowDesk/backend/pkg/auth"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
)

// LoginRequest estructura para datos de inicio de sesión
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest estructura para datos de registro
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

// AuthResponse estructura para respuesta de autenticación
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Login maneja el proceso de inicio de sesión
func Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Login error: invalid request format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de inicio de sesión inválidos", "details": err.Error()})
		return
	}

	log.Printf("Intento de inicio de sesión para: %s", request.Email)

	db := database.GetDB()
	if db == nil {
		log.Printf("Login: database connection not available, using mock user for: %s", request.Email)

		// En modo desarrollo sin DB, permitir login con credenciales de prueba
		if request.Email == "admin@growdesk.com" && request.Password == "admin123" {
			mockUser := models.User{
				ID:        "admin-user-123",
				Email:     request.Email,
				FirstName: "Admin",
				LastName:  "User",
				Role:      "admin",
			}

			token, err := auth.GenerateToken(mockUser.ID, mockUser.Email, mockUser.Role)
			if err != nil {
				log.Printf("Login error: fallo al generar token para %s: %v", request.Email, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
				return
			}

			c.JSON(http.StatusOK, AuthResponse{
				Token: token,
				User:  mockUser,
			})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	var user models.User
	result := db.Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		log.Printf("Login error: usuario no encontrado: %s", request.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	if !user.CheckPassword(request.Password) {
		log.Printf("Login error: contraseña incorrecta para: %s", request.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		log.Printf("Login error: fallo al generar token para %s: %v", request.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	log.Printf("Login exitoso para: %s (ID: %s, Role: %s)", user.Email, user.ID, user.Role)

	// No devolver la contraseña en la respuesta
	user.Password = ""

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

// Register maneja el registro de nuevos usuarios
func Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Register error: invalid request format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de registro inválidos", "details": err.Error()})
		return
	}

	log.Printf("Intento de registro para: %s", request.Email)

	db := database.GetDB()
	if db == nil {
		log.Printf("Register: database connection not available, using mock registration for: %s", request.Email)

		// En modo desarrollo sin DB, simular registro exitoso
		mockUser := models.User{
			ID:        "user-" + time.Now().Format("20060102150405"),
			Email:     request.Email,
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Role:      "customer",
		}

		token, err := auth.GenerateToken(mockUser.ID, mockUser.Email, mockUser.Role)
		if err != nil {
			log.Printf("Register error: fallo al generar token para %s: %v", request.Email, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
			return
		}

		c.JSON(http.StatusCreated, AuthResponse{
			Token: token,
			User:  mockUser,
		})
		return
	}

	var existingUser models.User
	result := db.Where("email = ?", request.Email).First(&existingUser)
	if result.Error == nil {
		log.Printf("Register error: email ya registrado: %s", request.Email)
		c.JSON(http.StatusConflict, gin.H{"error": "El correo electrónico ya está registrado"})
		return
	}

	user := models.User{
		Email:     request.Email,
		Password:  request.Password,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Role:      "customer",
	}

	tx := db.Begin()

	result = tx.Create(&user)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("Register error: fallo al crear usuario %s: %v", request.Email, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		tx.Rollback()
		log.Printf("Register error: fallo al generar token para %s: %v", request.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token de autenticación"})
		return
	}

	tx.Commit()

	log.Printf("Registro exitoso para: %s (ID: %s)", user.Email, user.ID)

	// No devolver la contraseña en la respuesta
	user.Password = ""

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}

// GetMe obtiene información del usuario autenticado actual
func GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
		return
	}

	db := database.GetDB()
	if db == nil {
		// En modo desarrollo sin DB, simular un usuario
		c.JSON(http.StatusOK, gin.H{
			"id":        userID,
			"email":     c.GetString("email"),
			"role":      c.GetString("role"),
			"firstName": "Usuario",
			"lastName":  "Desarrollo",
		})
		return
	}

	var user models.User
	result := db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		log.Printf("GetMe error: usuario no encontrado: %s", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// No devolver la contraseña en la respuesta
	user.Password = ""

	c.JSON(http.StatusOK, user)
}
