package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// GET /api/v1/employees
func listEmployees(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT emp_no, birth_date, first_name, last_name, gender, hire_date FROM employees LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listEmployees query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	employees := []Employee{}
	for rows.Next() {
		var e Employee
		err := rows.Scan(&e.EmpNo, &e.BirthDate, &e.FirstName, &e.LastName, &e.Gender, &e.HireDate)
		if err != nil {
			log.Printf("[ERROR] listEmployees scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		employees = append(employees, e)
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: employees})
}

// GET /api/v1/employees/:id
func getEmployee(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/employees/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var e Employee
	err = db.QueryRow("SELECT emp_no, birth_date, first_name, last_name, gender, hire_date FROM employees WHERE emp_no = ?", empNo).
		Scan(&e.EmpNo, &e.BirthDate, &e.FirstName, &e.LastName, &e.Gender, &e.HireDate)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Employee not found"})
		} else {
			log.Printf("[ERROR] getEmployee query failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		}
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: e})
}

// GET /api/v1/departments
func listDepartments(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT dept_no, dept_name FROM departments")
	if err != nil {
		log.Printf("[ERROR] listDepartments query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	departments := []Department{}
	for rows.Next() {
		var d Department
		err := rows.Scan(&d.DeptNo, &d.DeptName)
		if err != nil {
			log.Printf("[ERROR] listDepartments scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		departments = append(departments, d)
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: departments})
}

// GET /api/v1/departments/:id
func getDepartment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/departments/")

	var d Department
	err := db.QueryRow("SELECT dept_no, dept_name FROM departments WHERE dept_no = ?", id).
		Scan(&d.DeptNo, &d.DeptName)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Department not found"})
		} else {
			log.Printf("[ERROR] getDepartment query failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		}
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: d})
}

// GET /api/v1/salaries
func listSalaries(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT emp_no, salary, from_date, to_date FROM salaries LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listSalaries query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	salaries := []Salary{}
	for rows.Next() {
		var s Salary
		err := rows.Scan(&s.EmpNo, &s.Salary, &s.FromDate, &s.ToDate)
		if err != nil {
			log.Printf("[ERROR] listSalaries scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		salaries = append(salaries, s)
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: salaries})
}

// GET /api/v1/salaries/:id
func getSalary(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/salaries/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	rows, err := db.Query("SELECT emp_no, salary, from_date, to_date FROM salaries WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] getSalary query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	salaries := []Salary{}
	for rows.Next() {
		var s Salary
		err := rows.Scan(&s.EmpNo, &s.Salary, &s.FromDate, &s.ToDate)
		if err != nil {
			log.Printf("[ERROR] getSalary scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		salaries = append(salaries, s)
	}

	if len(salaries) == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "No salaries found"})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: salaries})
}

// GET /api/v1/titles
func listTitles(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT emp_no, title, from_date, to_date FROM titles LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listTitles query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	titles := []Title{}
	for rows.Next() {
		var t Title
		err := rows.Scan(&t.EmpNo, &t.Title, &t.FromDate, &t.ToDate)
		if err != nil {
			log.Printf("[ERROR] listTitles scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		titles = append(titles, t)
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: titles})
}

// GET /api/v1/titles/:id
func getTitle(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/titles/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	rows, err := db.Query("SELECT emp_no, title, from_date, to_date FROM titles WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] getTitle query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	titles := []Title{}
	for rows.Next() {
		var t Title
		err := rows.Scan(&t.EmpNo, &t.Title, &t.FromDate, &t.ToDate)
		if err != nil {
			log.Printf("[ERROR] getTitle scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		titles = append(titles, t)
	}

	if len(titles) == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "No titles found"})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: titles})
}

// GET /api/v1/dept_emp
func listDeptEmp(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT emp_no, dept_no, from_date, to_date FROM dept_emp LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listDeptEmp query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	deptEmps := []DeptEmp{}
	for rows.Next() {
		var de DeptEmp
		err := rows.Scan(&de.EmpNo, &de.DeptNo, &de.FromDate, &de.ToDate)
		if err != nil {
			log.Printf("[ERROR] listDeptEmp scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		deptEmps = append(deptEmps, de)
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: deptEmps})
}

// GET /api/v1/dept_emp/:id
func getDeptEmp(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/dept_emp/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	rows, err := db.Query("SELECT emp_no, dept_no, from_date, to_date FROM dept_emp WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] getDeptEmp query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	deptEmps := []DeptEmp{}
	for rows.Next() {
		var de DeptEmp
		err := rows.Scan(&de.EmpNo, &de.DeptNo, &de.FromDate, &de.ToDate)
		if err != nil {
			log.Printf("[ERROR] getDeptEmp scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		deptEmps = append(deptEmps, de)
	}

	if len(deptEmps) == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "No department assignments found"})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: deptEmps})
}

// GET /api/v1/dept_manager
func listDeptManager(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT emp_no, dept_no, from_date, to_date FROM dept_manager LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listDeptManager query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	deptMgrs := []DeptManager{}
	for rows.Next() {
		var dm DeptManager
		err := rows.Scan(&dm.EmpNo, &dm.DeptNo, &dm.FromDate, &dm.ToDate)
		if err != nil {
			log.Printf("[ERROR] listDeptManager scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		deptMgrs = append(deptMgrs, dm)
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: deptMgrs})
}

// GET /api/v1/dept_manager/:id
func getDeptManager(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/dept_manager/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	rows, err := db.Query("SELECT emp_no, dept_no, from_date, to_date FROM dept_manager WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] getDeptManager query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	deptMgrs := []DeptManager{}
	for rows.Next() {
		var dm DeptManager
		err := rows.Scan(&dm.EmpNo, &dm.DeptNo, &dm.FromDate, &dm.ToDate)
		if err != nil {
			log.Printf("[ERROR] getDeptManager scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		deptMgrs = append(deptMgrs, dm)
	}

	if len(deptMgrs) == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "No manager assignments found"})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: deptMgrs})
}

// GET /health
func healthCheck(w http.ResponseWriter, r *http.Request) {
	err := db.Ping()
	if err != nil {
		log.Printf("[ERROR] healthcheck failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"status": "ok"}})
}

// POST /api/v1/employees
func createEmployee(w http.ResponseWriter, r *http.Request) {
	var e Employee
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := db.Exec(
		"INSERT INTO employees (emp_no, birth_date, first_name, last_name, gender, hire_date) VALUES (?, ?, ?, ?, ?, ?)",
		e.EmpNo, e.BirthDate, e.FirstName, e.LastName, e.Gender, e.HireDate,
	)
	if err != nil {
		log.Printf("[ERROR] createEmployee insert failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create employee"})
		return
	}

	id, _ := result.LastInsertId()
	e.EmpNo = int(id)
	writeJSON(w, http.StatusCreated, APIResponse{Success: true, Data: e})
}

// POST /api/v1/departments
func createDepartment(w http.ResponseWriter, r *http.Request) {
	var d Department
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	_, err = db.Exec(
		"INSERT INTO departments (dept_no, dept_name) VALUES (?, ?)",
		d.DeptNo, d.DeptName,
	)
	if err != nil {
		log.Printf("[ERROR] createDepartment insert failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create department"})
		return
	}

	writeJSON(w, http.StatusCreated, APIResponse{Success: true, Data: d})
}

// POST /api/v1/salaries
func createSalary(w http.ResponseWriter, r *http.Request) {
	var s Salary
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	_, err = db.Exec(
		"INSERT INTO salaries (emp_no, salary, from_date, to_date) VALUES (?, ?, ?, ?)",
		s.EmpNo, s.Salary, s.FromDate, s.ToDate,
	)
	if err != nil {
		log.Printf("[ERROR] createSalary insert failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create salary"})
		return
	}

	writeJSON(w, http.StatusCreated, APIResponse{Success: true, Data: s})
}

