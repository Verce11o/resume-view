package handler

import (
	"context"
	"github.com/Verce11o/resume-view/employee-service/api"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type PositionService interface {
	CreatePosition(ctx context.Context, request api.CreatePosition) (models.Position, error)
	GetPosition(ctx context.Context, id string) (models.Position, error)
	UpdatePosition(ctx context.Context, id string, request api.UpdatePosition) (models.Position, error)
	DeletePosition(ctx context.Context, id string) error
}

type EmployeeService interface {
	CreateEmployee(ctx context.Context, request api.CreateEmployee) (models.Employee, error)
	GetEmployee(ctx context.Context, id string) (models.Employee, error)
	UpdateEmployee(ctx context.Context, id string, request api.UpdateEmployee) (models.Employee, error)
	DeleteEmployee(ctx context.Context, id string) error
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

	employee, err := h.employeeService.CreateEmployee(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *Handler) GetEmployeeByID(c *gin.Context, id string) {
	employee, err := h.employeeService.GetEmployee(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employee)
}

func (h *Handler) UpdateEmployeeByID(c *gin.Context, id string) {
	var input api.UpdateEmployee

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	employee, err := h.employeeService.UpdateEmployee(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *Handler) DeleteEmployeeByID(c *gin.Context, id string) {
	err := h.employeeService.DeleteEmployee(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) CreatePosition(c *gin.Context) {
	var input api.CreatePosition

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	position, err := h.positionService.CreatePosition(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, position)
}

func (h *Handler) GetPositionByID(c *gin.Context, id string) {
	position, err := h.positionService.GetPosition(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, position)
}

func (h *Handler) UpdatePositionByID(c *gin.Context, id string) {
	var input api.UpdatePosition

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	position, err := h.positionService.UpdatePosition(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, position)
}

func (h *Handler) DeletePositionByID(c *gin.Context, id string) {
	err := h.positionService.DeletePosition(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
