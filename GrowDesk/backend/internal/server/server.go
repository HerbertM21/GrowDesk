package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/internal/routes"
	"github.com/hmdev/GrowDesk/backend/pkg/auth"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
	"github.com/hmdev/GrowDesk/backend/pkg/websocket"
)

// Server representa el servidor de la aplicación
type Server struct {
	router *gin.Engine
	cfg    *Config
}

// Config contiene la configuración del servidor
type Config struct {
	Port            int
	JWTSecret       string
	DatabaseConfig  database.Config
	MaxRequestSize  int64
	TrustedProxies  []string
	ProductionMode  bool
	ShutdownTimeout int
}

// DefaultConfig retorna la configuración por defecto del servidor
func DefaultConfig() *Config {
	port, _ := strconv.Atoi(getEnvOrDefault("PORT", "8080"))

	return &Config{
		Port:            port,
		JWTSecret:       getEnvOrDefault("JWT_SECRET", "growdesk-secret-key"),
		DatabaseConfig:  database.LoadConfigFromEnv(),
		MaxRequestSize:  32 << 20, // 32MB
		TrustedProxies:  []string{},
		ProductionMode:  getEnvOrDefault("APP_ENV", "development") == "production",
		ShutdownTimeout: 10, // segundos
	}
}

// New crea una nueva instancia del servidor
func New(cfg *Config) *Server {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Configurar modo de Gin según entorno
	if cfg.ProductionMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Crear router
	router := gin.New()

	// Configurar middleware global
	router.Use(gin.Recovery())
	if !cfg.ProductionMode {
		router.Use(gin.Logger())
	}

	// Configurar tamaño máximo de request
	router.MaxMultipartMemory = cfg.MaxRequestSize

	// Configurar proxies confiables
	if len(cfg.TrustedProxies) > 0 {
		router.SetTrustedProxies(cfg.TrustedProxies)
	}

	return &Server{
		router: router,
		cfg:    cfg,
	}
}

// Initialize inicializa los componentes del servidor
func (s *Server) Initialize() error {
	// Inicializar base de datos
	err := database.InitDB(s.cfg.DatabaseConfig)
	if err != nil {
		return fmt.Errorf("error inicializando base de datos: %v", err)
	}

	// Inicializar autenticación
	auth.SetJWTSecret(s.cfg.JWTSecret)

	// Inicializar hub de WebSocket
	websocket.NewHub()

	// Configurar rutas
	routes.SetupRoutes(s.router)

	return nil
}

// Start inicia el servidor HTTP
func (s *Server) Start() error {
	// Crear servidor HTTP
	addr := fmt.Sprintf(":%d", s.cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	// Canal para manejar señales de sistema operativo
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en goroutine separada
	go func() {
		log.Printf("Servidor iniciado en http://localhost:%d", s.cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar servidor: %v", err)
		}
	}()

	// Esperar señal de apagado
	<-quit
	log.Println("Apagando servidor...")

	// Establecer timeout para cierre controlado
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.ShutdownTimeout)*time.Second)
	defer cancel()

	// Cerrar servicios
	database.CloseDB()

	// Cerrar servidor
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error al apagar servidor: %v", err)
	}

	log.Println("Servidor apagado correctamente")
	return nil
}

// getEnvOrDefault obtiene una variable de entorno o devuelve un valor por defecto
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
