package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Verce11o/resume-view/employee-service/api"
	serviceMock "github.com/Verce11o/resume-view/employee-service/internal/handler/mocks"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
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

	log := zap.NewNop().Sugar()

	for _, tt := range tests {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = &http.Request{
			Method: http.MethodPost,
			Header: make(http.Header),
		}

		MockJSONPost(ctx, tt.input)

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			positionService := serviceMock.NewMockPositionService(ctrl)
			employeeService := serviceMock.NewMockEmployeeService(ctrl)

			h := &Handler{
				log:             log,
				positionService: positionService,
				employeeService: employeeService,
			}

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

	log := zap.NewNop().Sugar()

	for _, tt := range tests {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = &http.Request{
			Method: http.MethodGet,
			Header: make(http.Header),
		}

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			positionService := serviceMock.NewMockPositionService(ctrl)
			employeeService := serviceMock.NewMockEmployeeService(ctrl)

			h := &Handler{
				log:             log,
				positionService: positionService,
				employeeService: employeeService,
			}

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			h.GetEmployeeByID(ctx, tt.id)
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

	log := zap.NewNop().Sugar()

	for _, tt := range tests {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = &http.Request{
			Method: http.MethodGet,
			Header: make(http.Header),
		}

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			positionService := serviceMock.NewMockPositionService(ctrl)
			employeeService := serviceMock.NewMockEmployeeService(ctrl)

			h := &Handler{
				log:             log,
				positionService: positionService,
				employeeService: employeeService,
			}

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			h.GetEmployeeList(ctx, api.GetEmployeeListParams{Cursor: tt.cursor})
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
			name:       "Invalid ID",
			id:         "invalid",
			input:      `{"first_name":"NewJohn","last_name":"NewDoe", "position_id":"` + uuid.NewString() + `"}`,
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

	log := zap.NewNop().Sugar()

	for _, tt := range tests {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = &http.Request{
			Method: http.MethodPut,
			Header: make(http.Header),
		}

		MockJSONPost(ctx, tt.input)

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			positionService := serviceMock.NewMockPositionService(ctrl)
			employeeService := serviceMock.NewMockEmployeeService(ctrl)

			h := &Handler{
				log:             log,
				positionService: positionService,
				employeeService: employeeService,
			}

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			h.UpdateEmployeeByID(ctx, tt.id)
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

	log := zap.NewNop().Sugar()

	for _, tt := range tests {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = &http.Request{
			Method: http.MethodPost,
			Header: make(http.Header),
		}

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			positionService := serviceMock.NewMockPositionService(ctrl)
			employeeService := serviceMock.NewMockEmployeeService(ctrl)

			h := &Handler{
				log:             log,
				positionService: positionService,
				employeeService: employeeService,
			}

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			h.DeleteEmployeeByID(ctx, tt.id)
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

func MockJSONPost(c *gin.Context, body string) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(body)))
}
