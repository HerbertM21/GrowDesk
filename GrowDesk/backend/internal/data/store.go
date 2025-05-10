package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
)

// Store representa el almacén de datos en memoria
type Store struct {
	Tickets    []models.Ticket
	Users      []models.User
	Categories []models.Category
	FAQs       []models.FAQ

	// Conexiones WebSocket por ID de ticket
	// Map de ID de ticket a lista de conexiones
	TicketConnections      map[string][]WebSocketConnection
	AlternateConnectionMap map[string]string // Para compatibilidad de conexiones

	mu sync.RWMutex // Para seguridad de hilos

	// Rutas de archivo para persistencia de datos
	TicketsFile    string
	UsersFile      string
	CategoriesFile string
	FAQsFile       string
}

// WebSocketConnection representa una conexión WebSocket
type WebSocketConnection struct {
	ID          string
	Socket      *websocket.Conn
	ConnectedAt time.Time
}

// NewStore crea un nuevo almacén de datos y carga datos iniciales
func NewStore(dataDir string) *Store {
	// Crear directorio de datos si no existe
	os.MkdirAll(dataDir, 0755)

	store := &Store{
		Tickets:                make([]models.Ticket, 0),
		Users:                  make([]models.User, 0),
		Categories:             make([]models.Category, 0),
		FAQs:                   make([]models.FAQ, 0),
		TicketConnections:      make(map[string][]WebSocketConnection),
		AlternateConnectionMap: make(map[string]string),
		TicketsFile:            filepath.Join(dataDir, "tickets.json"),
		UsersFile:              filepath.Join(dataDir, "users.json"),
		CategoriesFile:         filepath.Join(dataDir, "categories.json"),
		FAQsFile:               filepath.Join(dataDir, "faqs.json"),
	}

	// Cargar datos desde archivos o inicializar con valores por defecto
	store.loadTickets()
	store.loadUsers()
	store.loadCategories()
	store.loadFAQs()

	return store
}

// LoadTickets carga tickets desde archivo o inicializa con lista vacía
func (s *Store) loadTickets() {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.TicketsFile)
	if err != nil {
		fmt.Printf("No se encontró archivo de tickets, iniciando con lista vacía: %v\n", err)
		s.Tickets = make([]models.Ticket, 0)
		return
	}

	if err := json.Unmarshal(data, &s.Tickets); err != nil {
		fmt.Printf("Error al analizar archivo de tickets, iniciando con lista vacía: %v\n", err)
		s.Tickets = make([]models.Ticket, 0)
		return
	}

	fmt.Printf("Cargados %d tickets desde archivo\n", len(s.Tickets))
}

// LoadUsers carga usuarios desde archivo o inicializa con usuarios por defecto
func (s *Store) loadUsers() {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.UsersFile)
	if err != nil {
		fmt.Println("No se encontró archivo de usuarios, inicializando con usuarios por defecto")
		s.initializeDefaultUsers()
		return
	}

	if err := json.Unmarshal(data, &s.Users); err != nil {
		fmt.Printf("Error al analizar archivo de usuarios, inicializando con usuarios por defecto: %v\n", err)
		s.initializeDefaultUsers()
		return
	}

	fmt.Printf("Cargados %d usuarios desde archivo\n", len(s.Users))
}

