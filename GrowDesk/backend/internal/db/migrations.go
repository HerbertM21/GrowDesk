package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
)

// MigrateUsersFromJSON migra los usuarios desde el archivo JSON a la base de datos
func MigrateUsersFromJSON(db *sql.DB, jsonPath string) error {
	// Leer el archivo JSON
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("error al leer archivo de usuarios: %v", err)
	}

	// Decodificar JSON a slice de usuarios
	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return fmt.Errorf("error al decodificar JSON de usuarios: %v", err)
	}

	// Iniciar una transacción
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Insertar usuarios en la base de datos
	stmt, err := tx.Prepare(`
		INSERT INTO users (
			id, email, first_name, last_name, role, department, active, 
			password, created_at, updated_at, position, phone, language
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) ON CONFLICT (email) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("error al preparar statement: %v", err)
	}
	defer stmt.Close()

	for _, user := range users {
		// Generar UUID si no existe
		if user.ID == "" {
			user.ID = uuid.New().String()
		}

		_, err := stmt.Exec(
			user.ID,
			user.Email,
			user.FirstName,
			user.LastName,
			user.Role,
			user.Department,
			user.Active,
			user.Password,
			user.CreatedAt,
			user.UpdatedAt,
			user.Position,
			user.Phone,
			user.Language,
		)
		if err != nil {
			return fmt.Errorf("error al insertar usuario %s: %v", user.Email, err)
		}
	}

	// Commit de la transacción
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar transacción: %v", err)
	}

	log.Printf("Migrados %d usuarios desde JSON a PostgreSQL", len(users))
	return nil
}

// MigrateTicketsFromJSON migra los tickets desde el archivo JSON a la base de datos
func MigrateTicketsFromJSON(db *sql.DB, jsonPath string) error {
	// Leer el archivo JSON
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("error al leer archivo de tickets: %v", err)
	}

	// Decodificar JSON a slice de tickets
	var tickets []models.Ticket
	if err := json.Unmarshal(data, &tickets); err != nil {
		return fmt.Errorf("error al decodificar JSON de tickets: %v", err)
	}

	// Iniciar una transacción
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Insertar tickets en la base de datos
	ticketStmt, err := tx.Prepare(`
		INSERT INTO tickets (
			id, title, subject, description, status, priority, created_at, updated_at,
			category_id, assigned_to, created_by, customer_name, customer_email, source, 
			widget_id, department
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		) ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("error al preparar statement de tickets: %v", err)
	}
	defer ticketStmt.Close()

	// Preparar statement para mensajes
	messageStmt, err := tx.Prepare(`
		INSERT INTO messages (
			id, ticket_id, content, is_client, is_internal, created_at, user_id, user_name, user_email
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("error al preparar statement de mensajes: %v", err)
	}
	defer messageStmt.Close()

	for _, ticket := range tickets {
		// Asegurarse de que el ID exista
		if ticket.ID == "" {
			ticket.ID = uuid.New().String()
		}

		// Corregir el valor de source si no está definido
		if ticket.Source == "" {
			ticket.Source = "manual" // valor por defecto
		}

		// Ejecutar la inserción del ticket
		_, err := ticketStmt.Exec(
			ticket.ID,
			ticket.Title,
			ticket.Subject,
			ticket.Description,
			ticket.Status,
			ticket.Priority,
			ticket.CreatedAt,
			ticket.UpdatedAt,
			stringOrNil(ticket.Category),
			stringOrNil(ticket.AssignedTo),
			stringOrNil(ticket.CreatedBy),
			ticket.Customer.Name,
			ticket.Customer.Email,
			ticket.Source,
			stringOrNil(ticket.WidgetID),
			stringOrNil(ticket.Department),
		)
		if err != nil {
			return fmt.Errorf("error al insertar ticket %s: %v", ticket.ID, err)
		}

		// Insertar mensajes del ticket
		for _, message := range ticket.Messages {
			// Asegurarse de que el ID del mensaje exista
			if message.ID == "" {
				message.ID = uuid.New().String()
			}

			// Usar timestamp como createdAt si está disponible
			createdAt := message.CreatedAt
			if createdAt.IsZero() {
				createdAt = message.Timestamp
			}

			_, err := messageStmt.Exec(
				message.ID,
				ticket.ID,
				message.Content,
				message.IsClient,
				message.IsInternal,
				createdAt,
				stringOrNil(message.UserID),
				stringOrNil(message.UserName),
				stringOrNil(message.UserEmail),
			)
			if err != nil {
				return fmt.Errorf("error al insertar mensaje %s: %v", message.ID, err)
			}
		}
	}

	// Commit de la transacción
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar transacción: %v", err)
	}

	log.Printf("Migrados %d tickets desde JSON a PostgreSQL", len(tickets))
	return nil
}

// MigrateCategoriesFromJSON migra las categorías desde el archivo JSON a la base de datos
func MigrateCategoriesFromJSON(db *sql.DB, jsonPath string) error {
	// Leer el archivo JSON
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("error al leer archivo de categorías: %v", err)
	}

	// Decodificar JSON a slice de categorías
	var categories []models.Category
	if err := json.Unmarshal(data, &categories); err != nil {
		return fmt.Errorf("error al decodificar JSON de categorías: %v", err)
	}

	// Iniciar una transacción
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Insertar categorías en la base de datos
	stmt, err := tx.Prepare(`
		INSERT INTO categories (
			id, name, description, color, icon, active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("error al preparar statement: %v", err)
	}
	defer stmt.Close()

	for _, category := range categories {
		// Generar UUID si no existe
		if category.ID == "" {
			category.ID = uuid.New().String()
		}

		_, err := stmt.Exec(
			category.ID,
			category.Name,
			category.Description,
			category.Color,
			category.Icon,
			category.Active,
			category.CreatedAt,
			category.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("error al insertar categoría %s: %v", category.Name, err)
		}
	}

	// Commit de la transacción
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar transacción: %v", err)
	}

	log.Printf("Migradas %d categorías desde JSON a PostgreSQL", len(categories))
	return nil
}

