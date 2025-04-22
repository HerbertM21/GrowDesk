package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hmdev/GrowDesk/backend/internal/routes"
	"github.com/hmdev/GrowDesk/backend/internal/server"
	"github.com/hmdev/GrowDesk/backend/models"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Println("ADVERTENCIA: NO SE ENCONTRÓ .ENV")
	}

	// Configurar modo de Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializar base de datos
	database.Initialize()

	// Migrar modelos
	database.MigrateDB(
		&models.User{},
		&models.Ticket{},
		&models.Message{},
		&models.Attachment{},
	)

	// Crear router
	router := gin.Default()

	// Configurar CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Widget-API-Key", "X-Widget-ID", "X-Widget-Token"}
	router.Use(cors.New(config))

	// Configurar rutas principales
	server.SetupRoutes(router)

	// Configurar rutas del widget
	routes.SetupWidgetRoutes(router)

	// Obtener puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Iniciar servidor
	log.Printf("Servidor iniciando en el puerto %s", port)
	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
