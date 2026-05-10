package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func scanCountry(rows interface {
	Scan(dest ...interface{}) error
}, c *Country) error {
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

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// GET /api/v1/employees
func listEmployees(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT emp_no, employee_id, date_of_birth, first_name, last_name, middle_names, gender, date_of_hiring, date_of_termination, date_of_probation_end FROM employees LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listEmployees query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	employees := []Employee{}
	for rows.Next() {
		var e Employee
		err := rows.Scan(&e.EmpNo, &e.EmployeeID, &e.DateOfBirth, &e.FirstName, &e.LastName, &e.MiddleNames, &e.Gender, &e.DateOfHiring, &e.DateOfTermination, &e.DateOfProbationEnd)
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
	err = db.QueryRow("SELECT emp_no, employee_id, date_of_birth, first_name, last_name, middle_names, gender, date_of_hiring, date_of_termination, date_of_probation_end FROM employees WHERE emp_no = ?", empNo).
		Scan(&e.EmpNo, &e.EmployeeID, &e.DateOfBirth, &e.FirstName, &e.LastName, &e.MiddleNames, &e.Gender, &e.DateOfHiring, &e.DateOfTermination, &e.DateOfProbationEnd)
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
	rows, err := db.Query("SELECT dept_no, dept_name FROM department")
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
	err := db.QueryRow("SELECT dept_no, dept_name FROM department WHERE dept_no = ?", id).
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
	rows, err := db.Query("SELECT emp_no, salary, from_date, to_date FROM salary LIMIT 100")
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

	rows, err := db.Query("SELECT emp_no, salary, from_date, to_date FROM salary WHERE emp_no = ?", empNo)
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
		"INSERT INTO employees (employee_id, date_of_birth, first_name, last_name, middle_names, gender, date_of_hiring, date_of_termination, date_of_probation_end) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		e.EmployeeID, e.DateOfBirth, e.FirstName, e.LastName, e.MiddleNames, e.Gender, e.DateOfHiring, e.DateOfTermination, e.DateOfProbationEnd,
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
		"INSERT INTO department (dept_no, dept_name) VALUES (?, ?)",
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
		"INSERT INTO salary (emp_no, salary, from_date, to_date) VALUES (?, ?, ?, ?)",
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
		"UPDATE employees SET employee_id = ?, date_of_birth = ?, first_name = ?, last_name = ?, middle_names = ?, gender = ?, date_of_hiring = ?, date_of_termination = ?, date_of_probation_end = ? WHERE emp_no = ?",
		e.EmployeeID, e.DateOfBirth, e.FirstName, e.LastName, e.MiddleNames, e.Gender, e.DateOfHiring, e.DateOfTermination, e.DateOfProbationEnd, empNo,
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
		"UPDATE department SET dept_name = ? WHERE dept_no = ?",
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
		"UPDATE salary SET salary = ?, from_date = ?, to_date = ? WHERE emp_no = ? AND from_date = ?",
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

	result, err := db.Exec("DELETE FROM department WHERE dept_no = ?", id)
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

	result, err := db.Exec("DELETE FROM salary WHERE emp_no = ?", empNo)
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

// GET /api/v1/salary_groups
func listSalaryGroups(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT sg_no, sg_name, base_salary, from_date, to_date FROM salary_group LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listSalaryGroups query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []SalaryGroup{}
	for rows.Next() {
		var sg SalaryGroup
		if err := rows.Scan(&sg.SgNo, &sg.SgName, &sg.BaseSalary, &sg.FromDate, &sg.ToDate); err != nil {
			log.Printf("[ERROR] listSalaryGroups scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, sg)
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

// GET /api/v1/salary_groups/:id
func getSalaryGroup(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/salary_groups/")
	sgNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid salary group ID"})
		return
	}

	rows, err := db.Query("SELECT sg_no, sg_name, base_salary, from_date, to_date FROM salary_group WHERE sg_no = ?", sgNo)
	if err != nil {
		log.Printf("[ERROR] getSalaryGroup query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []SalaryGroup{}
	for rows.Next() {
		var sg SalaryGroup
		if err := rows.Scan(&sg.SgNo, &sg.SgName, &sg.BaseSalary, &sg.FromDate, &sg.ToDate); err != nil {
			log.Printf("[ERROR] getSalaryGroup scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, sg)
	}
	if len(items) == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Salary group not found"})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func createSalaryGroup(w http.ResponseWriter, r *http.Request) {
	var sg SalaryGroup
	if err := json.NewDecoder(r.Body).Decode(&sg); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}
	result, err := db.Exec("INSERT INTO salary_group (sg_name, base_salary, from_date, to_date) VALUES (?, ?, ?, ?)", sg.SgName, sg.BaseSalary, sg.FromDate, sg.ToDate)
	if err != nil {
		log.Printf("[ERROR] createSalaryGroup insert failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create salary group"})
		return
	}
	id, _ := result.LastInsertId()
	sg.SgNo = int(id)
	writeJSON(w, http.StatusCreated, APIResponse{Success: true, Data: sg})
}

func updateSalaryGroup(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/salary_groups/")
	sgNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid salary group ID"})
		return
	}
	var sg SalaryGroup
	if err := json.NewDecoder(r.Body).Decode(&sg); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}
	result, err := db.Exec("UPDATE salary_group SET sg_name = ?, base_salary = ?, from_date = ?, to_date = ? WHERE sg_no = ? AND from_date = ?", sg.SgName, sg.BaseSalary, sg.FromDate, sg.ToDate, sgNo, sg.FromDate)
	if err != nil {
		log.Printf("[ERROR] updateSalaryGroup query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update salary group"})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Salary group not found"})
		return
	}
	sg.SgNo = sgNo
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: sg})
}

func deleteSalaryGroup(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/salary_groups/")
	sgNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid salary group ID"})
		return
	}
	result, err := db.Exec("DELETE FROM salary_group WHERE sg_no = ?", sgNo)
	if err != nil {
		log.Printf("[ERROR] deleteSalaryGroup query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete salary group"})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Salary group not found"})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Salary group deleted"}})
}

