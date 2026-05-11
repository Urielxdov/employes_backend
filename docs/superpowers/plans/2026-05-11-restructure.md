# Restructure + DB Fix Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Separar handlers en un archivo por entidad dentro de `handlers/`, corregir el nombre de tabla `employee` (singular) en todas las queries, y simplificar `main.go`.

**Architecture:** Paquete `handlers` con un struct por entidad que recibe `*sql.DB`. `router.go` en `package main` instancia todos los handlers y registra rutas. `models.go` se mueve a `handlers/models.go`.

**Tech Stack:** Go 1.21, `database/sql`, `github.com/go-sql-driver/mysql`, `net/http` stdlib. Módulo: `github.com/user/employees-api`.

---

## File Map

| Acción | Archivo |
|---|---|
| Crear | `handlers/models.go` |
| Crear | `handlers/common.go` |
| Crear | `handlers/employees.go` |
| Crear | `handlers/departments.go` |
| Crear | `handlers/salaries.go` |
| Crear | `handlers/titles.go` |
| Crear | `handlers/dept_emp.go` |
| Crear | `handlers/dept_manager.go` |
| Crear | `handlers/salary_groups.go` |
| Crear | `handlers/sg_emp.go` |
| Crear | `handlers/countries.go` |
| Crear | `handlers/regions.go` |
| Crear | `handlers/region_emp.go` |
| Crear | `router.go` |
| Modificar | `main.go` |
| Eliminar | `handlers.go` |
| Eliminar | `models.go` |

---

### Task 1: handlers/models.go + handlers/common.go

**Files:**
- Create: `handlers/models.go`
- Create: `handlers/common.go`

- [ ] **Step 1: Crear `handlers/models.go`**

```go
package handlers

import "time"

type Employee struct {
	EmpNo              int       `json:"emp_no"`
	EmployeeID         string    `json:"employee_id"`
	DateOfBirth        time.Time `json:"date_of_birth"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	MiddleNames        string    `json:"middle_names"`
	Gender             string    `json:"gender"`
	DateOfHiring       time.Time `json:"date_of_hiring"`
	DateOfTermination  time.Time `json:"date_of_termination"`
	DateOfProbationEnd time.Time `json:"date_of_probation_end"`
}

type Department struct {
	DeptNo   string `json:"dept_no"`
	DeptName string `json:"dept_name"`
}

type Salary struct {
	EmpNo    int       `json:"emp_no"`
	Salary   int       `json:"salary"`
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}

type Title struct {
	EmpNo    int       `json:"emp_no"`
	Title    string    `json:"title"`
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}

type DeptEmp struct {
	EmpNo    int       `json:"emp_no"`
	DeptNo   string    `json:"dept_no"`
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}

type DeptManager struct {
	EmpNo    int       `json:"emp_no"`
	DeptNo   string    `json:"dept_no"`
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}

type SalaryGroup struct {
	SgNo       int       `json:"sg_no"`
	SgName     string    `json:"sg_name"`
	BaseSalary float64   `json:"base_salary"`
	FromDate   time.Time `json:"from_date"`
	ToDate     time.Time `json:"to_date"`
}

type SgEmp struct {
	EmpNo    int       `json:"emp_no"`
	SgNo     int       `json:"sg_no"`
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}

type Country struct {
	ID        int     `json:"id"`
	ISO       string  `json:"iso"`
	Name      string  `json:"name"`
	NiceName  string  `json:"nicename"`
	ISO3      *string `json:"iso3"`
	NumCode   *int    `json:"numcode"`
	PhoneCode int     `json:"phonecode"`
}

type Region struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	NiceName string `json:"nicename"`
	Note     string `json:"note"`
	Country  int    `json:"country"`
}

