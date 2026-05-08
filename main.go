package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func router(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	// Health check
	if path == "/health" && method == "GET" {
		healthCheck(w, r)
		return
	}

	// Route requests by entity and method
	if strings.HasPrefix(path, "/api/v1/employees") {
		handleEmployees(w, r)
		return
	}
	if strings.HasPrefix(path, "/api/v1/departments") {
		handleDepartments(w, r)
		return
	}
	if strings.HasPrefix(path, "/api/v1/salaries") {
		handleSalaries(w, r)
		return
	}
	if strings.HasPrefix(path, "/api/v1/titles") {
		handleTitles(w, r)
		return
	}
	if strings.HasPrefix(path, "/api/v1/dept_emp") {
		handleDeptEmp(w, r)
		return
	}
	if strings.HasPrefix(path, "/api/v1/dept_manager") {
		handleDeptManager(w, r)
		return
	}

	writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Route not found"})
}

func handleEmployees(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/employees"), "/")
	hasID := len(parts) > 1 && parts[1] != ""

	switch r.Method {
	case "GET":
		if hasID {
			getEmployee(w, r)
		} else {
			listEmployees(w, r)
		}
	case "POST":
		createEmployee(w, r)
	case "PUT":
		if !hasID {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for PUT"})
			return
		}
		updateEmployee(w, r)
	case "DELETE":
		if !hasID {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for DELETE"})
			return
		}
		deleteEmployee(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "Method not allowed"})
	}
}

func handleDepartments(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/departments"), "/")
	hasID := len(parts) > 1 && parts[1] != ""

	switch r.Method {
	case "GET":
		if hasID {
			getDepartment(w, r)
		} else {
			listDepartments(w, r)
		}
	case "POST":
		createDepartment(w, r)
	case "PUT":
		if !hasID {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for PUT"})
			return
		}
		updateDepartment(w, r)
	case "DELETE":
		if !hasID {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for DELETE"})
			return
		}
		deleteDepartment(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "Method not allowed"})
	}
}

func handleSalaries(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/salaries"), "/")
	hasID := len(parts) > 1 && parts[1] != ""

	switch r.Method {
	case "GET":
		if hasID {
			getSalary(w, r)
		} else {
			listSalaries(w, r)
		}
	case "POST":
		createSalary(w, r)
	case "PUT":
		if !hasID {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for PUT"})
			return
		}
		updateSalary(w, r)
	case "DELETE":
		if !hasID {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for DELETE"})
			return
		}
		deleteSalary(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "Method not allowed"})
	}
}

func handleTitles(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/titles"), "/")
	hasID := len(parts) > 1 && parts[1] != ""

	switch r.Method {
	case "GET":
		if hasID {
			getTitle(w, r)
		} else {
			listTitles(w, r)
		}
	case "POST":
		createTitle(w, r)
	case "PUT":
		if !hasID {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for PUT"})
			return
		}
		updateTitle(w, r)
	case "DELETE":
		if !hasID {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for DELETE"})
			return
		}
		deleteTitle(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "Method not allowed"})
	}
}

func handleDeptEmp(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/dept_emp"), "/")
	hasID := len(parts) > 1 && parts[1] != ""

	switch r.Method {
	case "GET":
		if hasID {
			getDeptEmp(w, r)
		} else {
			listDeptEmp(w, r)
		}
	case "POST":
		createDeptEmp(w, r)
	case "PUT":
		if !hasID {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for PUT"})
			return
		}
		updateDeptEmp(w, r)
	case "DELETE":
		if !hasID {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for DELETE"})
			return
		}
		deleteDeptEmp(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "Method not allowed"})
	}
}

func handleDeptManager(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/dept_manager"), "/")
	hasID := len(parts) > 1 && parts[1] != ""

	switch r.Method {
	case "GET":
		if hasID {
			getDeptManager(w, r)
		} else {
			listDeptManager(w, r)
		}
	case "POST":
		createDeptManager(w, r)
	case "PUT":
		if !hasID {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for PUT"})
			return
		}
		updateDeptManager(w, r)
	case "DELETE":
		if !hasID {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for DELETE"})
			return
		}
		deleteDeptManager(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "Method not allowed"})
	}
}

func main() {
	// Load environment
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbUser == "" {
		dbUser = "root"
	}

	// Initialize database
	log.Println("[INIT] Connecting to database...")
	err := initDB()
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize database: %v\n", err)
	}
	defer closeDB()

	// Start server
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[INFO] Starting server on port %s\n", port)
	http.HandleFunc("/", router)

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("[FATAL] Server failed: %v\n", err)
	}
}
