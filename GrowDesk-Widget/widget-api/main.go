package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

// Estructuras para los mensajes y tickets
type TicketCreationData struct {
	Name     string                 `json:"name" binding:"required"`
	Email    string                 `json:"email" binding:"required,email"`
	Message  string                 `json:"message" binding:"required"`
	Metadata map[string]interface{} `json:"metadata"`
}

type TicketCreationResponse struct {
	Success           bool   `json:"success"`
	TicketID          string `json:"ticketId"`
	LiveChatAvailable bool   `json:"liveChatAvailable"`
}

type MessageData struct {
	TicketID string `json:"ticketId" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

type MessageResponse struct {
	MessageID string `json:"messageId"`
	Message   string `json:"message"`
}

// GrowDeskTicket es la estructura para enviar tickets al sistema GrowDesk
type GrowDeskTicket struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Subject     string                 `json:"subject"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	Priority    string                 `json:"priority"`
	Email       string                 `json:"email"`
	Name        string                 `json:"name"`
	ClientName  string                 `json:"clientName"`
	ClientEmail string                 `json:"clientEmail"`
	Department  string                 `json:"department"`
	Source      string                 `json:"source"`
	WidgetID    string                 `json:"widgetId"`
	CreatedAt   string                 `json:"createdAt"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// GrowDeskMessage es la estructura para enviar mensajes al sistema GrowDesk
type GrowDeskMessage struct {
	TicketID  string `json:"ticketId"`
	Content   string `json:"content"`
	UserID    string `json:"userId"`
	IsClient  bool   `json:"isClient"`
	UserName  string `json:"userName,omitempty"`
	UserEmail string `json:"userEmail,omitempty"`
}

// Ticket representa un ticket de soporte
type Ticket struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Messages    []Message `json:"messages"`
	UserEmail   string    `json:"userEmail"`
	UserName    string    `json:"userName"`
	ClientName  string    `json:"clientName"`
	ClientEmail string    `json:"clientEmail"`
	WidgetID    string    `json:"widgetId"`
	Department  string    `json:"department"`
	Metadata    Metadata  `json:"metadata"`
}

// Message representa un mensaje en un ticket
type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	IsClient  bool      `json:"isClient"`
	CreatedAt time.Time `json:"createdAt"`
	UserName  string    `json:"userName"`
	UserEmail string    `json:"userEmail"`
}

// Metadata contiene información adicional
type Metadata struct {
	URL        string `json:"url"`
	Referrer   string `json:"referrer"`
	UserAgent  string `json:"userAgent"`
	ScreenSize string `json:"screenSize"`
	ExternalID string `json:"externalId"`
}

// TicketRequest se utiliza para crear un nuevo ticket
type TicketRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Message  string   `json:"message"`
	Metadata Metadata `json:"metadata"`
}

// TicketResponse es la respuesta del servidor después de crear un ticket
type TicketResponse struct {
	TicketID string `json:"ticketId"`
	Message  string `json:"message"`
}

// MessageRequest se utiliza para enviar un mensaje a un ticket
type MessageRequest struct {
	TicketID  string `json:"ticketId"`
	Message   string `json:"message"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
}

// AgentMessageRequest es la estructura para mensajes enviados por agentes de soporte
type AgentMessageRequest struct {
	TicketID  string `json:"ticketId" binding:"required"`
	Content   string `json:"content" binding:"required"`
	UserID    string `json:"userId"`
	AgentName string `json:"agentName"`
}

// Configuración WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Permitir cualquier origen para desarrollo
	CheckOrigin: func(r *http.Request) bool {
		// Permitir todas las conexiones en modo desarrollo
		return true
	},
}

// Mapa para almacenar las conexiones WebSocket activas, agrupadas por ticketId
var wsConnections = make(map[string][]*websocket.Conn)
var wsConnectionsMutex = sync.Mutex{}

// Estructura para mensajes WebSocket
type WebSocketMessage struct {
	Type     string      `json:"type"`
	Message  interface{} `json:"message"`
	TicketID string      `json:"ticketId"`
}

