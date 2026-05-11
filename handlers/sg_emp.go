package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type SgEmpHandler struct{ DB *sql.DB }

func NewSgEmpHandler(db *sql.DB) *SgEmpHandler {
	return &SgEmpHandler{DB: db}
}

func (h *SgEmpHandler) Handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/sg_emp"), "/")
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

func (h *SgEmpHandler) list(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.DB.Query("SELECT emp_no, sg_no, from_date, to_date FROM sg_emp LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listSgEmp query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []SgEmp{}
	for rows.Next() {
		var se SgEmp
		if err := rows.Scan(&se.EmpNo, &se.SgNo, &se.FromDate, &se.ToDate); err != nil {
			log.Printf("[ERROR] listSgEmp scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, se)
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func (h *SgEmpHandler) get(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/sg_emp/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	rows, err := h.DB.Query("SELECT emp_no, sg_no, from_date, to_date FROM sg_emp WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] getSgEmp query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []SgEmp{}
	for rows.Next() {
		var se SgEmp
		if err := rows.Scan(&se.EmpNo, &se.SgNo, &se.FromDate, &se.ToDate); err != nil {
			log.Printf("[ERROR] getSgEmp scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, se)
	}

	if len(items) == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "No salary group assignments found"})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func (h *SgEmpHandler) create(w http.ResponseWriter, r *http.Request) {
	var se SgEmp
	if err := json.NewDecoder(r.Body).Decode(&se); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	_, err := h.DB.Exec("INSERT INTO sg_emp (emp_no, sg_no, from_date, to_date) VALUES (?, ?, ?, ?)", se.EmpNo, se.SgNo, se.FromDate, se.ToDate)
	if err != nil {
		log.Printf("[ERROR] createSgEmp insert failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create salary group assignment"})
		return
	}
	WriteJSON(w, http.StatusCreated, APIResponse{Success: true, Data: se})
}

func (h *SgEmpHandler) update(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/sg_emp/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var se SgEmp
	if err := json.NewDecoder(r.Body).Decode(&se); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec(
		"UPDATE sg_emp SET sg_no = ?, from_date = ?, to_date = ? WHERE emp_no = ? AND from_date = ?",
		se.SgNo, se.FromDate, se.ToDate, empNo, se.FromDate,
	)
	if err != nil {
		log.Printf("[ERROR] updateSgEmp query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update salary group assignment"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Salary group assignment not found"})
		return
	}

	se.EmpNo = empNo
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: se})
}

func (h *SgEmpHandler) delete(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/sg_emp/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	result, err := h.DB.Exec("DELETE FROM sg_emp WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteSgEmp query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete salary group assignment"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Salary group assignment not found"})
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Salary group assignment deleted"}})
}
