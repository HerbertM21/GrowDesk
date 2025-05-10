package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/hmdev/GrowDeskV2/GrowDesk/backend/internal/models"
)

// UserRepository maneja las operaciones de base de datos para los usuarios
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository crea un nuevo repositorio de usuarios
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetAll obtiene todos los usuarios
func (r *UserRepository) GetAll() ([]models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, role, department, active, 
		       position, phone, language, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar usuarios: %v", err)
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Role,
			&user.Department,
			&user.Active,
			&user.Position,
			&user.Phone,
			&user.Language,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear usuario: %v", err)
		}

		user.CreatedAt = createdAt
		user.UpdatedAt = updatedAt
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar usuarios: %v", err)
	}

	return users, nil
}

// GetByID obtiene un usuario por su ID
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, role, department, active,
		       position, phone, language, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.Department,
		&user.Active,
		&user.Position,
		&user.Phone,
		&user.Language,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario con ID %s no encontrado", id)
		}
		return nil, fmt.Errorf("error al consultar usuario: %v", err)
	}

	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	return &user, nil
}

// GetByEmail obtiene un usuario por su correo electrónico
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, role, department, active,
		       position, phone, language, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.Department,
		&user.Active,
		&user.Position,
		&user.Phone,
		&user.Language,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario con email %s no encontrado", email)
		}
		return nil, fmt.Errorf("error al consultar usuario por email: %v", err)
	}

	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	return &user, nil
}

// Create crea un nuevo usuario
func (r *UserRepository) Create(user models.User) (*models.User, error) {
	query := `
		INSERT INTO users (id, first_name, last_name, email, password, role, department,
		                  active, position, phone, language, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`

	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	err := r.db.QueryRow(
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Role,
		user.Department,
		user.Active,
		user.Position,
		user.Phone,
		user.Language,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return nil, fmt.Errorf("error al crear usuario: %v", err)
	}

	return &user, nil
}

// Update actualiza un usuario existente
func (r *UserRepository) Update(user models.User) error {
	query := `
		UPDATE users
		SET first_name = $2, last_name = $3, email = $4, role = $5,
		    department = $6, active = $7, position = $8, phone = $9,
		    language = $10, updated_at = $11
		WHERE id = $1
	`

	now := time.Now()
	user.UpdatedAt = now

	result, err := r.db.Exec(
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Role,
		user.Department,
		user.Active,
		user.Position,
		user.Phone,
		user.Language,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error al actualizar usuario: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuario con ID %s no encontrado", user.ID)
	}

	return nil
}

// Delete elimina un usuario por su ID
func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar usuario: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuario con ID %s no encontrado", id)
	}

	return nil
}

// UpdatePassword actualiza la contraseña de un usuario
func (r *UserRepository) UpdatePassword(id string, password string) error {
	query := `
		UPDATE users
		SET password = $2, updated_at = $3
		WHERE id = $1
	`

	now := time.Now()

	result, err := r.db.Exec(query, id, password, now)
	if err != nil {
		return fmt.Errorf("error al actualizar contraseña: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuario con ID %s no encontrado", id)
	}

	return nil
}
