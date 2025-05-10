package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
)

// TicketRepository maneja las operaciones de base de datos relacionadas con tickets
type TicketRepository struct {
	DB *sql.DB
}

// NewTicketRepository crea una nueva instancia del repositorio de tickets
func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{
		DB: db,
	}
}

// GetAll obtiene todos los tickets de la base de datos
func (r *TicketRepository) GetAll() ([]models.Ticket, error) {
	query := `
		SELECT t.id, t.title, t.subject, t.description, t.status, t.priority, 
		       t.category, t.category_id, t.assigned_to, t.created_by, t.user_id,
		       t.source, t.widget_id, t.department, t.metadata,
		       t.created_at, t.updated_at
		FROM tickets t
		ORDER BY t.created_at DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar tickets: %v", err)
	}
	defer rows.Close()

	tickets := make([]models.Ticket, 0)
	for rows.Next() {
		var ticket models.Ticket
		var categoryID, assignedTo, createdBy, userID, metadataJSON sql.NullString
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&ticket.ID,
			&ticket.Title,
			&ticket.Subject,
			&ticket.Description,
			&ticket.Status,
			&ticket.Priority,
			&ticket.Category,
			&categoryID,
			&assignedTo,
			&createdBy,
			&userID,
			&ticket.Source,
			&ticket.WidgetID,
			&ticket.Department,
			&metadataJSON,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear ticket: %v", err)
		}

		// Asignar valores nulos
		if categoryID.Valid {
			ticket.CategoryID = categoryID.String
		}
		if assignedTo.Valid {
			ticket.AssignedTo = assignedTo.String
		}
		if createdBy.Valid {
			ticket.CreatedBy = createdBy.String
		}
		if userID.Valid {
			ticket.UserID = userID.String
		}

		// Parsear metadata JSON si existe
		if metadataJSON.Valid && metadataJSON.String != "" {
			var metadata models.Metadata
			if err := json.Unmarshal([]byte(metadataJSON.String), &metadata); err == nil {
				ticket.Metadata = &metadata
			}
		}

		ticket.CreatedAt = createdAt
		ticket.UpdatedAt = updatedAt

		// Cargar mensajes del ticket
		messages, err := r.getMessagesForTicket(ticket.ID)
		if err != nil {
			return nil, fmt.Errorf("error al obtener mensajes para ticket %s: %v", ticket.ID, err)
		}
		ticket.Messages = messages

		tickets = append(tickets, ticket)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar tickets: %v", err)
	}

	return tickets, nil
}

// GetByID obtiene un ticket por su ID
func (r *TicketRepository) GetByID(id string) (*models.Ticket, error) {
	query := `
		SELECT t.id, t.title, t.subject, t.description, t.status, t.priority, 
		       t.category, t.category_id, t.assigned_to, t.created_by, t.user_id,
		       t.source, t.widget_id, t.department, t.metadata,
		       t.created_at, t.updated_at
		FROM tickets t
		WHERE t.id = $1
	`

	var ticket models.Ticket
	var categoryID, assignedTo, createdBy, userID, metadataJSON sql.NullString
	var createdAt, updatedAt time.Time

	err := r.DB.QueryRow(query, id).Scan(
		&ticket.ID,
		&ticket.Title,
		&ticket.Subject,
		&ticket.Description,
		&ticket.Status,
		&ticket.Priority,
		&ticket.Category,
		&categoryID,
		&assignedTo,
		&createdBy,
		&userID,
		&ticket.Source,
		&ticket.WidgetID,
		&ticket.Department,
		&metadataJSON,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ticket con ID %s no encontrado", id)
		}
		return nil, fmt.Errorf("error al consultar ticket: %v", err)
	}

	// Asignar valores nulos
	if categoryID.Valid {
		ticket.CategoryID = categoryID.String
	}
	if assignedTo.Valid {
		ticket.AssignedTo = assignedTo.String
	}
	if createdBy.Valid {
		ticket.CreatedBy = createdBy.String
	}
	if userID.Valid {
		ticket.UserID = userID.String
	}

	// Parsear metadata JSON si existe
	if metadataJSON.Valid && metadataJSON.String != "" {
		var metadata models.Metadata
		if err := json.Unmarshal([]byte(metadataJSON.String), &metadata); err == nil {
			ticket.Metadata = &metadata
		}
	}

	ticket.CreatedAt = createdAt
	ticket.UpdatedAt = updatedAt

	// Cargar mensajes del ticket
	messages, err := r.getMessagesForTicket(ticket.ID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener mensajes para ticket %s: %v", ticket.ID, err)
	}
	ticket.Messages = messages

	return &ticket, nil
}

// Create crea un nuevo ticket en la base de datos
func (r *TicketRepository) Create(ticket models.Ticket) (*models.Ticket, error) {
	// Iniciar transacción
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("error al iniciar transacción: %v", err)
	}
	defer tx.Rollback()

	// Generar ID si no existe
	if ticket.ID == "" {
		ticket.ID = fmt.Sprintf("TICKET-%s", time.Now().Format("20060102-150405"))
	}

	// Establecer timestamps si no están definidos
	now := time.Now()
	if ticket.CreatedAt.IsZero() {
		ticket.CreatedAt = now
	}
	if ticket.UpdatedAt.IsZero() {
		ticket.UpdatedAt = now
	}

	// Convertir metadata a JSON si existe
	var metadataJSON sql.NullString
	if ticket.Metadata != nil {
		data, err := json.Marshal(ticket.Metadata)
		if err != nil {
			return nil, fmt.Errorf("error al serializar metadata: %v", err)
		}
		metadataJSON = sql.NullString{String: string(data), Valid: true}
	}

	// Insertar ticket
	query := `
		INSERT INTO tickets (
			id, title, subject, description, status, priority, category, category_id,
			assigned_to, created_by, user_id, source, widget_id, department, metadata,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
		RETURNING id
	`

	err = tx.QueryRow(
		query,
		ticket.ID,
		ticket.Title,
		ticket.Subject,
		ticket.Description,
		ticket.Status,
		ticket.Priority,
		ticket.Category,
		nullString(ticket.CategoryID),
		nullString(ticket.AssignedTo),
		nullString(ticket.CreatedBy),
		nullString(ticket.UserID),
		ticket.Source,
		ticket.WidgetID,
		ticket.Department,
		metadataJSON,
		ticket.CreatedAt,
		ticket.UpdatedAt,
	).Scan(&ticket.ID)

	if err != nil {
		return nil, fmt.Errorf("error al crear ticket: %v", err)
	}

	// Insertar mensajes si existen
	for i, message := range ticket.Messages {
		// Generar ID para mensaje si no existe
		if message.ID == "" {
			message.ID = uuid.New().String()
			ticket.Messages[i].ID = message.ID
		}

		// Establecer timestamps si no están definidos
		if message.CreatedAt.IsZero() {
			message.CreatedAt = now
			ticket.Messages[i].CreatedAt = now
		}
		if message.Timestamp.IsZero() {
			message.Timestamp = now
			ticket.Messages[i].Timestamp = now
		}

		// Insertar mensaje
		messageQuery := `
			INSERT INTO messages (
				id, ticket_id, content, is_client, is_internal, user_id,
				user_name, user_email, timestamp, created_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
			)
		`

		_, err = tx.Exec(
			messageQuery,
			message.ID,
			ticket.ID,
			message.Content,
			message.IsClient,
			message.IsInternal,
			nullString(message.UserID),
			nullString(message.UserName),
			nullString(message.UserEmail),
			message.Timestamp,
			message.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error al crear mensaje para ticket: %v", err)
		}
	}

	// Confirmar transacción
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error al confirmar transacción: %v", err)
	}

	log.Printf("Ticket creado con ID: %s", ticket.ID)
	return &ticket, nil
}

// Update actualiza un ticket existente
func (r *TicketRepository) Update(ticket models.Ticket) error {
	// Iniciar transacción
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %v", err)
	}
	defer tx.Rollback()

	// Actualizar timestamp
	ticket.UpdatedAt = time.Now()

	// Convertir metadata a JSON si existe
	var metadataJSON sql.NullString
	if ticket.Metadata != nil {
		data, err := json.Marshal(ticket.Metadata)
		if err != nil {
			return fmt.Errorf("error al serializar metadata: %v", err)
		}
		metadataJSON = sql.NullString{String: string(data), Valid: true}
	}

	// Actualizar ticket
	query := `
		UPDATE tickets
		SET title = $2, subject = $3, description = $4, status = $5,
		    priority = $6, category = $7, category_id = $8, assigned_to = $9,
		    created_by = $10, user_id = $11, source = $12, widget_id = $13,
		    department = $14, metadata = $15, updated_at = $16
		WHERE id = $1
	`

	result, err := tx.Exec(
		query,
		ticket.ID,
		ticket.Title,
		ticket.Subject,
		ticket.Description,
		ticket.Status,
		ticket.Priority,
		ticket.Category,
		nullString(ticket.CategoryID),
		nullString(ticket.AssignedTo),
		nullString(ticket.CreatedBy),
		nullString(ticket.UserID),
		ticket.Source,
		ticket.WidgetID,
		ticket.Department,
		metadataJSON,
		ticket.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar ticket: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ticket con ID %s no encontrado", ticket.ID)
	}

	// Confirmar transacción
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar transacción: %v", err)
	}

	return nil
}

// Delete elimina un ticket y sus mensajes
func (r *TicketRepository) Delete(id string) error {
	// Iniciar transacción
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %v", err)
	}
	defer tx.Rollback()

	// Eliminar mensajes asociados
	_, err = tx.Exec("DELETE FROM messages WHERE ticket_id = $1", id)
	if err != nil {
		return fmt.Errorf("error al eliminar mensajes del ticket: %v", err)
	}

	// Eliminar metadata del ticket
	_, err = tx.Exec("DELETE FROM ticket_metadata WHERE ticket_id = $1", id)
	if err != nil {
		return fmt.Errorf("error al eliminar metadata del ticket: %v", err)
	}

	// Eliminar ticket
	result, err := tx.Exec("DELETE FROM tickets WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error al eliminar ticket: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ticket con ID %s no encontrado", id)
	}

	// Confirmar transacción
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar transacción: %v", err)
	}

	return nil
}

// AddMessage añade un nuevo mensaje a un ticket existente
func (r *TicketRepository) AddMessage(ticketID string, message models.Message) (*models.Message, error) {
	// Iniciar transacción
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("error al iniciar transacción: %v", err)
	}
	defer tx.Rollback()

	// Verificar que el ticket existe
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM tickets WHERE id = $1)", ticketID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error al verificar ticket: %v", err)
	}

	if !exists {
		return nil, fmt.Errorf("ticket con ID %s no encontrado", ticketID)
	}

	// Generar ID para mensaje si no existe
	if message.ID == "" {
		message.ID = uuid.New().String()
	}

	// Establecer timestamps si no están definidos
	now := time.Now()
	if message.CreatedAt.IsZero() {
		message.CreatedAt = now
	}
	if message.Timestamp.IsZero() {
		message.Timestamp = now
	}

	// Insertar mensaje
	query := `
		INSERT INTO messages (
			id, ticket_id, content, is_client, is_internal, user_id,
			user_name, user_email, timestamp, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
		RETURNING id
	`

	err = tx.QueryRow(
		query,
		message.ID,
		ticketID,
		message.Content,
		message.IsClient,
		message.IsInternal,
		nullString(message.UserID),
		nullString(message.UserName),
		nullString(message.UserEmail),
		message.Timestamp,
		message.CreatedAt,
	).Scan(&message.ID)

	if err != nil {
		return nil, fmt.Errorf("error al crear mensaje: %v", err)
	}

	// Actualizar timestamp del ticket
	_, err = tx.Exec("UPDATE tickets SET updated_at = $1 WHERE id = $2", now, ticketID)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar timestamp del ticket: %v", err)
	}

	// Confirmar transacción
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error al confirmar transacción: %v", err)
	}

	log.Printf("Mensaje añadido con ID: %s", message.ID)
	return &message, nil
}

// getMessagesForTicket obtiene todos los mensajes para un ticket
func (r *TicketRepository) getMessagesForTicket(ticketID string) ([]models.Message, error) {
	query := `
		SELECT id, content, is_client, is_internal, user_id, user_name, user_email,
		       timestamp, created_at
		FROM messages
		WHERE ticket_id = $1
		ORDER BY timestamp ASC
	`

	rows, err := r.DB.Query(query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("error al consultar mensajes: %v", err)
	}
	defer rows.Close()

	messages := make([]models.Message, 0)
	for rows.Next() {
		var message models.Message
		var userID, userName, userEmail sql.NullString
		var timestamp, createdAt time.Time

		err := rows.Scan(
			&message.ID,
			&message.Content,
			&message.IsClient,
			&message.IsInternal,
			&userID,
			&userName,
			&userEmail,
			&timestamp,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear mensaje: %v", err)
		}

		// Asignar valores nulos
		if userID.Valid {
			message.UserID = userID.String
		}
		if userName.Valid {
			message.UserName = userName.String
		}
		if userEmail.Valid {
			message.UserEmail = userEmail.String
		}

		message.Timestamp = timestamp
		message.CreatedAt = createdAt

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar mensajes: %v", err)
	}

	return messages, nil
}

// Helper para convertir string a sql.NullString
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
