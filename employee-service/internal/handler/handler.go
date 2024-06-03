package handler

import (
	"net/http"
)

type EmployeeHandler interface {
	SignIn(w http.ResponseWriter, r *http.Request)
	CreateEmployee(w http.ResponseWriter, r *http.Request)
	GetEmployeeByID(w http.ResponseWriter, r *http.Request)
	GetEmployeeList(w http.ResponseWriter, r *http.Request)
	UpdateEmployeeByID(w http.ResponseWriter, r *http.Request)
	DeleteEmployeeByID(w http.ResponseWriter, r *http.Request)
}

type PositionHandler interface {
	CreatePosition(w http.ResponseWriter, r *http.Request)
	GetPositionByID(w http.ResponseWriter, r *http.Request)
	GetPositionList(w http.ResponseWriter, r *http.Request)
	UpdatePositionByID(w http.ResponseWriter, r *http.Request)
	DeletePositionByID(w http.ResponseWriter, r *http.Request)
}
