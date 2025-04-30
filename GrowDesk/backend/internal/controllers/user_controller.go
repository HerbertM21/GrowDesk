package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hmdev/GrowDesk/backend/models"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
)

// GetAllUsers obtiene lista de todos los usuarios (para administradores)
func GetAllUsers(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		// En modo desarrollo sin DB, devolver usuarios de ejemplo
		mockUsers := []models.User{
			{
				ID:         "admin-123",
				Email:      "admin@growdesk.com",
				FirstName:  "Admin",
				LastName:   "Usuario",
				Role:       "admin",
				Department: "Tecnología",
				Active:     true,
				CreatedAt:  time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt:  time.Now().Add(-2 * 24 * time.Hour),
			},
			{
				ID:         "agent-456",
				Email:      "agente@growdesk.com",
				FirstName:  "Agente",
				LastName:   "Soporte",
				Role:       "agent",
				Department: "Soporte",
				Active:     true,
				CreatedAt:  time.Now().Add(-20 * 24 * time.Hour),
				UpdatedAt:  time.Now().Add(-1 * 24 * time.Hour),
			},
			{
				ID:        "customer-789",
				Email:     "cliente@ejemplo.com",
				FirstName: "Cliente",
				LastName:  "Ejemplo",
				Role:      "customer",
				Active:    true,
				CreatedAt: time.Now().Add(-10 * 24 * time.Hour),
				UpdatedAt: time.Now(),
			},
		}

		c.JSON(http.StatusOK, mockUsers)
		return
	}

	var users []models.User
	result := db.Find(&users)
	if result.Error != nil {
		log.Printf("GetAllUsers error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
		return
	}

	// Eliminar contraseñas antes de devolver
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}

// GetUser obtiene los detalles de un usuario específico
func GetUser(c *gin.Context) {
	id := c.Param("id")

	// Verificar si el usuario actual tiene permiso para ver este usuario
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
		return
	}

	currentRole, _ := c.Get("role")
	isAdmin := currentRole == "admin"
	isSameUser := currentUserID.(string) == id

	// Solo el propio usuario o un admin puede ver los detalles
	if !isAdmin && !isSameUser {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permiso para ver este usuario"})
		return
	}

	db := database.GetDB()
	if db == nil {
		// En modo desarrollo sin DB, devolver usuario de ejemplo
		if id == "admin-123" {
			c.JSON(http.StatusOK, models.User{
				ID:         "admin-123",
				Email:      "admin@growdesk.com",
				FirstName:  "Admin",
				LastName:   "Usuario",
				Role:       "admin",
				Department: "Tecnología",
				Active:     true,
				CreatedAt:  time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt:  time.Now().Add(-2 * 24 * time.Hour),
			})
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	var user models.User
	result := db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// No devolver la contraseña
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// CreateUser crea un nuevo usuario (solo para admins)
func CreateUser(c *gin.Context) {
	var request struct {
		Email      string `json:"email" binding:"required,email"`
		Password   string `json:"password" binding:"required,min=6"`
		FirstName  string `json:"firstName" binding:"required"`
		LastName   string `json:"lastName" binding:"required"`
		Role       string `json:"role" binding:"required"`
		Department string `json:"department"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	// Verificar rol válido
	validRoles := map[string]bool{
		"admin":    true,
		"agent":    true,
		"customer": true,
	}

	if !validRoles[request.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rol inválido"})
		return
	}

	db := database.GetDB()
	if db == nil {
		// En modo desarrollo sin DB, simular creación
		mockUser := models.User{
			ID:         uuid.New().String(),
			Email:      request.Email,
			FirstName:  request.FirstName,
			LastName:   request.LastName,
			Role:       request.Role,
			Department: request.Department,
			Active:     true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		c.JSON(http.StatusCreated, mockUser)
		return
	}

	// Verificar si el email ya existe
	var existingUser models.User
	if db.Where("email = ?", request.Email).First(&existingUser).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "El correo electrónico ya está registrado"})
		return
	}

	// Crear nuevo usuario
	user := models.User{
		Email:      request.Email,
		Password:   request.Password,
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		Role:       request.Role,
		Department: request.Department,
		Active:     true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if result := db.Create(&user); result.Error != nil {
		log.Printf("CreateUser error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario"})
		return
	}

	// No devolver la contraseña
	user.Password = ""

	c.JSON(http.StatusCreated, user)
}

// UpdateUser actualiza la información de un usuario
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	// Verificar permisos
	currentUserID, _ := c.Get("userID")
	currentRole, _ := c.Get("role")
	isAdmin := currentRole == "admin"
	isSameUser := currentUserID.(string) == id

	// Solo el propio usuario o un admin puede actualizar
	if !isAdmin && !isSameUser {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permiso para modificar este usuario"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// No permitir actualizar campos sensibles excepto para admins
	delete(updateData, "id")
	delete(updateData, "password") // Contraseña se actualiza en otro endpoint
	delete(updateData, "createdAt")

	if !isAdmin {
		delete(updateData, "role")
		delete(updateData, "active")
	}

	db := database.GetDB()
	if db == nil {
		// En modo desarrollo sin DB, simular actualización
		c.JSON(http.StatusOK, gin.H{
			"id":      id,
			"message": "Usuario actualizado correctamente (simulado)",
			"updated": updateData,
		})
		return
	}

	var user models.User
	if result := db.Where("id = ?", id).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Actualizar usuario
	updateData["updated_at"] = time.Now()
	if result := db.Model(&user).Updates(updateData); result.Error != nil {
		log.Printf("UpdateUser error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar usuario"})
		return
	}

	// Obtener usuario actualizado
	if result := db.Where("id = ?", id).First(&user); result.Error != nil {
		log.Printf("UpdateUser error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuario actualizado"})
		return
	}

	// No devolver la contraseña
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// UpdateUserRole actualiza el rol de un usuario (solo admin)
func UpdateUserRole(c *gin.Context) {
	id := c.Param("id")

	var request struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rol no especificado"})
		return
	}

	// Verificar rol válido
	validRoles := map[string]bool{
		"admin":    true,
		"agent":    true,
		"customer": true,
	}

	if !validRoles[request.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rol inválido"})
		return
	}

	db := database.GetDB()
	if db == nil {
		// En modo desarrollo sin DB, simular actualización
		c.JSON(http.StatusOK, gin.H{
			"id":      id,
			"message": "Rol actualizado correctamente (simulado)",
			"role":    request.Role,
		})
		return
	}

	var user models.User
	if result := db.Where("id = ?", id).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Actualizar rol
	updates := map[string]interface{}{
		"role":       request.Role,
		"updated_at": time.Now(),
	}

	if result := db.Model(&user).Updates(updates); result.Error != nil {
		log.Printf("UpdateUserRole error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar rol"})
		return
	}

	// Obtener usuario actualizado
	if result := db.Where("id = ?", id).First(&user); result.Error != nil {
		log.Printf("UpdateUserRole error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuario actualizado"})
		return
	}

	// No devolver la contraseña
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// DeleteUser elimina un usuario (solo admin)
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// No permitir eliminar al propio usuario
	currentUserID, _ := c.Get("userID")
	if currentUserID.(string) == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No puede eliminar su propia cuenta"})
		return
	}

	db := database.GetDB()
	if db == nil {
		// En modo desarrollo sin DB, simular eliminación
		c.JSON(http.StatusOK, gin.H{
			"message": "Usuario eliminado correctamente (simulado)",
			"id":      id,
		})
		return
	}

	// Eliminar usuario
	result := db.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		log.Printf("DeleteUser error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar usuario"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado correctamente"})
}