type RegionEmp struct {
	EmpNo    int       `json:"emp_no"`
	RegionID int       `json:"region_id"`
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
```

- [ ] **Step 2: Crear `handlers/common.go`**

```go
package handlers

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
```

- [ ] **Step 3: Verificar compilación parcial**

```bash
cd /home/uhernand/employes_backend && go build ./handlers/...
```

Esperado: sin output (compila bien). Si hay errores de tipos, corregir antes de continuar.

---

### Task 2: handlers/employees.go

**Files:**
- Create: `handlers/employees.go`

- [ ] **Step 1: Crear `handlers/employees.go`**

```go
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
```

- [ ] **Step 2: Verificar compilación**

```bash
cd /home/uhernand/employes_backend && go build ./handlers/...
```

Esperado: sin output.

---

### Task 3: handlers/departments.go

**Files:**
- Create: `handlers/departments.go`

- [ ] **Step 1: Crear `handlers/departments.go`**

```go
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
```

- [ ] **Step 2: Verificar compilación**

```bash
cd /home/uhernand/employes_backend && go build ./handlers/...
```

Esperado: sin output.

---

### Task 4: handlers/salaries.go

**Files:**
- Create: `handlers/salaries.go`

- [ ] **Step 1: Crear `handlers/salaries.go`**

```go
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
```

- [ ] **Step 2: Verificar compilación**

```bash
cd /home/uhernand/employes_backend && go build ./handlers/...
```

Esperado: sin output.

---

### Task 5: handlers/titles.go

**Files:**
- Create: `handlers/titles.go`

- [ ] **Step 1: Crear `handlers/titles.go`**

```go
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type TitleHandler struct{ DB *sql.DB }

func NewTitleHandler(db *sql.DB) *TitleHandler {
	return &TitleHandler{DB: db}
}

func (h *TitleHandler) Handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/titles"), "/")
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

func (h *TitleHandler) list(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.DB.Query("SELECT emp_no, title, from_date, to_date FROM titles LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listTitles query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	titles := []Title{}
	for rows.Next() {
		var t Title
		if err := rows.Scan(&t.EmpNo, &t.Title, &t.FromDate, &t.ToDate); err != nil {
			log.Printf("[ERROR] listTitles scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		titles = append(titles, t)
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: titles})
}

func (h *TitleHandler) get(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/titles/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	rows, err := h.DB.Query("SELECT emp_no, title, from_date, to_date FROM titles WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] getTitle query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	titles := []Title{}
	for rows.Next() {
		var t Title
		if err := rows.Scan(&t.EmpNo, &t.Title, &t.FromDate, &t.ToDate); err != nil {
			log.Printf("[ERROR] getTitle scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		titles = append(titles, t)
	}

	if len(titles) == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "No titles found"})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: titles})
}

func (h *TitleHandler) create(w http.ResponseWriter, r *http.Request) {
	var t Title
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	_, err := h.DB.Exec("INSERT INTO titles (emp_no, title, from_date, to_date) VALUES (?, ?, ?, ?)", t.EmpNo, t.Title, t.FromDate, t.ToDate)
	if err != nil {
		log.Printf("[ERROR] createTitle insert failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create title"})
		return
	}
	WriteJSON(w, http.StatusCreated, APIResponse{Success: true, Data: t})
}

func (h *TitleHandler) update(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/titles/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var t Title
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec(
		"UPDATE titles SET title = ?, from_date = ?, to_date = ? WHERE emp_no = ? AND from_date = ?",
		t.Title, t.FromDate, t.ToDate, empNo, t.FromDate,
	)
	if err != nil {
		log.Printf("[ERROR] updateTitle query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update title"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Title record not found"})
		return
	}

	t.EmpNo = empNo
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: t})
}

