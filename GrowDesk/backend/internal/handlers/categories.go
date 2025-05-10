package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/data"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/middleware"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/utils"
)

// CategoryHandler contiene manejadores para categorías
type CategoryHandler struct {
	Store data.DataStore
}

// GetAllCategories devuelve todas las categorías
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	// Esta función solo maneja solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Establecer CORS
	utils.SetCORS(w)

	// Obtener categorías del almacén
	categories, err := h.Store.GetCategories()
	if err != nil {
		http.Error(w, "Error al obtener categorías", http.StatusInternalServerError)
		return
	}

	// Devolver categorías como JSON
	utils.WriteJSON(w, http.StatusOK, categories)
}

// GetCategory obtiene una categoría por ID
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	// Esta función solo maneja solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Establecer CORS
	utils.SetCORS(w)

	// Obtener ID de la URL
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 4 {
		http.Error(w, "URL de categoría inválida", http.StatusBadRequest)
		return
	}

	categoryID := segments[3]

	// Obtener categoría del almacén
	category, err := h.Store.GetCategory(categoryID)
	if err != nil {
		http.Error(w, "Categoría no encontrada", http.StatusNotFound)
		return
	}

	// Devolver categoría como JSON
	utils.WriteJSON(w, http.StatusOK, category)
}

// CreateCategory crea una nueva categoría
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Esta función solo maneja solicitudes POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Establecer CORS
	utils.SetCORS(w)

	// Verificar permisos de administrador
	role, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	// Decodificar cuerpo de la solicitud
	var category models.Category
	if err := utils.DecodeJSON(r, &category); err != nil {
		http.Error(w, "Error al leer datos de la categoría", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if category.Name == "" {
		http.Error(w, "El nombre de la categoría es requerido", http.StatusBadRequest)
		return
	}

	// Establecer valores por defecto si no se proporcionan
	if category.Color == "" {
		category.Color = "#3498db" // Color azul por defecto
	}
	if category.Icon == "" {
		category.Icon = "category" // Icono por defecto
	}

	// Establecer ID y marcas de tiempo
	category.ID = uuid.New().String()
	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now
	category.Active = true

	// Guardar en el almacén
	if err := h.Store.CreateCategory(category); err != nil {
		http.Error(w, "Error al crear categoría", http.StatusInternalServerError)
		return
	}

	// Devolver categoría creada
	utils.WriteJSON(w, http.StatusCreated, category)
}

// UpdateCategory actualiza una categoría existente
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	// Esta función solo maneja solicitudes PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Establecer CORS
	utils.SetCORS(w)

	// Verificar permisos de administrador
	role, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	// Obtener ID de la URL
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 4 {
		http.Error(w, "URL de categoría inválida", http.StatusBadRequest)
		return
	}

	categoryID := segments[3]

	// Obtener categoría existente
	existingCategory, err := h.Store.GetCategory(categoryID)
	if err != nil {
		http.Error(w, "Categoría no encontrada", http.StatusNotFound)
		return
	}

	// Decodificar cuerpo de la solicitud
	var updates models.Category
	if err := utils.DecodeJSON(r, &updates); err != nil {
		http.Error(w, "Error al leer datos de actualización", http.StatusBadRequest)
		return
	}

	// Actualizar campos de la categoría existente
	if updates.Name != "" {
		existingCategory.Name = updates.Name
	}
	if updates.Description != "" {
		existingCategory.Description = updates.Description
	}
	if updates.Color != "" {
		existingCategory.Color = updates.Color
	}
	if updates.Icon != "" {
		existingCategory.Icon = updates.Icon
	}

	// Actualizar campo active explícitamente si se proporciona
	existingCategory.Active = updates.Active

	// Actualizar marca de tiempo
	existingCategory.UpdatedAt = time.Now()

	// Guardar en el almacén
	if err := h.Store.UpdateCategory(*existingCategory); err != nil {
		http.Error(w, "Error al actualizar categoría", http.StatusInternalServerError)
		return
	}

	// Devolver categoría actualizada
	utils.WriteJSON(w, http.StatusOK, existingCategory)
}

// DeleteCategory elimina una categoría existente
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Esta función solo maneja solicitudes DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Establecer CORS
	utils.SetCORS(w)

	// Verificar permisos de administrador
	role, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	// Obtener ID de la URL
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 4 {
		http.Error(w, "URL de categoría inválida", http.StatusBadRequest)
		return
	}

	categoryID := segments[3]

	// Eliminar la categoría
	if err := h.Store.DeleteCategory(categoryID); err != nil {
		http.Error(w, "Error al eliminar categoría", http.StatusInternalServerError)
		return
	}

	// Devolver respuesta exitosa
	utils.WriteJSON(w, http.StatusOK, map[string]bool{"success": true})
}
