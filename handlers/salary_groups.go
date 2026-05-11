package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type SalaryGroupHandler struct{ DB *sql.DB }

func NewSalaryGroupHandler(db *sql.DB) *SalaryGroupHandler {
	return &SalaryGroupHandler{DB: db}
}

func (h *SalaryGroupHandler) Handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/salary_groups"), "/")
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

func (h *SalaryGroupHandler) list(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.DB.Query("SELECT sg_no, sg_name, base_salary, from_date, to_date FROM salary_group LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listSalaryGroups query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []SalaryGroup{}
	for rows.Next() {
		var sg SalaryGroup
		if err := rows.Scan(&sg.SgNo, &sg.SgName, &sg.BaseSalary, &sg.FromDate, &sg.ToDate); err != nil {
			log.Printf("[ERROR] listSalaryGroups scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, sg)
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func (h *SalaryGroupHandler) get(w http.ResponseWriter, r *http.Request) {
	sgNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/salary_groups/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid salary group ID"})
		return
	}

	rows, err := h.DB.Query("SELECT sg_no, sg_name, base_salary, from_date, to_date FROM salary_group WHERE sg_no = ?", sgNo)
	if err != nil {
		log.Printf("[ERROR] getSalaryGroup query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []SalaryGroup{}
	for rows.Next() {
		var sg SalaryGroup
		if err := rows.Scan(&sg.SgNo, &sg.SgName, &sg.BaseSalary, &sg.FromDate, &sg.ToDate); err != nil {
			log.Printf("[ERROR] getSalaryGroup scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, sg)
	}

	if len(items) == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Salary group not found"})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func (h *SalaryGroupHandler) create(w http.ResponseWriter, r *http.Request) {
	var sg SalaryGroup
	if err := json.NewDecoder(r.Body).Decode(&sg); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec("INSERT INTO salary_group (sg_name, base_salary, from_date, to_date) VALUES (?, ?, ?, ?)", sg.SgName, sg.BaseSalary, sg.FromDate, sg.ToDate)
	if err != nil {
		log.Printf("[ERROR] createSalaryGroup insert failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create salary group"})
		return
	}

	id, _ := result.LastInsertId()
	sg.SgNo = int(id)
	WriteJSON(w, http.StatusCreated, APIResponse{Success: true, Data: sg})
}

func (h *SalaryGroupHandler) update(w http.ResponseWriter, r *http.Request) {
	sgNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/salary_groups/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid salary group ID"})
		return
	}

	var sg SalaryGroup
	if err := json.NewDecoder(r.Body).Decode(&sg); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec(
		"UPDATE salary_group SET sg_name = ?, base_salary = ?, from_date = ?, to_date = ? WHERE sg_no = ? AND from_date = ?",
		sg.SgName, sg.BaseSalary, sg.FromDate, sg.ToDate, sgNo, sg.FromDate,
	)
	if err != nil {
		log.Printf("[ERROR] updateSalaryGroup query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update salary group"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Salary group not found"})
		return
	}

	sg.SgNo = sgNo
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: sg})
}

func (h *SalaryGroupHandler) delete(w http.ResponseWriter, r *http.Request) {
	sgNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/salary_groups/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid salary group ID"})
		return
	}

	result, err := h.DB.Exec("DELETE FROM salary_group WHERE sg_no = ?", sgNo)
	if err != nil {
		log.Printf("[ERROR] deleteSalaryGroup query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete salary group"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Salary group not found"})
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Salary group deleted"}})
}
