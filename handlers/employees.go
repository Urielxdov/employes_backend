package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type EmployeeHandler struct{ DB *sql.DB }

func NewEmployeeHandler(db *sql.DB) *EmployeeHandler {
	return &EmployeeHandler{DB: db}
}

func (h *EmployeeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/employees"), "/")
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

func (h *EmployeeHandler) list(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.DB.Query("SELECT emp_no, employee_id, date_of_birth, first_name, last_name, middle_names, gender, date_of_hiring, date_of_termination, date_of_probation_end FROM employee LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listEmployees query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	employees := []Employee{}
	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.EmpNo, &e.EmployeeID, &e.DateOfBirth, &e.FirstName, &e.LastName, &e.MiddleNames, &e.Gender, &e.DateOfHiring, &e.DateOfTermination, &e.DateOfProbationEnd); err != nil {
			log.Printf("[ERROR] listEmployees scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		employees = append(employees, e)
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: employees})
}

func (h *EmployeeHandler) get(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/employees/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var e Employee
	err = h.DB.QueryRow("SELECT emp_no, employee_id, date_of_birth, first_name, last_name, middle_names, gender, date_of_hiring, date_of_termination, date_of_probation_end FROM employee WHERE emp_no = ?", empNo).
		Scan(&e.EmpNo, &e.EmployeeID, &e.DateOfBirth, &e.FirstName, &e.LastName, &e.MiddleNames, &e.Gender, &e.DateOfHiring, &e.DateOfTermination, &e.DateOfProbationEnd)
	if err == sql.ErrNoRows {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Employee not found"})
		return
	}
	if err != nil {
		log.Printf("[ERROR] getEmployee query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: e})
}

func (h *EmployeeHandler) create(w http.ResponseWriter, r *http.Request) {
	var e Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec(
		"INSERT INTO employee (employee_id, date_of_birth, first_name, last_name, middle_names, gender, date_of_hiring, date_of_termination, date_of_probation_end) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		e.EmployeeID, e.DateOfBirth, e.FirstName, e.LastName, e.MiddleNames, e.Gender, e.DateOfHiring, e.DateOfTermination, e.DateOfProbationEnd,
	)
	if err != nil {
		log.Printf("[ERROR] createEmployee insert failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create employee"})
		return
	}

	id, _ := result.LastInsertId()
	e.EmpNo = int(id)
	WriteJSON(w, http.StatusCreated, APIResponse{Success: true, Data: e})
}

func (h *EmployeeHandler) update(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/employees/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var e Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec(
		"UPDATE employee SET employee_id = ?, date_of_birth = ?, first_name = ?, last_name = ?, middle_names = ?, gender = ?, date_of_hiring = ?, date_of_termination = ?, date_of_probation_end = ? WHERE emp_no = ?",
		e.EmployeeID, e.DateOfBirth, e.FirstName, e.LastName, e.MiddleNames, e.Gender, e.DateOfHiring, e.DateOfTermination, e.DateOfProbationEnd, empNo,
	)
	if err != nil {
		log.Printf("[ERROR] updateEmployee query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update employee"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Employee not found"})
		return
	}

	e.EmpNo = empNo
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: e})
}

func (h *EmployeeHandler) delete(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/employees/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	result, err := h.DB.Exec("DELETE FROM employee WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteEmployee query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete employee"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Employee not found"})
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Employee deleted"}})
}
