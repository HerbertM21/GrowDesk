package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

var db *sql.DB

// InitDB inicializa la conexión a la base de datos PostgreSQL
func InitDB() (*sql.DB, error) {
	// Obtener variables de entorno de la base de datos
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	user := getEnvOrDefault("DB_USER", "postgres")
	password := getEnvOrDefault("DB_PASSWORD", "postgres")
	dbname := getEnvOrDefault("DB_NAME", "growdesk")
	sslmode := getEnvOrDefault("DB_SSLMODE", "disable")

	// Cadena de conexión
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	// Abrir conexión
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión: %v", err)
	}

	// Verificar conexión
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error al verificar conexión: %v", err)
	}

	log.Println("Conexión a PostgreSQL establecida exitosamente")
	return db, nil
}

// InitializeSchema inicializa el esquema de la base de datos
func InitializeSchema(db *sql.DB) error {
	// Leer archivo de esquema
	schemaPath := filepath.Join("internal", "db", "schema.sql")
	schemaContent, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("error al leer archivo de esquema: %v", err)
	}

	// Ejecutar script SQL
	_, err = db.Exec(string(schemaContent))
	if err != nil {
		return fmt.Errorf("error al ejecutar script de esquema: %v", err)
	}

	log.Println("Esquema de base de datos inicializado exitosamente")
	return nil
}

// Close cierra la conexión a la base de datos
func Close() {
	if db != nil {
		db.Close()
		log.Println("Conexión a PostgreSQL cerrada")
	}
}

// GetDB obtiene la conexión a la base de datos
func GetDB() *sql.DB {
	return db
}

// Helper para obtener variables de entorno con valor por defecto
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
