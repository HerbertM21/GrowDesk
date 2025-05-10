package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// User representa la estructura de un usuario
type User struct {
	ID         string `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department,omitempty"`
	Active     bool   `json:"active"`
	Password   string `json:"password,omitempty"`
	CreatedAt  string `json:"createdAt,omitempty"`
	UpdatedAt  string `json:"updatedAt,omitempty"`
}

// SyncUsersHandler maneja la sincronización de usuarios desde el frontend con localStorage
func SyncUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS para permitir solicitudes desde el frontend
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Manejar solicitudes OPTIONS (preflight CORS)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Método no permitido")
		return
	}

	// Lee el cuerpo de la solicitud
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error al leer datos: %v", err)
		return
	}
	defer r.Body.Close()

	// Parsea los usuarios JSON
	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error al parsear JSON: %v", err)
		return
	}

	// Convertir los usuarios a JSON formateado
	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error al generar JSON: %v", err)
		return
	}

	// Definir las rutas de directorios donde se guardarán los datos
	locations := []struct {
		dir  string
		desc string
	}{
		{os.Getenv("DATA_DIR"), "volumen Docker"},         // Volumen Docker
		{os.Getenv("LOCAL_DATA_DIR"), "directorio local"}, // Directorio local montado
	}

	var successCount int
	var allPaths []string

	// Guardar en todas las ubicaciones
	for _, loc := range locations {
		if loc.dir == "" {
			continue // Saltar ubicaciones vacías
		}

		// Crear el directorio si no existe
		if err := os.MkdirAll(loc.dir, 0755); err != nil {
			log.Printf("Advertencia: No se pudo crear directorio en %s (%s): %v", loc.dir, loc.desc, err)
			continue
		}

		// Ruta al archivo de usuarios
		usersFile := filepath.Join(loc.dir, "users.json")

		// Escribir los datos al archivo
		if err := os.WriteFile(usersFile, jsonData, 0644); err != nil {
			log.Printf("Error al escribir en %s (%s): %v", usersFile, loc.desc, err)
			continue
		}

		log.Printf("Usuarios sincronizados: %d usuarios guardados en %s (%s)", len(users), usersFile, loc.desc)
		successCount++
		allPaths = append(allPaths, usersFile)
	}

	// Verificar si se guardó en al menos un lugar
	if successCount == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"success": false, "message": "No se pudo guardar en ninguna ubicación"}`)
		return
	}

	// Responder con éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	pathsJSON, _ := json.Marshal(allPaths)
	fmt.Fprintf(w, `{"success": true, "message": "Usuarios sincronizados correctamente en %d ubicaciones", "count": %d, "paths": %s}`,
		successCount, len(users), pathsJSON)
}

func main() {
	// Obtener el puerto desde variables de entorno o usar valor por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Puerto por defecto
	}

	// Imprimir información sobre la configuración
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "/app/data" // Valor por defecto para Docker
	}
	log.Printf("Directorio de datos: %s", dataDir)

	// Configurar rutas
	http.HandleFunc("/api/sync/users", SyncUsersHandler)

	// Iniciar servidor
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Servidor de sincronización iniciado en http://0.0.0.0%s", addr)
	log.Printf("Endpoint de sincronización: http://localhost%s/api/sync/users", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
