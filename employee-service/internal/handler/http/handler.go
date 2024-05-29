package http

import (
	"net/http"

	"github.com/Verce11o/resume-view/employee-service/api"
	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	log             *zap.SugaredLogger
	positionService service.Position
	employeeService service.Employee
	authService     service.Auth
}

func NewHandler(log *zap.SugaredLogger, positionService service.Position, employeeService service.Employee,
	authService service.Auth) *Handler {
	return &Handler{log: log, positionService: positionService, employeeService: employeeService,
		authService: authService}
}

func (h *Handler) SignIn(c *gin.Context) {
	var input api.SignInEmployee

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	employeeID, err := uuid.Parse(input.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	token, err := h.authService.SignIn(c.Request.Context(), employeeID)

	if err != nil {
		h.log.Errorf("error while sign in: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": token,
	})
}

func (h *Handler) CreateEmployee(c *gin.Context) {
	var input api.CreateEmployee

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	employee, err := h.employeeService.CreateEmployee(c.Request.Context(), domain.CreateEmployee{
		EmployeeID:   uuid.New(),
		PositionID:   uuid.New(),
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		PositionName: input.PositionName,
		Salary:       input.Salary,
	})

	if err != nil {
		h.log.Errorf("error creating employee: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *Handler) GetEmployeeByID(c *gin.Context, id string) {
	employeeID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	employee, err := h.employeeService.GetEmployee(c.Request.Context(), employeeID)
	if err != nil {
		h.log.Errorf("error getting employee: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *Handler) GetEmployeeList(c *gin.Context, params api.GetEmployeeListParams) {
	var cursor string
	if params.Cursor != nil {
		cursor = *params.Cursor
	}

	employee, err := h.employeeService.GetEmployeeList(c.Request.Context(), cursor)
	if err != nil {
		h.log.Errorf("error getting employee: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *Handler) UpdateEmployeeByID(c *gin.Context, id string) {
	employeeID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	var input api.UpdateEmployee

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	positionID, err := uuid.Parse(input.PositionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	employee, err := h.employeeService.UpdateEmployee(c.Request.Context(), domain.UpdateEmployee{
		EmployeeID: employeeID,
		PositionID: positionID,
		FirstName:  input.FirstName,
		LastName:   input.LastName,
	})

	if err != nil {
		h.log.Errorf("error updating employee: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *Handler) DeleteEmployeeByID(c *gin.Context, id string) {
	employeeID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	err = h.employeeService.DeleteEmployee(c.Request.Context(), employeeID)
	if err != nil {
		h.log.Errorf("error deleting employee: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) CreatePosition(c *gin.Context) {
	var input api.CreatePosition

	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Errorf("error parsing position: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	position, err := h.positionService.CreatePosition(c.Request.Context(), domain.CreatePosition{
		ID:     uuid.New(),
		Name:   input.Name,
		Salary: input.Salary,
	})
	if err != nil {
		h.log.Errorf("error creating position: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

		return
	}

	c.JSON(http.StatusOK, position)
}

func (h *Handler) GetPositionByID(c *gin.Context, id string) {
	positionID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	position, err := h.positionService.GetPosition(c.Request.Context(), positionID)
	if err != nil {
		h.log.Errorf("error getting position: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

		return
	}

	c.JSON(http.StatusOK, position)
}

func (h *Handler) GetPositionList(c *gin.Context, params api.GetPositionListParams) {
	var cursor string
	if params.Cursor != nil {
		cursor = *params.Cursor
	}

	employee, err := h.positionService.GetPositionList(c.Request.Context(), cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *Handler) UpdatePositionByID(c *gin.Context, id string) {
	positionID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	var input api.UpdatePosition

	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Errorf("error parsing position: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	position, err := h.positionService.UpdatePosition(c.Request.Context(), domain.UpdatePosition{
		ID:     positionID,
		Name:   input.Name,
		Salary: input.Salary,
	})

	if err != nil {
		h.log.Errorf("error updating position: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

		return
	}

	c.JSON(http.StatusOK, position)
}

func (h *Handler) DeletePositionByID(c *gin.Context, id string) {
	positionID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	err = h.positionService.DeletePosition(c.Request.Context(), positionID)
	if err != nil {
		h.log.Errorf("error deleting position: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