// LoadCategories carga categorías desde archivo o inicializa con valores por defecto
func (s *Store) loadCategories() {
	// Leemos el archivo fuera del mutex para evitar deadlocks
	data, err := os.ReadFile(s.CategoriesFile)
	if err != nil {
		fmt.Println("No se encontró archivo de categorías, inicializando con categorías por defecto")
		now := time.Now()
		defaultCategories := []models.Category{
			{
				ID:          "1",
				Name:        "Soporte Técnico",
				Description: "Problemas técnicos con hardware o software",
				Color:       "#4CAF50",
				Icon:        "computer",
				Active:      true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "2",
				Name:        "Consultas Generales",
				Description: "Preguntas sobre productos o servicios",
				Color:       "#2196F3",
				Icon:        "help",
				Active:      true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "3",
				Name:        "Facturación",
				Description: "Problemas o dudas sobre facturación",
				Color:       "#FFC107",
				Icon:        "credit_card",
				Active:      true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		// Actualizar categorías con lock
		s.mu.Lock()
		s.Categories = defaultCategories
		s.mu.Unlock()

		// Guardar directamente a archivo sin usar SaveCategories
		dataToSave, err := json.MarshalIndent(defaultCategories, "", "  ")
		if err != nil {
			fmt.Printf("Error al serializar categorías: %v\n", err)
			return
		}

		if err := os.WriteFile(s.CategoriesFile, dataToSave, 0644); err != nil {
			fmt.Printf("Error al escribir archivo de categorías: %v\n", err)
			return
		}

		fmt.Printf("Guardadas %d categorías por defecto en archivo\n", len(defaultCategories))
		return
	}

	// Si el archivo existe, lo deserializamos
	var categories []models.Category
	if err := json.Unmarshal(data, &categories); err != nil {
		fmt.Printf("Error al analizar archivo de categorías, inicializando con valores por defecto: %v\n", err)
		// Llamamos recursivamente para inicializar con valores por defecto
		// Esta llamada es segura porque el archivo no se podrá leer en la próxima iteración
		s.loadCategories()
		return
	}

	// Actualizamos el store con lock
	s.mu.Lock()
	s.Categories = categories
	s.mu.Unlock()

	fmt.Printf("Cargados %d categorías desde archivo\n", len(categories))
}

// LoadFAQs carga FAQs desde archivo o inicializa con valores por defecto
func (s *Store) loadFAQs() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.FAQsFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Si el archivo no existe, inicializar con FAQs por defecto
			s.FAQs = []models.FAQ{
				{
					ID:          1,
					Question:    "¿Cómo puedo crear un ticket?",
					Answer:      "Para crear un ticket, haz clic en el botón 'Nuevo Ticket' y completa el formulario.",
					IsPublished: true,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					ID:          2,
					Question:    "¿Cómo puedo ver mis tickets?",
					Answer:      "Puedes ver tus tickets en la sección 'Mis Tickets' del panel de control.",
					IsPublished: true,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					ID:          3,
					Question:    "¿Cómo puedo actualizar mi perfil?",
					Answer:      "Para actualizar tu perfil, ve a la sección 'Configuración' y haz clic en 'Editar Perfil'.",
					IsPublished: true,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			}
			// Guardar las FAQs por defecto
			return s.SaveFAQs()
		}
		return fmt.Errorf("error al leer archivo de FAQs: %v", err)
	}

	var faqs []models.FAQ
	if err := json.Unmarshal(data, &faqs); err != nil {
		return fmt.Errorf("error al analizar archivo de FAQs: %v", err)
	}

	s.FAQs = faqs
	return nil
}

// SaveTickets guarda tickets en archivo
func (s *Store) SaveTickets() error {
	s.mu.Lock()
	tickets := make([]models.Ticket, len(s.Tickets))
	copy(tickets, s.Tickets)
	s.mu.Unlock()

	data, err := json.MarshalIndent(tickets, "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar tickets: %w", err)
	}

	if err := os.WriteFile(s.TicketsFile, data, 0644); err != nil {
		return fmt.Errorf("error al escribir archivo de tickets: %w", err)
	}

	fmt.Printf("Guardados %d tickets en archivo\n", len(tickets))
	return nil
}

// SaveUsers guarda usuarios en archivo
func (s *Store) SaveUsers() error {
	s.mu.Lock()
	users := make([]models.User, len(s.Users))
	copy(users, s.Users)
	s.mu.Unlock()

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar usuarios: %w", err)
	}

	if err := os.WriteFile(s.UsersFile, data, 0644); err != nil {
		return fmt.Errorf("error al escribir archivo de usuarios: %w", err)
	}

	fmt.Printf("Guardados %d usuarios en archivo\n", len(users))
	return nil
}

// SaveCategories guarda categorías en archivo
func (s *Store) SaveCategories() error {
	s.mu.Lock()
	categories := make([]models.Category, len(s.Categories))
	copy(categories, s.Categories)
	s.mu.Unlock()

	data, err := json.MarshalIndent(categories, "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar categorías: %w", err)
	}

	if err := os.WriteFile(s.CategoriesFile, data, 0644); err != nil {
		return fmt.Errorf("error al escribir archivo de categorías: %w", err)
	}

	fmt.Printf("Guardados %d categorías en archivo\n", len(categories))
	return nil
}