// TicketData estructura para los datos recibidos al crear un ticket desde el widget
type TicketData struct {
	Subject     string `json:"subject" binding:"required"`
	Message     string `json:"message"`
	Priority    string `json:"priority"`
	Department  string `json:"department"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	WidgetID    string `json:"widgetId"`
	ClientName  string `json:"clientName"`
	ClientEmail string `json:"clientEmail"`
	Metadata    struct {
		URL        string `json:"url"`
		UserAgent  string `json:"userAgent"`
		Referrer   string `json:"referrer"`
		ScreenSize string `json:"screenSize"`
	} `json:"metadata"`
}

// GetUserInfo extrae información de usuario de los headers o el cuerpo de la solicitud
func GetUserInfo(c *gin.Context, req interface{}) (string, string) {
	// Primero intentar obtener de los headers
	userName := c.GetHeader("X-User-Name")
	userEmail := c.GetHeader("X-User-Email")

	// Si no están en los headers, intentar obtenerlos del cuerpo
	if userName == "" || userEmail == "" {
		switch r := req.(type) {
		case *TicketRequest:
			userName = r.Name
			userEmail = r.Email
		case *MessageRequest:
			userName = r.UserName
			userEmail = r.UserEmail
		}
	}

	return userName, userEmail
}

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Println("Archivo .env no encontrado, usando variables de entorno del sistema")
	}

	// Crear directorio de datos si no existe
	os.MkdirAll("data", 0755)

	// Configuración del router con CORS habilitado
	router := gin.Default()

	// Middleware para CORS
	router.Use(func(c *gin.Context) {
		// Obtener los orígenes permitidos desde variables de entorno
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "*" // Permitir todos los orígenes por defecto
		}

		// Establecer el origen en la respuesta
		origin := c.Request.Header.Get("Origin")
		if origin != "" && (allowedOrigins == "*" || strings.Contains(allowedOrigins, origin)) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Widget-ID, X-Widget-Token, X-User-Name, X-User-Email, X-Source, X-Client-Created, X-Widget-Ticket-ID, X-Message-Source, X-From-Client, X-Client-Message, X-Ticket-ID, Origin, Accept")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// Manejar solicitudes OPTIONS
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Endpoint de salud
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "GrowDesk Widget API está en funcionamiento",
		})
	})

	// Endpoint para verificar la conexión con el backend
	router.GET("/check-backend", func(c *gin.Context) {
		// Obtener URL del backend
		apiURL := os.Getenv("GROWDESK_API_URL")
		if apiURL == "" {
			apiURL = "http://localhost:8080"
		}

		// Normalizar URL
		baseURL := strings.TrimSuffix(apiURL, "/")
		healthURL := fmt.Sprintf("%s/health", baseURL)

		// Realizar solicitud al backend
		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(healthURL)

		if err != nil {
			log.Printf("Error al conectar con el backend: %v", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":      "ERROR",
				"message":     "No se pudo conectar con el backend: " + err.Error(),
				"backend_url": apiURL,
			})
			return
		}
		defer resp.Body.Close()

		// Leer respuesta
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error al leer respuesta del backend: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "ERROR",
				"message": "Error al leer respuesta del backend",
			})
			return
		}

		// Retornar información sobre la conexión
		c.JSON(http.StatusOK, gin.H{
			"status":           "OK",
			"backend_url":      apiURL,
			"backend_status":   resp.StatusCode,
			"backend_response": string(body),
			"message":          "Conexión exitosa con el backend",
		})
	})

	// Agregar endpoint de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "widget-api",
		})
	})

	// API de Widget - Incluir todas las rutas bajo /widget
	widgetAPI := router.Group("/widget")
	{
		// Endpoint para verificar estado del servicio
		widgetAPI.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "OK",
				"message": "GrowDesk Widget API está funcionando correctamente",
				"version": "1.0.0",
			})
		})

		// Tickets y mensajes
		widgetAPI.POST("/tickets", createTicket)
		widgetAPI.POST("/messages", sendMessage)
		widgetAPI.GET("/tickets/:ticketId/messages", getMessages)

		// Ruta para FAQs
		widgetAPI.GET("/faqs", getFaqs)
	}

	// WebSocket y API para agentes - Estas rutas no van bajo /widget
	router.GET("/api/ws/chat/:ticketId", handleWebSocketConnection)
	router.POST("/api/agent/messages", handleAgentMessage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Puerto por defecto
	}

	log.Printf("Servidor iniciado en el puerto %s", port)
	router.Run(":" + port)
}

// getFaqs retorna las preguntas frecuentes disponibles
func getFaqs(c *gin.Context) {
	log.Printf("Solicitando FAQs del backend")

	// Obtener URL del backend
	apiURL := os.Getenv("GROWDESK_API_URL")
	apiKey := os.Getenv("GROWDESK_API_KEY")

	if apiURL == "" {
		apiURL = "http://localhost:8080"
		log.Printf("GROWDESK_API_URL no definido, usando valor por defecto: %s", apiURL)
	}

	if apiKey == "" {
		apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiJhZG1pbi0xMjMiLCJlbWFpbCI6ImFkbWluQGdyb3dkZXNrLmNvbSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTcyNDA4ODQwMH0.8J5ayPvA4B-1vF5NaqRXycW1pmIl9qjKP6Ddj4Ot_Cw"
		log.Printf("GROWDESK_API_KEY no definido, usando valor por defecto")
	}

	// Widget ID para filtrar si está disponible
	widgetID := c.GetHeader("X-Widget-ID")

	// Construir URL
	baseURL := strings.TrimSuffix(apiURL, "/")
	faqsURL := fmt.Sprintf("%s/api/faqs", baseURL)

	// Añadir parámetro de widget si está disponible
	if widgetID != "" {
		faqsURL = fmt.Sprintf("%s?widgetId=%s", faqsURL, widgetID)
	}

	// Crear la solicitud
	req, err := http.NewRequest("GET", faqsURL, nil)
	if err != nil {
		log.Printf("Error al crear solicitud para FAQs: %v", err)
		// Devolver FAQs fallback
		c.JSON(http.StatusOK, getFallbackFaqs())
		return
	}

	// Configurar cabeceras
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	if widgetID != "" {
		req.Header.Set("X-Widget-ID", widgetID)
	}

	// Realizar solicitud
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error al obtener FAQs del backend: %v", err)
		// Devolver FAQs fallback
		c.JSON(http.StatusOK, getFallbackFaqs())
		return
	}
	defer resp.Body.Close()

	// Comprobar código de respuesta
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Error del backend al obtener FAQs. Código: %d, Respuesta: %s",
			resp.StatusCode, string(body))
		// Devolver FAQs fallback
		c.JSON(http.StatusOK, getFallbackFaqs())
		return
	}

	// Leer y reenviar la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error al leer respuesta de FAQs: %v", err)
		// Devolver FAQs fallback
		c.JSON(http.StatusOK, getFallbackFaqs())
		return
	}

	// Establecer el tipo de contenido y devolver la respuesta tal cual
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(body))
}

// getFallbackFaqs devuelve FAQs predeterminadas cuando no se pueden obtener del backend
func getFallbackFaqs() []gin.H {
	return []gin.H{
		{
			"id":          1,
			"question":    "¿Cómo puedo restablecer mi contraseña?",
			"answer":      "Puedes restablecer tu contraseña haciendo clic en el enlace 'Olvidé mi contraseña' en la página de inicio de sesión.",
			"category":    "Cuenta",
			"isPublished": true,
		},
		{
			"id":          2,
			"question":    "¿Cómo contactar con soporte?",
			"answer":      "Puedes contactar con nuestro equipo de soporte mediante este chat o enviando un correo a soporte@empresa.com",
			"category":    "Ayuda",
			"isPublished": true,
		},
		{
			"id":          3,
			"question":    "¿Cuáles son los horarios de atención?",
			"answer":      "Nuestro equipo está disponible de lunes a viernes de 9:00 a 18:00 horas.",
			"category":    "Ayuda",
			"isPublished": true,
		},
	}
}

// sendToGrowDesk envía datos al sistema GrowDesk
func sendToGrowDesk(url string, jsonData []byte, apiKey string, ticketID string) {
	// Crear una solicitud HTTP
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error al crear solicitud HTTP: %v", err)
		return
	}

	// Añadir cabeceras
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("X-Message-Source", "widget-client")
	req.Header.Set("X-Widget-ID", "true")

	// Crear cliente HTTP con timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Enviar solicitud
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error al enviar datos a GrowDesk: %v", err)
		return
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error al leer respuesta de GrowDesk: %v", err)
		return
	}

	// Comprobar si la respuesta fue exitosa
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("Datos enviados correctamente a GrowDesk para ticket %s. Respuesta: %s", ticketID, body)
	} else {
		log.Printf("Error de respuesta de GrowDesk para ticket %s. Código: %d, Respuesta: %s",
			ticketID, resp.StatusCode, body)
	}
}

// generateEmbedCode crea el código HTML para incrustar el widget
func generateEmbedCode(widgetId, widgetToken, brandName, welcomeMessage, primaryColor, position string) string {
	// Valores por defecto si no se proporcionan
	if widgetId == "" {
		widgetId = "widget-" + fmt.Sprintf("%d", time.Now().Unix())
	}

	if widgetToken == "" {
		widgetToken = "token-" + fmt.Sprintf("%d", time.Now().Unix())
	}

	if brandName == "" {
		brandName = "GrowDesk"
	}

	if welcomeMessage == "" {
		welcomeMessage = "¡Hola! ¿En qué podemos ayudarte hoy?"
	}

	if primaryColor == "" {
		primaryColor = "#3498db"
	}

	if position == "" {
		position = "right"
	}

	// URLs de producción
	baseUrl := os.Getenv("WIDGET_BASE_URL")
	apiUrl := os.Getenv("WIDGET_API_URL")

	// URLs por defecto (desarrollo)
	if baseUrl == "" {
		baseUrl = "http://localhost:3030"
	}

	if apiUrl == "" {
		apiUrl = "http://localhost:3000"
	}

	return `<script src="` + baseUrl + `/widget.js" id="growdesk-widget"
  data-widget-id="` + widgetId + `"
  data-widget-token="` + widgetToken + `"
  data-api-url="` + apiUrl + `"
  data-brand-name="` + brandName + `"
  data-welcome-message="` + welcomeMessage + `"
  data-primary-color="` + primaryColor + `"
  data-position="` + position + `">
</script>`
}

// SaveTicket guarda un ticket en el almacenamiento local
func SaveTicket(ticket Ticket) error {
	// Verificar si el directorio data existe
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		log.Printf("El directorio data no existe, creándolo...")
		if err := os.Mkdir("data", 0755); err != nil {
			log.Printf("Error al crear directorio data: %v", err)
			return err
		}
	}

	// Serializar ticket a JSON con indentación
	data, err := json.MarshalIndent(ticket, "", "  ")
	if err != nil {
		log.Printf("Error al serializar ticket a JSON: %v", err)
		return err
	}

	// Obtener ruta absoluta para el archivo de ticket
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Error al obtener directorio de trabajo: %v", err)
		wd = "."
	}

	// Guardar en archivo
	filename := fmt.Sprintf("data/ticket_%s.json", ticket.ID)
	absFilename := path.Join(wd, filename)
	log.Printf("Guardando ticket en archivo: %s", absFilename)

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Printf("Error al escribir archivo de ticket: %v", err)
		return err
	}

	log.Printf("Ticket guardado exitosamente: %s", filename)

	// Verificar que el archivo existe
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("¡Error! Archivo no encontrado después de guardar: %s", filename)
	} else {
		log.Printf("Verificación: archivo existe después de guardar: %s", filename)
	}

	return nil
}

// LoadTicket carga un ticket desde el almacenamiento local
func LoadTicket(ticketID string) (Ticket, error) {
	filename := fmt.Sprintf("data/ticket_%s.json", ticketID)
	data, err := os.ReadFile(filename)
	if err != nil {
		// Si no se encuentra el archivo, puede ser un ticket del sistema GrowDesk
		// Verificar si tiene formato de ticket de GrowDesk (TICKET-YYYYMMDDHHMMSS)
		if strings.HasPrefix(ticketID, "TICKET-") {
			// Intentar buscar en el sistema GrowDesk
			log.Printf("Ticket %s no encontrado localmente, buscando en GrowDesk...", ticketID)
			growdeskTicket, err := getTicketFromGrowDesk(ticketID)
			if err != nil {
				log.Printf("Error al buscar ticket en GrowDesk: %v", err)
				return Ticket{}, err
			}
			return growdeskTicket, nil
		}
		return Ticket{}, err
	}

	var ticket Ticket
	err = json.Unmarshal(data, &ticket)
	return ticket, err
}

// getMessages obtiene los mensajes de un ticket
func getMessages(c *gin.Context) {
	ticketId := c.Param("ticketId")
	if ticketId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ticket no proporcionado", "success": false})
		return
	}

	log.Printf("Solicitando mensajes para ticket: %s", ticketId)

	// Intentar obtener mensajes del backend Go primero
	apiURL := os.Getenv("GROWDESK_API_URL")
	apiKey := os.Getenv("GROWDESK_API_KEY")
	widgetID := c.GetHeader("X-Widget-ID")

	if apiURL == "" {
		apiURL = "http://localhost:8080"
		log.Printf("GROWDESK_API_URL no definido, usando valor por defecto: %s", apiURL)
	}

	if apiKey == "" {
		apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiJhZG1pbi0xMjMiLCJlbWFpbCI6ImFkbWluQGdyb3dkZXNrLmNvbSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTcyNDA4ODQwMH0.8J5ayPvA4B-1vF5NaqRXycW1pmIl9qjKP6Ddj4Ot_Cw"
	}

	// Construir URL
	baseURL := strings.TrimSuffix(apiURL, "/")
	messagesURL := fmt.Sprintf("%s/api/tickets/%s/messages", baseURL, ticketId)

	// Intenta obtener mensajes del backend
	if gotMessagesFromBackend := tryGetMessagesFromBackend(c, messagesURL, apiKey, widgetID); gotMessagesFromBackend {
		return // Si tuvo éxito, terminamos
	}

	// Si llegamos aquí, fallamos al obtener mensajes del backend
	loadLocalMessages(c, ticketId)
}

// tryGetMessagesFromBackend intenta obtener mensajes del backend Go
// Devuelve true si tuvo éxito y ya envió la respuesta al cliente
func tryGetMessagesFromBackend(c *gin.Context, messagesURL, apiKey, widgetID string) bool {
	// Crear la solicitud
	req, err := http.NewRequest("GET", messagesURL, nil)
	if err != nil {
		log.Printf("Error al crear solicitud para mensajes: %v", err)
		return false
	}

	// Configurar cabeceras
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	if widgetID != "" {
		req.Header.Set("X-Widget-ID", widgetID)
	}

	// Realizar solicitud al backend Go
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error al obtener mensajes del backend: %v", err)
		return false
	}
	defer resp.Body.Close()

	// Verificar código de respuesta
	if resp.StatusCode == http.StatusOK {
		// Leer respuesta
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			// Enviar respuesta del backend directo al cliente
			c.Header("Content-Type", "application/json")
			c.String(http.StatusOK, string(body))
			return true
		}
		log.Printf("Error al leer respuesta de mensajes: %v", err)
	} else {
		log.Printf("Error del backend al obtener mensajes. Código: %d", resp.StatusCode)
	}

	return false
}

// loadLocalMessages carga mensajes de un ticket almacenado localmente
func loadLocalMessages(c *gin.Context, ticketId string) {
	log.Printf("Usando fallback para cargar mensajes localmente")

	// Cargar el ticket localmente
	ticket, err := LoadTicket(ticketId)
	if err != nil {
		log.Printf("Error al cargar ticket localmente: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Ticket no encontrado",
			"success": false,
		})
		return
	}

	// Devolver mensajes al cliente
	messages := ticket.Messages
	if messages == nil {
		messages = []Message{} // Asegurar que devolvemos al menos un array vacío
	}

	c.JSON(http.StatusOK, messages)
}

// getTicketFromGrowDesk obtiene un ticket desde el sistema principal GrowDesk
func getTicketFromGrowDesk(ticketID string) (Ticket, error) {
	var ticket Ticket

	// Verificar que tenemos la URL y API key
	apiURL := os.Getenv("GROWDESK_API_URL")
	apiKey := os.Getenv("GROWDESK_API_KEY")

	if apiURL == "" {
		apiURL = "http://localhost:8080"
		log.Printf("GROWDESK_API_URL no definido, usando valor por defecto: %s", apiURL)
	}

	if apiKey == "" {
		apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiJhZG1pbi0xMjMiLCJlbWFpbCI6ImFkbWluQGdyb3dkZXNrLmNvbSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTcyNDA4ODQwMH0.8J5ayPvA4B-1vF5NaqRXycW1pmIl9qjKP6Ddj4Ot_Cw" // Token por defecto para desarrollo
		log.Printf("GROWDESK_API_KEY no definido, usando valor por defecto")
	}

	// Construir URL para obtener el ticket
	ticketURL := fmt.Sprintf("%s/api/tickets/%s", apiURL, ticketID)
	log.Printf("Solicitando ticket a: %s", ticketURL)

	// Crear cliente HTTP con timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Crear request
	req, err := http.NewRequest("GET", ticketURL, nil)
	if err != nil {
		log.Printf("Error al crear request para obtener ticket: %v", err)
		return ticket, err
	}

	// Añadir headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Enviar request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error al solicitar ticket a GrowDesk: %v", err)
		return ticket, err
	}
	defer resp.Body.Close()

	// Verificar código de respuesta
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error al obtener ticket, código de estado: %d", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Respuesta: %s", string(body))
		return ticket, fmt.Errorf("error al obtener ticket, código: %d", resp.StatusCode)
	}

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error al leer respuesta: %v", err)
		return ticket, err
	}

	log.Printf("Respuesta del servidor GrowDesk: %s", string(body))

	// Formato esperado de la respuesta para un ticket de GrowDesk
	type GrowDeskTicketResponse struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		CreatedAt   string `json:"createdAt"`
		Customer    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"customer"`
		Messages []struct {
			ID        string `json:"id"`
			Content   string `json:"content"`
			IsClient  bool   `json:"isClient"`
			Timestamp string `json:"timestamp"`
		} `json:"messages"`
	}

	// Deserializar respuesta
	var growdeskTicket GrowDeskTicketResponse
	if err := json.Unmarshal(body, &growdeskTicket); err != nil {
		log.Printf("Error al deserializar respuesta: %v", err)
		return ticket, err
	}

	// Validar que el ID coincide
	if growdeskTicket.ID != ticketID {
		log.Printf("Error: ID del ticket recibido (%s) no coincide con el solicitado (%s)",
			growdeskTicket.ID, ticketID)
		return ticket, fmt.Errorf("id de ticket no coincide")
	}

	// Convertir a formato de ticket local
	ticket = Ticket{
		ID:          growdeskTicket.ID,
		Title:       growdeskTicket.Title,
		Description: growdeskTicket.Description,
		Status:      growdeskTicket.Status,
		CreatedBy:   growdeskTicket.Customer.Email,
		UserEmail:   growdeskTicket.Customer.Email,
		UserName:    growdeskTicket.Customer.Name,
		Metadata:    Metadata{},
	}

	// Convertir timestamp
	createdAt, err := time.Parse(time.RFC3339, growdeskTicket.CreatedAt)
	if err == nil {
		ticket.CreatedAt = createdAt
		ticket.UpdatedAt = time.Now() // Actualizar fecha de actualización a ahora
	} else {
		log.Printf("Error al parsear timestamp: %v", err)
		ticket.CreatedAt = time.Now()
		ticket.UpdatedAt = time.Now()
	}

	// Convertir mensajes
	for _, msg := range growdeskTicket.Messages {
		newMsg := Message{
			ID:       msg.ID,
			Content:  msg.Content,
			IsClient: msg.IsClient,
		}

		// Convertir timestamp
		msgTime, err := time.Parse(time.RFC3339, msg.Timestamp)
		if err == nil {
			newMsg.CreatedAt = msgTime
		} else {
			log.Printf("Error al parsear timestamp del mensaje: %v", err)
			newMsg.CreatedAt = time.Now()
		}

		// Añadir mensaje a la lista
		ticket.Messages = append(ticket.Messages, newMsg)
	}

	log.Printf("Ticket obtenido correctamente de GrowDesk, contiene %d mensajes",
		len(ticket.Messages))
	return ticket, nil
}

