package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Cargar ticket
func LoadTicket(ticketID string) (Ticket, error) {
	filename := fmt.Sprintf("data/ticket_%s.json", ticketID)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		// Si el archivo no existe y tiene formato GrowDesk, intentar buscar en el sistema principal
		if os.IsNotExist(err) && strings.HasPrefix(ticketID, "TICKET-") {
			log.Printf("Ticket %s no encontrado localmente, intentando buscar en GrowDesk...", ticketID)
			return GetTicketFromGrowDesk(ticketID)
		}
		return Ticket{}, err
	}

	var ticket Ticket
	err = json.Unmarshal(data, &ticket)
	return ticket, err
}

// GetTicketFromGrowDesk busca un ticket en el sistema principal GrowDesk
func GetTicketFromGrowDesk(ticketID string) (Ticket, error) {
	log.Printf("Buscando ticket %s en GrowDesk...", ticketID)

	// Configurar URL y clave API
	growDeskApiUrl := os.Getenv("GROWDESK_API_URL")
	if growDeskApiUrl == "" {
		growDeskApiUrl = "http://localhost:8000/api"
	}

	apiKey := os.Getenv("GROWDESK_API_KEY")
	if apiKey == "" {
		apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiJhZG1pbi0xMjMiLCJlbWFpbCI6ImFkbWluQGdyb3dkZXNrLmNvbSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTcyNDA4ODQwMH0.8J5ayPvA4B-1vF5NaqRXycW1pmIl9qjKP6Ddj4Ot_Cw"
	}

	// Construir URL completa
	url := fmt.Sprintf("%s/tickets/%s", growDeskApiUrl, ticketID)

	// Crear cliente HTTP con timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Crear solicitud
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error al crear solicitud HTTP para GrowDesk: %v", err)
		return Ticket{}, err
	}

	// Establecer headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Ejecutar solicitud
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error al solicitar ticket a GrowDesk: %v", err)
		return Ticket{}, err
	}
	defer resp.Body.Close()

	// Verificar código de respuesta
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error de GrowDesk: código %d", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Detalle del error: %s", string(body))
		return Ticket{}, fmt.Errorf("error al obtener ticket de GrowDesk: código %d", resp.StatusCode)
	}

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error al leer respuesta de GrowDesk: %v", err)
		return Ticket{}, err
	}

	// Parsear respuesta JSON
	var growDeskTicket map[string]interface{}
	if err := json.Unmarshal(body, &growDeskTicket); err != nil {
		log.Printf("Error al parsear respuesta JSON de GrowDesk: %v", err)
		return Ticket{}, err
	}

	// Convertir a nuestro formato de ticket
	ticket := Ticket{
		ID:          ticketID,
		Title:       fmt.Sprintf("%v", growDeskTicket["title"]),
		Description: fmt.Sprintf("%v", growDeskTicket["description"]),
		Status:      fmt.Sprintf("%v", growDeskTicket["status"]),
		CreatedBy:   fmt.Sprintf("%v", growDeskTicket["createdBy"]),
		Messages:    []Message{},
	}

	// Intentar parsear timestamps
	if createdAt, ok := growDeskTicket["createdAt"].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339, createdAt)
		if err == nil {
			ticket.CreatedAt = parsedTime
		} else {
			ticket.CreatedAt = time.Now()
		}
	} else {
		ticket.CreatedAt = time.Now()
	}

	if updatedAt, ok := growDeskTicket["updatedAt"].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339, updatedAt)
		if err == nil {
			ticket.UpdatedAt = parsedTime
		} else {
			ticket.UpdatedAt = time.Now()
		}
	} else {
		ticket.UpdatedAt = time.Now()
	}

	// Obtener información del cliente
	if customer, ok := growDeskTicket["customer"].(map[string]interface{}); ok {
		if name, ok := customer["name"].(string); ok {
			ticket.UserName = name
		}
		if email, ok := customer["email"].(string); ok {
			ticket.UserEmail = email
		}
	}

	// Obtener mensajes
	if messages, ok := growDeskTicket["messages"].([]interface{}); ok {
		for _, msgInterface := range messages {
			if msg, ok := msgInterface.(map[string]interface{}); ok {
				message := Message{
					ID:      fmt.Sprintf("%v", msg["id"]),
					Content: fmt.Sprintf("%v", msg["content"]),
				}

				// Determinar si es mensaje de cliente
				if isClient, ok := msg["isClient"].(bool); ok {
					message.IsClient = isClient
				}

				// Parsear timestamp
				if timestamp, ok := msg["timestamp"].(string); ok {
					parsedTime, err := time.Parse(time.RFC3339, timestamp)
					if err == nil {
						message.CreatedAt = parsedTime
					} else {
						message.CreatedAt = time.Now()
					}
				} else if createdAt, ok := msg["createdAt"].(string); ok {
					parsedTime, err := time.Parse(time.RFC3339, createdAt)
					if err == nil {
						message.CreatedAt = parsedTime
					} else {
						message.CreatedAt = time.Now()
					}
				} else {
					message.CreatedAt = time.Now()
				}

				// Agregar mensaje a la lista
				ticket.Messages = append(ticket.Messages, message)
			}
		}
	}

	// Guardar ticket localmente para futuras consultas
	if err := SaveTicket(ticket); err != nil {
		log.Printf("Advertencia: No se pudo guardar el ticket localmente: %v", err)
	} else {
		log.Printf("Ticket %s obtenido de GrowDesk y guardado localmente", ticketID)
	}

	return ticket, nil
}

