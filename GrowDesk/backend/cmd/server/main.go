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
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/data"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/db"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/handlers"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/middleware"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se encontró archivo .env, usando variables de entorno directamente")
	}

	// Parsear flags
	var (
		port        = flag.Int("port", getEnvInt("PORT", 8080), "HTTP service port")
		dataDir     = flag.String("data-dir", getEnv("DATA_DIR", "./data"), "Directory for data storage")
		useMock     = flag.Bool("mock-auth", getEnvBool("MOCK_AUTH", true), "Use mock authentication for development")
		usePostgres = flag.Bool("postgres", getEnvBool("USE_POSTGRES", true), "Use PostgreSQL database instead of file storage")
		migrateData = flag.Bool("migrate", getEnvBool("MIGRATE_DATA", true), "Migrate data from JSON files to PostgreSQL")
	)
	flag.Parse()

	// Crear directorio de datos si no existe
	if err := os.MkdirAll(*dataDir, 0755); err != nil {
		log.Fatalf("Error al crear directorio de datos: %v", err)
	}

	// Inicializar el almacén de datos (store)
	var store data.DataStore

	// Decidir si usar PostgreSQL o almacenamiento en archivos basado en la flag
	if *usePostgres {
		log.Println("Usando PostgreSQL como almacén de datos")

		// Inicializar conexión a PostgreSQL
		database, err := db.InitDB()
		if err != nil {
			log.Fatalf("Error al conectar a PostgreSQL: %v", err)
		}
		defer db.Close()

		// Inicializar esquema de base de datos
		if err := db.InitializeSchema(database); err != nil {
			log.Fatalf("Error al inicializar esquema de base de datos: %v", err)
		}

		// Migrar datos desde JSON si se solicita
		if *migrateData {
			log.Println("Migrando datos de archivos JSON a PostgreSQL...")
			if err := db.MigrateAllFromJSON(database, *dataDir); err != nil {
				log.Printf("Advertencia: Error durante la migración de datos: %v", err)
			}
		}

		// Crear store PostgreSQL
		store = db.NewPostgreSQLStore(database)
	} else {
		log.Println("Usando almacenamiento en archivos")
		store = data.NewStore(*dataDir)
	}

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

	// Rutas de usuarios (autenticadas)
	mux.Handle("/api/users", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configurar CORS explícitamente
		utils.SetCORS(w)

		// Responder a preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Manejar basado en el método HTTP
		switch r.Method {
		case http.MethodGet:
			// Obtener todos los usuarios
			users, err := store.GetUsers()
			if err != nil {
				http.Error(w, "Error al obtener usuarios", http.StatusInternalServerError)
				return
			}
			utils.WriteJSON(w, http.StatusOK, users)
		case http.MethodPost:
			// Crear un nuevo usuario
			var user struct {
				FirstName  string `json:"firstName"`
				LastName   string `json:"lastName"`
				Email      string `json:"email"`
				Password   string `json:"password"`
				Role       string `json:"role"`
				Department string `json:"department"`
				Active     bool   `json:"active"`
			}

			// Leer el cuerpo de la solicitud
			if err := utils.DecodeJSON(r, &user); err != nil {
				http.Error(w, "Error al leer datos del usuario", http.StatusBadRequest)
				return
			}

			// Generar ID único
			id := fmt.Sprintf("%d", time.Now().UnixNano())

			// Crear nuevo usuario
			newUser := models.User{
				ID:         id,
				FirstName:  user.FirstName,
				LastName:   user.LastName,
				Email:      user.Email,
				Password:   user.Password, // esto debería ser hasheado
				Role:       user.Role,
				Department: user.Department,
				Active:     user.Active,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			// Agregar el usuario al store
			if err := store.CreateUser(newUser); err != nil {
				http.Error(w, "Error al guardar usuario", http.StatusInternalServerError)
				return
			}

			// Devolver el usuario creado
			utils.WriteJSON(w, http.StatusCreated, newUser)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})))

	// Rutas de usuarios individuales
	mux.Handle("/api/users/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configurar CORS explícitamente
		utils.SetCORS(w)

		// Responder a preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Obtener ID del usuario de la URL
		path := r.URL.Path
		segments := strings.Split(path, "/")
		if len(segments) < 4 {
			http.Error(w, "URL de usuario inválida", http.StatusBadRequest)
			return
		}

		userID := segments[3]

		// Manejar basado en el método HTTP
		switch r.Method {
		case http.MethodGet:
			// Obtener un usuario específico
			user, err := store.GetUser(userID)
			if err != nil {
				http.Error(w, "Usuario no encontrado", http.StatusNotFound)
				return
			}

			utils.WriteJSON(w, http.StatusOK, user)
		case http.MethodPut:
			// Actualizar un usuario
			var updates models.User
			if err := utils.DecodeJSON(r, &updates); err != nil {
				http.Error(w, "Error al leer datos de actualización", http.StatusBadRequest)
				return
			}

			// Obtener usuario existente
			user, err := store.GetUser(userID)
			if err != nil {
				http.Error(w, "Usuario no encontrado", http.StatusNotFound)
				return
			}

			// Actualizar campos
			if updates.FirstName != "" {
				user.FirstName = updates.FirstName
			}
			if updates.LastName != "" {
				user.LastName = updates.LastName
			}
			if updates.Email != "" {
				user.Email = updates.Email
			}
			if updates.Role != "" {
				user.Role = updates.Role
			}
			if updates.Department != "" {
				user.Department = updates.Department
			}
			if updates.Password != "" {
				user.Password = updates.Password // En producción real, esto debería ser hasheado
			}

			// Marcar como actualizado
			user.UpdatedAt = time.Now()

			// Actualizar en el store
			if err := store.UpdateUser(*user); err != nil {
				http.Error(w, "Error al actualizar usuario", http.StatusInternalServerError)
				return
			}

			// Devolver la respuesta actualizada
			utils.WriteJSON(w, http.StatusOK, user)
		case http.MethodDelete:
			// Eliminar un usuario
			if err := store.DeleteUser(userID); err != nil {
				http.Error(w, "Error al eliminar usuario", http.StatusInternalServerError)
				return
			}

			utils.WriteJSON(w, http.StatusOK, map[string]bool{"success": true})
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})))

	// Ruta de WebSocket para el chat de tickets
	mux.HandleFunc("/api/ws/chat/", func(w http.ResponseWriter, r *http.Request) {
		// Configurar CORS para WebSocket
		utils.SetCORS(w)

		// Extraer el ID del ticket de la URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "ID de ticket inválido", http.StatusBadRequest)
			return
		}
		ticketID := parts[len(parts)-1]

		// Actualizar a WebSocket
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Permitir todas las solicitudes en desarrollo
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error al actualizar conexión a WebSocket: %v", err)
			return
		}

		// Registrar la conexión
		connID := store.AddWSConnection(ticketID, conn)

		// Cerrar la conexión cuando finalice
		defer func() {
			conn.Close()
			store.RemoveWSConnection(ticketID, connID)
		}()

		// Enviar mensajes existentes del ticket
		ticket, err := store.GetTicket(ticketID)
		if err == nil && ticket != nil {
			// Enviar mensajes existentes
			wsMessage := models.WebSocketMessage{
				Type:     "init_messages",
				TicketID: ticketID,
				Messages: ticket.Messages,
			}

			conn.WriteJSON(wsMessage)
		}

		// Mantener la conexión abierta y escuchar mensajes
		for {
			messageType, _, err := conn.ReadMessage()
			if err != nil {
				break // Salir si hay error (cliente desconectado)
			}

			// Ping-pong para mantener la conexión activa
			if messageType == websocket.PingMessage {
				conn.WriteMessage(websocket.PongMessage, []byte{})
			}
		}
	})

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

// Helper para obtener variables de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper para obtener variables de entorno numéricas con valor por defecto
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	var result int
	_, err := fmt.Sscanf(value, "%d", &result)
	if err != nil {
		return defaultValue
	}
	return result
}

// Helper para obtener variables de entorno como booleano
func getEnvBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Printf("Advertencia: Variable de entorno %s no es un booleano válido, usando valor por defecto: %v", key, defaultValue)
		return defaultValue
	}

	return value
}
