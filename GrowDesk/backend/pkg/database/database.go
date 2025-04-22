package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Iiniciar la bd
func Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: .env no encontrado.")
	}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "growdesk")
	dbSSLMode := getEnv("DB_SSL_MODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// logger de GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Umbral para consultas lentas
			LogLevel:                  logger.Info, // Nivel de log
			IgnoreRecordNotFoundError: true,        // Ignorar errores de registro no encontrado
			Colorful:                  true,        // Activar color
		},
	)

	// Intentar conectar a la base de datos con reintentos
	maxRetries := 5
	var retryCount int

	for retryCount < maxRetries {
		log.Printf("Conectar la bd (attempt %d/%d)...", retryCount+1, maxRetries)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})

		if err == nil {
			DB = db
			log.Println("Conexión de base de datos éxitosa")

			// Verificar conexión a la base de datos
			sqlDB, err := db.DB()
			if err != nil {
				log.Printf("Error al conectar con la base de datos: %v", err)
				continue
			}

			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)

			// Verificar que la conexión está viva con ping
			err = sqlDB.Ping()
			if err != nil {
				log.Printf("Error al hacer ping a la base de datos: %v", err)
				retryCount++
				time.Sleep(time.Second * 3)
				continue
			}

			return
		}

		log.Printf("Error al conectar con la base de datos: %v", err)
		retryCount++
		time.Sleep(time.Second * 3) // Esperar 3 segundos antes de reintentar
	}

	log.Printf("ADVERTENCIA: Error al conectar con la base de datos en %d intentos.", maxRetries)

}

// GetDB retorna la instancia de la bd
func GetDB() *gorm.DB {
	return DB
}

func MigrateDB(models ...interface{}) {
	log.Println("Iniciando migración...")
	err := DB.AutoMigrate(models...)
	if err != nil {
		log.Fatalf("Error al migrar BD: %v", err)
	}
	log.Println("Migración de la base de datos éxitosa")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
