package chi

import (
	"encoding/json"
	"net/http"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/service"
	"github.com/go-chi/chi"
	chiRender "github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

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
		chiRender.Status(r, http.StatusBadRequest)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	employeeID, err := uuid.Parse(input.ID)
	if err != nil {
		chiRender.Status(r, http.StatusBadRequest)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	token, err := h.authService.SignIn(r.Context(), employeeID)

	if err != nil {
		h.log.Errorf("error while sign in: %v", err)
		chiRender.Status(r, http.StatusUnauthorized)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	chiRender.Status(r, http.StatusOK)
	chiRender.JSON(w, r, chiRender.M{
		"message": token,
	})
}

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateEmployeeRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		chiRender.Status(r, http.StatusBadRequest)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

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
		chiRender.Status(r, http.StatusInternalServerError)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	chiRender.Status(r, http.StatusOK)
	chiRender.JSON(w, r, employee)
}

func (h *Handler) GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	employeeID, err := uuid.Parse(id)

	if err != nil {
		chiRender.Status(r, http.StatusBadRequest)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	employee, err := h.employeeService.GetEmployee(r.Context(), employeeID)
	if err != nil {
		h.log.Errorf("error getting employee: %v", err)
		chiRender.Status(r, http.StatusInternalServerError)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	chiRender.Status(r, http.StatusOK)
	chiRender.JSON(w, r, employee)
}

func (h *Handler) GetEmployeeList(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")

	employee, err := h.employeeService.GetEmployeeList(r.Context(), cursor)
	if err != nil {
		h.log.Errorf("error getting employee: %v", err)
		chiRender.Status(r, http.StatusInternalServerError)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	chiRender.Status(r, http.StatusOK)
	chiRender.JSON(w, r, employee)
}

func (h *Handler) UpdateEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	employeeID, err := uuid.Parse(id)
	if err != nil {
		chiRender.Status(r, http.StatusBadRequest)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	var input domain.UpdateEmployeeRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		chiRender.Status(r, http.StatusBadRequest)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	positionID, err := uuid.Parse(input.PositionID)
	if err != nil {
		chiRender.Status(r, http.StatusBadRequest)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

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
		chiRender.Status(r, http.StatusInternalServerError)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	chiRender.Status(r, http.StatusOK)
	chiRender.JSON(w, r, employee)
}

func (h *Handler) DeleteEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	employeeID, err := uuid.Parse(id)
	if err != nil {
		chiRender.Status(r, http.StatusBadRequest)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	err = h.employeeService.DeleteEmployee(r.Context(), employeeID)
	if err != nil {
		h.log.Errorf("error deleting employee: %v", err)
		chiRender.Status(r, http.StatusInternalServerError)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	chiRender.Status(r, http.StatusOK)
	chiRender.JSON(w, r, chiRender.M{
		"message": "success",
	})
}

func (h *Handler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	var input domain.CreatePositionRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		chiRender.Status(r, http.StatusBadRequest)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	position, err := h.positionService.CreatePosition(r.Context(), domain.CreatePosition{
		ID:     uuid.New(),
		Name:   input.Name,
		Salary: input.Salary,
	})

	if err != nil {
		h.log.Errorf("error creating position: %v", err)
		chiRender.Status(r, http.StatusInternalServerError)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	chiRender.Status(r, http.StatusOK)
	chiRender.JSON(w, r, position)
}
func (h *Handler) GetPositionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	positionID, err := uuid.Parse(id)
	if err != nil {
		chiRender.Status(r, http.StatusBadRequest)
		chiRender.JSON(w, r, chiRender.M{
			"message": "invalid ID",
		})

		return
	}

	position, err := h.positionService.GetPosition(r.Context(), positionID)
	if err != nil {
		h.log.Errorf("error getting position: %v", err)
		chiRender.Status(r, http.StatusInternalServerError)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	chiRender.Status(r, http.StatusOK)
	chiRender.JSON(w, r, position)
}
func (h *Handler) GetPositionList(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")

	position, err := h.positionService.GetPositionList(r.Context(), cursor)
	if err != nil {
		chiRender.Status(r, http.StatusInternalServerError)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	chiRender.Status(r, http.StatusOK)
	chiRender.JSON(w, r, position)
}

func (h *Handler) UpdatePositionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	positionID, err := uuid.Parse(id)
	if err != nil {
		chiRender.Status(r, http.StatusBadRequest)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	var input domain.UpdatePositionRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		chiRender.Status(r, http.StatusBadRequest)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	position, err := h.positionService.UpdatePosition(r.Context(), domain.UpdatePosition{
		ID:     positionID,
		Name:   input.Name,
		Salary: input.Salary,
	})

	if err != nil {
		h.log.Errorf("error updating position: %v", err)
		chiRender.Status(r, http.StatusInternalServerError)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	chiRender.Status(r, http.StatusOK)
	chiRender.JSON(w, r, position)
}

func (h *Handler) DeletePositionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	positionID, err := uuid.Parse(id)
	if err != nil {
		chiRender.Status(r, http.StatusBadRequest)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	err = h.positionService.DeletePosition(r.Context(), positionID)
	if err != nil {
		h.log.Errorf("error deleting position: %v", err)
		chiRender.Status(r, http.StatusInternalServerError)
		chiRender.JSON(w, r, chiRender.M{
			"message": err.Error(),
		})

		return
	}

	chiRender.Status(r, http.StatusOK)
	chiRender.JSON(w, r, chiRender.M{
		"message": "success",
	})
}