// createTicket handles creation of new support tickets
func createTicket(c *gin.Context) {
	log.Printf("===== INICIO CREACIÓN TICKET WIDGET =====")

	// Verificar token de widget si está configurado
	widgetToken := c.GetHeader("X-Widget-Token")
	widgetID := c.GetHeader("X-Widget-ID")

	log.Printf("Widget ID: %s", widgetID)
	log.Printf("Token recibido: %s", widgetToken)

	// Obtener el cuerpo de la solicitud para depuración
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error al leer el cuerpo de la solicitud: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer la solicitud"})
		return
	}

	// Restaurar el cuerpo para que pueda ser leído nuevamente
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Loguear el cuerpo para depuración
	log.Printf("Cuerpo de la solicitud: %s", string(bodyBytes))

	// Primero intentar unmarshall a una estructura intermedia más flexible
	var dataRaw map[string]interface{}
	if jsonErr := json.Unmarshal(bodyBytes, &dataRaw); jsonErr != nil {
		log.Printf("Error al parsear JSON básico: %v", jsonErr)
	} else {
		// Verificar si falta el campo subject y añadirlo
		if _, hasSubject := dataRaw["subject"]; !hasSubject {
			log.Printf("SUBJECT NO ENCONTRADO en solicitud, añadiéndolo")

			// Obtener el nombre para generar un subject
			var name string
			if nameVal, hasName := dataRaw["name"]; hasName {
				name = fmt.Sprintf("%v", nameVal)
			} else {
				name = "Cliente"
			}

			// Añadir el subject al mapa
			dataRaw["subject"] = fmt.Sprintf("Solicitud de soporte - %s", name)

			// Reserializar a JSON
			modifiedBytes, jsonErr := json.Marshal(dataRaw)
			if jsonErr == nil {
				log.Printf("Cuerpo de solicitud modificado con subject: %s", string(modifiedBytes))
				// Reemplazar el body
				c.Request.Body = io.NopCloser(bytes.NewBuffer(modifiedBytes))
				// Actualizar bodyBytes para logging
				bodyBytes = modifiedBytes
			} else {
				log.Printf("Error al reserializar JSON con subject: %v", jsonErr)
			}
		}
	}

	var ticketData TicketData
	if err := c.ShouldBindJSON(&ticketData); err != nil {
		log.Printf("ERROR DETALLADO EN DATOS DE ENTRADA: %v", err)
		log.Printf("Tipo de error: %T", err)

		// Errores de validación mostrados más detalladamente
		if ve, ok := err.(validator.ValidationErrors); ok {
			errorDetails := []string{}
			for _, e := range ve {
				errorDetails = append(errorDetails, fmt.Sprintf("Campo '%s' falló validación '%s'", e.Field(), e.Tag()))
				log.Printf("Campo con error: '%s', valor recibido: '%v', validación fallida: '%s'",
					e.Field(), e.Value(), e.Tag())
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validación fallida", "details": errorDetails})
		} else {
			// Otros tipos de errores
			log.Printf("Error de JSON no validación: %v", err)

			// Intentar realizar un parse manual como fallback
			var rawData map[string]interface{}
			if jsonErr := json.Unmarshal(bodyBytes, &rawData); jsonErr == nil {
				log.Printf("Parse manual exitoso: %v", rawData)

				// Crear una estructura TicketData mínima
				ticketData = TicketData{
					Subject: fmt.Sprintf("%v", rawData["subject"]),
				}

				// Extraer otros campos si existen
				if name, ok := rawData["name"].(string); ok {
					ticketData.Name = name
				}

				if email, ok := rawData["email"].(string); ok {
					ticketData.Email = email
				}

				if message, ok := rawData["message"].(string); ok {
					ticketData.Message = message
				}

				if priority, ok := rawData["priority"].(string); ok {
					ticketData.Priority = priority
				}

				if dept, ok := rawData["department"].(string); ok {
					ticketData.Department = dept
				}

				if widgetID, ok := rawData["widgetId"].(string); ok {
					ticketData.WidgetID = widgetID
				}

				// Intenta extraer metadata si existe
				if meta, ok := rawData["metadata"].(map[string]interface{}); ok {
					if url, ok := meta["url"].(string); ok {
						ticketData.Metadata.URL = url
					}
					if ua, ok := meta["userAgent"].(string); ok {
						ticketData.Metadata.UserAgent = ua
					}
					if ref, ok := meta["referrer"].(string); ok {
						ticketData.Metadata.Referrer = ref
					}
					if ss, ok := meta["screenSize"].(string); ok {
						ticketData.Metadata.ScreenSize = ss
					}
				}

				// Verificar que al menos tengamos un subject
				if ticketData.Subject == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "El campo 'subject' es obligatorio"})
					return
				}

				// Continuar con el proceso normal
				log.Printf("Recuperación exitosa con datos mínimos: %+v", ticketData)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
	}

	// Verificar campos obligatorios
	if ticketData.Subject == "" {
		log.Printf("Error: Subject es obligatorio")
		c.JSON(http.StatusBadRequest, gin.H{"error": "El asunto del ticket es obligatorio"})
		return
	}

	// Obtener información del usuario
	userName := ticketData.Name
	userEmail := ticketData.Email
	clientName := ticketData.ClientName
	clientEmail := ticketData.ClientEmail

	// Si no hay cliente específico, usar los valores de name/email
	if clientName == "" {
		clientName = userName
	}

	if clientEmail == "" {
		clientEmail = userEmail
	}

	if userName == "" {
		userName = c.GetHeader("X-User-Name")
	}

	if userEmail == "" {
		userEmail = c.GetHeader("X-User-Email")
	}

	// Si aún no tenemos nombre/email del usuario, usar valores por defecto
	if userName == "" {
		userName = "Anónimo"
	}

	if userEmail == "" {
		// Generar un email temporal basado en timestamp
		userEmail = fmt.Sprintf("user_%d@temporary.com", time.Now().Unix())
	}

	// Asegurar que cliente tenga valores válidos
	if clientName == "" {
		clientName = userName
	}

	if clientEmail == "" {
		clientEmail = userEmail
	}

	log.Printf("Datos validados: Subject='%s', Name='%s', Email='%s', ClientName='%s', ClientEmail='%s'",
		ticketData.Subject, userName, userEmail, clientName, clientEmail)

	// Generar ID de ticket único
	now := time.Now()
	ticketID := fmt.Sprintf("TICKET-%s", now.Format("20060102-150405"))

	log.Printf("ID de ticket generado: %s", ticketID)

	// Crear el ticket
	ticket := Ticket{
		ID:          ticketID,
		Title:       ticketData.Subject, // Usamos Subject como Title
		Subject:     ticketData.Subject,
		Status:      "open",
		Priority:    ticketData.Priority,
		Description: ticketData.Message,
		CreatedAt:   now,
		UpdatedAt:   now,
		UserEmail:   userEmail,
		UserName:    userName,
		ClientName:  clientName,
		ClientEmail: clientEmail,
		WidgetID:    ticketData.WidgetID,
		Department:  ticketData.Department,
		Messages:    []Message{},
		Metadata: Metadata{
			URL:        ticketData.Metadata.URL,
			Referrer:   ticketData.Metadata.Referrer,
			UserAgent:  ticketData.Metadata.UserAgent,
			ScreenSize: ticketData.Metadata.ScreenSize,
		},
	}

	// Si hay un mensaje inicial, añadirlo al ticket
	if ticketData.Message != "" {
		messageID := fmt.Sprintf("MSG-%d", time.Now().UnixNano())
		message := Message{
			ID:        messageID,
			Content:   ticketData.Message,
			IsClient:  true, // Es un mensaje del cliente
			CreatedAt: now,
			UserName:  userName,
			UserEmail: userEmail,
		}
		ticket.Messages = append(ticket.Messages, message)
	}

	// Guardar el ticket localmente
	if err := SaveTicket(ticket); err != nil {
		log.Printf("Error al guardar ticket localmente: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el ticket"})
		return
	}

	// Enviar el ticket al sistema GrowDesk en una goroutine separada
	go func() {
		// Configuración del API de GrowDesk
		apiURL := os.Getenv("GROWDESK_API_URL")
		apiKey := os.Getenv("GROWDESK_API_KEY")

		if apiURL == "" {
			apiURL = "http://localhost:8080"
			log.Printf("GROWDESK_API_URL no definido, usando valor por defecto: %s", apiURL)
		}

		if apiKey == "" {
			apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiJhZG1pbi0xMjMiLCJlbWFpbCI6ImFkbWluQGdyb3dkZXNrLmNvbSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTcyNDA4ODQwMH0.8J5ayPvA4B-1vF5NaqRXycW1pmIl9qjKP6Ddj4Ot_Cw"
			log.Printf("GROWDESK_API_KEY no definido, usando valor por defecto")
		}

		// Preparar datos para GrowDesk
		// La estructura debe coincidir con lo que espera el backend de GrowDesk
		growDeskTicket := GrowDeskTicket{
			ID:          ticketID,
			Title:       ticketData.Subject,
			Subject:     ticketData.Subject,
			Description: ticketData.Message,
			Status:      "open",
			Priority:    ticketData.Priority,
			Name:        userName,
			Email:       userEmail,
			ClientName:  clientName,
			ClientEmail: clientEmail,
			Department:  ticketData.Department,
			Source:      "widget",
			WidgetID:    ticketData.WidgetID,
			CreatedAt:   now.Format(time.RFC3339),
		}

		// Convertir a JSON
		jsonData, err := json.Marshal(growDeskTicket)
		if err != nil {
			log.Printf("Error al convertir ticket a JSON: %v", err)
			return
		}

		// Normalizar URL base
		baseURL := strings.TrimSuffix(apiURL, "/")

		// Construir URL para la API de tickets
		url := fmt.Sprintf("%s/widget/tickets", baseURL)

		log.Printf("Enviando ticket al sistema GrowDesk en URL: %s", url)
		log.Printf("Payload JSON: %s", string(jsonData))

		headers := map[string]string{
			"Content-Type":     "application/json",
			"Authorization":    "Bearer " + apiKey,
			"X-Source":         "widget",
			"X-Widget-ID":      ticketData.WidgetID,
			"X-Client-Created": "true",
		}

		// Enviar el ticket con reintentos
		resp, err := sendHttpRequestWithRetry(url, jsonData, headers, 3)
		if err != nil {
			log.Printf("Error en todas las llamadas al sistema GrowDesk: %v", err)
		} else if resp != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			log.Printf("Respuesta de GrowDesk: Status %d, Body: %s", resp.StatusCode, string(body))

			// Si el ticket ya existe en GrowDesk, actualizar su estado
			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				log.Printf("Ticket creado correctamente en GrowDesk")

				// Procesar la respuesta si es necesario
				var respData map[string]interface{}
				if err := json.Unmarshal(body, &respData); err == nil {
					// Si GrowDesk devuelve un ID, actualizar el ticket local
					if growdeskID, ok := respData["id"].(string); ok && growdeskID != "" && growdeskID != ticketID {
						log.Printf("GrowDesk asignó un ID diferente al ticket: %s (local: %s)", growdeskID, ticketID)
						// Podría ser necesario actualizar la referencia local con este ID
					}
				}
			} else {
				log.Printf("Error al crear ticket en GrowDesk. Código: %d", resp.StatusCode)
			}
		}
	}()

	// Responder al cliente con el ID del ticket creado
	// IMPORTANTE: incluir "id" en la respuesta ya que el widget lo espera
	c.JSON(http.StatusCreated, gin.H{
		"ticketId":          ticketID,
		"message":           "Ticket creado correctamente",
		"success":           true,
		"id":                ticketID, // Campo importante para el widget
		"liveChatAvailable": true,
	})

	log.Printf("===== FIN CREACIÓN TICKET WIDGET =====")
}

