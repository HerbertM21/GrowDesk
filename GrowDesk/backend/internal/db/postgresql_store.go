package db

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/data"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/db/repository"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
)

// PostgreSQLStore implementa la interfaz DataStore con PostgreSQL
type PostgreSQLStore struct {
	db             *sql.DB
	userRepo       *repository.UserRepository
	ticketRepo     *repository.TicketRepository
	categoryRepo   *repository.CategoryRepository
	faqRepo        *repository.FAQRepository
	wsConnections  map[string]map[string]*websocket.Conn
	wsConnectionMu sync.Mutex
}

// NewPostgreSQLStore crea una nueva instancia de PostgreSQLStore
func NewPostgreSQLStore(db *sql.DB) *PostgreSQLStore {
	return &PostgreSQLStore{
		db:            db,
		userRepo:      repository.NewUserRepository(db),
		ticketRepo:    repository.NewTicketRepository(db),
		categoryRepo:  repository.NewCategoryRepository(db),
		faqRepo:       repository.NewFAQRepository(db),
		wsConnections: make(map[string]map[string]*websocket.Conn),
	}
}

// Implementación de métodos para usuarios
func (s *PostgreSQLStore) GetUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

func (s *PostgreSQLStore) GetUser(id string) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *PostgreSQLStore) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *PostgreSQLStore) CreateUser(user models.User) error {
	_, err := s.userRepo.Create(user)
	return err
}

func (s *PostgreSQLStore) UpdateUser(user models.User) error {
	return s.userRepo.Update(user)
}

func (s *PostgreSQLStore) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}

// Implementación de métodos para tickets
func (s *PostgreSQLStore) GetTickets() ([]models.Ticket, error) {
	return s.ticketRepo.GetAll()
}

func (s *PostgreSQLStore) GetTicket(id string) (*models.Ticket, error) {
	return s.ticketRepo.GetByID(id)
}

func (s *PostgreSQLStore) CreateTicket(ticket models.Ticket) error {
	_, err := s.ticketRepo.Create(ticket)
	return err
}

func (s *PostgreSQLStore) UpdateTicket(ticket models.Ticket) error {
	return s.ticketRepo.Update(ticket)
}

func (s *PostgreSQLStore) DeleteTicket(id string) error {
	return s.ticketRepo.Delete(id)
}

func (s *PostgreSQLStore) AddTicketMessage(ticketID string, message models.Message) error {
	_, err := s.ticketRepo.AddMessage(ticketID, message)
	return err
}

// Implementación de métodos para categorías
func (s *PostgreSQLStore) GetCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

func (s *PostgreSQLStore) GetCategory(id string) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *PostgreSQLStore) CreateCategory(category models.Category) error {
	_, err := s.categoryRepo.Create(category)
	return err
}

func (s *PostgreSQLStore) UpdateCategory(category models.Category) error {
	return s.categoryRepo.Update(category)
}

func (s *PostgreSQLStore) DeleteCategory(id string) error {
	return s.categoryRepo.Delete(id)
}

// Implementación de métodos para FAQs
func (s *PostgreSQLStore) GetFAQs() ([]models.FAQ, error) {
	return s.faqRepo.GetAll()
}

func (s *PostgreSQLStore) GetFAQsByStatus(published bool) ([]models.FAQ, error) {
	return s.faqRepo.GetByStatus(published)
}

func (s *PostgreSQLStore) GetFAQ(id int) (*models.FAQ, error) {
	return s.faqRepo.GetByID(id)
}

func (s *PostgreSQLStore) CreateFAQ(faq models.FAQ) error {
	_, err := s.faqRepo.Create(faq)
	return err
}

func (s *PostgreSQLStore) UpdateFAQ(faq models.FAQ) error {
	return s.faqRepo.Update(faq)
}

func (s *PostgreSQLStore) DeleteFAQ(id int) error {
	return s.faqRepo.Delete(id)
}

func (s *PostgreSQLStore) ToggleFAQPublish(id int) error {
	return s.faqRepo.TogglePublish(id)
}

// Implementación de métodos para WebSocket
func (s *PostgreSQLStore) AddWSConnection(ticketID string, conn *websocket.Conn) string {
	s.wsConnectionMu.Lock()
	defer s.wsConnectionMu.Unlock()

	connectionID := fmt.Sprintf("conn-%d", len(s.wsConnections))
	if _, exists := s.wsConnections[ticketID]; !exists {
		s.wsConnections[ticketID] = make(map[string]*websocket.Conn)
	}
	s.wsConnections[ticketID][connectionID] = conn
	return connectionID
}

func (s *PostgreSQLStore) RemoveWSConnection(ticketID, connectionID string) {
	s.wsConnectionMu.Lock()
	defer s.wsConnectionMu.Unlock()

	if conns, exists := s.wsConnections[ticketID]; exists {
		delete(conns, connectionID)
		if len(conns) == 0 {
			delete(s.wsConnections, ticketID)
		}
	}
}

func (s *PostgreSQLStore) BroadcastMessage(ticketID string, message models.Message) {
	s.wsConnectionMu.Lock()
	defer s.wsConnectionMu.Unlock()

	if conns, exists := s.wsConnections[ticketID]; exists {
		for _, conn := range conns {
			err := conn.WriteJSON(message)
			if err != nil {
				// Si hay error al enviar, simplemente registramos y continuamos
				fmt.Printf("Error al enviar mensaje por WebSocket: %v\n", err)
			}
		}
	}
}

// Verifica que PostgreSQLStore implementa la interfaz DataStore
var _ data.DataStore = (*PostgreSQLStore)(nil)
