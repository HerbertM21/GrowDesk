package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/models"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
)

// FAQRequest estructura para solicitudes de creación/actualización
type FAQRequest struct {
	Question    string `json:"question" binding:"required"`
	Answer      string `json:"answer" binding:"required"`
	Category    string `json:"category" binding:"required"`
	IsPublished *bool  `json:"isPublished"`
}

// GetAllFaqs obtiene todas las preguntas frecuentes
func GetAllFaqs(c *gin.Context) {
	// Parámetros opcionales de filtro
	category := c.Query("category")
	onlyPublished := c.Query("published") == "true"

	// Obtener rol del usuario para determinar si mostrar no publicadas
	_, exists := c.Get("role")
	isAgent := false
	if exists {
		role := c.GetString("role")
		isAgent = role == "agent" || role == "admin"
	}

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, devolver FAQs de ejemplo
		mockFaqs := []models.FAQ{
			{
				ID:          1,
				Question:    "¿Cómo puedo cambiar mi contraseña?",
				Answer:      "Puede cambiar su contraseña accediendo a su perfil y seleccionando la opción 'Cambiar contraseña'.",
				Category:    "Cuenta",
				IsPublished: true,
				CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt:   time.Now().Add(-5 * 24 * time.Hour),
			},
			{
				ID:          2,
				Question:    "¿Cómo puedo crear un nuevo ticket?",
				Answer:      "Para crear un nuevo ticket, vaya al panel principal y haga clic en el botón 'Nuevo Ticket'.",
				Category:    "Tickets",
				IsPublished: true,
				CreatedAt:   time.Now().Add(-20 * 24 * time.Hour),
				UpdatedAt:   time.Now().Add(-2 * 24 * time.Hour),
			},
			{
				ID:          3,
				Question:    "¿Qué información necesito proporcionar para un ticket?",
				Answer:      "Es importante incluir un título descriptivo, detallar el problema y adjuntar capturas de pantalla si es posible.",
				Category:    "Tickets",
				IsPublished: false, // Esta no está publicada
				CreatedAt:   time.Now().Add(-10 * 24 * time.Hour),
				UpdatedAt:   time.Now().Add(-1 * 24 * time.Hour),
			},
		}

		// Aplicar filtros
		var filtered []models.FAQ
		for _, faq := range mockFaqs {
			include := true

			// Filtrar por categoría
			if category != "" && faq.Category != category {
				include = false
			}

			// Filtrar por publicación
			if onlyPublished && !faq.IsPublished {
				include = false
			}

			// Si no es agente y la FAQ no está publicada, no mostrar
			if !isAgent && !faq.IsPublished {
				include = false
			}

			if include {
				filtered = append(filtered, faq)
			}
		}

		c.JSON(http.StatusOK, filtered)
		return
	}

	// Construir consulta con filtros
	query := db.Model(&models.FAQ{})

	// Filtrar por categoría si se especifica
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Si solo se quieren publicadas o el usuario no es agente/admin
	if onlyPublished || !isAgent {
		query = query.Where("is_published = ?", true)
	}

	var faqs []models.FAQ
	if err := query.Order("category, id").Find(&faqs).Error; err != nil {
		log.Printf("GetAllFaqs error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener preguntas frecuentes"})
		return
	}

	c.JSON(http.StatusOK, faqs)
}

// GetFaq obtiene una pregunta frecuente específica
func GetFaq(c *gin.Context) {
	idParam := c.Param("id")

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, devolver FAQ de ejemplo
		if idParam == "1" {
			mockFaq := models.FAQ{
				ID:          1,
				Question:    "¿Cómo puedo cambiar mi contraseña?",
				Answer:      "Puede cambiar su contraseña accediendo a su perfil y seleccionando la opción 'Cambiar contraseña'.",
				Category:    "Cuenta",
				IsPublished: true,
				CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt:   time.Now().Add(-5 * 24 * time.Hour),
			}

			c.JSON(http.StatusOK, mockFaq)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Pregunta frecuente no encontrada"})
		return
	}

	var faq models.FAQ
	if err := db.Where("id = ?", idParam).First(&faq).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pregunta frecuente no encontrada"})
		return
	}

	// Verificar si es público o si el usuario tiene permisos
	role, exists := c.Get("role")
	if !faq.IsPublished && (!exists || (role != "admin" && role != "agent")) {
		c.JSON(http.StatusForbidden, gin.H{"error": "No tiene permiso para ver esta pregunta frecuente"})
		return
	}

	c.JSON(http.StatusOK, faq)
}