// sendHttpRequestWithRetry envía una solicitud HTTP con reintentos
func sendHttpRequestWithRetry(url string, jsonData []byte, headers map[string]string, maxRetries int) (*http.Response, error) {
	var resp *http.Response
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		// Crear solicitud
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error al crear solicitud HTTP (intento %d): %v", i+1, err)
			lastErr = err
			continue
		}

		// Añadir cabeceras
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		// Crear cliente con timeout
		client := &http.Client{
			Timeout: 30 * time.Second, // Aumentado de 10 a 30 segundos
		}

		// Enviar solicitud
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Éxito
			return resp, nil
		}

		if err != nil {
			log.Printf("Error en intento %d: %v", i+1, err)
			lastErr = err

			// Intenta resolver problemas de DNS o hosts
			if strings.Contains(err.Error(), "lookup backend") {
				log.Printf("Error de resolución DNS detectado. Intentando con nombre completo 'growdesk-backend'")
				if strings.Contains(url, "backend:") {
					url = strings.Replace(url, "backend:", "growdesk-backend:", 1)
					log.Printf("URL ajustada a: %s", url)
				}
			}
		} else {
			log.Printf("Respuesta no exitosa en intento %d: %d", i+1, resp.StatusCode)
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			log.Printf("Contenido de respuesta: %s", string(bodyBytes))
			resp.Body.Close()
			lastErr = fmt.Errorf("código de respuesta no exitoso: %d", resp.StatusCode)
		}

		// Esperar antes de reintentar (backoff exponencial)
		waitTime := time.Duration(300*(i+1)) * time.Millisecond
		log.Printf("Esperando %v antes del siguiente intento...", waitTime)
		time.Sleep(waitTime)
	}

	return nil, fmt.Errorf("fallo después de %d intentos: %v", maxRetries, lastErr)
}

