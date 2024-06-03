//go:build !integration

package gin

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Verce11o/resume-view/employee-service/internal/models"
	serviceMock "github.com/Verce11o/resume-view/employee-service/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestHandler_CreateEmployee(t *testing.T) {
	t.Parallel()

	type fields struct {
		employeeService *serviceMock.MockEmployeeService
		positionService *serviceMock.MockPositionService
	}

	tests := []struct {
		name       string
		input      string
		response   any
		mockFunc   func(f *fields)
		statusCode int
		err        error
	}{
		{
			name:  "Valid Input",
			input: `{"first_name":"John","last_name":"Doe","position_name":"Developer","salary":60000}`,
			response: models.Employee{ // I'm not checking for IDs, created_at & updated_at since I can't predict them.
				FirstName: "John",
				LastName:  "Doe",
			},
			mockFunc: func(f *fields) {
				f.employeeService.EXPECT().CreateEmployee(gomock.Any(), gomock.Any()).
					Return(models.Employee{
						FirstName: "John",
						LastName:  "Doe",
					}, nil)
			},
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid Input",
			input:      `{"name":"John","surname":"Doe"}`,
			response:   models.Employee{},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		ctx, w := createTestContext(http.MethodPost, tt.input)

		MockJSONRequest(ctx, tt.input, http.MethodPost)

		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			h.CreateEmployee(ctx)
			assert.EqualValues(t, tt.statusCode, w.Code)

			if tt.response != nil {
				var responseBody models.Employee
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.EqualValues(t, tt.response, responseBody)
			}
		})
	}
}

func TestHandler_GetEmployeeByID(t *testing.T) {
	t.Parallel()

	type fields struct {
		employeeService *serviceMock.MockEmployeeService
		positionService *serviceMock.MockPositionService
	}

	tests := []struct {
		name       string
		id         string
		response   any
		mockFunc   func(f *fields)
		statusCode int
		err        error
	}{
		{
			name: "Valid ID",
			id:   uuid.New().String(),
			response: models.Employee{
				FirstName: "John",
				LastName:  "Doe",
			},
			mockFunc: func(f *fields) {
				f.employeeService.EXPECT().GetEmployee(gomock.Any(), gomock.Any()).
					Return(models.Employee{
						FirstName: "John",
						LastName:  "Doe",
					}, nil)
			},
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid ID",
			id:         "invalid",
			response:   models.Employee{},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		ctx, w := createTestContext(http.MethodPost, "")

		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			ctx.AddParam("id", tt.id)

			h.GetEmployeeByID(ctx)
			assert.EqualValues(t, tt.statusCode, w.Code)

			if tt.response != nil {
				var responseBody models.Employee
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.EqualValues(t, tt.response, responseBody)
			}
		})
	}
}

func TestHandler_GetEmployeeList(t *testing.T) {
	t.Parallel()

	type fields struct {
		employeeService *serviceMock.MockEmployeeService
		positionService *serviceMock.MockPositionService
	}

	cursor := ""
	employees := []models.Employee{
		{
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			FirstName: "Steven",
			LastName:  "Hocking",
		},
	}

	tests := []struct {
		name       string
		cursor     *string
		response   any
		mockFunc   func(f *fields)
		statusCode int
		err        error
	}{
		{
			name:   "Valid empty cursor",
			cursor: &cursor,
			response: models.EmployeeList{
				Cursor:    "example",
				Employees: employees,
			},
			mockFunc: func(f *fields) {
				f.employeeService.EXPECT().GetEmployeeList(gomock.Any(), gomock.Any()).
					Return(models.EmployeeList{
						Cursor:    "example",
						Employees: employees,
					}, nil)
			},
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		ctx, w := createTestContext(http.MethodPost, "")

		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			ctx.Request.URL.Query().Set("cursor", cursor)

			h.GetEmployeeList(ctx)
			assert.EqualValues(t, tt.statusCode, w.Code)

			if tt.response != nil {
				var responseBody models.EmployeeList
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.EqualValues(t, tt.response, responseBody)
			}
		})
	}
}