func listSgEmp(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT emp_no, sg_no, from_date, to_date FROM sg_emp LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listSgEmp query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []SgEmp{}
	for rows.Next() {
		var se SgEmp
		if err := rows.Scan(&se.EmpNo, &se.SgNo, &se.FromDate, &se.ToDate); err != nil {
			log.Printf("[ERROR] listSgEmp scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, se)
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func getSgEmp(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/sg_emp/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}
	rows, err := db.Query("SELECT emp_no, sg_no, from_date, to_date FROM sg_emp WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] getSgEmp query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []SgEmp{}
	for rows.Next() {
		var se SgEmp
		if err := rows.Scan(&se.EmpNo, &se.SgNo, &se.FromDate, &se.ToDate); err != nil {
			log.Printf("[ERROR] getSgEmp scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, se)
	}
	if len(items) == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "No salary group assignments found"})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func createSgEmp(w http.ResponseWriter, r *http.Request) {
	var se SgEmp
	if err := json.NewDecoder(r.Body).Decode(&se); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}
	_, err := db.Exec("INSERT INTO sg_emp (emp_no, sg_no, from_date, to_date) VALUES (?, ?, ?, ?)", se.EmpNo, se.SgNo, se.FromDate, se.ToDate)
	if err != nil {
		log.Printf("[ERROR] createSgEmp insert failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create salary group assignment"})
		return
	}
	writeJSON(w, http.StatusCreated, APIResponse{Success: true, Data: se})
}

func updateSgEmp(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/sg_emp/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}
	var se SgEmp
	if err := json.NewDecoder(r.Body).Decode(&se); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}
	result, err := db.Exec("UPDATE sg_emp SET sg_no = ?, from_date = ?, to_date = ? WHERE emp_no = ? AND from_date = ?", se.SgNo, se.FromDate, se.ToDate, empNo, se.FromDate)
	if err != nil {
		log.Printf("[ERROR] updateSgEmp query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update salary group assignment"})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Salary group assignment not found"})
		return
	}
	se.EmpNo = empNo
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: se})
}

func deleteSgEmp(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/sg_emp/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}
	result, err := db.Exec("DELETE FROM sg_emp WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteSgEmp query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete salary group assignment"})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Salary group assignment not found"})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Salary group assignment deleted"}})
}

func listCountries(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, iso, name, nicename, iso3, numcode, phonecode FROM country LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listCountries query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []Country{}
	for rows.Next() {
		var c Country
		if err := scanCountry(rows, &c); err != nil {
			log.Printf("[ERROR] listCountries scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, c)
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func getCountry(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/countries/")
	countryID, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid country ID"})
		return
	}
	var c Country
	err = scanCountry(db.QueryRow("SELECT id, iso, name, nicename, iso3, numcode, phonecode FROM country WHERE id = ?", countryID), &c)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Country not found"})
		} else {
			log.Printf("[ERROR] getCountry query failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		}
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: c})
}

func createCountry(w http.ResponseWriter, r *http.Request) {
	var c Country
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}
	result, err := db.Exec("INSERT INTO country (iso, name, nicename, iso3, numcode, phonecode) VALUES (?, ?, ?, ?, ?, ?)", c.ISO, c.Name, c.NiceName, c.ISO3, c.NumCode, c.PhoneCode)
	if err != nil {
		log.Printf("[ERROR] createCountry insert failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create country"})
		return
	}
	id, _ := result.LastInsertId()
	c.ID = int(id)
	writeJSON(w, http.StatusCreated, APIResponse{Success: true, Data: c})
}

func updateCountry(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/countries/")
	countryID, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid country ID"})
		return
	}
	var c Country
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}
	result, err := db.Exec("UPDATE country SET iso = ?, name = ?, nicename = ?, iso3 = ?, numcode = ?, phonecode = ? WHERE id = ?", c.ISO, c.Name, c.NiceName, c.ISO3, c.NumCode, c.PhoneCode, countryID)
	if err != nil {
		log.Printf("[ERROR] updateCountry query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update country"})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Country not found"})
		return
	}
	c.ID = countryID
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: c})
}

