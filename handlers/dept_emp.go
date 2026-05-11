package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type DeptEmpHandler struct{ DB *sql.DB }

func NewDeptEmpHandler(db *sql.DB) *DeptEmpHandler {
	return &DeptEmpHandler{DB: db}
}

func (h *DeptEmpHandler) Handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/dept_emp"), "/")
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

func (h *DeptEmpHandler) list(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.DB.Query("SELECT emp_no, dept_no, from_date, to_date FROM dept_emp LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listDeptEmp query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []DeptEmp{}
	for rows.Next() {
		var de DeptEmp
		if err := rows.Scan(&de.EmpNo, &de.DeptNo, &de.FromDate, &de.ToDate); err != nil {
			log.Printf("[ERROR] listDeptEmp scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, de)
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func (h *DeptEmpHandler) get(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/dept_emp/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	rows, err := h.DB.Query("SELECT emp_no, dept_no, from_date, to_date FROM dept_emp WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] getDeptEmp query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []DeptEmp{}
	for rows.Next() {
		var de DeptEmp
		if err := rows.Scan(&de.EmpNo, &de.DeptNo, &de.FromDate, &de.ToDate); err != nil {
			log.Printf("[ERROR] getDeptEmp scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, de)
	}

	if len(items) == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "No department assignments found"})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func (h *DeptEmpHandler) create(w http.ResponseWriter, r *http.Request) {
	var de DeptEmp
	if err := json.NewDecoder(r.Body).Decode(&de); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	_, err := h.DB.Exec("INSERT INTO dept_emp (emp_no, dept_no, from_date, to_date) VALUES (?, ?, ?, ?)", de.EmpNo, de.DeptNo, de.FromDate, de.ToDate)
	if err != nil {
		log.Printf("[ERROR] createDeptEmp insert failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create department assignment"})
		return
	}
	WriteJSON(w, http.StatusCreated, APIResponse{Success: true, Data: de})
}

func (h *DeptEmpHandler) update(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/dept_emp/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var de DeptEmp
	if err := json.NewDecoder(r.Body).Decode(&de); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec(
		"UPDATE dept_emp SET dept_no = ?, from_date = ?, to_date = ? WHERE emp_no = ? AND from_date = ?",
		de.DeptNo, de.FromDate, de.ToDate, empNo, de.FromDate,
	)
	if err != nil {
		log.Printf("[ERROR] updateDeptEmp query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update department assignment"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Department assignment not found"})
		return
	}

	de.EmpNo = empNo
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: de})
}

func (h *DeptEmpHandler) delete(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/dept_emp/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	result, err := h.DB.Exec("DELETE FROM dept_emp WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteDeptEmp query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete department assignment"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Department assignment not found"})
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Department assignment deleted"}})
}