// handleWebSocketConnection maneja las conexiones WebSocket para chat en tiempo real
func handleWebSocketConnection(c *gin.Context) {
	ticketId := c.Param("ticketId")
	if ticketId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ticket no proporcionado"})
		return
	}

	log.Printf("Intentando establecer conexión WebSocket para ticket: %s", ticketId)

	// Verificar que el ticket existe antes de permitir una conexión WebSocket
	ticket, err := LoadTicket(ticketId)
	if err != nil {
		log.Printf("Error al cargar ticket para WebSocket: %v. Se intentará crear o utilizar un ticket existente.", err)

		// Si el ticketId tiene un formato específico (como TICKET-YYYYMMDD-HHMMSS)
		// Podemos intentar usar otro ticket activo del mismo usuario
		session := getSessionFromRequest(c)
		if session.ticketId != "" && session.ticketId != ticketId {
			log.Printf("Se encontró un ticket alternativo %s en la sesión del usuario", session.ticketId)
			// Redirigir a la conexión WebSocket con el ticket correcto
			c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/api/ws/chat/%s", session.ticketId))
			return
		}
	} else {
		log.Printf("Ticket encontrado: %s - %s", ticket.ID, ticket.Title)
	}

	// Configuración para la actualización
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true // Permitir conexiones de cualquier origen para desarrollo
	}

	// Mejorar a WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error al mejorar a WebSocket: %v", err)
		return
	}

	// Registrar la conexión
	wsConnectionsMutex.Lock()
	if _, exists := wsConnections[ticketId]; !exists {
		wsConnections[ticketId] = make([]*websocket.Conn, 0)
	}
	wsConnections[ticketId] = append(wsConnections[ticketId], ws)
	wsConnectionsMutex.Unlock()

	log.Printf("Nueva conexión WebSocket establecida para ticket: %s", ticketId)

	// Enviar mensaje de bienvenida/confirmación de conexión
	welcomeMsg := map[string]interface{}{
		"type": "connection_established",
		"data": map[string]interface{}{
			"message":  "Conexión establecida",
			"ticketId": ticketId,
			"status":   "connected",
		},
	}

	welcomeBytes, _ := json.Marshal(welcomeMsg)
	if err := ws.WriteMessage(websocket.TextMessage, welcomeBytes); err != nil {
		log.Printf("Error al enviar mensaje de bienvenida: %v", err)
	}

	// Manejar desconexión
	defer func() {
		ws.Close()
		wsConnectionsMutex.Lock()
		// Eliminar la conexión del slice
		for i, conn := range wsConnections[ticketId] {
			if conn == ws {
				wsConnections[ticketId] = append(wsConnections[ticketId][:i], wsConnections[ticketId][i+1:]...)
				break
			}
		}
		// Si no hay más conexiones para este ticket, eliminar la entrada
		if len(wsConnections[ticketId]) == 0 {
			delete(wsConnections, ticketId)
		}
		wsConnectionsMutex.Unlock()
		log.Printf("Conexión WebSocket cerrada para ticket: %s", ticketId)
	}()

	// Mantener la conexión y escuchar mensajes
	for {
		// Leer mensaje (puede ser un ping o mensaje real)
		messageType, messageBytes, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error al leer mensaje WebSocket: %v", err)
			break // Salir si hay un error (cliente desconectado)
		}

		// Solo procesar mensajes de texto
		if messageType != websocket.TextMessage {
			continue
		}

		// Intentar parsear JSON
		var message map[string]interface{}
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("Mensaje no es JSON válido: %v", err)
			continue
		}

		// Determinar tipo de mensaje
		msgType, ok := message["type"].(string)
		if !ok {
			log.Printf("Mensaje sin tipo definido")
			continue
		}

		log.Printf("Mensaje WebSocket recibido - Tipo: %s", msgType)

		// Manejar diferentes tipos de mensajes
		switch msgType {
		case "ping":
			// Responder con pong
			pongMsg := map[string]interface{}{
				"type": "pong",
				"data": map[string]string{
					"time": time.Now().Format(time.RFC3339),
				},
			}
			pongBytes, _ := json.Marshal(pongMsg)
			if err := ws.WriteMessage(websocket.TextMessage, pongBytes); err != nil {
				log.Printf("Error al enviar pong: %v", err)
			}

		case "client_message":
			// Extraer contenido del mensaje
			var content string
			data, ok := message["data"].(map[string]interface{})
			if !ok {
				data, ok = message["message"].(map[string]interface{})
				if !ok {
					log.Printf("Formato incorrecto para mensaje de cliente")
					continue
				}
			}

			// Obtener contenido del mensaje
			if contentVal, ok := data["content"].(string); ok {
				content = contentVal
			} else if contentVal, ok := data["message"].(string); ok {
				content = contentVal
			} else {
				log.Printf("Contenido del mensaje no encontrado")
				continue
			}

			// Crear estructura para el mensaje
			messageData := MessageData{
				TicketID: ticketId,
				Message:  content,
			}

			// Simular solicitud para procesamiento
			mockContext := &gin.Context{}
			reqBody, _ := json.Marshal(messageData)
			mockContext.Request = &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(reqBody)),
			}

			// Usar la función existente para procesar el mensaje
			sendMessage(mockContext)

		default:
			log.Printf("Tipo de mensaje no manejado: %s", msgType)
		}
	}
}