// POST /api/v1/titles
func createTitle(w http.ResponseWriter, r *http.Request) {
	var t Title
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	_, err = db.Exec(
		"INSERT INTO titles (emp_no, title, from_date, to_date) VALUES (?, ?, ?, ?)",
		t.EmpNo, t.Title, t.FromDate, t.ToDate,
	)
	if err != nil {
		log.Printf("[ERROR] createTitle insert failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create title"})
		return
	}

	writeJSON(w, http.StatusCreated, APIResponse{Success: true, Data: t})
}

// POST /api/v1/dept_emp
func createDeptEmp(w http.ResponseWriter, r *http.Request) {
	var de DeptEmp
	err := json.NewDecoder(r.Body).Decode(&de)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	_, err = db.Exec(
		"INSERT INTO dept_emp (emp_no, dept_no, from_date, to_date) VALUES (?, ?, ?, ?)",
		de.EmpNo, de.DeptNo, de.FromDate, de.ToDate,
	)
	if err != nil {
		log.Printf("[ERROR] createDeptEmp insert failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create department assignment"})
		return
	}

	writeJSON(w, http.StatusCreated, APIResponse{Success: true, Data: de})
}

// POST /api/v1/dept_manager
func createDeptManager(w http.ResponseWriter, r *http.Request) {
	var dm DeptManager
	err := json.NewDecoder(r.Body).Decode(&dm)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	_, err = db.Exec(
		"INSERT INTO dept_manager (emp_no, dept_no, from_date, to_date) VALUES (?, ?, ?, ?)",
		dm.EmpNo, dm.DeptNo, dm.FromDate, dm.ToDate,
	)
	if err != nil {
		log.Printf("[ERROR] createDeptManager insert failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create manager assignment"})
		return
	}

	writeJSON(w, http.StatusCreated, APIResponse{Success: true, Data: dm})
}