// SaveFAQs guarda FAQs en archivo
func (s *Store) SaveFAQs() error {
	s.mu.Lock()
	faqsCopy := make([]models.FAQ, len(s.FAQs))
	copy(faqsCopy, s.FAQs)
	s.mu.Unlock()

	// Asegurarse de que el directorio existe
	if err := os.MkdirAll(filepath.Dir(s.FAQsFile), 0755); err != nil {
		return fmt.Errorf("Error al crear directorio: %v", err)
	}

	data, err := json.MarshalIndent(faqsCopy, "", "  ")
	if err != nil {
		return fmt.Errorf("Error al serializar FAQs: %v", err)
	}

	if err := os.WriteFile(s.FAQsFile, data, 0644); err != nil {
		return fmt.Errorf("Error al escribir archivo de FAQs: %v", err)
	}

	return nil
}

// InitializeDefaultUsers inicializa el almacén con usuarios por defecto
func (s *Store) initializeDefaultUsers() {
	s.Users = []models.User{
		{
			ID:         "1",
			FirstName:  "Admin",
			LastName:   "Usuario",
			Email:      "admin@example.com",
			Role:       "admin",
			Department: "Tecnología",
			Active:     true,
			Password:   "password", // Sería hasheado en una implementación real
		},
		{
			ID:         "2",
			FirstName:  "Asistente",
			LastName:   "Soporte",
			Email:      "asistente@example.com",
			Role:       "assistant",
			Department: "Soporte",
			Active:     true,
			Password:   "password",
		},
		{
			ID:         "3",
			FirstName:  "Empleado",
			LastName:   "Regular",
			Email:      "empleado@example.com",
			Role:       "employee",
			Department: "Ventas",
			Active:     true,
			Password:   "password",
		},
	}

	// Guardar en archivo
	s.SaveUsers()
}

// InitializeDefaultFAQs inicializa el almacén con FAQs por defecto
func (s *Store) initializeDefaultFAQs() {
	now := time.Now()
	s.FAQs = []models.FAQ{
		{
			ID:          1,
			Question:    "¿Cómo puedo restablecer mi contraseña?",
			Answer:      "Para restablecer su contraseña, haga clic en el enlace \"¿Olvidó su contraseña?\" en la pantalla de inicio de sesión y siga las instrucciones enviadas a su correo electrónico.",
			Category:    "Cuenta",
			IsPublished: true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          2,
			Question:    "¿Cómo puedo actualizar mi información de contacto?",
			Answer:      "Inicie sesión en su cuenta, vaya a \"Mi Perfil\" y haga clic en \"Editar Información\". Allí podrá actualizar su dirección de correo electrónico, número de teléfono y dirección postal.",
			Category:    "Cuenta",
			IsPublished: true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          3,
			Question:    "¿Cómo puedo reportar un problema técnico?",
			Answer:      "Para reportar un problema técnico, vaya a la sección \"Soporte\", haga clic en \"Crear Ticket\" y complete el formulario con los detalles del problema. Un técnico se pondrá en contacto con usted lo antes posible.",
			Category:    "Soporte Técnico",
			IsPublished: true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          4,
			Question:    "¿Cuáles son los horarios de atención al cliente?",
			Answer:      "Nuestro equipo de atención al cliente está disponible de lunes a viernes, de 9:00 a.m. a 6:00 p.m. (hora local). Para asistencia de emergencia fuera de este horario, por favor envíe un correo electrónico a soporte@ejemplo.com.",
			Category:    "General",
			IsPublished: false,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          5,
			Question:    "¿Cómo funciona el sistema de tickets?",
			Answer:      "El sistema de tickets le permite enviar consultas o reportar problemas. Cada ticket recibe un número único de seguimiento. Puede ver el estado de sus tickets en cualquier momento iniciando sesión en su cuenta y accediendo a la sección \"Mis Tickets\".",
			Category:    "Soporte Técnico",
			IsPublished: true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          6,
			Question:    "¿Puedo cambiar el idioma de la plataforma?",
			Answer:      "Sí, puede cambiar el idioma de la plataforma accediendo a \"Configuración\" desde su perfil. Actualmente ofrecemos soporte para inglés y español.",
			Category:    "Configuración",
			IsPublished: true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	// Guardar en archivo
	s.SaveFAQs()
}

// AddTicket agrega un nuevo ticket al almacén
func (s *Store) AddTicket(ticket models.Ticket) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ticket.ID == "" {
		ticket.ID = fmt.Sprintf("TICKET-%s", time.Now().Format("20060102-150405"))
	}

	s.Tickets = append(s.Tickets, ticket)
	s.SaveTickets()
}

// GetTicket recupera un ticket por ID
func (s *Store) GetTicket(id string) (*models.Ticket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i, ticket := range s.Tickets {
		if ticket.ID == id {
			return &s.Tickets[i], nil
		}
	}

	return nil, fmt.Errorf("Ticket no encontrado: %s", id)
}

// UpdateTicket actualiza un ticket existente
func (s *Store) UpdateTicket(ticket models.Ticket) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existingTicket := range s.Tickets {
		if existingTicket.ID == ticket.ID {
			// Preservar mensajes existentes si no se proporcionan
			if len(ticket.Messages) == 0 {
				ticket.Messages = existingTicket.Messages
			}

			// Asegurar que la fecha de actualización se establece
			ticket.UpdatedAt = time.Now()

			// Actualizar el ticket
			s.Tickets[i] = ticket
			return s.SaveTickets()
		}
	}

	return fmt.Errorf("Ticket no encontrado: %s", ticket.ID)
}

