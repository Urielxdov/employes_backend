package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type RegionEmpHandler struct{ DB *sql.DB }

func NewRegionEmpHandler(db *sql.DB) *RegionEmpHandler {
	return &RegionEmpHandler{DB: db}
}

func (h *RegionEmpHandler) Handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/region_emp"), "/")
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

func (h *RegionEmpHandler) list(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.DB.Query("SELECT emp_no, region_id, from_date, to_date FROM region_emp LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listRegionEmp query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []RegionEmp{}
	for rows.Next() {
		var re RegionEmp
		if err := rows.Scan(&re.EmpNo, &re.RegionID, &re.FromDate, &re.ToDate); err != nil {
			log.Printf("[ERROR] listRegionEmp scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, re)
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func (h *RegionEmpHandler) get(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/region_emp/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	rows, err := h.DB.Query("SELECT emp_no, region_id, from_date, to_date FROM region_emp WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] getRegionEmp query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []RegionEmp{}
	for rows.Next() {
		var re RegionEmp
		if err := rows.Scan(&re.EmpNo, &re.RegionID, &re.FromDate, &re.ToDate); err != nil {
			log.Printf("[ERROR] getRegionEmp scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, re)
	}

	if len(items) == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "No region assignments found"})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func (h *RegionEmpHandler) create(w http.ResponseWriter, r *http.Request) {
	var re RegionEmp
	if err := json.NewDecoder(r.Body).Decode(&re); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	_, err := h.DB.Exec("INSERT INTO region_emp (emp_no, region_id, from_date, to_date) VALUES (?, ?, ?, ?)", re.EmpNo, re.RegionID, re.FromDate, re.ToDate)
	if err != nil {
		log.Printf("[ERROR] createRegionEmp insert failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create region assignment"})
		return
	}
	WriteJSON(w, http.StatusCreated, APIResponse{Success: true, Data: re})
}

func (h *RegionEmpHandler) update(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/region_emp/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var re RegionEmp
	if err := json.NewDecoder(r.Body).Decode(&re); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec(
		"UPDATE region_emp SET region_id = ?, from_date = ?, to_date = ? WHERE emp_no = ? AND from_date = ?",
		re.RegionID, re.FromDate, re.ToDate, empNo, re.FromDate,
	)
	if err != nil {
		log.Printf("[ERROR] updateRegionEmp query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update region assignment"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Region assignment not found"})
		return
	}

	re.EmpNo = empNo
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: re})
}

func (h *RegionEmpHandler) delete(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/region_emp/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	result, err := h.DB.Exec("DELETE FROM region_emp WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteRegionEmp query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete region assignment"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Region assignment not found"})
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Region assignment deleted"}})
}