// PUT /api/v1/employees/:id
func updateEmployee(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/employees/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var e Employee
	err = json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := db.Exec(
		"UPDATE employees SET birth_date = ?, first_name = ?, last_name = ?, gender = ?, hire_date = ? WHERE emp_no = ?",
		e.BirthDate, e.FirstName, e.LastName, e.Gender, e.HireDate, empNo,
	)
	if err != nil {
		log.Printf("[ERROR] updateEmployee query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update employee"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Employee not found"})
		return
	}

	e.EmpNo = empNo
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: e})
}

// PUT /api/v1/departments/:id
func updateDepartment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/departments/")

	var d Department
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := db.Exec(
		"UPDATE departments SET dept_name = ? WHERE dept_no = ?",
		d.DeptName, id,
	)
	if err != nil {
		log.Printf("[ERROR] updateDepartment query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update department"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Department not found"})
		return
	}

	d.DeptNo = id
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: d})
}

// PUT /api/v1/salaries/:id
func updateSalary(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/salaries/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var s Salary
	err = json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := db.Exec(
		"UPDATE salaries SET salary = ?, from_date = ?, to_date = ? WHERE emp_no = ? AND from_date = ?",
		s.Salary, s.FromDate, s.ToDate, empNo, s.FromDate,
	)
	if err != nil {
		log.Printf("[ERROR] updateSalary query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update salary"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Salary record not found"})
		return
	}

	s.EmpNo = empNo
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: s})
}

// PUT /api/v1/titles/:id
func updateTitle(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/titles/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var t Title
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := db.Exec(
		"UPDATE titles SET title = ?, from_date = ?, to_date = ? WHERE emp_no = ? AND from_date = ?",
		t.Title, t.FromDate, t.ToDate, empNo, t.FromDate,
	)
	if err != nil {
		log.Printf("[ERROR] updateTitle query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update title"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Title record not found"})
		return
	}

	t.EmpNo = empNo
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: t})
}

// PUT /api/v1/dept_emp/:id
func updateDeptEmp(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/dept_emp/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var de DeptEmp
	err = json.NewDecoder(r.Body).Decode(&de)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := db.Exec(
		"UPDATE dept_emp SET dept_no = ?, from_date = ?, to_date = ? WHERE emp_no = ? AND from_date = ?",
		de.DeptNo, de.FromDate, de.ToDate, empNo, de.FromDate,
	)
	if err != nil {
		log.Printf("[ERROR] updateDeptEmp query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update department assignment"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Department assignment not found"})
		return
	}

	de.EmpNo = empNo
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: de})
}

// PUT /api/v1/dept_manager/:id
func updateDeptManager(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/dept_manager/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var dm DeptManager
	err = json.NewDecoder(r.Body).Decode(&dm)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := db.Exec(
		"UPDATE dept_manager SET dept_no = ?, from_date = ?, to_date = ? WHERE emp_no = ? AND from_date = ?",
		dm.DeptNo, dm.FromDate, dm.ToDate, empNo, dm.FromDate,
	)
	if err != nil {
		log.Printf("[ERROR] updateDeptManager query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update manager assignment"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Manager assignment not found"})
		return
	}

	dm.EmpNo = empNo
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: dm})
}

// DELETE /api/v1/employees/:id
func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/employees/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	result, err := db.Exec("DELETE FROM employees WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteEmployee query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete employee"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Employee not found"})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Employee deleted"}})
}

// DELETE /api/v1/departments/:id
func deleteDepartment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/departments/")

	result, err := db.Exec("DELETE FROM departments WHERE dept_no = ?", id)
	if err != nil {
		log.Printf("[ERROR] deleteDepartment query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete department"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Department not found"})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Department deleted"}})
}

// DELETE /api/v1/salaries/:id
func deleteSalary(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/salaries/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	result, err := db.Exec("DELETE FROM salaries WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteSalary query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete salary"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Salary record not found"})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Salary deleted"}})
}

// DELETE /api/v1/titles/:id
func deleteTitle(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/titles/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	result, err := db.Exec("DELETE FROM titles WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteTitle query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete title"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Title record not found"})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Title deleted"}})
}

// DELETE /api/v1/dept_emp/:id
func deleteDeptEmp(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/dept_emp/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	result, err := db.Exec("DELETE FROM dept_emp WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteDeptEmp query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete department assignment"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Department assignment not found"})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Department assignment deleted"}})
}

// DELETE /api/v1/dept_manager/:id
func deleteDeptManager(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/dept_manager/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	result, err := db.Exec("DELETE FROM dept_manager WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteDeptManager query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete manager assignment"})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Manager assignment not found"})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Manager assignment deleted"}})
}