// AddMessageToTicket agrega un mensaje a un ticket
func (s *Store) AddMessageToTicket(ticketID string, message models.Message) (*models.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, ticket := range s.Tickets {
		if ticket.ID == ticketID {
			if message.ID == "" {
				message.ID = fmt.Sprintf("MSG-%s", uuid.New().String())
			}
			if message.Timestamp.IsZero() {
				message.Timestamp = time.Now()
			}
			if message.CreatedAt.IsZero() {
				message.CreatedAt = time.Now()
			}

			s.Tickets[i].Messages = append(s.Tickets[i].Messages, message)
			s.Tickets[i].UpdatedAt = time.Now()
			s.SaveTickets()
			return &message, nil
		}
	}

	return nil, fmt.Errorf("Ticket no encontrado: %s", ticketID)
}

// GetAllFAQs devuelve todas las FAQs
func (s *Store) GetAllFAQs() []models.FAQ {
	s.mu.RLock()
	faqs := s.FAQs
	s.mu.RUnlock()

	// Si no hay FAQs, intentar recargar del archivo
	if len(faqs) == 0 {
		// Liberar el mutex antes de recargar
		s.loadFAQs()

		// Obtener las FAQs actualizadas
		s.mu.RLock()
		faqs = s.FAQs
		s.mu.RUnlock()

		// Si aún no hay FAQs, inicializar con valores por defecto
		if len(faqs) == 0 {
			s.mu.Lock()
			s.initializeDefaultFAQs()
			faqs = s.FAQs
			s.mu.Unlock()
		}
	}

	// Crear una copia para evitar problemas de concurrencia
	result := make([]models.FAQ, len(faqs))
	copy(result, faqs)
	return result
}

// GetPublishedFAQs devuelve solo las FAQs publicadas
func (s *Store) GetPublishedFAQs() []models.FAQ {
	s.mu.RLock()
	if len(s.FAQs) == 0 {
		s.mu.RUnlock()
		if err := s.loadFAQs(); err != nil {
			fmt.Printf("Error al cargar FAQs: %v\n", err)
			return []models.FAQ{}
		}
		s.mu.RLock()
	}

	// Crear una copia de las FAQs publicadas
	publishedFAQs := make([]models.FAQ, 0)
	for _, faq := range s.FAQs {
		if faq.IsPublished {
			publishedFAQs = append(publishedFAQs, faq)
		}
	}
	s.mu.RUnlock()

	return publishedFAQs
}