// CreateFaq crea una nueva pregunta frecuente
func CreateFaq(c *gin.Context) {
	var request FAQRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	isPublished := false
	if request.IsPublished != nil {
		isPublished = *request.IsPublished
	}

	// Crear nueva FAQ
	faq := models.FAQ{
		Question:    request.Question,
		Answer:      request.Answer,
		Category:    request.Category,
		IsPublished: isPublished,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular creación
		faq.ID = uint(time.Now().Unix() % 1000)

		c.JSON(http.StatusCreated, faq)
		return
	}

	if err := db.Create(&faq).Error; err != nil {
		log.Printf("CreateFaq error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear pregunta frecuente"})
		return
	}

	c.JSON(http.StatusCreated, faq)
}

// UpdateFaq actualiza una pregunta frecuente existente
func UpdateFaq(c *gin.Context) {
	idParam := c.Param("id")

	var request FAQRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular actualización
		isPublished := false
		if request.IsPublished != nil {
			isPublished = *request.IsPublished
		}

		id := 0
		if idParam == "1" {
			id = 1
		}

		mockFaq := models.FAQ{
			ID:          uint(id),
			Question:    request.Question,
			Answer:      request.Answer,
			Category:    request.Category,
			IsPublished: isPublished,
			UpdatedAt:   time.Now(),
		}

		c.JSON(http.StatusOK, mockFaq)
		return
	}

	var faq models.FAQ
	if err := db.Where("id = ?", idParam).First(&faq).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pregunta frecuente no encontrada"})
		return
	}

	// Actualizar campos
	faq.Question = request.Question
	faq.Answer = request.Answer
	faq.Category = request.Category

	if request.IsPublished != nil {
		faq.IsPublished = *request.IsPublished
	}

	faq.UpdatedAt = time.Now()

	if err := db.Save(&faq).Error; err != nil {
		log.Printf("UpdateFaq error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar pregunta frecuente"})
		return
	}

	c.JSON(http.StatusOK, faq)
}

// ToggleFaqPublish cambia el estado de publicación de una FAQ
func ToggleFaqPublish(c *gin.Context) {
	idParam := c.Param("id")

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular cambio
		c.JSON(http.StatusOK, gin.H{
			"message": "Estado de publicación cambiado correctamente (simulado)",
			"id":      idParam,
		})
		return
	}

	var faq models.FAQ
	if err := db.Where("id = ?", idParam).First(&faq).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pregunta frecuente no encontrada"})
		return
	}

	// Invertir estado de publicación
	faq.IsPublished = !faq.IsPublished
	faq.UpdatedAt = time.Now()

	if err := db.Save(&faq).Error; err != nil {
		log.Printf("ToggleFaqPublish error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cambiar estado de publicación"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Estado de publicación cambiado correctamente",
		"id":          faq.ID,
		"isPublished": faq.IsPublished,
	})
}

// DeleteFaq elimina una pregunta frecuente
func DeleteFaq(c *gin.Context) {
	idParam := c.Param("id")

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular eliminación
		c.JSON(http.StatusOK, gin.H{
			"message": "Pregunta frecuente eliminada correctamente (simulado)",
			"id":      idParam,
		})
		return
	}

	result := db.Delete(&models.FAQ{}, "id = ?", idParam)
	if result.Error != nil {
		log.Printf("DeleteFaq error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar pregunta frecuente"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pregunta frecuente no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pregunta frecuente eliminada correctamente"})
}