// Función auxiliar para obtener información de sesión de las cookies o headers
type sessionInfo struct {
	name     string
	email    string
	ticketId string
}

func getSessionFromRequest(c *gin.Context) sessionInfo {
	result := sessionInfo{}

	// Intentar obtener del encabezado
	result.name = c.GetHeader("X-User-Name")
	result.email = c.GetHeader("X-User-Email")
	result.ticketId = c.GetHeader("X-Ticket-ID")

	// Intentar obtener de las cookies
	cookieValue, err := c.Cookie("growdesk_session")
	if err == nil && cookieValue != "" {
		var sessionData map[string]interface{}
		if err := json.Unmarshal([]byte(cookieValue), &sessionData); err == nil {
			if name, ok := sessionData["name"].(string); ok {
				result.name = name
			}
			if email, ok := sessionData["email"].(string); ok {
				result.email = email
			}
			if ticketId, ok := sessionData["ticketId"].(string); ok {
				result.ticketId = ticketId
			}
		}
	}

	return result
}

// sendMessageToWebSocketClients envía un mensaje a todos los clientes WebSocket conectados a un ticket
func sendMessageToWebSocketClients(ticketId string, message Message) {
	wsConnectionsMutex.Lock()
	defer wsConnectionsMutex.Unlock()

	connections, exists := wsConnections[ticketId]
	if !exists || len(connections) == 0 {
		log.Printf("No hay conexiones WebSocket activas para el ticket: %s", ticketId)
		return
	}

	// IMPORTANTE: Asegurarse de que el mensaje tiene la estructura esperada
	// Crear un mapa explícito con los campos exactos que espera el cliente
	messageObj := map[string]interface{}{
		"id":        message.ID,
		"content":   message.Content,
		"isClient":  message.IsClient,
		"createdAt": message.CreatedAt.Format(time.RFC3339),
		"timestamp": message.CreatedAt.Format(time.RFC3339), // Para compatibilidad
		"userName":  message.UserName,
		"userEmail": message.UserEmail,
	}

	// Estructura compatible con ambos backends (JS y Go)
	wsMessage := map[string]interface{}{
		"type":     "new_message",
		"data":     messageObj, // Para compatibilidad con Go
		"message":  messageObj, // Para compatibilidad con JS
		"ticketId": ticketId,
	}

	// Serializar a JSON
	msgBytes, err := json.Marshal(wsMessage)
	if err != nil {
		log.Printf("Error al serializar mensaje WebSocket: %v", err)
		return
	}

	log.Printf("Enviando mensaje a %d cliente(s) del ticket: %s", len(connections), ticketId)
	log.Printf("Mensaje: %s", message.Content)

	// Enviar a todas las conexiones
	for _, conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
			log.Printf("Error al enviar mensaje WebSocket: %v", err)
		} else {
			log.Printf("Mensaje enviado con éxito")
		}
	}
}

