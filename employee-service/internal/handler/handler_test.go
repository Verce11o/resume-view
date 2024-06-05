//go:build !integration

package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	chiHandler "github.com/Verce11o/resume-view/employee-service/internal/handler/http/chi"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	serviceMock "github.com/Verce11o/resume-view/employee-service/internal/service/mocks"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type m map[string]any

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
			input:      `{{{{`,
			response:   models.Employee{},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			req, err := http.NewRequest(http.MethodPost, "/employees", bytes.NewBufferString(tt.input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Post("/employees", h.CreateEmployee)

			asd := mux.NewRouter()

			asd.HandleFunc("/employees", h.CreateEmployee)

			r.ServeHTTP(rr, req)

			assert.EqualValues(t, tt.statusCode, rr.Code)

			if tt.response != nil {
				var responseBody models.Employee
				err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
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
		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			req, err := http.NewRequest(http.MethodGet, "/employees/"+tt.id, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Get("/employees/{id}", h.GetEmployeeByID)

			r.ServeHTTP(rr, req)

			assert.EqualValues(t, tt.statusCode, rr.Code)

			if tt.response != nil {
				var responseBody models.Employee
				err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
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
		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			req, err := http.NewRequest(http.MethodGet, "/employees?cursor="+*tt.cursor, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Get("/employees", h.GetEmployeeList)

			r.ServeHTTP(rr, req)

			assert.EqualValues(t, tt.statusCode, rr.Code)

			if tt.response != nil {
				var responseBody models.EmployeeList
				err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
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
		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			req, err := http.NewRequest(http.MethodPut, "/employees/"+tt.id, bytes.NewBufferString(tt.input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Put("/employees/{id}", h.UpdateEmployeeByID)

			r.ServeHTTP(rr, req)

			assert.EqualValues(t, tt.statusCode, rr.Code)

			if tt.response != nil {
				var responseBody models.Employee
				err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
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
			response: map[string]string{
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
			response: map[string]string{
				"message": errors.New("invalid UUID length: 7").Error(),
			},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			req, err := http.NewRequest(http.MethodDelete, "/employees/"+tt.id, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Delete("/employees/{id}", h.DeleteEmployeeByID)

			r.ServeHTTP(rr, req)

			assert.EqualValues(t, tt.statusCode, rr.Code)

			if tt.response != nil {
				var responseBody map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
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
			input:      `{{{`,
			response:   models.Position{},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			req, err := http.NewRequest(http.MethodPost, "/positions", bytes.NewBufferString(tt.input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Post("/positions", h.CreatePosition)

			r.ServeHTTP(rr, req)

			assert.EqualValues(t, tt.statusCode, rr.Code)

			if tt.response != nil {
				var responseBody models.Position
				err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
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
		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			req, err := http.NewRequest(http.MethodGet, "/positions/"+tt.id, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Get("/positions/{id}", h.GetPositionByID)

			r.ServeHTTP(rr, req)

			assert.EqualValues(t, tt.statusCode, rr.Code)

			if tt.response != nil {
				var responseBody models.Position

				err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
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
			Name:   "Python Developer",
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
		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			req, err := http.NewRequest(http.MethodGet, "/positions?cursor="+*tt.cursor, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Get("/positions", h.GetPositionList)

			r.ServeHTTP(rr, req)

			assert.EqualValues(t, tt.statusCode, rr.Code)

			if tt.response != nil {
				var responseBody models.PositionList
				err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
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
			input:      `{{{`,
			response:   models.Position{},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			req, err := http.NewRequest(http.MethodPut, "/positions/"+tt.id, bytes.NewBufferString(tt.input))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Put("/positions/{id}", h.UpdatePositionByID)

			r.ServeHTTP(rr, req)

			assert.EqualValues(t, tt.statusCode, rr.Code)

			if tt.response != nil {
				var responseBody models.Position
				err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
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
			response: map[string]string{
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
			response: map[string]string{
				"message": "invalid UUID length: 7",
			},
			mockFunc:   func(_ *fields) {},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, employeeService, positionService, h := initMocks(t)
			defer ctrl.Finish()

			tt.mockFunc(&fields{
				employeeService: employeeService,
				positionService: positionService,
			})

			req, err := http.NewRequest(http.MethodDelete, "/positions/"+tt.id, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Delete("/positions/{id}", h.DeletePositionByID)

			r.ServeHTTP(rr, req)

			assert.EqualValues(t, tt.statusCode, rr.Code)

			if tt.response != nil {
				var responseBody map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.EqualValues(t, tt.response, responseBody)
			}
		})
	}
}

func initMocks(t *testing.T) (*gomock.Controller, *serviceMock.MockEmployeeService,
	*serviceMock.MockPositionService, Handler) {
	ctrl := gomock.NewController(t)
	positionService := serviceMock.NewMockPositionService(ctrl)
	employeeService := serviceMock.NewMockEmployeeService(ctrl)
	authService := serviceMock.NewMockAuthService(ctrl)

	log := zap.NewNop().Sugar()

	h := chiHandler.New(log, positionService, employeeService, authService)

	return ctrl, employeeService, positionService, h
}
