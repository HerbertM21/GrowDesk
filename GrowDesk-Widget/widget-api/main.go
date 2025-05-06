package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// Ticket represents a support ticket
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

// Message represents a message in a ticket
type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	IsClient  bool      `json:"isClient"`
	CreatedAt time.Time `json:"createdAt"`
	UserName  string    `json:"userName"`
	UserEmail string    `json:"userEmail"`
}

// Metadata contains additional information
type Metadata struct {
	URL        string `json:"url"`
	Referrer   string `json:"referrer"`
	UserAgent  string `json:"userAgent"`
	ScreenSize string `json:"screenSize"`
	ExternalID string `json:"externalId"`
}

// TicketRequest is used to create a new ticket
type TicketRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Message  string   `json:"message"`
	Metadata Metadata `json:"metadata"`
}

// TicketResponse is the server response after creating a ticket
type TicketResponse struct {
	TicketID string `json:"ticketId"`
	Message  string `json:"message"`
}

// MessageRequest is used to send a message to a ticket
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

// GetUserInfo extracts user information from headers or request body
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
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Widget-ID, X-Widget-Token, X-User-Name, X-User-Email, X-Source, X-Client-Created, X-Widget-Ticket-ID, X-Message-Source, X-From-Client, X-Client-Message")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

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
	// Intentar cargar FAQs desde alguna fuente (archivo, base de datos, etc.)
	// Por ahora devolvemos FAQs de ejemplo
	faqs := []gin.H{
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
	}

	c.JSON(http.StatusOK, faqs)
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
	// URL base donde se aloja el JS del widget (puerto 3030)
	baseUrl := "http://localhost:3030"

	// URL de la API del widget (puerto 8082)
	apiUrl := "http://localhost:8082"

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

// getMessages obtiene todos los mensajes para un ticket
func getMessages(c *gin.Context) {
	ticketID := c.Param("ticketId")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is required"})
		return
	}

	log.Printf("Obteniendo mensajes para ticket: %s", ticketID)

	// Intentar obtener mensajes directamente desde GrowDesk primero
	apiURL := os.Getenv("GROWDESK_API_URL")
	apiKey := os.Getenv("GROWDESK_API_KEY")

	if apiURL != "" && apiKey != "" {
		log.Printf("Intentando obtener mensajes desde GrowDesk para ticket: %s", ticketID)

		// Solicitar datos actualizados del ticket a GrowDesk
		updatedTicket, err := getTicketFromGrowDesk(ticketID)
		if err == nil {
			// Si obtenemos correctamente el ticket actualizado, usar sus mensajes
			log.Printf("Ticket actualizado correctamente desde GrowDesk, usando mensajes actualizados")

			// Guardar el ticket actualizado localmente
			if err := SaveTicket(updatedTicket); err != nil {
				log.Printf("Error al guardar el ticket actualizado: %v", err)
			}

			// Devolver mensajes actualizados
			c.JSON(http.StatusOK, gin.H{
				"ticketId": updatedTicket.ID,
				"messages": updatedTicket.Messages,
			})
			return
		} else {
			log.Printf("No se pudo obtener ticket desde GrowDesk: %v. Intentando cargar localmente.", err)
		}
	}

	// Si no se puede obtener desde GrowDesk, cargar ticket localmente
	ticket, err := LoadTicket(ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Return messages from local storage
	c.JSON(http.StatusOK, gin.H{
		"ticketId": ticket.ID,
		"messages": ticket.Messages,
	})
}

// getTicketFromGrowDesk obtiene un ticket desde el sistema principal GrowDesk
func getTicketFromGrowDesk(ticketID string) (Ticket, error) {
	var ticket Ticket

	// Verificar que tenemos la URL y API key
	apiURL := os.Getenv("GROWDESK_API_URL")
	apiKey := os.Getenv("GROWDESK_API_KEY")

	if apiURL == "" {
		apiURL = "http://172.21.0.5:8080"
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

	var ticketData TicketData
	if err := c.ShouldBindJSON(&ticketData); err != nil {
		log.Printf("Error en datos de entrada: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
			apiURL = "http://172.21.0.5:8080"
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
		url := fmt.Sprintf("%s/api/tickets", baseURL)

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
	c.JSON(http.StatusCreated, gin.H{
		"ticketId": ticketID,
		"message":  "Ticket creado correctamente",
	})

	log.Printf("===== FIN CREACIÓN TICKET WIDGET =====")
}

// sendHttpRequestWithRetry envía una solicitud HTTP con reintentos
func sendHttpRequestWithRetry(url string, jsonData []byte, headers map[string]string, maxRetries int) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < maxRetries; i++ {
		// Crear solicitud
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error al crear solicitud HTTP (intento %d): %v", i+1, err)
			continue
		}

		// Añadir cabeceras
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		// Crear cliente con timeout
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		// Enviar solicitud
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Éxito
			return resp, nil
		}

		if err != nil {
			log.Printf("Error en intento %d: %v", i+1, err)
		} else {
			log.Printf("Respuesta no exitosa en intento %d: %d", i+1, resp.StatusCode)
			resp.Body.Close()
		}

		// Esperar antes de reintentar (backoff exponencial)
		time.Sleep(time.Duration(300*(i+1)) * time.Millisecond)
	}

	return nil, fmt.Errorf("fallo después de %d intentos: %v", maxRetries, err)
}

