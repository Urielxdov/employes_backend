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
