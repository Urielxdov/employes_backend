package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type SalaryHandler struct{ DB *sql.DB }

func NewSalaryHandler(db *sql.DB) *SalaryHandler {
	return &SalaryHandler{DB: db}
}

func (h *SalaryHandler) Handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/salaries"), "/")
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

func (h *SalaryHandler) list(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.DB.Query("SELECT emp_no, salary, from_date, to_date FROM salary LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listSalaries query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	salaries := []Salary{}
	for rows.Next() {
		var s Salary
		if err := rows.Scan(&s.EmpNo, &s.Salary, &s.FromDate, &s.ToDate); err != nil {
			log.Printf("[ERROR] listSalaries scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		salaries = append(salaries, s)
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: salaries})
}

func (h *SalaryHandler) get(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/salaries/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	rows, err := h.DB.Query("SELECT emp_no, salary, from_date, to_date FROM salary WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] getSalary query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	salaries := []Salary{}
	for rows.Next() {
		var s Salary
		if err := rows.Scan(&s.EmpNo, &s.Salary, &s.FromDate, &s.ToDate); err != nil {
			log.Printf("[ERROR] getSalary scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		salaries = append(salaries, s)
	}

	if len(salaries) == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "No salaries found"})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: salaries})
}

func (h *SalaryHandler) create(w http.ResponseWriter, r *http.Request) {
	var s Salary
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	_, err := h.DB.Exec("INSERT INTO salary (emp_no, salary, from_date, to_date) VALUES (?, ?, ?, ?)", s.EmpNo, s.Salary, s.FromDate, s.ToDate)
	if err != nil {
		log.Printf("[ERROR] createSalary insert failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create salary"})
		return
	}
	WriteJSON(w, http.StatusCreated, APIResponse{Success: true, Data: s})
}

func (h *SalaryHandler) update(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/salaries/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var s Salary
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec(
		"UPDATE salary SET salary = ?, from_date = ?, to_date = ? WHERE emp_no = ? AND from_date = ?",
		s.Salary, s.FromDate, s.ToDate, empNo, s.FromDate,
	)
	if err != nil {
		log.Printf("[ERROR] updateSalary query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update salary"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Salary record not found"})
		return
	}

	s.EmpNo = empNo
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: s})
}

func (h *SalaryHandler) delete(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/salaries/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	result, err := h.DB.Exec("DELETE FROM salary WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteSalary query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete salary"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Salary record not found"})
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Salary deleted"}})
}
