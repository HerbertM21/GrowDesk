package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/data"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/handlers"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/middleware"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/utils"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/websocket"
)

func main() {
	// Parsear flags
	var (
		port    = flag.Int("port", 8080, "HTTP service port")
		dataDir = flag.String("data-dir", "./data", "Directory for data storage")
		useMock = flag.Bool("mock-auth", true, "Use mock authentication for development")
	)
	flag.Parse()

	// Crear directorio de datos si no existe
	if err := os.MkdirAll(*dataDir, 0755); err != nil {
		log.Fatalf("Error al crear directorio de datos: %v", err)
	}

	// Inicializar el almacén de datos
	store := data.NewStore(*dataDir)

	// Crear handlers
	authHandler := &handlers.AuthHandler{Store: store}
	ticketHandler := &handlers.TicketHandler{Store: store}
	categoryHandler := &handlers.CategoryHandler{Store: store}
	faqHandler := &handlers.FAQHandler{Store: store}

	// Crear enrutador (usando http.ServeMux básico para simplicidad)
	mux := http.NewServeMux()

	// Middleware de autenticación
	var authMiddleware func(http.Handler) http.Handler
	if *useMock {
		authMiddleware = middleware.MockAuth
		log.Println("Usando autenticación de prueba para desarrollo")
	} else {
		authMiddleware = middleware.Auth
	}

	// Rutas de comprobación de estado
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "Servidor funcionando",
		})
	})

	// Rutas de autenticación
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.Handle("/api/auth/me", authMiddleware(http.HandlerFunc(authHandler.Me)))

	// Rutas de tickets (autenticadas)
	mux.Handle("/api/tickets", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Manejar basado en el método HTTP
		switch r.Method {
		case http.MethodGet:
			ticketHandler.GetAllTickets(w, r)
		case http.MethodPost:
			ticketHandler.CreateTicket(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})))

	// Rutas de tickets individuales
	mux.Handle("/api/tickets/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// Manejar la ruta de mensajes de tickets
		if filepath.Base(filepath.Dir(path)) == "tickets" && filepath.Ext(path) == "" {
			// Esta es una ruta para un ID de ticket específico como /api/tickets/:id
			switch r.Method {
			case http.MethodGet:
				ticketHandler.GetTicket(w, r)
			case http.MethodPut:
				ticketHandler.UpdateTicket(w, r)
			default:
				http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			}
		} else if filepath.Base(path) == "messages" {
			// Esta es una ruta para mensajes de tickets como /api/tickets/:id/messages
			switch r.Method {
			case http.MethodGet:
				ticketHandler.GetTicketMessages(w, r)
			case http.MethodPost:
				ticketHandler.AddTicketMessage(w, r)
			default:
				http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			}
		} else {
			http.NotFound(w, r)
		}
	})))

	// Rutas de widget (públicas)
	mux.HandleFunc("/widget/tickets", ticketHandler.CreateWidgetTicket)

	// Ruta para obtener mensajes de un ticket desde el widget
	mux.HandleFunc("/widget/tickets/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if filepath.Base(path) == "messages" {
			ticketHandler.GetTicketMessages(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	// Rutas de categorías (autenticadas)
	mux.Handle("/api/categories", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Manejar basado en el método HTTP
		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetAllCategories(w, r)
		case http.MethodPost:
			categoryHandler.CreateCategory(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})))

	// Rutas de categorías individuales
	mux.Handle("/api/categories/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Manejar basado en el método HTTP
		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetCategory(w, r)
		case http.MethodPut:
			categoryHandler.UpdateCategory(w, r)
		case http.MethodDelete:
			categoryHandler.DeleteCategory(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})))

	// Rutas de FAQ (autenticadas para operaciones de administrador)
	mux.Handle("/api/faqs", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Manejar basado en el método HTTP
		switch r.Method {
		case http.MethodGet:
			faqHandler.GetAllFAQs(w, r)
		case http.MethodPost:
			faqHandler.CreateFAQ(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})))

	// Rutas de FAQ individuales
	mux.Handle("/api/faqs/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// Comprobar si esta es una ruta para un endpoint de toggle-publish
		if filepath.Base(path) == "toggle-publish" {
			if r.Method == http.MethodPatch {
				faqHandler.TogglePublishFAQ(w, r)
			} else {
				http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			}
		} else {
			// Esta es una ruta para un ID de FAQ específico
			switch r.Method {
			case http.MethodGet:
				faqHandler.GetFAQ(w, r)
			case http.MethodPut:
				faqHandler.UpdateFAQ(w, r)
			case http.MethodDelete:
				faqHandler.DeleteFAQ(w, r)
			default:
				http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			}
		}
	})))

	// Rutas de FAQ públicas
	mux.HandleFunc("/widget/faqs", faqHandler.GetPublishedFAQs)
	mux.HandleFunc("/faqs", faqHandler.GetPublishedFAQs) // Endpoint alternativo

	// Rutas de compatibilidad de widget (para manejar las rutas duplicadas /widget/widget/...)
	mux.HandleFunc("/widget/widget/faqs", faqHandler.GetPublishedFAQs)

	// Ruta de WebSocket para el chat de tickets
	mux.HandleFunc("/api/ws/chat/", websocket.ChatHandler(store))

	// Middleware de CORS
	corsMiddleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Establecer encabezados CORS
			utils.SetCORS(w)

			// Manejar preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Llamar al manejador envolvente
			h.ServeHTTP(w, r)
		})
	}

	// Crear servidor con manejador envolvente
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      corsMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor en una goroutine
	go func() {
		log.Printf("Iniciando servidor en el puerto %d", *port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar servidor: %v", err)
		}
	}()

	// Esperar a la señal de interrupción para cerrar el servidor de maneragraceful
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Servidor se está cerrando...")

	// Crear contexto con timeout para cerrar el servidor
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Intentar cerrar el servidor de maneragraceful
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Servidor forzado a cerrarse: %v", err)
	}

	log.Println("Servidor cerrado de maneragraceful")
}
