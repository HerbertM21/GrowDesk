package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/data"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/middleware"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
)

// FAQHandler contiene manejadores para operaciones de FAQ
type FAQHandler struct {
	Store *data.Store
}

// GetAllFAQs devuelve todas las FAQs (ruta de administrador)
func (h *FAQHandler) GetAllFAQs(w http.ResponseWriter, r *http.Request) {
	// Establecer contexto de timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Solo maneja solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Crear canal para respuesta
	faqsChan := make(chan []models.FAQ, 1)
	errChan := make(chan error, 1)

	// Obtener FAQs en una goroutine
	go func() {
		faqs := h.Store.GetAllFAQs()
		select {
		case faqsChan <- faqs:
		case <-ctx.Done():
			errChan <- ctx.Err()
		}
	}()

	// Esperar respuesta o timeout
	select {
	case faqs := <-faqsChan:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(faqs)
	case err := <-errChan:
		if err == context.DeadlineExceeded {
			http.Error(w, "Tiempo de espera de solicitud agotado", http.StatusGatewayTimeout)
		} else {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
	case <-ctx.Done():
		http.Error(w, "Tiempo de espera de solicitud agotado", http.StatusGatewayTimeout)
	}
}

// GetPublishedFAQs devuelve solo las FAQs publicadas (ruta pública)
func (h *FAQHandler) GetPublishedFAQs(w http.ResponseWriter, r *http.Request) {
	// Establecer contexto de timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Solo maneja solicitudes GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Comprobar credenciales de widget
	widgetID := r.Header.Get("X-Widget-ID")
	widgetToken := r.Header.Get("X-Widget-Token")

	// Log credenciales de widget para depuración
	if widgetID != "" {
		fmt.Printf("Solicitud de widget - ID: %s, Token: %s\n", widgetID, widgetToken)
	}

	// Crear canal para respuesta
	faqsChan := make(chan []models.FAQ, 1)
	errChan := make(chan error, 1)

	// Obtener FAQs publicadas en una goroutine
	go func() {
		faqs := h.Store.GetPublishedFAQs()
		select {
		case faqsChan <- faqs:
		case <-ctx.Done():
			errChan <- ctx.Err()
		}
	}()

	// Esperar respuesta o timeout
	select {
	case faqs := <-faqsChan:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(faqs)
	case err := <-errChan:
		if err == context.DeadlineExceeded {
			http.Error(w, "Tiempo de espera de solicitud agotado", http.StatusGatewayTimeout)
		} else {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
	case <-ctx.Done():
		http.Error(w, "Tiempo de espera de solicitud agotado", http.StatusGatewayTimeout)
	}
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

	// Encontrar la FAQ
	var faq *models.FAQ
	for i, f := range h.Store.FAQs {
		if f.ID == id {
			faq = &h.Store.FAQs[i]
			break
		}
	}

	if faq == nil {
		http.Error(w, "FAQ no encontrada", http.StatusNotFound)
		return
	}

	// Devolver la FAQ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faq)
}

// CreateFAQ crea una nueva FAQ
func (h *FAQHandler) CreateFAQ(w http.ResponseWriter, r *http.Request) {
	// Establecer contexto de timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Solo maneja solicitudes POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear el cuerpo de la solicitud
	var faq models.FAQ
	if err := json.NewDecoder(r.Body).Decode(&faq); err != nil {
		http.Error(w, "El cuerpo de la solicitud es inválido", http.StatusBadRequest)
		return
	}

	// Crear canal para respuesta
	faqChan := make(chan *models.FAQ, 1)
	errChan := make(chan error, 1)

	// Crear FAQ en una goroutine
	go func() {
		newFAQ, err := h.Store.CreateFAQ(&faq)
		if err != nil {
			errChan <- err
			return
		}
		select {
		case faqChan <- newFAQ:
		case <-ctx.Done():
			errChan <- ctx.Err()
		}
	}()

	// Esperar respuesta o timeout
	select {
	case newFAQ := <-faqChan:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newFAQ)
	case err := <-errChan:
		if err == context.DeadlineExceeded {
			http.Error(w, "Tiempo de espera de solicitud agotado", http.StatusGatewayTimeout)
		} else {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
	case <-ctx.Done():
		http.Error(w, "Tiempo de espera de solicitud agotado", http.StatusGatewayTimeout)
	}
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

	// Encontrar la FAQ
	var faqIndex = -1
	for i, f := range h.Store.FAQs {
		if f.ID == id {
			faqIndex = i
			break
		}
	}

	if faqIndex == -1 {
		http.Error(w, "FAQ no encontrada", http.StatusNotFound)
		return
	}

	// Actualizar campos
	h.Store.FAQs[faqIndex].Question = updateReq.Question
	h.Store.FAQs[faqIndex].Answer = updateReq.Answer
	h.Store.FAQs[faqIndex].Category = updateReq.Category
	if updateReq.IsPublished != nil {
		h.Store.FAQs[faqIndex].IsPublished = *updateReq.IsPublished
	}
	h.Store.FAQs[faqIndex].UpdatedAt = time.Now()

	// Guardar cambios
	h.Store.SaveFAQs()

	// Devolver la FAQ actualizada
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.Store.FAQs[faqIndex])
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

	// Encontrar y eliminar la FAQ
	var found bool
	var updatedFAQs []models.FAQ
	for _, f := range h.Store.FAQs {
		if f.ID != id {
			updatedFAQs = append(updatedFAQs, f)
		} else {
			found = true
		}
	}

	if !found {
		http.Error(w, "FAQ no encontrada", http.StatusNotFound)
		return
	}

	// Actualizar lista de FAQs
	h.Store.FAQs = updatedFAQs
	h.Store.SaveFAQs()

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

	// Encontrar la FAQ
	var faqIndex = -1
	for i, f := range h.Store.FAQs {
		if f.ID == id {
			faqIndex = i
			break
		}
	}

	if faqIndex == -1 {
		http.Error(w, "FAQ no encontrada", http.StatusNotFound)
		return
	}

	// Alternar estado publicado
	h.Store.FAQs[faqIndex].IsPublished = !h.Store.FAQs[faqIndex].IsPublished
	h.Store.FAQs[faqIndex].UpdatedAt = time.Now()

	// Guardar cambios
	h.Store.SaveFAQs()

	// Devolver la FAQ actualizada
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.Store.FAQs[faqIndex])
}
