package gorilla

import (
	"encoding/json"
	"net/http"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type m map[string]string

type Handler struct {
	log             *zap.SugaredLogger
	positionService service.Position
	employeeService service.Employee
	authService     service.Auth
}

func New(log *zap.SugaredLogger, positionService service.Position, employeeService service.Employee,
	authService service.Auth) *Handler {
	return &Handler{log: log, positionService: positionService, employeeService: employeeService,
		authService: authService}
}
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var input domain.SignInEmployeeRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleErr(w, err.Error(), http.StatusBadRequest)

		return
	}

	employeeID, err := uuid.Parse(input.ID)
	if err != nil {
		handleErr(w, err.Error(), http.StatusBadRequest)

		return
	}

	token, err := h.authService.SignIn(r.Context(), employeeID)
	if err != nil {
		h.log.Errorf("error while sign in: %v", err)
		handleErr(w, err.Error(), http.StatusUnauthorized)

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(m{"message": token})

	if err != nil {
		h.log.Errorf("error while encode response: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateEmployeeRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleErr(w, err.Error(), http.StatusBadRequest)

		return
	}

	employee, err := h.employeeService.CreateEmployee(r.Context(), domain.CreateEmployee{
		EmployeeID:   uuid.New(),
		PositionID:   uuid.New(),
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		PositionName: input.PositionName,
		Salary:       input.Salary,
	})

	if err != nil {
		h.log.Errorf("error creating employee: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(employee)

	if err != nil {
		h.log.Errorf("error while encode response: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *Handler) GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	employeeID, err := uuid.Parse(id)

	if err != nil {
		handleErr(w, err.Error(), http.StatusBadRequest)

		return
	}

	employee, err := h.employeeService.GetEmployee(r.Context(), employeeID)
	if err != nil {
		h.log.Errorf("error getting employee: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(employee)

	if err != nil {
		h.log.Errorf("error while encode response: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *Handler) GetEmployeeList(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")

	employee, err := h.employeeService.GetEmployeeList(r.Context(), cursor)
	if err != nil {
		h.log.Errorf("error getting employee: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(employee)

	if err != nil {
		h.log.Errorf("error while encode response: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *Handler) UpdateEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	employeeID, err := uuid.Parse(id)
	if err != nil {
		handleErr(w, err.Error(), http.StatusBadRequest)

		return
	}

	var input domain.UpdateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleErr(w, err.Error(), http.StatusBadRequest)

		return
	}

	positionID, err := uuid.Parse(input.PositionID)
	if err != nil {
		handleErr(w, err.Error(), http.StatusBadRequest)

		return
	}

	employee, err := h.employeeService.UpdateEmployee(r.Context(), domain.UpdateEmployee{
		EmployeeID: employeeID,
		PositionID: positionID,
		FirstName:  input.FirstName,
		LastName:   input.LastName,
	})

	if err != nil {
		h.log.Errorf("error updating employee: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(employee)

	if err != nil {
		h.log.Errorf("error while encode response: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *Handler) DeleteEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	employeeID, err := uuid.Parse(id)
	if err != nil {
		handleErr(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = h.employeeService.DeleteEmployee(r.Context(), employeeID)
	if err != nil {
		h.log.Errorf("error deleting employee: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(m{"message": "success"})

	if err != nil {
		h.log.Errorf("error while encode response: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *Handler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	var input domain.CreatePosition

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleErr(w, err.Error(), http.StatusBadRequest)

		return
	}

	position, err := h.positionService.CreatePosition(r.Context(), domain.CreatePosition{
		ID:     uuid.New(),
		Name:   input.Name,
		Salary: input.Salary,
	})

	if err != nil {
		h.log.Errorf("error creating position: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(position)

	if err != nil {
		h.log.Errorf("error while encode response: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *Handler) GetPositionByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	positionID, err := uuid.Parse(id)
	if err != nil {
		handleErr(w, err.Error(), http.StatusBadRequest)

		return
	}

	position, err := h.positionService.GetPosition(r.Context(), positionID)
	if err != nil {
		h.log.Errorf("error getting position: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(position)

	if err != nil {
		h.log.Errorf("error while encode response: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *Handler) GetPositionList(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")

	position, err := h.positionService.GetPositionList(r.Context(), cursor)
	if err != nil {
		h.log.Errorf("error getting position: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(position)

	if err != nil {
		h.log.Errorf("error while encode response: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *Handler) UpdatePositionByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	positionID, err := uuid.Parse(id)
	if err != nil {
		handleErr(w, err.Error(), http.StatusBadRequest)

		return
	}

	var input domain.UpdatePosition
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleErr(w, err.Error(), http.StatusBadRequest)

		return
	}

	position, err := h.positionService.UpdatePosition(r.Context(), domain.UpdatePosition{
		ID:     positionID,
		Name:   input.Name,
		Salary: input.Salary,
	})

	if err != nil {
		h.log.Errorf("error updating position: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(position)

	if err != nil {
		h.log.Errorf("error while encode response: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *Handler) DeletePositionByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	positionID, err := uuid.Parse(id)
	if err != nil {
		handleErr(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = h.positionService.DeletePosition(r.Context(), positionID)
	if err != nil {
		h.log.Errorf("error deleting position: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(m{"message": "success"})

	if err != nil {
		h.log.Errorf("error while encode response: %v", err)
		handleErr(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func handleErr(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(m{"message": msg})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