func deleteCountry(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/countries/")
	countryID, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid country ID"})
		return
	}
	result, err := db.Exec("DELETE FROM country WHERE id = ?", countryID)
	if err != nil {
		log.Printf("[ERROR] deleteCountry query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete country"})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Country not found"})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Country deleted"}})
}

func listRegions(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, nicename, note, country FROM region LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listRegions query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []Region{}
	for rows.Next() {
		var rg Region
		if err := rows.Scan(&rg.ID, &rg.Name, &rg.NiceName, &rg.Note, &rg.Country); err != nil {
			log.Printf("[ERROR] listRegions scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, rg)
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func getRegion(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/regions/")
	regionID, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid region ID"})
		return
	}
	var rg Region
	err = db.QueryRow("SELECT id, name, nicename, note, country FROM region WHERE id = ?", regionID).Scan(&rg.ID, &rg.Name, &rg.NiceName, &rg.Note, &rg.Country)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Region not found"})
		} else {
			log.Printf("[ERROR] getRegion query failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		}
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: rg})
}

func createRegion(w http.ResponseWriter, r *http.Request) {
	var rg Region
	if err := json.NewDecoder(r.Body).Decode(&rg); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}
	result, err := db.Exec("INSERT INTO region (name, nicename, note, country) VALUES (?, ?, ?, ?)", rg.Name, rg.NiceName, rg.Note, rg.Country)
	if err != nil {
		log.Printf("[ERROR] createRegion insert failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create region"})
		return
	}
	id, _ := result.LastInsertId()
	rg.ID = int(id)
	writeJSON(w, http.StatusCreated, APIResponse{Success: true, Data: rg})
}

func updateRegion(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/regions/")
	regionID, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid region ID"})
		return
	}
	var rg Region
	if err := json.NewDecoder(r.Body).Decode(&rg); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}
	result, err := db.Exec("UPDATE region SET name = ?, nicename = ?, note = ?, country = ? WHERE id = ?", rg.Name, rg.NiceName, rg.Note, rg.Country, regionID)
	if err != nil {
		log.Printf("[ERROR] updateRegion query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update region"})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Region not found"})
		return
	}
	rg.ID = regionID
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: rg})
}

func deleteRegion(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/regions/")
	regionID, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid region ID"})
		return
	}
	result, err := db.Exec("DELETE FROM region WHERE id = ?", regionID)
	if err != nil {
		log.Printf("[ERROR] deleteRegion query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete region"})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Region not found"})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Region deleted"}})
}

func listRegionEmp(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT emp_no, region_id, from_date, to_date FROM region_emp LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listRegionEmp query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []RegionEmp{}
	for rows.Next() {
		var re RegionEmp
		if err := rows.Scan(&re.EmpNo, &re.RegionID, &re.FromDate, &re.ToDate); err != nil {
			log.Printf("[ERROR] listRegionEmp scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, re)
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func getRegionEmp(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/region_emp/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}
	rows, err := db.Query("SELECT emp_no, region_id, from_date, to_date FROM region_emp WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] getRegionEmp query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []RegionEmp{}
	for rows.Next() {
		var re RegionEmp
		if err := rows.Scan(&re.EmpNo, &re.RegionID, &re.FromDate, &re.ToDate); err != nil {
			log.Printf("[ERROR] getRegionEmp scan failed: %v\n", err)
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, re)
	}
	if len(items) == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "No region assignments found"})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func createRegionEmp(w http.ResponseWriter, r *http.Request) {
	var re RegionEmp
	if err := json.NewDecoder(r.Body).Decode(&re); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}
	_, err := db.Exec("INSERT INTO region_emp (emp_no, region_id, from_date, to_date) VALUES (?, ?, ?, ?)", re.EmpNo, re.RegionID, re.FromDate, re.ToDate)
	if err != nil {
		log.Printf("[ERROR] createRegionEmp insert failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create region assignment"})
		return
	}
	writeJSON(w, http.StatusCreated, APIResponse{Success: true, Data: re})
}

func updateRegionEmp(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/region_emp/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}
	var re RegionEmp
	if err := json.NewDecoder(r.Body).Decode(&re); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}
	result, err := db.Exec("UPDATE region_emp SET region_id = ?, from_date = ?, to_date = ? WHERE emp_no = ? AND from_date = ?", re.RegionID, re.FromDate, re.ToDate, empNo, re.FromDate)
	if err != nil {
		log.Printf("[ERROR] updateRegionEmp query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update region assignment"})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Region assignment not found"})
		return
	}
	re.EmpNo = empNo
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: re})
}

func deleteRegionEmp(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/region_emp/")
	empNo, err := strconv.Atoi(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}
	result, err := db.Exec("DELETE FROM region_emp WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteRegionEmp query failed: %v\n", err)
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete region assignment"})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Region assignment not found"})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Region assignment deleted"}})
}
