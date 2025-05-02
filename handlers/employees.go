package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"wider-circle-backend/db"
	"wider-circle-backend/models"

	"gorm.io/gorm"
)

type EmployeeHandler struct {
	DB  *gorm.DB
	URL string
}

func NewEmployeeHandler(db *gorm.DB, sourceURL string) *EmployeeHandler {
	return &EmployeeHandler{
		DB:  db,
		URL: sourceURL,
	}
}

func (h *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	exists, err := db.EmployeeExists(h.DB)
	if err != nil {
		log.Printf("Database check failed: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Println("No employee table or data found. Fetching from external source...")

		// Change the URL to the GitHub gist URL
		jsonURL := "https://gist.githubusercontent.com/chancock09/6d2a5a4436dcd488b8287f3e3e4fc73d/raw/fa47d64c6d5fc860fabd3033a1a4e3c59336324e/employees.json"

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(jsonURL)
		if err != nil {
			log.Printf("HTTP request failed: %v", err)
			http.Error(w, "Failed to fetch employees", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Unexpected status code: %d", resp.StatusCode)
			http.Error(w, "Failed to fetch employees", http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			http.Error(w, "Failed to read employee data", http.StatusInternalServerError)
			return
		}

		if len(body) == 0 {
			log.Println("Response body is empty")
			http.Error(w, "Empty response from employee source", http.StatusInternalServerError)
			return
		}

		var employees []models.Employee
		if err := json.Unmarshal(body, &employees); err != nil {
			log.Printf("JSON unmarshal failed: %v", err)
			http.Error(w, "Invalid JSON format", http.StatusInternalServerError)
			return
		}

		// Save the employees to the database
		if err := db.SaveEmployees(h.DB, employees); err != nil {
			log.Printf("Failed to save employees: %v", err)
			// Continue anyway since we want to return the JSON even if saving fails
		} else {
			log.Printf("Saved %d employees to database", len(employees))
		}

		// Return the fetched employees directly
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body) // Return the original JSON
		return
	}

	// If database exists and has data, fetch from database as before
	employees, err := db.GetEmployees(h.DB)
	if err != nil {
		log.Printf("Failed to retrieve employees: %v", err)
		http.Error(w, "Failed to retrieve employees", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Failed to encode employees to JSON: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Printf("Returned %d employees", len(employees))
}

func (h *EmployeeHandler) fetchEmployeesFromSource() ([]models.Employee, error) {
	log.Printf("Fetching employees from URL: %s", h.URL)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(h.URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("response body is empty")
	}

	var employees []models.Employee
	if err := json.Unmarshal(body, &employees); err != nil {
		log.Printf("JSON unmarshal failed: %v", err)
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	log.Printf("Fetched %d employees from source", len(employees))
	for i, emp := range employees {
		log.Printf("Employee %d: ID=%d, Name=%s, Title=%s, ManagerID=%v", i, emp.ID, emp.Name, emp.Title, emp.ManagerID)
	}

	return employees, nil
}