// SyncMessagesFromGrowDesk obtiene los mensajes más recientes de un ticket desde GrowDesk
// y los sincroniza con la copia local
func SyncMessagesFromGrowDesk(ticketID string) error {
	log.Printf("Sincronizando mensajes del ticket %s desde GrowDesk...", ticketID)

	// Verificar que el ticket existe localmente primero
	localTicket, err := LoadTicket(ticketID)
	if err != nil {
		return fmt.Errorf("error al cargar ticket local: %v", err)
	}

	// Obtener el ticket completo desde GrowDesk
	remoteTicket, err := GetTicketFromGrowDesk(ticketID)
	if err != nil {
		return fmt.Errorf("error al obtener ticket de GrowDesk: %v", err)
	}

	// Verificar si hay mensajes nuevos para sincronizar
	if len(remoteTicket.Messages) <= len(localTicket.Messages) {
		// No hay mensajes nuevos
		log.Printf("No hay mensajes nuevos para sincronizar en el ticket %s", ticketID)
		return nil
	}

	// Actualizar el ticket local con los nuevos mensajes
	localTicket.Messages = remoteTicket.Messages
	localTicket.UpdatedAt = time.Now()

	// Guardar el ticket actualizado
	if err := SaveTicket(localTicket); err != nil {
		return fmt.Errorf("error al guardar ticket local actualizado: %v", err)
	}

	log.Printf("Ticket %s sincronizado correctamente, %d mensajes", ticketID, len(localTicket.Messages))

	// Notificar a los clientes conectados vía WebSocket
	for _, message := range remoteTicket.Messages {
		// Solo notificar de mensajes recientes (últimos 5 minutos)
		if time.Since(message.CreatedAt) < 5*time.Minute {
			sendMessageToWebSocketClients(ticketID, message)
		}
	}

	return nil
}

// StartSyncWorker inicia un worker en segundo plano que sincroniza
// periódicamente los tickets activos con GrowDesk
func StartSyncWorker() {
	// Ejecutar en una goroutine separada
	go func() {
		syncInterval := 30 * time.Second
		ticker := time.NewTicker(syncInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Obtener lista de tickets activos
				tickets, err := listActiveTickets()
				if err != nil {
					log.Printf("Error al listar tickets activos: %v", err)
					continue
				}

				// Sincronizar cada ticket
				for _, ticket := range tickets {
					// Solo sincronizar tickets con formato GrowDesk (TICKET-*)
					if strings.HasPrefix(ticket.ID, "TICKET-") {
						if err := SyncMessagesFromGrowDesk(ticket.ID); err != nil {
							log.Printf("Error al sincronizar mensajes del ticket %s: %v", ticket.ID, err)
						}
					}
				}
			}
		}
	}()

	log.Println("Worker de sincronización iniciado")
}

// listActiveTickets obtiene la lista de tickets activos en el almacenamiento local
func listActiveTickets() ([]Ticket, error) {
	// Verificar que el directorio de datos existe
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		return []Ticket{}, nil
	}

	// Leer todos los archivos de ticket
	files, err := ioutil.ReadDir("data")
	if err != nil {
		return nil, err
	}

	var tickets []Ticket
	for _, file := range files {
		// Solo procesar archivos de tickets
		if !file.IsDir() && strings.HasPrefix(file.Name(), "ticket_") {
			// Extraer ID del ticket del nombre del archivo
			ticketID := strings.TrimPrefix(file.Name(), "ticket_")
			ticketID = strings.TrimSuffix(ticketID, ".json")

			// Cargar ticket
			ticket, err := LoadTicket(ticketID)
			if err != nil {
				log.Printf("Error al cargar ticket %s: %v", ticketID, err)
				continue
			}

			// Solo incluir tickets activos (no cerrados)
			if ticket.Status != "closed" {
				tickets = append(tickets, ticket)
			}
		}
	}

	return tickets, nil
}

// Función principal actualizada para iniciar el worker de sincronización
func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("Archivo .env no encontrado. Usando variables de entorno del sistema.")
	}

	// Configurar modo de Gin según entorno
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Iniciar worker de sincronización
	StartSyncWorker()

	// Inicializar router
	router := gin.Default()

	// Resto del código de main()...
}