func (h *TitleHandler) delete(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/titles/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	result, err := h.DB.Exec("DELETE FROM titles WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteTitle query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete title"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Title record not found"})
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Title deleted"}})
}
```

- [ ] **Step 2: Verificar compilación**

```bash
cd /home/uhernand/employes_backend && go build ./handlers/...
```

Esperado: sin output.

---

### Task 6: handlers/dept_emp.go

**Files:**
- Create: `handlers/dept_emp.go`

- [ ] **Step 1: Crear `handlers/dept_emp.go`**

```go
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
```

- [ ] **Step 2: Verificar compilación**

```bash
cd /home/uhernand/employes_backend && go build ./handlers/...
```

Esperado: sin output.

---

### Task 7: handlers/dept_manager.go

**Files:**
- Create: `handlers/dept_manager.go`

- [ ] **Step 1: Crear `handlers/dept_manager.go`**

```go
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
```

- [ ] **Step 2: Verificar compilación**

```bash
cd /home/uhernand/employes_backend && go build ./handlers/...
```

Esperado: sin output.

---

### Task 8: handlers/salary_groups.go

**Files:**
- Create: `handlers/salary_groups.go`

- [ ] **Step 1: Crear `handlers/salary_groups.go`**

```go
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
```

- [ ] **Step 2: Verificar compilación**

```bash
cd /home/uhernand/employes_backend && go build ./handlers/...
```

Esperado: sin output.

---

### Task 9: handlers/sg_emp.go

**Files:**
- Create: `handlers/sg_emp.go`

- [ ] **Step 1: Crear `handlers/sg_emp.go`**

```go
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
```

- [ ] **Step 2: Verificar compilación**

```bash
cd /home/uhernand/employes_backend && go build ./handlers/...
```

Esperado: sin output.

---

### Task 10: handlers/countries.go

**Files:**
- Create: `handlers/countries.go`

- [ ] **Step 1: Crear `handlers/countries.go`**

```go
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type CountryHandler struct{ DB *sql.DB }

func NewCountryHandler(db *sql.DB) *CountryHandler {
	return &CountryHandler{DB: db}
}

func (h *CountryHandler) Handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/countries"), "/")
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

func scanCountry(rows interface{ Scan(dest ...interface{}) error }, c *Country) error {
	var iso3 sql.NullString
	var numcode sql.NullInt64
	if err := rows.Scan(&c.ID, &c.ISO, &c.Name, &c.NiceName, &iso3, &numcode, &c.PhoneCode); err != nil {
		return err
	}
	if iso3.Valid {
		c.ISO3 = &iso3.String
	}
	if numcode.Valid {
		n := int(numcode.Int64)
		c.NumCode = &n
	}
	return nil
}

func (h *CountryHandler) list(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.DB.Query("SELECT id, iso, name, nicename, iso3, numcode, phonecode FROM country LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listCountries query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []Country{}
	for rows.Next() {
		var c Country
		if err := scanCountry(rows, &c); err != nil {
			log.Printf("[ERROR] listCountries scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, c)
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func (h *CountryHandler) get(w http.ResponseWriter, r *http.Request) {
	countryID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/countries/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid country ID"})
		return
	}

	var c Country
	err = scanCountry(h.DB.QueryRow("SELECT id, iso, name, nicename, iso3, numcode, phonecode FROM country WHERE id = ?", countryID), &c)
	if err == sql.ErrNoRows {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Country not found"})
		return
	}
	if err != nil {
		log.Printf("[ERROR] getCountry query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: c})
}

func (h *CountryHandler) create(w http.ResponseWriter, r *http.Request) {
	var c Country
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec("INSERT INTO country (iso, name, nicename, iso3, numcode, phonecode) VALUES (?, ?, ?, ?, ?, ?)", c.ISO, c.Name, c.NiceName, c.ISO3, c.NumCode, c.PhoneCode)
	if err != nil {
		log.Printf("[ERROR] createCountry insert failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create country"})
		return
	}

	id, _ := result.LastInsertId()
	c.ID = int(id)
	WriteJSON(w, http.StatusCreated, APIResponse{Success: true, Data: c})
}

func (h *CountryHandler) update(w http.ResponseWriter, r *http.Request) {
	countryID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/countries/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid country ID"})
		return
	}

	var c Country
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec("UPDATE country SET iso = ?, name = ?, nicename = ?, iso3 = ?, numcode = ?, phonecode = ? WHERE id = ?", c.ISO, c.Name, c.NiceName, c.ISO3, c.NumCode, c.PhoneCode, countryID)
	if err != nil {
		log.Printf("[ERROR] updateCountry query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update country"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Country not found"})
		return
	}

	c.ID = countryID
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: c})
}

func (h *CountryHandler) delete(w http.ResponseWriter, r *http.Request) {
	countryID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/countries/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid country ID"})
		return
	}

	result, err := h.DB.Exec("DELETE FROM country WHERE id = ?", countryID)
	if err != nil {
		log.Printf("[ERROR] deleteCountry query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete country"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Country not found"})
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Country deleted"}})
}
```

- [ ] **Step 2: Verificar compilación**

```bash
cd /home/uhernand/employes_backend && go build ./handlers/...
```

Esperado: sin output.

---

### Task 11: handlers/regions.go

**Files:**
- Create: `handlers/regions.go`

- [ ] **Step 1: Crear `handlers/regions.go`**

```go
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type RegionHandler struct{ DB *sql.DB }

func NewRegionHandler(db *sql.DB) *RegionHandler {
	return &RegionHandler{DB: db}
}

func (h *RegionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/regions"), "/")
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

func (h *RegionHandler) list(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.DB.Query("SELECT id, name, nicename, note, country FROM region LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listRegions query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []Region{}
	for rows.Next() {
		var rg Region
		if err := rows.Scan(&rg.ID, &rg.Name, &rg.NiceName, &rg.Note, &rg.Country); err != nil {
			log.Printf("[ERROR] listRegions scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, rg)
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func (h *RegionHandler) get(w http.ResponseWriter, r *http.Request) {
	regionID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/regions/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid region ID"})
		return
	}

	var rg Region
	err = h.DB.QueryRow("SELECT id, name, nicename, note, country FROM region WHERE id = ?", regionID).
		Scan(&rg.ID, &rg.Name, &rg.NiceName, &rg.Note, &rg.Country)
	if err == sql.ErrNoRows {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Region not found"})
		return
	}
	if err != nil {
		log.Printf("[ERROR] getRegion query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: rg})
}

func (h *RegionHandler) create(w http.ResponseWriter, r *http.Request) {
	var rg Region
	if err := json.NewDecoder(r.Body).Decode(&rg); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec("INSERT INTO region (name, nicename, note, country) VALUES (?, ?, ?, ?)", rg.Name, rg.NiceName, rg.Note, rg.Country)
	if err != nil {
		log.Printf("[ERROR] createRegion insert failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create region"})
		return
	}

	id, _ := result.LastInsertId()
	rg.ID = int(id)
	WriteJSON(w, http.StatusCreated, APIResponse{Success: true, Data: rg})
}

func (h *RegionHandler) update(w http.ResponseWriter, r *http.Request) {
	regionID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/regions/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid region ID"})
		return
	}

	var rg Region
	if err := json.NewDecoder(r.Body).Decode(&rg); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec("UPDATE region SET name = ?, nicename = ?, note = ?, country = ? WHERE id = ?", rg.Name, rg.NiceName, rg.Note, rg.Country, regionID)
	if err != nil {
		log.Printf("[ERROR] updateRegion query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update region"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Region not found"})
		return
	}

	rg.ID = regionID
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: rg})
}

func (h *RegionHandler) delete(w http.ResponseWriter, r *http.Request) {
	regionID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/regions/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid region ID"})
		return
	}

	result, err := h.DB.Exec("DELETE FROM region WHERE id = ?", regionID)
	if err != nil {
		log.Printf("[ERROR] deleteRegion query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete region"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Region not found"})
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Region deleted"}})
}
```

- [ ] **Step 2: Verificar compilación**

```bash
cd /home/uhernand/employes_backend && go build ./handlers/...
```

Esperado: sin output.

---

### Task 12: handlers/region_emp.go

**Files:**
- Create: `handlers/region_emp.go`

- [ ] **Step 1: Crear `handlers/region_emp.go`**

```go
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
```

- [ ] **Step 2: Verificar compilación**

```bash
cd /home/uhernand/employes_backend && go build ./handlers/...
```

Esperado: sin output.

---

### Task 13: router.go + health check

**Files:**
- Create: `router.go`

- [ ] **Step 1: Crear `router.go`**

```go
package main

import (
	"database/sql"
	"net/http"

	"github.com/user/employees-api/handlers"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func setupRouter(database *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if err := database.Ping(); err != nil {
			handlers.WriteJSON(w, http.StatusInternalServerError, handlers.APIResponse{Success: false, Error: "Database unavailable"})
			return
		}
		handlers.WriteJSON(w, http.StatusOK, handlers.APIResponse{Success: true, Data: map[string]string{"status": "ok"}})
	}))

	eh := handlers.NewEmployeeHandler(database)
	mux.HandleFunc("/api/v1/employees", corsMiddleware(eh.Handle))
	mux.HandleFunc("/api/v1/employees/", corsMiddleware(eh.Handle))

	dh := handlers.NewDepartmentHandler(database)
	mux.HandleFunc("/api/v1/departments", corsMiddleware(dh.Handle))
	mux.HandleFunc("/api/v1/departments/", corsMiddleware(dh.Handle))

	sh := handlers.NewSalaryHandler(database)
	mux.HandleFunc("/api/v1/salaries", corsMiddleware(sh.Handle))
	mux.HandleFunc("/api/v1/salaries/", corsMiddleware(sh.Handle))

	th := handlers.NewTitleHandler(database)
	mux.HandleFunc("/api/v1/titles", corsMiddleware(th.Handle))
	mux.HandleFunc("/api/v1/titles/", corsMiddleware(th.Handle))

	deh := handlers.NewDeptEmpHandler(database)
	mux.HandleFunc("/api/v1/dept_emp", corsMiddleware(deh.Handle))
	mux.HandleFunc("/api/v1/dept_emp/", corsMiddleware(deh.Handle))

	dmh := handlers.NewDeptManagerHandler(database)
	mux.HandleFunc("/api/v1/dept_manager", corsMiddleware(dmh.Handle))
	mux.HandleFunc("/api/v1/dept_manager/", corsMiddleware(dmh.Handle))

	sgh := handlers.NewSalaryGroupHandler(database)
	mux.HandleFunc("/api/v1/salary_groups", corsMiddleware(sgh.Handle))
	mux.HandleFunc("/api/v1/salary_groups/", corsMiddleware(sgh.Handle))

	seh := handlers.NewSgEmpHandler(database)
	mux.HandleFunc("/api/v1/sg_emp", corsMiddleware(seh.Handle))
	mux.HandleFunc("/api/v1/sg_emp/", corsMiddleware(seh.Handle))

	ch := handlers.NewCountryHandler(database)
	mux.HandleFunc("/api/v1/countries", corsMiddleware(ch.Handle))
	mux.HandleFunc("/api/v1/countries/", corsMiddleware(ch.Handle))

	rh := handlers.NewRegionHandler(database)
	mux.HandleFunc("/api/v1/regions", corsMiddleware(rh.Handle))
	mux.HandleFunc("/api/v1/regions/", corsMiddleware(rh.Handle))

	reh := handlers.NewRegionEmpHandler(database)
	mux.HandleFunc("/api/v1/region_emp", corsMiddleware(reh.Handle))
	mux.HandleFunc("/api/v1/region_emp/", corsMiddleware(reh.Handle))

	return mux
}
```

- [ ] **Step 2: Verificar compilación (aún fallará — main.go tiene conflictos)**

```bash
cd /home/uhernand/employes_backend && go build ./... 2>&1 | head -20
```

Esperado: errores por duplicados (`corsMiddleware`, `handleEmployees`, etc. definidos tanto en `main.go` como en el nuevo `router.go`). Esto es normal — se resuelve en el Task 14.

---

### Task 14: Reemplazar main.go + eliminar archivos viejos

**Files:**
- Modify: `main.go`
- Delete: `handlers.go`
- Delete: `models.go`

- [ ] **Step 1: Reemplazar contenido de `main.go`**

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("[INIT] Connecting to database...")
	if err := initDB(); err != nil {
		log.Fatalf("[FATAL] Failed to initialize database: %v", err)
	}
	defer closeDB()

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[INFO] Starting server on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), setupRouter(db)); err != nil {
		log.Fatalf("[FATAL] Server failed: %v", err)
	}
}
```

- [ ] **Step 2: Eliminar `handlers.go`**

```bash
rm /home/uhernand/employes_backend/handlers.go
```

- [ ] **Step 3: Eliminar `models.go`**

```bash
rm /home/uhernand/employes_backend/models.go
```

- [ ] **Step 4: Verificar compilación completa**

```bash
cd /home/uhernand/employes_backend && go build ./...
```

Esperado: sin output (compilación limpia). Si hay errores de símbolo no encontrado, verificar que todos los archivos en `handlers/` tienen `package handlers` y que `router.go` importa `github.com/user/employees-api/handlers`.

---

### Task 15: Commit final

- [ ] **Step 1: Verificar estado**

```bash
cd /home/uhernand/employes_backend && git status
```

Esperado: archivos nuevos en `handlers/`, `router.go` nuevo, `main.go` modificado, `handlers.go` y `models.go` eliminados.

- [ ] **Step 2: Verificar que el binario levanta**

```bash
cd /home/uhernand/employes_backend && go vet ./...
```

Esperado: sin output.

- [ ] **Step 3: Commit atómico**

```bash
cd /home/uhernand/employes_backend && git add -A && git commit -m "refactor: split handlers by entity, fix employee table name

- handlers/ package with one file per entity (11 handlers)
- Handler structs receive *sql.DB via constructor (no globals)
- router.go wires all handlers; main.go is entry-point only
- models.go moved to handlers/models.go (package handlers)
- Fix: queries used 'employees' (plural), schema defines 'employee'"
```