// handleAgentMessage procesa mensajes enviados por agentes y los reenvía a los clientes
func handleAgentMessage(c *gin.Context) {
	var req AgentMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar ticket ID
	if req.TicketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	log.Printf("Recibido mensaje de agente para ticket: %s, contenido: %s", req.TicketID, req.Content)

	// Cargar ticket
	ticket, err := LoadTicket(req.TicketID)
	if err != nil {
		log.Printf("Error al cargar ticket %s: %v", req.TicketID, err)

		// Si el ticket no existe, crear uno nuevo
		if os.IsNotExist(err) {
			log.Printf("El ticket no existe, creando uno nuevo con ID: %s", req.TicketID)

			// Crear ticket
			ticket = Ticket{
				ID:          req.TicketID,
				Title:       "Solicitud de soporte de agente",
				Description: "Este ticket fue creado automáticamente al recibir un mensaje de un agente.",
				Status:      "new",
				CreatedBy:   "agent",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Messages:    []Message{},
				UserEmail:   "client@example.com",
				UserName:    "Cliente",
				Metadata:    Metadata{},
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cargar el ticket", "details": err.Error()})
			return
		}
	}

	// Usar uuid o un ID de mensaje con formato adecuado
	messageID := fmt.Sprintf("MSG-%d", time.Now().UnixNano())

	// Si no se proporcionó nombre del agente, usar uno por defecto
	agentName := req.AgentName
	if agentName == "" {
		agentName = "Soporte"
	}

	// Crear nuevo mensaje (desde agente, no cliente)
	newMessage := Message{
		ID:        messageID,
		Content:   req.Content,
		IsClient:  false, // Mensaje de agente - EXPLÍCITAMENTE FALSE
		CreatedAt: time.Now(),
		UserName:  agentName, // Nombre del agente
	}

	// Agregar mensaje al ticket
	ticket.Messages = append(ticket.Messages, newMessage)
	ticket.UpdatedAt = time.Now()

	// Guardar ticket actualizado
	if err := SaveTicket(ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message", "details": err.Error()})
		return
	}

	log.Printf("Mensaje de agente guardado correctamente en ticket %s", req.TicketID)

	// IMPORTANTE: Asegurarse de que el mensaje tenga el formato correcto antes de enviarlo por WebSocket
	log.Printf("Enviando mensaje de agente a clientes via WebSocket. IsClient=%v", newMessage.IsClient)

	// Enviar a todos los clientes conectados por WebSocket
	sendMessageToWebSocketClients(req.TicketID, newMessage)

	// Enviar mensaje al sistema principal GrowDesk
	growDeskMessage := GrowDeskMessage{
		TicketID:  req.TicketID,
		Content:   req.Content,
		IsClient:  false, // FALSE para mensajes de agente
		UserID:    req.UserID,
		UserName:  agentName,
		UserEmail: "agent@growdesk.com",
	}

	// Serializar los datos a JSON
	jsonData, err := json.Marshal(growDeskMessage)
	if err != nil {
		log.Printf("Error al serializar mensaje de agente para GrowDesk: %v", err)
	} else {
		// Obtener API URL
		apiURL := os.Getenv("GROWDESK_API_URL")
		apiKey := os.Getenv("GROWDESK_API_KEY")

		if apiURL == "" {
			apiURL = "http://localhost:8080"
			log.Printf("GROWDESK_API_URL no definido, usando valor por defecto: %s", apiURL)
		}

		if apiKey == "" {
			apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiJhZG1pbi0xMjMiLCJlbWFpbCI6ImFkbWluQGdyb3dkZXNrLmNvbSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTcyNDA4ODQwMH0.8J5ayPvA4B-1vF5NaqRXycW1pmIl9qjKP6Ddj4Ot_Cw"
		}

		// Normalizar URL base
		baseURL := strings.TrimSuffix(apiURL, "/")

		// Construir URL para la API de mensajes
		url := fmt.Sprintf("%s/widget/tickets/%s/messages?from_client=true", baseURL, req.TicketID)

		log.Printf("Enviando mensaje al sistema GrowDesk en la URL: %s", url)

		headers := map[string]string{
			"Content-Type":       "application/json",
			"Authorization":      "Bearer " + apiKey,
			"X-Message-Source":   "widget-client",
			"X-Widget-ID":        req.TicketID,
			"X-Client-Message":   "true",
			"X-Widget-Ticket-ID": req.TicketID,
			"X-From-Client":      "true",
		}

		// Enviar con reintentos
		resp, err := sendHttpRequestWithRetry(url, jsonData, headers, 3)
		if err != nil {
			log.Printf("Error al enviar mensaje a GrowDesk: %v", err)
		} else if resp != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			log.Printf("Respuesta al enviar mensaje a GrowDesk: Status %d, Body: %s", resp.StatusCode, string(body))
		}
	}

	// Devolver respuesta de éxito
	c.JSON(http.StatusOK, gin.H{
		"messageId": newMessage.ID,
		"message":   "Agent message sent successfully",
		"success":   true,
	})
}

