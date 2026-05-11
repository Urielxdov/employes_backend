package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type DepartmentHandler struct{ DB *sql.DB }

func NewDepartmentHandler(db *sql.DB) *DepartmentHandler {
	return &DepartmentHandler{DB: db}
}

func (h *DepartmentHandler) Handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/departments"), "/")
	hasID := len(parts) > 1 && parts[1] != ""

	switch r.Method {
	case "GET":
		if hasID {
			h.get(w, r)
		} else {
			h.list(w, r)
		}
	case "POST":
		h.create(w, r)
	case "PUT":
		if !hasID {
			WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for PUT"})
			return
		}
		h.update(w, r)
	case "DELETE":
		if !hasID {
			WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for DELETE"})
			return
		}
		h.delete(w, r)
	default:
		WriteJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "Method not allowed"})
	}
}

func (h *DepartmentHandler) list(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.DB.Query("SELECT dept_no, dept_name FROM department")
	if err != nil {
		log.Printf("[ERROR] listDepartments query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	departments := []Department{}
	for rows.Next() {
		var d Department
		if err := rows.Scan(&d.DeptNo, &d.DeptName); err != nil {
			log.Printf("[ERROR] listDepartments scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		departments = append(departments, d)
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: departments})
}

func (h *DepartmentHandler) get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/departments/")

	var d Department
	err := h.DB.QueryRow("SELECT dept_no, dept_name FROM department WHERE dept_no = ?", id).
		Scan(&d.DeptNo, &d.DeptName)
	if err == sql.ErrNoRows {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Department not found"})
		return
	}
	if err != nil {
		log.Printf("[ERROR] getDepartment query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: d})
}

func (h *DepartmentHandler) create(w http.ResponseWriter, r *http.Request) {
	var d Department
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	_, err := h.DB.Exec("INSERT INTO department (dept_no, dept_name) VALUES (?, ?)", d.DeptNo, d.DeptName)
	if err != nil {
		log.Printf("[ERROR] createDepartment insert failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create department"})
		return
	}
	WriteJSON(w, http.StatusCreated, APIResponse{Success: true, Data: d})
}

func (h *DepartmentHandler) update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/departments/")

	var d Department
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec("UPDATE department SET dept_name = ? WHERE dept_no = ?", d.DeptName, id)
	if err != nil {
		log.Printf("[ERROR] updateDepartment query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update department"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Department not found"})
		return
	}

	d.DeptNo = id
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: d})
}

func (h *DepartmentHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/departments/")

	result, err := h.DB.Exec("DELETE FROM department WHERE dept_no = ?", id)
	if err != nil {
		log.Printf("[ERROR] deleteDepartment query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete department"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Department not found"})
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Department deleted"}})
}