func TestHandler_UpdateEmployeeByID(t *testing.T) {
	t.Parallel()

	type fields struct {
		employeeService *serviceMock.MockEmployeeService
		positionService *serviceMock.MockPositionService
	}

	tests := []struct {
		name       string
		id         string
		input      string
		response   any
		mockFunc   func(f *fields)
		statusCode int
		err        error
	}{
		{
			name:  "Valid Input",
			id:    uuid.NewString(),
			input: `{"first_name":"NewJohn","last_name":"NewDoe", "position_id":"` + uuid.NewString() + `"}`,
			response: models.Employee{
				FirstName: "NewJohn",
				LastName:  "NewDoe",
			},
			mockFunc: func(f *fields) {
				f.employeeService.EXPECT().UpdateEmployee(gomock.Any(), gomock.Any()).
					Return(models.Employee{
						FirstName: "NewJohn",
						LastName:  "NewDoe",
					}, nil)
			},
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid employee ID",
			id:         "invalid",
			input:      `{"first_name":"NewJohn","last_name":"NewDoe", "position_id":"` + uuid.NewString() + `"}`,
			response:   models.Employee{},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Invalid position ID",
			id:         uuid.NewString(),
			input:      `{"first_name":"NewJohn","last_name":"NewDoe", "position_id":"invalid"`,
			response:   models.Employee{},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Invalid input",
			id:         uuid.NewString(),
			input:      `{"name":"NewJohn","surname":"NewDoe"}`,
			response:   models.Employee{},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		ctx, w := createTestContext(http.MethodPut, tt.input)

		MockJSONRequest(ctx, tt.input, http.MethodPut)

		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			ctx.AddParam("id", tt.id)

			h.UpdateEmployeeByID(ctx)
			assert.EqualValues(t, tt.statusCode, w.Code)

			if tt.response != nil {
				var responseBody models.Employee
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.EqualValues(t, tt.response, responseBody)
			}
		})
	}
}

func TestHandler_DeleteEmployeeByID(t *testing.T) {
	t.Parallel()

	type fields struct {
		employeeService *serviceMock.MockEmployeeService
		positionService *serviceMock.MockPositionService
	}

	tests := []struct {
		name       string
		id         string
		response   any
		mockFunc   func(f *fields)
		statusCode int
		err        error
	}{
		{
			name: "Valid ID",
			id:   uuid.NewString(),
			response: gin.H{
				"message": "success",
			},
			mockFunc: func(f *fields) {
				f.employeeService.EXPECT().DeleteEmployee(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Invalid ID",
			id:   "invalid",
			response: gin.H{
				"message": errors.New("invalid UUID length: 7").Error(),
			},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		ctx, w := createTestContext(http.MethodPost, "")

		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			ctx.AddParam("id", tt.id)

			h.DeleteEmployeeByID(ctx)
			assert.EqualValues(t, tt.statusCode, w.Code)

			if tt.response != nil {
				var responseBody gin.H
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.EqualValues(t, tt.response, responseBody)
			}
		})
	}
}

func TestHandler_CreatePosition(t *testing.T) {
	t.Parallel()

	type fields struct {
		employeeService *serviceMock.MockEmployeeService
		positionService *serviceMock.MockPositionService
	}

	tests := []struct {
		name       string
		input      string
		response   any
		mockFunc   func(f *fields)
		statusCode int
		err        error
	}{
		{
			name:  "Valid Input",
			input: `{"name": "Go Developer", "salary": 30999}`,
			response: models.Position{
				Name:   "Go Developer",
				Salary: 30999,
			},
			mockFunc: func(f *fields) {
				f.positionService.EXPECT().CreatePosition(gomock.Any(), gomock.Any()).
					Return(models.Position{
						Name:   "Go Developer",
						Salary: 30999,
					}, nil)
			},
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid Input",
			input:      `{"name": "Go Developer", "pay_amount": 30999}`,
			response:   models.Position{},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		ctx, w := createTestContext(http.MethodPost, tt.input)

		MockJSONRequest(ctx, tt.input, http.MethodPost)

		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			h.CreatePosition(ctx)
			assert.EqualValues(t, tt.statusCode, w.Code)

			if tt.response != nil {
				var responseBody models.Position
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.EqualValues(t, tt.response, responseBody)
			}
		})
	}
}

func TestHandler_GetPositionByID(t *testing.T) {
	t.Parallel()

	type fields struct {
		employeeService *serviceMock.MockEmployeeService
		positionService *serviceMock.MockPositionService
	}

	tests := []struct {
		name       string
		id         string
		response   any
		mockFunc   func(f *fields)
		statusCode int
		err        error
	}{
		{
			name: "Valid ID",
			id:   uuid.New().String(),
			response: models.Position{
				Name:   "Go Developer",
				Salary: 30999,
			},
			mockFunc: func(f *fields) {
				f.positionService.EXPECT().GetPosition(gomock.Any(), gomock.Any()).
					Return(models.Position{
						Name:   "Go Developer",
						Salary: 30999,
					}, nil)
			},
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid ID",
			id:         "invalid",
			response:   models.Position{},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		ctx, w := createTestContext(http.MethodGet, "")

		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			ctx.AddParam("id", tt.id)

			h.GetPositionByID(ctx)
			assert.EqualValues(t, tt.statusCode, w.Code)

			if tt.response != nil {
				var responseBody models.Position
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.EqualValues(t, tt.response, responseBody)
			}
		})
	}
}