// sendMessage agrega un mensaje a un ticket existente
func sendMessage(c *gin.Context) {
	var messageData MessageData

	// Capturar el cuerpo original para depuración si es necesario
	var bodyBytes []byte
	if c.Request != nil && c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Loguear el cuerpo para depuración
		log.Printf("Cuerpo original de la solicitud: %s", string(bodyBytes))
	}

	if err := c.ShouldBindJSON(&messageData); err != nil {
		log.Printf("Error al procesar JSON del mensaje: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Formato de mensaje inválido: " + err.Error(),
			"success": false,
		})
		return
	}

	ticketID := messageData.TicketID
	messageContent := messageData.Message

	log.Printf("===== MENSAJE RECIBIDO DEL WIDGET =====")
	log.Printf("Ticket ID: %s", ticketID)
	log.Printf("Contenido: %s", messageContent)

	// Validaciones básicas
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Se requiere el ID del ticket",
			"success": false,
		})
		return
	}

	if messageContent == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "El mensaje no puede estar vacío",
			"success": false,
		})
		return
	}

	// Cargar el ticket existente
	ticket, err := LoadTicket(ticketID)
	if err != nil {
		log.Printf("Error al cargar ticket: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Ticket no encontrado",
			"success": false,
		})
		return
	}

	// Crear un ID único para el mensaje
	messageID := fmt.Sprintf("MSG-%d", time.Now().UnixNano())

	// Obtener información del usuario
	userName, userEmail := GetUserInfo(c, &messageData)
	log.Printf("Usuario: %s (%s)", userName, userEmail)

	// Si no hay nombre o email, intentar usar la información del ticket
	if userName == "" {
		userName = ticket.UserName
		if userName == "" {
			userName = "Cliente"
		}
	}

	if userEmail == "" {
		userEmail = ticket.UserEmail
		if userEmail == "" {
			userEmail = "client@example.com"
		}
	}

	// Crear nueva entrada de mensaje LOCAL
	// IMPORTANTE: Siempre con isClient=true para mensajes del widget
	message := Message{
		ID:        messageID,
		Content:   messageContent,
		IsClient:  true, // FORZAR isClient=true para mensajes del widget
		CreatedAt: time.Now(),
		UserName:  userName,
		UserEmail: userEmail,
	}

	// Añadir mensaje al ticket local
	ticket.Messages = append(ticket.Messages, message)
	ticket.UpdatedAt = time.Now()

	// Guardar ticket actualizado localmente
	if err := SaveTicket(ticket); err != nil {
		log.Printf("Error al guardar ticket localmente: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar mensaje en el ticket", "success": false})
		return
	}

	// Enviar a todos los clientes conectados por WebSocket
	go sendMessageToWebSocketClients(ticketID, message)

	// Enviar el mensaje al sistema GrowDesk en una goroutine separada
	go func() {
		// Configuración del API de GrowDesk
		apiURL := os.Getenv("GROWDESK_API_URL")
		apiKey := os.Getenv("GROWDESK_API_KEY")

		if apiURL == "" {
			apiURL = "http://localhost:8080"
			log.Printf("GROWDESK_API_URL no definido, usando valor por defecto: %s", apiURL)
		}

		if apiKey == "" {
			apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiJhZG1pbi0xMjMiLCJlbWFpbCI6ImFkbWluQGdyb3dkZXNrLmNvbSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTcyNDA4ODQwMH0.8J5ayPvA4B-1vF5NaqRXycW1pmIl9qjKP6Ddj4Ot_Cw" // Token por defecto
			log.Printf("GROWDESK_API_KEY no definido, usando valor por defecto")
		}

		// Obtener información del widget del ticket cargado
		widgetID := ticket.WidgetID
		if widgetID == "" {
			widgetID = c.GetHeader("X-Widget-ID")
		}

		// Preparar el mensaje explícitamente con isClient=true
		growDeskMsg := GrowDeskMessage{
			TicketID:  ticketID,
			Content:   messageContent,
			UserID:    userEmail,
			IsClient:  true, // Esto es CRUCIAL - siempre true para mensajes del widget
			UserName:  userName,
			UserEmail: userEmail,
		}

		// Convertir a JSON
		jsonData, err := json.Marshal(growDeskMsg)
		if err != nil {
			log.Printf("Error al convertir mensaje a JSON: %v", err)
			return
		}

		// Normalizar URL base
		baseURL := strings.TrimSuffix(apiURL, "/")

		// Construir URL para la API de tickets
		url := fmt.Sprintf("%s/widget/tickets/%s/messages?from_client=true", baseURL, ticketID)

		log.Printf("Enviando mensaje al sistema GrowDesk en la URL: %s", url)

		headers := map[string]string{
			"Content-Type":       "application/json",
			"Authorization":      "Bearer " + apiKey,
			"X-Message-Source":   "widget-client",
			"X-Widget-ID":        widgetID, // Usamos el ID real del widget
			"X-Client-Message":   "true",
			"X-Widget-Ticket-ID": ticketID,
			"X-From-Client":      "true",
		}

		// Enviar con reintentos
		resp, err := sendHttpRequestWithRetry(url, jsonData, headers, 3)
		if err != nil {
			log.Printf("Error al enviar mensaje a GrowDesk: %v", err)
		} else if resp != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			log.Printf("Respuesta al enviar mensaje a GrowDesk: Status %d, Body: %s", resp.StatusCode, string(body))
		}
	}()

	// Devolver respuesta exitosa
	c.JSON(http.StatusOK, gin.H{
		"messageId": messageID,
		"message":   "Mensaje enviado correctamente",
		"success":   true,
	})
}
