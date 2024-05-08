package handler

import (
	"context"
	"net/http"

	"github.com/Verce11o/resume-view/employee-service/api"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PositionService interface {
	CreatePosition(ctx context.Context, request api.CreatePosition) (models.Position, error)
	GetPosition(ctx context.Context, id uuid.UUID) (models.Position, error)
	GetPositionList(ctx context.Context, cursor string) (models.PositionList, error)
	UpdatePosition(ctx context.Context, id uuid.UUID, request api.UpdatePosition) (models.Position, error)
	DeletePosition(ctx context.Context, id uuid.UUID) error
}

type EmployeeService interface {
	CreateEmployee(ctx context.Context, employeeID uuid.UUID, positionID uuid.UUID, request api.CreateEmployee) (models.Employee, error)
	GetEmployee(ctx context.Context, id uuid.UUID) (models.Employee, error)
	GetEmployeeList(ctx context.Context, cursor string) (models.EmployeeList, error)
	UpdateEmployee(ctx context.Context, id uuid.UUID, request api.UpdateEmployee) (models.Employee, error)
	DeleteEmployee(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	log             *zap.SugaredLogger
	positionService PositionService
	employeeService EmployeeService
}

func NewHandler(log *zap.SugaredLogger, positionService PositionService, employeeService EmployeeService) *Handler {
	return &Handler{log: log, positionService: positionService, employeeService: employeeService}
}

func (h *Handler) CreateEmployee(c *gin.Context) {
	var input api.CreateEmployee

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	employee, err := h.employeeService.CreateEmployee(c.Request.Context(), uuid.New(), uuid.New(), input)
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

	employee, err := h.employeeService.UpdateEmployee(c.Request.Context(), employeeID, input)
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

	position, err := h.positionService.CreatePosition(c.Request.Context(), input)
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

	position, err := h.positionService.UpdatePosition(c.Request.Context(), positionID, input)
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