// MigrateFAQsFromJSON migra las FAQs desde el archivo JSON a la base de datos
func MigrateFAQsFromJSON(db *sql.DB, jsonPath string) error {
	// Leer el archivo JSON
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("error al leer archivo de FAQs: %v", err)
	}

	// Decodificar JSON a slice de FAQs
	var faqs []models.FAQ
	if err := json.Unmarshal(data, &faqs); err != nil {
		return fmt.Errorf("error al decodificar JSON de FAQs: %v", err)
	}

	// Iniciar una transacción
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Insertar FAQs en la base de datos
	stmt, err := tx.Prepare(`
		INSERT INTO faqs (
			question, answer, category, is_published, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id
	`)
	if err != nil {
		return fmt.Errorf("error al preparar statement: %v", err)
	}
	defer stmt.Close()

	for _, faq := range faqs {
		var id int
		err := stmt.QueryRow(
			faq.Question,
			faq.Answer,
			faq.Category,
			faq.IsPublished,
			faq.CreatedAt,
			faq.UpdatedAt,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("error al insertar FAQ %s: %v", faq.Question, err)
		}
	}

	// Commit de la transacción
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar transacción: %v", err)
	}

	log.Printf("Migradas %d FAQs desde JSON a PostgreSQL", len(faqs))
	return nil
}

// MigrateAllFromJSON migra todos los datos desde los archivos JSON a PostgreSQL
func MigrateAllFromJSON(db *sql.DB, dataDir string) error {
	// Verificar que el directorio existe
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return fmt.Errorf("el directorio de datos %s no existe", dataDir)
	}

	// Migrar usuarios
	usersPath := filepath.Join(dataDir, "users.json")
	if fileExists(usersPath) {
		if err := MigrateUsersFromJSON(db, usersPath); err != nil {
			log.Printf("ADVERTENCIA: Error al migrar usuarios: %v", err)
		}
	}

	// Migrar categorías
	categoriesPath := filepath.Join(dataDir, "categories.json")
	if fileExists(categoriesPath) {
		if err := MigrateCategoriesFromJSON(db, categoriesPath); err != nil {
			log.Printf("ADVERTENCIA: Error al migrar categorías: %v", err)
		}
	}

	// Migrar FAQs
	faqsPath := filepath.Join(dataDir, "faqs.json")
	if fileExists(faqsPath) {
		if err := MigrateFAQsFromJSON(db, faqsPath); err != nil {
			log.Printf("ADVERTENCIA: Error al migrar FAQs: %v", err)
		}
	}

	// Migrar tickets (debe hacerse después de las categorías y usuarios)
	ticketsPath := filepath.Join(dataDir, "tickets.json")
	if fileExists(ticketsPath) {
		if err := MigrateTicketsFromJSON(db, ticketsPath); err != nil {
			log.Printf("ADVERTENCIA: Error al migrar tickets: %v", err)
		}
	}

	// Migrar tickets de widget si existen (en el directorio data)
	widgetDataDir := filepath.Join(dataDir, "..", "..", "GrowDesk-Widget", "widget-api", "data")
	if dirExists(widgetDataDir) {
		log.Printf("Iniciando migración de tickets de widget desde %s...", widgetDataDir)
		MigrateWidgetTickets(db, widgetDataDir)
	}

	log.Println("Migración de datos desde JSON a PostgreSQL completada")
	return nil
}

