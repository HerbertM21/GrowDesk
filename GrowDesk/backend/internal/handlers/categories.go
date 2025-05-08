package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/data"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/middleware"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
)

// CategoryHandler contiene manejadores para operaciones de categoría
type CategoryHandler struct {
	Store *data.Store
}

// GetAllCategories devuelve todas las categorías
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Devolver todas las categorías
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.Store.Categories)
}

// GetCategory returns a specific category by ID
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer el ID de la categoría desde la URL
	// Formato de URL: /api/categories/:id
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	categoryID := parts[len(parts)-1]

	// Encontrar la categoría
	var category *models.Category
	for i, c := range h.Store.Categories {
		if c.ID == categoryID {
			category = &h.Store.Categories[i]
			break
		}
	}

	if category == nil {
		http.Error(w, "Categoría no encontrada", http.StatusNotFound)
		return
	}

	// Devolver la categoría
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// CreateCategory crea una nueva categoría
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Comprobar si el usuario es administrador
	role, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Prohibido: Solo los administradores pueden crear categorías", http.StatusForbidden)
		return
	}

	// Parsear el cuerpo de la solicitud
	var categoryReq struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
		Icon        string `json:"icon"`
		Active      bool   `json:"active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&categoryReq); err != nil {
		http.Error(w, "El cuerpo de la solicitud es inválido", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if categoryReq.Name == "" {
		http.Error(w, "El nombre de la categoría es requerido", http.StatusBadRequest)
		return
	}

	// Establecer valores predeterminados si no se proporcionan
	if categoryReq.Color == "" {
		categoryReq.Color = "#2196F3" // Valor predeterminado azul
	}
	if categoryReq.Icon == "" {
		categoryReq.Icon = "category" // Icono predeterminado
	}

	// Generar nuevo ID (ID incremental simple por ahora)
	newID := "1"
	if len(h.Store.Categories) > 0 {
		lastID, _ := strconv.Atoi(h.Store.Categories[len(h.Store.Categories)-1].ID)
		newID = strconv.Itoa(lastID + 1)
	}

	// Crear nueva categoría
	now := time.Now()
	newCategory := models.Category{
		ID:          newID,
		Name:        categoryReq.Name,
		Description: categoryReq.Description,
		Color:       categoryReq.Color,
		Icon:        categoryReq.Icon,
		Active:      categoryReq.Active,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Agregar al almacén y guardar
	h.Store.Categories = append(h.Store.Categories, newCategory)
	h.Store.SaveCategories()

	// Devolver la categoría creada
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

// UpdateCategory actualiza una categoría existente
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Comprobar si el usuario es administrador
	role, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Prohibido: Solo los administradores pueden actualizar categorías", http.StatusForbidden)
		return
	}

	// Extraer el ID de la categoría desde la URL
	// Formato de URL: /api/categories/:id
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	categoryID := parts[len(parts)-1]

	// Parsear el cuerpo de la solicitud
	var updateReq struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
		Icon        string `json:"icon"`
		Active      *bool  `json:"active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "El cuerpo de la solicitud es inválido", http.StatusBadRequest)
		return
	}

	// Encontrar la categoría
	var categoryIndex = -1
	for i, c := range h.Store.Categories {
		if c.ID == categoryID {
			categoryIndex = i
			break
		}
	}

	if categoryIndex == -1 {
		http.Error(w, "Categoría no encontrada", http.StatusNotFound)
		return
	}

	// Actualizar campos si se proporcionan
	if updateReq.Name != "" {
		h.Store.Categories[categoryIndex].Name = updateReq.Name
	}
	if updateReq.Description != "" {
		h.Store.Categories[categoryIndex].Description = updateReq.Description
	}
	if updateReq.Color != "" {
		h.Store.Categories[categoryIndex].Color = updateReq.Color
	}
	if updateReq.Icon != "" {
		h.Store.Categories[categoryIndex].Icon = updateReq.Icon
	}
	if updateReq.Active != nil {
		h.Store.Categories[categoryIndex].Active = *updateReq.Active
	}

	// Actualizar marca de tiempo
	h.Store.Categories[categoryIndex].UpdatedAt = time.Now()

	// Guardar cambios
	h.Store.SaveCategories()

	// Devolver la categoría actualizada
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.Store.Categories[categoryIndex])
}

// DeleteCategory elimina una categoría existente
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Comprobar si el usuario es administrador
	role, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Prohibido: Solo los administradores pueden eliminar categorías", http.StatusForbidden)
		return
	}

	// Extraer el ID de la categoría desde la URL
	// Formato de URL: /api/categories/:id
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	categoryID := parts[len(parts)-1]

	// Comprobar si alguna incidencia está usando esta categoría
	for _, ticket := range h.Store.Tickets {
		if ticket.Category == categoryID {
			http.Error(w, "La categoría está en uso por una o más incidencias", http.StatusBadRequest)
			return
		}
	}

	// Encontrar y eliminar la categoría
	var found bool
	var updatedCategories []models.Category
	for _, c := range h.Store.Categories {
		if c.ID != categoryID {
			updatedCategories = append(updatedCategories, c)
		} else {
			found = true
		}
	}

	if !found {
		http.Error(w, "Categoría no encontrada", http.StatusNotFound)
		return
	}

	// Actualizar lista de categorías
	h.Store.Categories = updatedCategories
	h.Store.SaveCategories()

	// Devolver éxito
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Categoría eliminada correctamente",
	})
}
