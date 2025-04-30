package database

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Variables globales para la gestión de la base de datos
var (
	db     *gorm.DB
	dbLock sync.Mutex
)

// DatabaseType representa el tipo de base de datos a utilizar
type DatabaseType string

const (
	DBTypeSQLite   DatabaseType = "sqlite"
	DBTypeMySQL    DatabaseType = "mysql"
	DBTypePostgres DatabaseType = "postgres"
)

// Config contiene la configuración de la base de datos
type Config struct {
	Type         DatabaseType
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	SSLMode      string
	FilePath     string // Para SQLite
	LogLevel     logger.LogLevel
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// LoadConfigFromEnv carga la configuración desde variables de entorno
func LoadConfigFromEnv() Config {
	dbType := DatabaseType(getEnvOrDefault("DB_TYPE", "sqlite"))

	return Config{
		Type:         dbType,
		Host:         getEnvOrDefault("DB_HOST", "localhost"),
		Port:         getEnvInt("DB_PORT", 5432),
		Username:     getEnvOrDefault("DB_USER", "postgres"),
		Password:     getEnvOrDefault("DB_PASSWORD", "password"),
		DatabaseName: getEnvOrDefault("DB_NAME", "growdesk"),
		SSLMode:      getEnvOrDefault("DB_SSL_MODE", "disable"),
		FilePath:     getEnvOrDefault("DB_FILE_PATH", "growdesk.db"),
		LogLevel:     logger.Error,
		MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
		MaxLifetime:  time.Duration(getEnvInt("DB_MAX_LIFETIME", 5)) * time.Minute,
	}
}

// InitDB inicializa la conexión a la base de datos
func InitDB(config Config) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	if db != nil {
		// Ya tenemos una conexión a la base de datos
		return nil
	}

	var err error
	var dialector gorm.Dialector

	// Configurar dialector según el tipo de base de datos
	switch config.Type {
	case DBTypeSQLite:
		// SQLite
		dialector = sqlite.Open(config.FilePath)
	case DBTypeMySQL:
		// MySQL
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Username, config.Password, config.Host, config.Port, config.DatabaseName)
		dialector = mysql.Open(dsn)
	case DBTypePostgres:
		// PostgreSQL
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.Username, config.Password, config.DatabaseName, config.SSLMode)
		dialector = postgres.Open(dsn)
	default:
		return fmt.Errorf("tipo de base de datos no soportado: %s", config.Type)
	}

	// Configuraciones de GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	}

	// Inicializar conexión
	db, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("error al conectar con la base de datos: %v", err)
	}

	// Configurar pool de conexiones
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("error al obtener el pool de conexiones: %v", err)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	// Verificar conexión
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("error al verificar conexión a la base de datos: %v", err)
	}

	// Log de éxito
	fmt.Printf("Conexión exitosa a la base de datos: %s\n", config.Type)

	return nil
}

// GetDB retorna la instancia de conexión a la base de datos
func GetDB() (*gorm.DB, error) {
	if db == nil {
		return nil, errors.New("la base de datos no ha sido inicializada, debe llamar a InitDB primero")
	}
	return db, nil
}

// CloseDB cierra la conexión a la base de datos
func CloseDB() {
	dbLock.Lock()
	defer dbLock.Unlock()

	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Error al obtener el pool de conexiones para cerrar: %v\n", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		fmt.Printf("Error al cerrar la conexión a la base de datos: %v\n", err)
		return
	}

	db = nil
	fmt.Println("Conexión a la base de datos cerrada correctamente")
}

// getEnvOrDefault obtiene una variable de entorno o devuelve un valor por defecto
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvInt obtiene una variable de entorno como entero o devuelve un valor por defecto
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var result int
	_, err := fmt.Sscanf(value, "%d", &result)
	if err != nil {
		return defaultValue
	}

	return result
}