// MigrateWidgetTickets migra los tickets del widget desde los archivos JSON a la base de datos
func MigrateWidgetTickets(db *sql.DB, widgetDataDir string) error {
	// Verificar que el directorio existe
	if _, err := os.Stat(widgetDataDir); os.IsNotExist(err) {
		return fmt.Errorf("el directorio de datos de widget %s no existe", widgetDataDir)
	}

	// Listar archivos JSON en el directorio
	files, err := filepath.Glob(filepath.Join(widgetDataDir, "ticket_*.json"))
	if err != nil {
		return fmt.Errorf("error al listar archivos de tickets de widget: %v", err)
	}

	log.Printf("Encontrados %d archivos de tickets de widget para migrar", len(files))

	// Iniciar una transacción
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Preparar statement para widget_tickets
	ticketStmt, err := tx.Prepare(`
		INSERT INTO widget_tickets (
			id, ticket_id, title, subject, description, status, priority,
			created_at, updated_at, client_name, client_email, widget_id, department, synced
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		) ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("error al preparar statement de widget_tickets: %v", err)
	}
	defer ticketStmt.Close()

	// Preparar statement para widget_messages
	messageStmt, err := tx.Prepare(`
		INSERT INTO widget_messages (
			id, widget_ticket_id, content, is_client, created_at, user_name, user_email, synced
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("error al preparar statement de widget_messages: %v", err)
	}
	defer messageStmt.Close()

	// Migrar cada archivo de ticket
	migratedCount := 0
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Error al leer archivo %s: %v", file, err)
			continue
		}

		var ticket map[string]interface{}
		if err := json.Unmarshal(data, &ticket); err != nil {
			log.Printf("Error al decodificar JSON del archivo %s: %v", file, err)
			continue
		}

		ticketID := getStringOrEmpty(ticket, "id")
		if ticketID == "" {
			// Extraer ID del nombre del archivo
			ticketID = strings.TrimPrefix(filepath.Base(file), "ticket_")
			ticketID = strings.TrimSuffix(ticketID, ".json")
		}

		// Generar un UUID para el ID de la tabla
		id := uuid.New().String()

		// Insertar el ticket en widget_tickets
		_, err = ticketStmt.Exec(
			id,
			ticketID,
			getStringOrEmpty(ticket, "title"),
			getStringOrEmpty(ticket, "subject"),
			getStringOrEmpty(ticket, "description"),
			getStringOrEmpty(ticket, "status"),
			getStringOrEmpty(ticket, "priority"),
			getTimeOrNow(ticket, "createdAt"),
			getTimeOrNow(ticket, "updatedAt"),
			getStringOrEmpty(ticket, "clientName"),
			getStringOrEmpty(ticket, "clientEmail"),
			getStringOrEmpty(ticket, "widgetId"),
			getStringOrEmpty(ticket, "department"),
			false, // synced = false
		)
		if err != nil {
			log.Printf("Error al insertar ticket de widget %s: %v", ticketID, err)
			continue
		}

		// Insertar mensajes del ticket
		messages, ok := ticket["messages"].([]interface{})
		if ok {
			for _, msgData := range messages {
				msg, ok := msgData.(map[string]interface{})
				if !ok {
					continue
				}

				msgID := getStringOrEmpty(msg, "id")
				if msgID == "" {
					msgID = uuid.New().String()
				}

				_, err = messageStmt.Exec(
					msgID,
					id,
					getStringOrEmpty(msg, "content"),
					getBoolOrDefault(msg, "isClient", true),
					getTimeOrNow(msg, "createdAt"),
					getStringOrEmpty(msg, "userName"),
					getStringOrEmpty(msg, "userEmail"),
					false, // synced = false
				)
				if err != nil {
					log.Printf("Error al insertar mensaje de widget %s: %v", msgID, err)
				}
			}
		}

		migratedCount++
	}

	// Commit de la transacción
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar transacción: %v", err)
	}

	log.Printf("Migrados %d tickets de widget desde JSON a PostgreSQL", migratedCount)
	return nil
}

// Funciones auxiliares

// stringOrNil devuelve el string como interface{} o nil si está vacío
func stringOrNil(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// fileExists comprueba si un archivo existe
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// dirExists comprueba si un directorio existe
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// getStringOrEmpty obtiene un valor string de un map o devuelve string vacío
func getStringOrEmpty(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

// getBoolOrDefault obtiene un valor bool de un map o devuelve el valor por defecto
func getBoolOrDefault(data map[string]interface{}, key string, defaultVal bool) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	return defaultVal
}

// getTimeOrNow intenta obtener un tiempo de un map o devuelve el tiempo actual
func getTimeOrNow(data map[string]interface{}, key string) interface{} {
	if val, ok := data[key].(string); ok && val != "" {
		return val // PostgreSQL puede parsear strings de tiempo ISO8601
	}
	return "NOW()" // Función SQL para tiempo actual
}
