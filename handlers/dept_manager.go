package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type DeptManagerHandler struct{ DB *sql.DB }

func NewDeptManagerHandler(db *sql.DB) *DeptManagerHandler {
	return &DeptManagerHandler{DB: db}
}

func (h *DeptManagerHandler) Handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/dept_manager"), "/")
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

func (h *DeptManagerHandler) list(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.DB.Query("SELECT emp_no, dept_no, from_date, to_date FROM dept_manager LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listDeptManager query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []DeptManager{}
	for rows.Next() {
		var dm DeptManager
		if err := rows.Scan(&dm.EmpNo, &dm.DeptNo, &dm.FromDate, &dm.ToDate); err != nil {
			log.Printf("[ERROR] listDeptManager scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, dm)
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func (h *DeptManagerHandler) get(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/dept_manager/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	rows, err := h.DB.Query("SELECT emp_no, dept_no, from_date, to_date FROM dept_manager WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] getDeptManager query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []DeptManager{}
	for rows.Next() {
		var dm DeptManager
		if err := rows.Scan(&dm.EmpNo, &dm.DeptNo, &dm.FromDate, &dm.ToDate); err != nil {
			log.Printf("[ERROR] getDeptManager scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, dm)
	}

	if len(items) == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "No manager assignments found"})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func (h *DeptManagerHandler) create(w http.ResponseWriter, r *http.Request) {
	var dm DeptManager
	if err := json.NewDecoder(r.Body).Decode(&dm); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	_, err := h.DB.Exec("INSERT INTO dept_manager (emp_no, dept_no, from_date, to_date) VALUES (?, ?, ?, ?)", dm.EmpNo, dm.DeptNo, dm.FromDate, dm.ToDate)
	if err != nil {
		log.Printf("[ERROR] createDeptManager insert failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create manager assignment"})
		return
	}
	WriteJSON(w, http.StatusCreated, APIResponse{Success: true, Data: dm})
}

func (h *DeptManagerHandler) update(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/dept_manager/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var dm DeptManager
	if err := json.NewDecoder(r.Body).Decode(&dm); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec(
		"UPDATE dept_manager SET dept_no = ?, from_date = ?, to_date = ? WHERE emp_no = ? AND from_date = ?",
		dm.DeptNo, dm.FromDate, dm.ToDate, empNo, dm.FromDate,
	)
	if err != nil {
		log.Printf("[ERROR] updateDeptManager query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update manager assignment"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Manager assignment not found"})
		return
	}

	dm.EmpNo = empNo
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: dm})
}

func (h *DeptManagerHandler) delete(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/dept_manager/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	result, err := h.DB.Exec("DELETE FROM dept_manager WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteDeptManager query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete manager assignment"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Manager assignment not found"})
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Manager assignment deleted"}})
}
