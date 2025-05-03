package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"wider-circle-backend/db"
	"wider-circle-backend/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func getEnv(key, defaultValue string) string {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	dbDriver := getEnv("DB_DRIVER", "sqlite")
	// sqlite example: "employees.db"
	// postgres url example: "host=localhost user=postgres password=postgres dbname=employees port=5432 sslmode=disable"
	dbDSN := getEnv("DB_DSN", "employees.db")

	port := getEnv("PORT", "8080")

	// base case employee URL
	// see env file other test case urls
	employeeURL := getEnv("EMPLOYEE_SEED_URL",
		"https://gist.githubusercontent.com/chancock09/6d2a5a4436dcd488b8287f3e3e4fc73d/raw/fa47d64c6d5fc860fabd3033a1a4e4c59336324e/employees.json")

	database, err := db.InitializeDatabase(dbDriver, dbDSN)
	if err != nil {
		log.Printf("Database initialization failed: %v", err)

		// Create a client with timeout
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		// Send GET request to /employees endpoint
		resp, err := client.Get("http://localhost:" + port + "/employee")
		if err != nil {
			log.Printf("Failed to send request to /employee: %v", err)
		} else {
			log.Printf("Request to /employee returned status: %s", resp.Status)
			resp.Body.Close()
		}

		// Continue execution with nil database, handlers should check for nil
		database = nil
	}

	employeeHandler := handlers.NewEmployeeHandler(database, employeeURL)

	router := mux.NewRouter()
	router.HandleFunc("/employee", employeeHandler.GetEmployees).Methods("GET")

	router.Use(corsMiddleware)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
