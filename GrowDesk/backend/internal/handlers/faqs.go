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
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/utils"
)

// FAQHandler contiene manejadores para FAQs
type FAQHandler struct {
	Store data.DataStore
}

// GetAllFAQs devuelve todas las FAQs
func (h *FAQHandler) GetAllFAQs(w http.ResponseWriter, r *http.Request) {
	// Esta función solo maneja solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Establecer CORS
	utils.SetCORS(w)

	// Determinar si se debe filtrar por estado de publicación
	published := r.URL.Query().Get("published")

	var faqs []models.FAQ
	var err error

	if published != "" {
		// Convertir string a bool
		isPublished := published == "true"
		faqs, err = h.Store.GetFAQsByStatus(isPublished)
	} else {
		// Obtener todas las FAQs
		faqs, err = h.Store.GetFAQs()
	}

	if err != nil {
		http.Error(w, "Error al obtener FAQs", http.StatusInternalServerError)
		return
	}

	// Devolver FAQs como JSON
	utils.WriteJSON(w, http.StatusOK, faqs)
}

// GetPublishedFAQs devuelve solo las FAQs publicadas
func (h *FAQHandler) GetPublishedFAQs(w http.ResponseWriter, r *http.Request) {
	// Esta función solo maneja solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Establecer CORS
	utils.SetCORS(w)

	// Obtener FAQs publicadas
	faqs, err := h.Store.GetFAQsByStatus(true)
	if err != nil {
		http.Error(w, "Error al obtener FAQs", http.StatusInternalServerError)
		return
	}

	// Devolver FAQs como JSON
	utils.WriteJSON(w, http.StatusOK, faqs)
}

// GetFAQ devuelve una FAQ específica por ID
func (h *FAQHandler) GetFAQ(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer el ID de la FAQ desde la URL
	// Formato de URL: /api/faqs/:id
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "ID de FAQ inválido", http.StatusBadRequest)
		return
	}

	// Parsear el ID como entero
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Formato de ID de FAQ inválido", http.StatusBadRequest)
		return
	}

	// Obtener FAQ por ID
	faq, err := h.Store.GetFAQ(id)
	if err != nil {
		http.Error(w, "FAQ no encontrada", http.StatusNotFound)
		return
	}

	// Devolver la FAQ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faq)
}

// CreateFAQ crea una nueva FAQ
func (h *FAQHandler) CreateFAQ(w http.ResponseWriter, r *http.Request) {
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
	var faq models.FAQ
	if err := utils.DecodeJSON(r, &faq); err != nil {
		http.Error(w, "Error al leer datos de la FAQ", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if faq.Question == "" || faq.Answer == "" {
		http.Error(w, "La pregunta y la respuesta son requeridas", http.StatusBadRequest)
		return
	}

	// Establecer marcas de tiempo
	now := time.Now()
	faq.CreatedAt = now
	faq.UpdatedAt = now

	// Guardar en el almacén
	if err := h.Store.CreateFAQ(faq); err != nil {
		http.Error(w, "Error al crear FAQ", http.StatusInternalServerError)
		return
	}

	// Devolver FAQ creada
	utils.WriteJSON(w, http.StatusCreated, faq)
}

// UpdateFAQ actualiza una FAQ existente
func (h *FAQHandler) UpdateFAQ(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Comprobar si el usuario es administrador o asistente
	role, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || (role != "admin" && role != "assistant") {
		http.Error(w, "Prohibido: Permisos insuficientes", http.StatusForbidden)
		return
	}

	// Extraer el ID de la FAQ desde la URL
	// Formato de URL: /api/faqs/:id
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Formato de ID de FAQ inválido", http.StatusBadRequest)
		return
	}

	// Parsear el ID como entero
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Formato de ID de FAQ inválido", http.StatusBadRequest)
		return
	}

	// Parsear el cuerpo de la solicitud
	var updateReq struct {
		Question    string `json:"question"`
		Answer      string `json:"answer"`
		Category    string `json:"category"`
		IsPublished *bool  `json:"isPublished"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "El cuerpo de la solicitud es inválido", http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if updateReq.Question == "" || updateReq.Answer == "" || updateReq.Category == "" {
		http.Error(w, "La pregunta, la respuesta y la categoría son requeridas", http.StatusBadRequest)
		return
	}

	// Obtener la FAQ existente
	faq, err := h.Store.GetFAQ(id)
	if err != nil {
		http.Error(w, "FAQ no encontrada", http.StatusNotFound)
		return
	}

	// Actualizar campos
	faq.Question = updateReq.Question
	faq.Answer = updateReq.Answer
	faq.Category = updateReq.Category
	if updateReq.IsPublished != nil {
		faq.IsPublished = *updateReq.IsPublished
	}
	faq.UpdatedAt = time.Now()

	// Guardar cambios
	if err := h.Store.UpdateFAQ(*faq); err != nil {
		http.Error(w, "Error al actualizar FAQ", http.StatusInternalServerError)
		return
	}

	// Devolver la FAQ actualizada
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faq)
}

// DeleteFAQ elimina una FAQ existente
func (h *FAQHandler) DeleteFAQ(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Comprobar si el usuario es administrador
	role, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || role != "admin" {
		http.Error(w, "Prohibido: Solo los administradores pueden eliminar FAQs", http.StatusForbidden)
		return
	}

	// Extraer el ID de la FAQ desde la URL
	// Formato de URL: /api/faqs/:id
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Formato de ID de FAQ inválido", http.StatusBadRequest)
		return
	}

	// Parsear el ID como entero
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Formato de ID de FAQ inválido", http.StatusBadRequest)
		return
	}

	// Eliminar la FAQ
	if err := h.Store.DeleteFAQ(id); err != nil {
		http.Error(w, "Error al eliminar FAQ", http.StatusInternalServerError)
		return
	}

	// Devolver no contenido para eliminación exitosa
	w.WriteHeader(http.StatusNoContent)
}

// TogglePublishFAQ alterna el estado publicado de una FAQ
func (h *FAQHandler) TogglePublishFAQ(w http.ResponseWriter, r *http.Request) {
	// Solo maneja solicitudes PATCH
	if r.Method != http.MethodPatch {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Comprobar si el usuario es administrador o asistente
	role, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok || (role != "admin" && role != "assistant") {
		http.Error(w, "Prohibido: Permisos insuficientes", http.StatusForbidden)
		return
	}

	// Extraer el ID de la FAQ desde la URL
	// Formato de URL: /api/faqs/:id/toggle-publish
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Formato de ID de FAQ inválido", http.StatusBadRequest)
		return
	}

	// Obtener la parte del ID (esperando formato /faqs/ID/toggle-publish)
	idPart := parts[len(parts)-2]
	id, err := strconv.Atoi(idPart)
	if err != nil {
		http.Error(w, "Formato de ID de FAQ inválido", http.StatusBadRequest)
		return
	}

	// Llamar al método para alternar estado de publicación
	if err := h.Store.ToggleFAQPublish(id); err != nil {
		http.Error(w, "Error al cambiar estado de publicación", http.StatusInternalServerError)
		return
	}

	// Obtener la FAQ actualizada
	faq, err := h.Store.GetFAQ(id)
	if err != nil {
		http.Error(w, "Error al obtener FAQ actualizada", http.StatusInternalServerError)
		return
	}

	// Devolver la FAQ actualizada
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faq)
}
