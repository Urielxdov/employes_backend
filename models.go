package main

import "time"

type Employee struct {
	EmpNo     int       `json:"emp_no"`
	BirthDate time.Time `json:"birth_date"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Gender    string    `json:"gender"`
	HireDate  time.Time `json:"hire_date"`
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

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
