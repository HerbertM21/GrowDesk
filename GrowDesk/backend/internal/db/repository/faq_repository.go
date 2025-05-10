package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
)

// FAQRepository maneja las operaciones de base de datos para las FAQs
type FAQRepository struct {
	db *sql.DB
}

// NewFAQRepository crea un nuevo repositorio de FAQs
func NewFAQRepository(db *sql.DB) *FAQRepository {
	return &FAQRepository{db: db}
}

// GetAll obtiene todas las FAQs
func (r *FAQRepository) GetAll() ([]models.FAQ, error) {
	query := `
		SELECT id, question, answer, category, is_published, created_at, updated_at
		FROM faqs
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar FAQs: %v", err)
	}
	defer rows.Close()

	faqs := make([]models.FAQ, 0)
	for rows.Next() {
		var faq models.FAQ
		var category sql.NullString
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&faq.ID,
			&faq.Question,
			&faq.Answer,
			&category,
			&faq.IsPublished,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear FAQ: %v", err)
		}

		// Asignar valores nulos
		if category.Valid {
			faq.Category = category.String
		}

		faq.CreatedAt = createdAt
		faq.UpdatedAt = updatedAt

		faqs = append(faqs, faq)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar FAQs: %v", err)
	}

	return faqs, nil
}

// GetByStatus obtiene FAQs filtradas por estado de publicación
func (r *FAQRepository) GetByStatus(published bool) ([]models.FAQ, error) {
	query := `
		SELECT id, question, answer, category, is_published, created_at, updated_at
		FROM faqs
		WHERE is_published = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, published)
	if err != nil {
		return nil, fmt.Errorf("error al consultar FAQs por estado: %v", err)
	}
	defer rows.Close()

	faqs := make([]models.FAQ, 0)
	for rows.Next() {
		var faq models.FAQ
		var category sql.NullString
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&faq.ID,
			&faq.Question,
			&faq.Answer,
			&category,
			&faq.IsPublished,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear FAQ: %v", err)
		}

		// Asignar valores nulos
		if category.Valid {
			faq.Category = category.String
		}

		faq.CreatedAt = createdAt
		faq.UpdatedAt = updatedAt

		faqs = append(faqs, faq)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar FAQs: %v", err)
	}

	return faqs, nil
}

// GetByID obtiene una FAQ por su ID
func (r *FAQRepository) GetByID(id int) (*models.FAQ, error) {
	query := `
		SELECT id, question, answer, category, is_published, created_at, updated_at
		FROM faqs
		WHERE id = $1
	`

	var faq models.FAQ
	var category sql.NullString
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(query, id).Scan(
		&faq.ID,
		&faq.Question,
		&faq.Answer,
		&category,
		&faq.IsPublished,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("FAQ con ID %d no encontrada", id)
		}
		return nil, fmt.Errorf("error al consultar FAQ: %v", err)
	}

	// Asignar valores nulos
	if category.Valid {
		faq.Category = category.String
	}

	faq.CreatedAt = createdAt
	faq.UpdatedAt = updatedAt

	return &faq, nil
}

// Create crea una nueva FAQ
func (r *FAQRepository) Create(faq models.FAQ) (*models.FAQ, error) {
	// Establecer timestamps si no están definidos
	now := time.Now()
	if faq.CreatedAt.IsZero() {
		faq.CreatedAt = now
	}
	if faq.UpdatedAt.IsZero() {
		faq.UpdatedAt = now
	}

	query := `
		INSERT INTO faqs (
			question, answer, category, is_published, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		faq.Question,
		faq.Answer,
		nullString(faq.Category),
		faq.IsPublished,
		faq.CreatedAt,
		faq.UpdatedAt,
	).Scan(&faq.ID)

	if err != nil {
		return nil, fmt.Errorf("error al crear FAQ: %v", err)
	}

	return &faq, nil
}

// Update actualiza una FAQ existente
func (r *FAQRepository) Update(faq models.FAQ) error {
	// Actualizar timestamp
	faq.UpdatedAt = time.Now()

	query := `
		UPDATE faqs
		SET question = $2, answer = $3, category = $4, 
		    is_published = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.Exec(
		query,
		faq.ID,
		faq.Question,
		faq.Answer,
		nullString(faq.Category),
		faq.IsPublished,
		faq.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar FAQ: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("FAQ con ID %d no encontrada", faq.ID)
	}

	return nil
}

// Delete elimina una FAQ por su ID
func (r *FAQRepository) Delete(id int) error {
	query := "DELETE FROM faqs WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar FAQ: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("FAQ con ID %d no encontrada", id)
	}

	return nil
}

// TogglePublish alterna el estado de publicación de una FAQ
func (r *FAQRepository) TogglePublish(id int) error {
	query := `
		UPDATE faqs
		SET is_published = NOT is_published, updated_at = $2
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id, time.Now())
	if err != nil {
		return fmt.Errorf("error al alternar estado de publicación: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("FAQ con ID %d no encontrada", id)
	}

	return nil
}
