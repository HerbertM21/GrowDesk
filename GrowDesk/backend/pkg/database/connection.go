package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Configuración de la base de datos
type Config struct {
	Type            string // mysql, postgres, sqlite
	Host            string
	Port            string
	Username        string
	Password        string
	Database        string
	AdditionalArgs  string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	Debug           bool
}

// InitDB inicializa la conexión a la base de datos
func InitDB(config Config) error {
	var err error
	var dialector gorm.Dialector

	// Configurar logger de GORM
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	if config.Debug {
		gormLogger.LogMode(logger.Info)
	}

	// Configurar dialector según el tipo de base de datos
	switch config.Type {
	case "mysql":
		dsn := config.Username + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
		if config.AdditionalArgs != "" {
			dsn += "&" + config.AdditionalArgs
		}
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn := "host=" + config.Host + " port=" + config.Port + " user=" + config.Username + " dbname=" + config.Database + " password=" + config.Password + " sslmode=disable TimeZone=UTC"
		if config.AdditionalArgs != "" {
			dsn += " " + config.AdditionalArgs
		}
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(config.Database)
	default:
		log.Println("Modo sin base de datos activo")
		return nil
	}

	// Establecer conexión
	db, err = gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		log.Printf("Error al conectar a la base de datos: %v", err)
		return err
	}

	// Configurar pool de conexiones
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error al obtener la conexión SQL: %v", err)
		return err
	}

	// Establecer configuraciones del pool de conexiones
	if config.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	} else {
		sqlDB.SetMaxIdleConns(10)
	}

	if config.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	} else {
		sqlDB.SetMaxOpenConns(100)
	}

	if config.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	} else {
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	// Verificar conexión
	if err = sqlDB.Ping(); err != nil {
		log.Printf("Error al verificar la conexión a la base de datos: %v", err)
		return err
	}

	log.Printf("Conectado a la base de datos %s en %s:%s", config.Database, config.Host, config.Port)

	return nil
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return db
}

// CloseDB cierra la conexión a la base de datos
func CloseDB() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// LoadConfigFromEnv carga la configuración de la base de datos desde variables de entorno
func LoadConfigFromEnv() Config {
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite" // Por defecto usar SQLite
	}

	maxIdleConns := 10
	maxOpenConns := 100
	connMaxLifetime := time.Hour

	return Config{
		Type:            dbType,
		Host:            os.Getenv("DB_HOST"),
		Port:            os.Getenv("DB_PORT"),
		Username:        os.Getenv("DB_USER"),
		Password:        os.Getenv("DB_PASSWORD"),
		Database:        os.Getenv("DB_NAME"),
		AdditionalArgs:  os.Getenv("DB_ARGS"),
		MaxIdleConns:    maxIdleConns,
		MaxOpenConns:    maxOpenConns,
		ConnMaxLifetime: connMaxLifetime,
		Debug:           os.Getenv("DB_DEBUG") == "true",
	}
}
