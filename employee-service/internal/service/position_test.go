//go:build !integration

package service

import (
	"context"
	"testing"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/Verce11o/resume-view/employee-service/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestPositionService_CreatePosition(t *testing.T) {
	t.Parallel()

	type fields struct {
		positionRepo *mocks.PositionRepository
	}

	positionID := uuid.New()

	tests := []struct {
		name     string
		input    domain.CreatePosition
		response models.Position
		mockFunc func(f *fields)
		wantErr  bool
	}{
		{
			name: "Valid",
			input: domain.CreatePosition{
				ID:     positionID,
				Name:   "Go Developer",
				Salary: 30999,
			},
			response: models.Position{
				ID:     positionID,
				Name:   "Go Developer",
				Salary: 30999,
			},
			mockFunc: func(f *fields) {
				f.positionRepo.On("CreatePosition", mock.Anything, mock.AnythingOfType("domain.CreatePosition")).
					Return(models.Position{
						ID:     positionID,
						Name:   "Go Developer",
						Salary: 30999,
					}, nil)
			},
		},
		{
			name: "Invalid position ID",
			input: domain.CreatePosition{
				ID:     uuid.Nil,
				Name:   "Go Developer",
				Salary: 30999,
			},
			response: models.Position{},
			mockFunc: func(f *fields) {
				f.positionRepo.On("CreatePosition", mock.Anything, mock.AnythingOfType("domain.CreatePosition")).
					Return(models.Position{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positionRepo := mocks.NewPositionRepository(t)
			cache := mocks.NewPositionCacheRepository(t)

			tt.mockFunc(&fields{
				positionRepo: positionRepo,
			})

			srv := &PositionService{
				log:   zap.NewNop().Sugar(),
				cache: cache,
				repo:  positionRepo,
			}

			position, err := srv.CreatePosition(context.TODO(), tt.input)

			assert.Equal(t, tt.wantErr, err != nil)

			assert.EqualValues(t, tt.response, position)
		})
	}
}

func TestPositionService_GetPosition(t *testing.T) {
	t.Parallel()

	type fields struct {
		positionRepo *mocks.PositionRepository
		cache        *mocks.PositionCacheRepository
	}

	positionID := uuid.New()
	tests := []struct {
		name     string
		id       uuid.UUID
		response models.Position
		mockFunc func(f *fields)
		wantErr  bool
	}{
		{
			name: "Valid from cache",
			id:   positionID,
			response: models.Position{
				ID:     positionID,
				Name:   "Go Developer",
				Salary: 30999,
			},
			mockFunc: func(f *fields) {
				f.cache.On("GetPosition", mock.Anything, mock.AnythingOfType("string")).
					Return(&models.Position{
						ID:     positionID,
						Name:   "Go Developer",
						Salary: 30999,
					}, nil)
			},
		},
		{
			name: "Valid from repo",
			id:   positionID,
			response: models.Position{
				ID:     positionID,
				Name:   "Go Developer",
				Salary: 30999,
			},
			mockFunc: func(f *fields) {
				f.cache.On("GetPosition", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, nil)

				f.positionRepo.On("GetPosition", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.Position{
						ID:     positionID,
						Name:   "Go Developer",
						Salary: 30999,
					}, nil)

				f.cache.On("SetPosition", mock.Anything, mock.AnythingOfType("string"),
					mock.AnythingOfType("*models.Position")).
					Return(nil)
			},
		},
		{
			name: "Cache error",
			id:   uuid.Nil,
			response: models.Position{
				ID:     positionID,
				Name:   "Go Developer",
				Salary: 30999,
			},
			mockFunc: func(f *fields) {
				f.cache.On("GetPosition", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, assert.AnError)

				f.positionRepo.On("GetPosition", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.Position{
						ID:     positionID,
						Name:   "Go Developer",
						Salary: 30999,
					}, nil)

				f.cache.On("SetPosition", mock.Anything, mock.AnythingOfType("string"),
					mock.AnythingOfType("*models.Position")).
					Return(assert.AnError)
			},
		},
		{
			name:     "Repository error",
			id:       uuid.Nil,
			response: models.Position{},
			mockFunc: func(f *fields) {
				f.cache.On("GetPosition", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, assert.AnError)

				f.positionRepo.On("GetPosition", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.Position{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positionRepo := mocks.NewPositionRepository(t)
			cache := mocks.NewPositionCacheRepository(t)

			tt.mockFunc(&fields{
				positionRepo: positionRepo,
				cache:        cache,
			})

			srv := &PositionService{
				log:   zap.NewNop().Sugar(),
				cache: cache,
				repo:  positionRepo,
			}

			position, err := srv.GetPosition(context.TODO(), tt.id)

			assert.Equal(t, tt.wantErr, err != nil)

			assert.EqualValues(t, tt.response, position)
		})
	}
}

func TestPositionService_GetPositionList(t *testing.T) {
	t.Parallel()

	type fields struct {
		positionRepo *mocks.PositionRepository
		cache        *mocks.PositionCacheRepository
	}

	positionID := uuid.New()

	tests := []struct {
		name     string
		cursor   string
		response models.PositionList
		mockFunc func(f *fields)
		wantErr  bool
	}{
		{
			name:   "Valid empty cursor",
			cursor: "",
			response: models.PositionList{
				Cursor: "cursorExample",
				Positions: []models.Position{
					{
						ID:     positionID,
						Name:   "Go Developer",
						Salary: 30999,
					},
				},
			},
			mockFunc: func(f *fields) {
				f.positionRepo.On("GetPositionList", mock.Anything, mock.AnythingOfType("string")).
					Return(models.PositionList{
						Cursor: "cursorExample",
						Positions: []models.Position{
							{
								ID:     positionID,
								Name:   "Go Developer",
								Salary: 30999,
							},
						},
					}, nil)
			},
		},
		{
			name:     "Invalid cursor",
			cursor:   "invalid",
			response: models.PositionList{},
			mockFunc: func(f *fields) {
				f.positionRepo.On("GetPositionList", mock.Anything, mock.AnythingOfType("string")).
					Return(models.PositionList{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positionRepo := mocks.NewPositionRepository(t)
			cache := mocks.NewPositionCacheRepository(t)

			tt.mockFunc(&fields{
				positionRepo: positionRepo,
				cache:        cache,
			})

			srv := &PositionService{
				log:   zap.NewNop().Sugar(),
				repo:  positionRepo,
				cache: cache,
			}

			position, err := srv.GetPositionList(context.TODO(), tt.cursor)

			assert.Equal(t, tt.wantErr, err != nil)

			assert.EqualValues(t, tt.response, position)
		})
	}
}

func TestPositionService_UpdatePosition(t *testing.T) {
	t.Parallel()

	type fields struct {
		positionRepo *mocks.PositionRepository
		cache        *mocks.PositionCacheRepository
	}

	positionID := uuid.New()

	tests := []struct {
		name     string
		input    domain.UpdatePosition
		response models.Position
		mockFunc func(f *fields)
		wantErr  bool
	}{
		{
			name: "Valid input",
			input: domain.UpdatePosition{
				ID:     positionID,
				Name:   "Go Developer",
				Salary: 30999,
			},
			response: models.Position{
				ID:     positionID,
				Name:   "Go Developer",
				Salary: 30999,
			},
			mockFunc: func(f *fields) {
				f.positionRepo.On("UpdatePosition", mock.Anything, mock.AnythingOfType("domain.UpdatePosition")).
					Return(models.Position{
						ID:     positionID,
						Name:   "Go Developer",
						Salary: 30999,
					}, nil)

				f.cache.On("DeletePosition", mock.Anything, mock.AnythingOfType("string")).
					Return(nil)
			},
		},
		{
			name: "Invalid input",
			input: domain.UpdatePosition{
				ID:     uuid.Nil,
				Name:   "Go Developer",
				Salary: 30999,
			},
			response: models.Position{},
			mockFunc: func(f *fields) {
				f.positionRepo.On("UpdatePosition", mock.Anything, mock.AnythingOfType("domain.UpdatePosition")).
					Return(models.Position{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Cache error",
			input: domain.UpdatePosition{
				ID:     positionID,
				Name:   "Go Developer",
				Salary: 30999,
			},
			response: models.Position{
				ID:     positionID,
				Name:   "Go Developer",
				Salary: 30999,
			},
			mockFunc: func(f *fields) {
				f.positionRepo.On("UpdatePosition", mock.Anything, mock.AnythingOfType("domain.UpdatePosition")).
					Return(models.Position{
						ID:     positionID,
						Name:   "Go Developer",
						Salary: 30999,
					}, nil)

				f.cache.On("DeletePosition", mock.Anything, mock.AnythingOfType("string")).
					Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positionRepo := mocks.NewPositionRepository(t)
			cache := mocks.NewPositionCacheRepository(t)
			tt.mockFunc(&fields{
				positionRepo: positionRepo,
				cache:        cache,
			})

			srv := &PositionService{
				log:   zap.NewNop().Sugar(),
				repo:  positionRepo,
				cache: cache,
			}
			position, err := srv.UpdatePosition(context.TODO(), tt.input)

			assert.Equal(t, tt.wantErr, err != nil)

			assert.EqualValues(t, tt.response, position)
		})
	}
}

func TestPositionService_DeletePosition(t *testing.T) {
	t.Parallel()

	type fields struct {
		positionRepo *mocks.PositionRepository
		cache        *mocks.PositionCacheRepository
	}

	positionID := uuid.New()

	tests := []struct {
		name     string
		id       uuid.UUID
		mockFunc func(f *fields)
		wantErr  bool
	}{
		{
			name: "Valid input",
			id:   positionID,
			mockFunc: func(f *fields) {
				f.positionRepo.On("GetPosition", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.Position{
						ID:     positionID,
						Name:   "Go Developer",
						Salary: 30999,
					}, nil)

				f.positionRepo.On("DeletePosition", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(nil)

				f.cache.On("DeletePosition", mock.Anything, mock.AnythingOfType("string")).
					Return(nil)
			},
		},
		{
			name: "Non-existing position",
			id:   uuid.New(),
			mockFunc: func(f *fields) {
				f.positionRepo.On("GetPosition", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.Position{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Repository error",
			id:   positionID,
			mockFunc: func(f *fields) {
				f.positionRepo.On("GetPosition", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.Position{
						ID:     positionID,
						Name:   "Go Developer",
						Salary: 30999,
					}, nil)

				f.positionRepo.On("DeletePosition", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Cache error",
			id:   positionID,
			mockFunc: func(f *fields) {
				f.positionRepo.On("GetPosition", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.Position{
						ID:     positionID,
						Name:   "Go Developer",
						Salary: 30999,
					}, nil)

				f.positionRepo.On("DeletePosition", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(nil)

				f.cache.On("DeletePosition", mock.Anything, mock.AnythingOfType("string")).
					Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positionRepo := mocks.NewPositionRepository(t)
			cache := mocks.NewPositionCacheRepository(t)
			tt.mockFunc(&fields{
				positionRepo: positionRepo,
				cache:        cache,
			})

			srv := &PositionService{
				log:   zap.NewNop().Sugar(),
				repo:  positionRepo,
				cache: cache,
			}

			err := srv.DeletePosition(context.TODO(), tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
