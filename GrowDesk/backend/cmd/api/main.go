package main

import (
	"log"
	"os"

	"github.com/hmdev/GrowDesk/backend/internal/server"
	"github.com/hmdev/GrowDesk/backend/pkg/auth"
	"github.com/hmdev/GrowDesk/backend/pkg/database"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró el archivo .env, usando variables de entorno del sistema")
	}

	// Inicializar sistema de autenticación JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("ADVERTENCIA: JWT_SECRET no está configurado, usando valor por defecto")
		jwtSecret = "clave_segura_por_defecto_solo_para_desarrollo"
	}
	auth.InitJWT(jwtSecret)

	// Inicializar conexión a la base de datos
	database.Initialize()

	// Verificar si la base de datos está inicializada
	if database.GetDB() == nil {
		log.Println("ADVERTENCIA: Base de datos no disponible, continuando en modo de prueba sin persistencia")
	}

	// Obtener puerto de las variables de entorno o usar puerto por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Puerto por defecto
	}

	// Iniciar el servidor HTTP
	if err := server.StartServer(port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