// handleWebSocketConnection maneja las conexiones WebSocket para chat en tiempo real
func handleWebSocketConnection(c *gin.Context) {
	ticketId := c.Param("ticketId")
	if ticketId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ticket no proporcionado"})
		return
	}

	log.Printf("Intentando establecer conexión WebSocket para ticket: %s", ticketId)
	log.Printf("Headers de la solicitud: %v", c.Request.Header)

	// Verificar que el ticket existe antes de permitir una conexión WebSocket
	_, err := LoadTicket(ticketId)
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
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error al leer mensaje WebSocket: %v", err)
			break // Salir si hay un error (cliente desconectado)
		}

		// Log del mensaje recibido
		log.Printf("Mensaje recibido de tipo %d: %s", messageType, string(message))

		// Procesamiento del mensaje si es necesario
		// proximamente
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
		return
	}

	// IMPORTANTE: Asegurarse de que el mensaje tiene la estructura esperada
	// Crear un mapa explícito con los campos exactos que espera el cliente
	messageObj := map[string]interface{}{
		"id":        message.ID,
		"content":   message.Content,
		"isClient":  message.IsClient, // Importante: mantener el nombre exacto del campo
		"createdAt": message.CreatedAt.Format(time.RFC3339),
		"timestamp": message.CreatedAt.Format(time.RFC3339), // Agregar timestamp para compatibilidad
		"userName":  message.UserName,
		"userEmail": message.UserEmail,
	}

	// Usar la estructura que espera GrowDesk (data en lugar de message)
	wsMessage := map[string]interface{}{
		"type":     "new_message",
		"data":     messageObj,
		"ticketId": ticketId,
	}

	// Serializar a JSON
	msgBytes, err := json.Marshal(wsMessage)
	if err != nil {
		log.Printf("Error al serializar mensaje WebSocket: %v", err)
		return
	}

	log.Printf("Enviando mensaje a los clientes del ticket: %s", ticketId)
	log.Printf("Estructura del mensaje WebSocket: %+v", wsMessage)

	// Imprimir JSON para depuración
	prettyJson, _ := json.MarshalIndent(wsMessage, "", "  ")
	log.Printf("JSON del mensaje a enviar: %s", string(prettyJson))

	// Enviar a todas las conexiones
	for _, conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
			log.Printf("Error al enviar mensaje WebSocket: %v", err)
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

	// Crear nuevo mensaje (desde agente, no cliente)
	newMessage := Message{
		ID:        uuid.New().String(),
		Content:   req.Content,
		IsClient:  false, // Mensaje de agente - EXPLÍCITAMENTE FALSE
		CreatedAt: time.Now(),
		UserName:  req.AgentName, // Nombre del agente
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
		TicketID: req.TicketID,
		Content:  req.Content,
		IsClient: false, // FALSE para mensajes de agente
		UserID:   req.UserID,
	}

	// Serializar los datos a JSON
	jsonData, err := json.Marshal(growDeskMessage)
	if err != nil {
		log.Printf("Error al serializar mensaje de agente para GrowDesk: %v", err)
	} else {
		// Enviar mensaje al sistema GrowDesk
		url := fmt.Sprintf("http://172.21.0.5:8080/api/tickets/%s/messages", req.TicketID)
		go sendToGrowDesk(url, jsonData, "", req.TicketID)
	}

	// Devolver respuesta de éxito
	c.JSON(http.StatusOK, gin.H{
		"messageId": newMessage.ID,
		"message":   "Agent message sent successfully",
	})
}

// sendMessage adds a message to an existing ticket
func sendMessage(c *gin.Context) {
	var messageData MessageData

	if err := c.ShouldBindJSON(&messageData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticketID := messageData.TicketID
	messageContent := messageData.Message

	log.Printf("===== MENSAJE RECIBIDO DEL WIDGET =====")
	log.Printf("Ticket ID: %s", ticketID)
	log.Printf("Contenido: %s", messageContent)

	// Cargar el ticket existente
	ticket, err := LoadTicket(ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
		return
	}

	// Crear un ID único para el mensaje
	messageID := fmt.Sprintf("MSG-%d", time.Now().UnixNano())

	// Obtener información del usuario
	userName, userEmail := GetUserInfo(c, &messageData)
	log.Printf("Usuario: %s (%s)", userName, userEmail)

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar mensaje en el ticket"})
		return
	}

	// Enviar el mensaje al sistema GrowDesk en una goroutine separada
	go func() {
		// Configuración del API de GrowDesk
		apiURL := os.Getenv("GROWDESK_API_URL")
		apiKey := os.Getenv("GROWDESK_API_KEY")

		if apiURL == "" {
			apiURL = "http://172.21.0.5:8080"
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

		// Construir URL correcta para enviar mensaje
		url := fmt.Sprintf("%s/api/tickets/%s/messages?from_client=true", baseURL, ticketID)

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

		// Enviar el mensaje con reintentos
		resp, err := sendHttpRequestWithRetry(url, jsonData, headers, 3)
		if err == nil && resp != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			log.Printf("Mensaje enviado exitosamente al sistema GrowDesk. Respuesta: %s", string(body))
		} else {
			log.Printf("Error al enviar mensaje al sistema GrowDesk: %v", err)
		}
	}()

	// Enviar mensaje a clientes WebSocket
	sendMessageToWebSocketClients(ticketID, message)

	// Responder al cliente
	c.JSON(http.StatusOK, gin.H{
		"message":   "Mensaje añadido correctamente",
		"messageId": messageID,
	})

	log.Printf("===== FIN MENSAJE WIDGET =====")
}
