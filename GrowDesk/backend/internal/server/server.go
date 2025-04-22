package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// StartServer inicia el servidor HTTP con todas las rutas configuradas
func StartServer(port string) error {

	router := gin.Default()

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "X-Widget-ID", "X-Widget-Token"}
	router.Use(cors.New(config))

	// Configurar rutas
	SetupRoutes(router)

	// Iniciar servidor
	return router.Run(":" + port)
}