// AddWSConnection agrega una conexión WebSocket para un ticket
func (s *Store) AddWSConnection(ticketID string, conn *websocket.Conn) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Inicializar el slice de conexiones si no existe
	if _, exists := s.TicketConnections[ticketID]; !exists {
		s.TicketConnections[ticketID] = make([]WebSocketConnection, 0)
	}

	// Generar un ID de conexión único
	connectionID := fmt.Sprintf("%d", time.Now().UnixNano())

	// Agregar la conexión
	s.TicketConnections[ticketID] = append(s.TicketConnections[ticketID], WebSocketConnection{
		ID:          connectionID,
		Socket:      conn,
		ConnectedAt: time.Now(),
	})

	fmt.Printf("Añadido WebSocket conexión %s para ticket %s\n", connectionID, ticketID)
	return connectionID
}

// RemoveWSConnection elimina una conexión WebSocket
func (s *Store) RemoveWSConnection(ticketID, connectionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if connections, exists := s.TicketConnections[ticketID]; exists {
		// Filtrar la conexión con el ID dado
		updatedConnections := make([]WebSocketConnection, 0)
		for _, conn := range connections {
			if conn.ID != connectionID {
				updatedConnections = append(updatedConnections, conn)
			}
		}

		if len(updatedConnections) > 0 {
			s.TicketConnections[ticketID] = updatedConnections
		} else {
			// Si no hay conexiones, eliminar la entrada
			delete(s.TicketConnections, ticketID)
		}

		fmt.Printf("Eliminada WebSocket conexión %s para ticket %s\n", connectionID, ticketID)
	}
}

// BroadcastMessage envía un mensaje a todos los clientes conectados para un ticket
func (s *Store) BroadcastMessage(ticketID string, message models.Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Printf("Enviando mensaje a clientes para ticket %s\n", ticketID)

	// Preparar el mensaje para WebSocket
	wsMessage := models.WebSocketMessage{
		Type:     "new_message",
		TicketID: ticketID,
		Data: map[string]interface{}{
			"id":        message.ID,
			"content":   message.Content,
			"isClient":  message.IsClient,
			"timestamp": message.Timestamp,
			"userName":  message.UserName,
		},
	}

	// Serializar el mensaje
	messageData, err := json.Marshal(wsMessage)
	if err != nil {
		fmt.Printf("Error al serializar mensaje WebSocket: %v\n", err)
		return
	}

	connectionFound := false

	// Enviar a conexiones directas
	if connections, exists := s.TicketConnections[ticketID]; exists && len(connections) > 0 {
		fmt.Printf("Encontrados %d conexiones directas para ticket %s\n", len(connections), ticketID)
		sendCount := 0

		for _, conn := range connections {
			if conn.Socket != nil {
				if err := conn.Socket.WriteMessage(websocket.TextMessage, messageData); err != nil {
					fmt.Printf("Error al enviar mensaje WebSocket: %v\n", err)
				} else {
					sendCount++
				}
			}
		}

		fmt.Printf("Enviado mensaje a %d/%d conexiones\n", sendCount, len(connections))
		connectionFound = true
	} else {
		fmt.Printf("No se encontraron conexiones directas para ticket %s\n", ticketID)
	}

	// Comprobar conexiones alternativas
	if altTicketID, exists := s.AlternateConnectionMap[ticketID]; exists {
		if connections, exists := s.TicketConnections[altTicketID]; exists && len(connections) > 0 {
			fmt.Printf("Encontrados %d conexiones alternativas para ticket %s a través de %s\n",
				len(connections), ticketID, altTicketID)

			sendCount := 0
			for _, conn := range connections {
				if conn.Socket != nil {
					if err := conn.Socket.WriteMessage(websocket.TextMessage, messageData); err != nil {
						fmt.Printf("Error al enviar mensaje WebSocket: %v\n", err)
					} else {
						sendCount++
					}
				}
			}

			fmt.Printf("Enviado mensaje a %d/%d conexiones alternativas\n", sendCount, len(connections))
			connectionFound = true
		}
	}

	if !connectionFound {
		fmt.Printf("No se encontraron conexiones para ticket %s\n", ticketID)
	}
}

