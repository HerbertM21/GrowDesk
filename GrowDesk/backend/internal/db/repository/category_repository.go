package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
)

// CategoryRepository maneja las operaciones de base de datos para las categorías
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository crea un nuevo repositorio de categorías
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetAll obtiene todas las categorías
func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	query := `
		SELECT id, name, description, color, icon, active, created_at, updated_at
		FROM categories
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar categorías: %v", err)
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var category models.Category
		var description, color, icon sql.NullString
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&description,
			&color,
			&icon,
			&category.Active,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear categoría: %v", err)
		}

		// Asignar valores nulos
		if description.Valid {
			category.Description = description.String
		}
		if color.Valid {
			category.Color = color.String
		}
		if icon.Valid {
			category.Icon = icon.String
		}

		category.CreatedAt = createdAt
		category.UpdatedAt = updatedAt

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar categorías: %v", err)
	}

	return categories, nil
}

// GetByID obtiene una categoría por su ID
func (r *CategoryRepository) GetByID(id string) (*models.Category, error) {
	query := `
		SELECT id, name, description, color, icon, active, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	var category models.Category
	var description, color, icon sql.NullString
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&description,
		&color,
		&icon,
		&category.Active,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("categoría con ID %s no encontrada", id)
		}
		return nil, fmt.Errorf("error al consultar categoría: %v", err)
	}

	// Asignar valores nulos
	if description.Valid {
		category.Description = description.String
	}
	if color.Valid {
		category.Color = color.String
	}
	if icon.Valid {
		category.Icon = icon.String
	}

	category.CreatedAt = createdAt
	category.UpdatedAt = updatedAt

	return &category, nil
}

// Create crea una nueva categoría
func (r *CategoryRepository) Create(category models.Category) (*models.Category, error) {
	// Generar ID si no existe
	if category.ID == "" {
		category.ID = fmt.Sprintf("CAT-%s", time.Now().Format("20060102-150405"))
	}

	// Establecer timestamps si no están definidos
	now := time.Now()
	if category.CreatedAt.IsZero() {
		category.CreatedAt = now
	}
	if category.UpdatedAt.IsZero() {
		category.UpdatedAt = now
	}

	query := `
		INSERT INTO categories (
			id, name, description, color, icon, active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		category.ID,
		category.Name,
		nullString(category.Description),
		nullString(category.Color),
		nullString(category.Icon),
		category.Active,
		category.CreatedAt,
		category.UpdatedAt,
	).Scan(&category.ID)

	if err != nil {
		return nil, fmt.Errorf("error al crear categoría: %v", err)
	}

	return &category, nil
}

// Update actualiza una categoría existente
func (r *CategoryRepository) Update(category models.Category) error {
	// Actualizar timestamp
	category.UpdatedAt = time.Now()

	query := `
		UPDATE categories
		SET name = $2, description = $3, color = $4, icon = $5,
		    active = $6, updated_at = $7
		WHERE id = $1
	`

	result, err := r.db.Exec(
		query,
		category.ID,
		category.Name,
		nullString(category.Description),
		nullString(category.Color),
		nullString(category.Icon),
		category.Active,
		category.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error al actualizar categoría: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("categoría con ID %s no encontrada", category.ID)
	}

	return nil
}

// Delete elimina una categoría por su ID
func (r *CategoryRepository) Delete(id string) error {
	query := "DELETE FROM categories WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar categoría: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("categoría con ID %s no encontrada", id)
	}

	return nil
}
