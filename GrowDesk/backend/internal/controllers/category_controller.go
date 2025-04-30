package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/models"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
)

// CategoryRequest estructura para solicitudes de creación/actualización
type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
	Active      *bool  `json:"active"`
}

// GetAllCategories obtiene todas las categorías
func GetAllCategories(c *gin.Context) {
	// Parámetro opcional para filtrar solo categorías activas
	onlyActive := c.Query("active") == "true"

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, devolver categorías de ejemplo
		mockCategories := []models.Category{
			{
				ID:          "CAT-001",
				Name:        "Soporte Técnico",
				Description: "Problemas técnicos y errores",
				Color:       "#2196F3",
				Icon:        "computer",
				Active:      true,
				CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt:   time.Now().Add(-5 * 24 * time.Hour),
			},
			{
				ID:          "CAT-002",
				Name:        "Facturación",
				Description: "Consultas sobre pagos y facturas",
				Color:       "#4CAF50",
				Icon:        "payments",
				Active:      true,
				CreatedAt:   time.Now().Add(-25 * 24 * time.Hour),
				UpdatedAt:   time.Now().Add(-2 * 24 * time.Hour),
			},
			{
				ID:          "CAT-003",
				Name:        "Ventas",
				Description: "Consultas sobre productos y servicios",
				Color:       "#FFC107",
				Icon:        "shopping_cart",
				Active:      true,
				CreatedAt:   time.Now().Add(-20 * 24 * time.Hour),
				UpdatedAt:   time.Now().Add(-1 * 24 * time.Hour),
			},
			{
				ID:          "CAT-004",
				Name:        "Características",
				Description: "Solicitudes de nuevas características",
				Color:       "#9C27B0",
				Icon:        "new_releases",
				Active:      false,
				CreatedAt:   time.Now().Add(-15 * 24 * time.Hour),
				UpdatedAt:   time.Now().Add(-15 * 24 * time.Hour),
			},
		}

		// Filtrar si es necesario
		if onlyActive {
			var activeCategories []models.Category
			for _, cat := range mockCategories {
				if cat.Active {
					activeCategories = append(activeCategories, cat)
				}
			}
			c.JSON(http.StatusOK, activeCategories)
			return
		}

		c.JSON(http.StatusOK, mockCategories)
		return
	}

	var categories []models.Category
	query := db.Order("name asc")

	// Filtrar sólo categorías activas si se solicita
	if onlyActive {
		query = query.Where("active = ?", true)
	}

	if err := query.Find(&categories).Error; err != nil {
		log.Printf("GetAllCategories error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener categorías"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategory obtiene una categoría específica
func GetCategory(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, devolver categoría de ejemplo
		if id == "CAT-001" {
			mockCategory := models.Category{
				ID:          "CAT-001",
				Name:        "Soporte Técnico",
				Description: "Problemas técnicos y errores",
				Color:       "#2196F3",
				Icon:        "computer",
				Active:      true,
				CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt:   time.Now().Add(-5 * 24 * time.Hour),
			}

			c.JSON(http.StatusOK, mockCategory)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
		return
	}

	var category models.Category
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// CreateCategory crea una nueva categoría
func CreateCategory(c *gin.Context) {
	var request CategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	active := true
	if request.Active != nil {
		active = *request.Active
	}

	// Crear nueva categoría
	category := models.Category{
		Name:        request.Name,
		Description: request.Description,
		Color:       request.Color,
		Icon:        request.Icon,
		Active:      active,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Valores predeterminados
	if category.Color == "" {
		category.Color = "#2196F3"
	}

	if category.Icon == "" {
		category.Icon = "category"
	}

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular creación
		category.ID = "CAT-" + time.Now().Format("20060102150405")

		c.JSON(http.StatusCreated, category)
		return
	}

	if err := db.Create(&category).Error; err != nil {
		log.Printf("CreateCategory error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear categoría"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// UpdateCategory actualiza una categoría existente
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var request CategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular actualización
		active := true
		if request.Active != nil {
			active = *request.Active
		}

		mockCategory := models.Category{
			ID:          id,
			Name:        request.Name,
			Description: request.Description,
			Color:       request.Color,
			Icon:        request.Icon,
			Active:      active,
			UpdatedAt:   time.Now(),
		}

		c.JSON(http.StatusOK, mockCategory)
		return
	}

	var category models.Category
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
		return
	}

	// Actualizar campos
	category.Name = request.Name
	category.Description = request.Description

	if request.Color != "" {
		category.Color = request.Color
	}

	if request.Icon != "" {
		category.Icon = request.Icon
	}

	if request.Active != nil {
		category.Active = *request.Active
	}

	category.UpdatedAt = time.Now()

	if err := db.Save(&category).Error; err != nil {
		log.Printf("UpdateCategory error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar categoría"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory elimina una categoría
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	if db == nil {
		// En modo de desarrollo sin DB, simular eliminación
		c.JSON(http.StatusOK, gin.H{
			"message": "Categoría eliminada correctamente (simulado)",
			"id":      id,
		})
		return
	}

	// Verificar uso en tickets
	var count int64
	db.Model(&models.Ticket{}).Where("category = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "No se puede eliminar la categoría porque está siendo utilizada en tickets",
			"count": count,
		})
		return
	}

	result := db.Delete(&models.Category{}, "id = ?", id)
	if result.Error != nil {
		log.Printf("DeleteCategory error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar categoría"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categoría no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Categoría eliminada correctamente"})
}