// CreateFAQ crea una nueva FAQ
func (s *Store) createFAQInternal(faq *models.FAQ) (*models.FAQ, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generar nuevo ID
	newID := 1
	if len(s.FAQs) > 0 {
		newID = s.FAQs[len(s.FAQs)-1].ID + 1
	}

	// Establecer marcas de tiempo
	now := time.Now()
	faq.ID = newID
	faq.CreatedAt = now
	faq.UpdatedAt = now

	// Agregar al almacén
	s.FAQs = append(s.FAQs, *faq)

	// Guardar en archivo
	if err := s.SaveFAQs(); err != nil {
		return nil, fmt.Errorf("error al guardar FAQs: %v", err)
	}

	return faq, nil
}

// CreateFAQ crea una nueva FAQ (implementación para la interfaz)
func (s *Store) CreateFAQ(faq models.FAQ) error {
	_, err := s.createFAQInternal(&faq)
	return err
}

// GetUsers devuelve todos los usuarios
func (s *Store) GetUsers() ([]models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Crear una copia para evitar problemas de concurrencia
	usersCopy := make([]models.User, len(s.Users))
	copy(usersCopy, s.Users)

	return usersCopy, nil
}

// GetUser obtiene un usuario por ID
func (s *Store) GetUser(id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.Users {
		if user.ID == id {
			// Crear una copia para evitar problemas de concurrencia
			userCopy := user
			return &userCopy, nil
		}
	}

	return nil, fmt.Errorf("usuario con ID %s no encontrado", id)
}

// AddUser agrega un nuevo usuario
func (s *Store) AddUser(user models.User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Users = append(s.Users, user)
}

// UpdateUser actualiza un usuario existente
func (s *Store) UpdateUser(user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existingUser := range s.Users {
		if existingUser.ID == user.ID {
			// Update user fields
			user.UpdatedAt = time.Now()
			s.Users[i] = user
			return s.SaveUsers()
		}
	}

	return fmt.Errorf("usuario con ID %s no encontrado", user.ID)
}

// DeleteUser elimina un usuario por ID
func (s *Store) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, user := range s.Users {
		if user.ID == id {
			// Eliminar usuario
			s.Users = append(s.Users[:i], s.Users[i+1:]...)
			return s.SaveUsers()
		}
	}

	return fmt.Errorf("usuario con ID %s no encontrado", id)
}

// GetUserByEmail busca un usuario por email
func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.Users {
		if user.Email == email {
			// Crear una copia para evitar problemas de concurrencia
			userCopy := user
			return &userCopy, nil
		}
	}

	return nil, fmt.Errorf("user with email %s not found", email)
}

// CreateUser agrega un nuevo usuario
func (s *Store) CreateUser(user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Asegurar que el usuario tiene un ID
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Establecer marcas de tiempo si no están ya establecidas
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}

	s.Users = append(s.Users, user)
	return s.SaveUsers()
}

// GetTickets devuelve todos los tickets
func (s *Store) GetTickets() ([]models.Ticket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Crear una copia para evitar problemas de concurrencia
	tickets := make([]models.Ticket, len(s.Tickets))
	copy(tickets, s.Tickets)

	return tickets, nil
}

// CreateTicket agrega un nuevo ticket
func (s *Store) CreateTicket(ticket models.Ticket) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Asegurar que el ticket tiene un ID
	if ticket.ID == "" {
		ticket.ID = uuid.New().String()
	}

	// Establecer marcas de tiempo si no están ya establecidas
	if ticket.CreatedAt.IsZero() {
		ticket.CreatedAt = time.Now()
	}
	if ticket.UpdatedAt.IsZero() {
		ticket.UpdatedAt = time.Now()
	}

	s.Tickets = append(s.Tickets, ticket)
	return s.SaveTickets()
}

// DeleteTicket elimina un ticket por ID
func (s *Store) DeleteTicket(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Encontrar el índice del ticket
	for i, ticket := range s.Tickets {
		if ticket.ID == id {
			// Eliminar ticket
			s.Tickets = append(s.Tickets[:i], s.Tickets[i+1:]...)
			return s.SaveTickets()
		}
	}

	return fmt.Errorf("ticket with ID %s not found", id)
}

// AddTicketMessage agrega un mensaje a un ticket
func (s *Store) AddTicketMessage(ticketID string, message models.Message) error {
	_, err := s.AddMessageToTicket(ticketID, message)
	return err
}

