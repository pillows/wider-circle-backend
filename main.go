package main

import (
	"log"
	"net/http"
	"os"
	"wider-circle-backend/db"
	"wider-circle-backend/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func getEnv(key, defaultValue string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
	employeeURL := getEnv("EMPLOYEE_SEED_URL",
		"https://gist.githubusercontent.com/chancock09/6d2a5a4436dcd488b8287f3e3e4fc73d/raw/fa47d64c6d5fc860fabd3033a1a4e3c59336324e/employees.json")

	// test case 1 with an extra nested employee
	// employeeURL := getEnv("EMPLOYEE_SEED_URL",
	// 	"https://gist.githubusercontent.com/pillows/dded29423234b68f3cdde47aeafb0d63/raw/ac706aa7ef6c28e2d5bcc78b1560f621c98778e3/gistfile1.json")

	// test case 2 with no employees
	// employeeURL := getEnv("EMPLOYEE_SEED_URL",
	// "https://gist.githubusercontent.com/pillows/dbca826d3821805a842efb39315591e1/raw/9ffcd84a87536977082a68a99513270b75714a39/empty.json")

	// test case 3 with a single employee
	// employeeURL := getEnv("EMPLOYEE_SEED_URL",
	// 	"https://gist.githubusercontent.com/pillows/dded29423234b68f3cdde47aeafb0d63/raw/28d7fd2b3535d88d56738c459dac630036d676a8/gistfile1.json")

	database, err := db.InitializeDatabase(dbDriver, dbDSN)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	employeeHandler := handlers.NewEmployeeHandler(database, employeeURL)

	router := mux.NewRouter()
	router.HandleFunc("/employees", employeeHandler.GetEmployees).Methods("GET")

	router.Use(corsMiddleware)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