func TestHandler_GetPositionList(t *testing.T) {
	t.Parallel()

	type fields struct {
		employeeService *serviceMock.MockEmployeeService
		positionService *serviceMock.MockPositionService
	}

	cursor := ""
	positions := []models.Position{
		{
			Name:   "Go Developer",
			Salary: 30999,
		},
		{
			Name:   "Python developer",
			Salary: 20845,
		},
	}

	tests := []struct {
		name       string
		cursor     *string
		response   any
		mockFunc   func(f *fields)
		statusCode int
		err        error
	}{
		{
			name:   "Valid empty cursor",
			cursor: &cursor,
			response: models.PositionList{
				Cursor:    "example",
				Positions: positions,
			},
			mockFunc: func(f *fields) {
				f.positionService.EXPECT().GetPositionList(gomock.Any(), gomock.Any()).
					Return(models.PositionList{
						Cursor:    "example",
						Positions: positions,
					}, nil)
			},
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		ctx, w := createTestContext(http.MethodGet, "")

		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			ctx.Request.URL.Query().Set("cursor", cursor)

			h.GetPositionList(ctx)
			assert.EqualValues(t, tt.statusCode, w.Code)

			if tt.response != nil {
				var responseBody models.PositionList
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.EqualValues(t, tt.response, responseBody)
			}
		})
	}
}

func TestHandler_UpdatePosition(t *testing.T) {
	t.Parallel()

	type fields struct {
		employeeService *serviceMock.MockEmployeeService
		positionService *serviceMock.MockPositionService
	}

	tests := []struct {
		name       string
		id         string
		input      string
		response   any
		mockFunc   func(f *fields)
		statusCode int
		err        error
	}{
		{
			name:  "Valid Input",
			id:    uuid.NewString(),
			input: `{"name":"NewGoDeveloper", "salary": 30999}`,
			response: models.Position{
				Name:   "NewGoDeveloper",
				Salary: 30999,
			},
			mockFunc: func(f *fields) {
				f.positionService.EXPECT().UpdatePosition(gomock.Any(), gomock.Any()).
					Return(models.Position{
						Name:   "NewGoDeveloper",
						Salary: 30999,
					}, nil)
			},
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid ID",
			id:         "invalid",
			input:      `{"name":"NewGoDeveloper", "salary": 30999}`,
			response:   models.Position{},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Invalid input",
			id:         uuid.NewString(),
			input:      `{"positionName":"NewGoDeveloper", "pay_amount": 30999}`,
			response:   models.Position{},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		ctx, w := createTestContext(http.MethodPut, tt.input)

		MockJSONRequest(ctx, tt.input, http.MethodPut)

		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			ctx.AddParam("id", tt.id)

			h.UpdatePositionByID(ctx)
			assert.EqualValues(t, tt.statusCode, w.Code)

			if tt.response != nil {
				var responseBody models.Position
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.EqualValues(t, tt.response, responseBody)
			}
		})
	}
}

func TestHandler_DeletePositionByID(t *testing.T) {
	t.Parallel()

	type fields struct {
		employeeService *serviceMock.MockEmployeeService
		positionService *serviceMock.MockPositionService
	}

	tests := []struct {
		name       string
		id         string
		response   any
		mockFunc   func(f *fields)
		statusCode int
		err        error
	}{
		{
			name: "Valid ID",
			id:   uuid.NewString(),
			response: gin.H{
				"message": "success",
			},
			mockFunc: func(f *fields) {
				f.positionService.EXPECT().DeletePosition(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Invalid ID",
			id:   "invalid",
			response: gin.H{
				"message": errors.New("invalid UUID length: 7").Error(),
			},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		ctx, w := createTestContext(http.MethodPost, "")

		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})
			ctx.AddParam("id", tt.id)

			h.DeletePositionByID(ctx)
			assert.EqualValues(t, tt.statusCode, w.Code)

			if tt.response != nil {
				var responseBody gin.H
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.EqualValues(t, tt.response, responseBody)
			}
		})
	}
}

func MockJSONRequest(c *gin.Context, body string, method string) {
	c.Request.Method = method
	c.Request.Header.Set("Content-Type", "application/json")

	c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(body)))
}

func createTestContext(method, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Method: method,
		Header: make(http.Header),
	}
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(body)))

	return ctx, w
}

func initMocks(t *testing.T) (*gomock.Controller, *serviceMock.MockEmployeeService,
	*serviceMock.MockPositionService, *Handler) {
	ctrl := gomock.NewController(t)
	positionService := serviceMock.NewMockPositionService(ctrl)
	employeeService := serviceMock.NewMockEmployeeService(ctrl)
	log := zap.NewNop().Sugar()

	h := &Handler{
		log:             log,
		positionService: positionService,
		employeeService: employeeService,
	}

	return ctrl, employeeService, positionService, h
}