// GetCategories devuelve todas las categorías
func (s *Store) GetCategories() ([]models.Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Crear una copia para evitar problemas de concurrencia
	categoriesCopy := make([]models.Category, len(s.Categories))
	copy(categoriesCopy, s.Categories)

	return categoriesCopy, nil
}

// GetCategory obtiene una categoría por ID
func (s *Store) GetCategory(id string) (*models.Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, category := range s.Categories {
		if category.ID == id {
			// Crear una copia para evitar problemas de concurrencia
			categoryCopy := category
			return &categoryCopy, nil
		}
	}

	return nil, fmt.Errorf("categoría con ID %s no encontrada", id)
}

// CreateCategory crea una nueva categoría
func (s *Store) CreateCategory(category models.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Asegurar que la categoría tiene un ID
	if category.ID == "" {
		category.ID = uuid.New().String()
	}

	// Establecer marcas de tiempo si no están ya establecidas
	now := time.Now()
	if category.CreatedAt.IsZero() {
		category.CreatedAt = now
	}
	if category.UpdatedAt.IsZero() {
		category.UpdatedAt = now
	}

	s.Categories = append(s.Categories, category)
	return s.SaveCategories()
}

// UpdateCategory actualiza una categoría existente
func (s *Store) UpdateCategory(category models.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existingCategory := range s.Categories {
		if existingCategory.ID == category.ID {
			// Actualizar marca de tiempo
			category.UpdatedAt = time.Now()
			s.Categories[i] = category
			return s.SaveCategories()
		}
	}

	return fmt.Errorf("categoría con ID %s no encontrada", category.ID)
}

// DeleteCategory elimina una categoría por ID
func (s *Store) DeleteCategory(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, category := range s.Categories {
		if category.ID == id {
			// Eliminar categoría
			s.Categories = append(s.Categories[:i], s.Categories[i+1:]...)
			return s.SaveCategories()
		}
	}

	return fmt.Errorf("categoría con ID %s no encontrada", id)
}

// GetFAQs devuelve todas las FAQs
func (s *Store) GetFAQs() ([]models.FAQ, error) {
	faqs := s.GetAllFAQs()
	return faqs, nil
}

// GetFAQsByStatus devuelve FAQs filtradas por estado de publicación
func (s *Store) GetFAQsByStatus(published bool) ([]models.FAQ, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Crear una copia filtrada
	filteredFAQs := make([]models.FAQ, 0)
	for _, faq := range s.FAQs {
		if faq.IsPublished == published {
			filteredFAQs = append(filteredFAQs, faq)
		}
	}

	return filteredFAQs, nil
}

// GetFAQ obtiene una FAQ por ID
func (s *Store) GetFAQ(id int) (*models.FAQ, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, faq := range s.FAQs {
		if faq.ID == id {
			// Crear una copia para evitar problemas de concurrencia
			faqCopy := faq
			return &faqCopy, nil
		}
	}

	return nil, fmt.Errorf("FAQ con ID %d no encontrada", id)
}

// UpdateFAQ actualiza una FAQ existente
func (s *Store) UpdateFAQ(faq models.FAQ) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existingFAQ := range s.FAQs {
		if existingFAQ.ID == faq.ID {
			// Actualizar marca de tiempo
			faq.UpdatedAt = time.Now()
			s.FAQs[i] = faq
			return s.SaveFAQs()
		}
	}

	return fmt.Errorf("FAQ con ID %d no encontrada", faq.ID)
}

// DeleteFAQ elimina una FAQ por ID
func (s *Store) DeleteFAQ(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, faq := range s.FAQs {
		if faq.ID == id {
			// Eliminar FAQ
			s.FAQs = append(s.FAQs[:i], s.FAQs[i+1:]...)
			return s.SaveFAQs()
		}
	}

	return fmt.Errorf("FAQ con ID %d no encontrada", id)
}

// ToggleFAQPublish cambia el estado de publicación de una FAQ
func (s *Store) ToggleFAQPublish(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, faq := range s.FAQs {
		if faq.ID == id {
			// Cambiar estado de publicación
			s.FAQs[i].IsPublished = !s.FAQs[i].IsPublished
			s.FAQs[i].UpdatedAt = time.Now()
			return s.SaveFAQs()
		}
	}

	return fmt.Errorf("FAQ con ID %d no encontrada", id)
}
