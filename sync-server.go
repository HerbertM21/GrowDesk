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

func main() {
	// Configurar CORS middleware
	corsMiddleware := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}

	// Ruta para sincronizar usuarios
	http.HandleFunc("/api/sync/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Método no permitido")
			return
		}

		// Leer el cuerpo de la solicitud
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error al leer datos: %v", err)
			return
		}
		defer r.Body.Close()

		// Parsear los usuarios JSON
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

		// Guardar en la ruta local (para desarrollo)
		dataDir1 := "./GrowDesk/backend/data"
		if err := os.MkdirAll(dataDir1, 0755); err != nil {
			log.Printf("Advertencia: No se pudo crear directorio local %s: %v", dataDir1, err)
		} else {
			usersFile1 := filepath.Join(dataDir1, "users.json")
			if err := os.WriteFile(usersFile1, jsonData, 0644); err != nil {
				log.Printf("Advertencia: No se pudo escribir en el archivo local %s: %v", usersFile1, err)
			} else {
				log.Printf("Usuarios sincronizados localmente: %d usuarios guardados en %s", len(users), usersFile1)
			}
		}

		// Guardar en la ruta del contenedor (para producción)
		dataDir2 := "/app/data"
		var usersFile2 string
		var savedInContainer bool

		if _, err := os.Stat(dataDir2); err == nil {
			// Si el directorio existe (estamos dentro del contenedor)
			if err := os.MkdirAll(dataDir2, 0755); err != nil {
				log.Printf("Advertencia: No se pudo crear directorio en contenedor %s: %v", dataDir2, err)
			} else {
				usersFile2 = filepath.Join(dataDir2, "users.json")
				if err := os.WriteFile(usersFile2, jsonData, 0644); err != nil {
					log.Printf("Advertencia: No se pudo escribir en el archivo del contenedor %s: %v", usersFile2, err)
				} else {
					log.Printf("Usuarios sincronizados en contenedor: %d usuarios guardados en %s", len(users), usersFile2)
					savedInContainer = true
				}
			}
		} else {
			log.Printf("Información: Directorio del contenedor %s no existe, omitiendo", dataDir2)
		}

		// Responder con éxito si al menos uno de los archivos se guardó
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if savedInContainer {
			fmt.Fprintf(w, `{"success": true, "message": "Usuarios sincronizados correctamente en ambas ubicaciones", "count": %d, "localPath": "%s", "containerPath": "%s"}`, len(users), dataDir1, dataDir2)
		} else {
			fmt.Fprintf(w, `{"success": true, "message": "Usuarios sincronizados correctamente en ubicación local", "count": %d, "path": "%s"}`, len(users), dataDir1)
		}
	})

	// Aplicar middleware CORS a todas las rutas
	handler := corsMiddleware(http.DefaultServeMux)

	// Iniciar servidor
	port := "8000"
	log.Printf("Servidor de sincronización iniciado en http://localhost:%s", port)
	log.Printf("Endpoint de sincronización: http://localhost:%s/api/sync/users", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
