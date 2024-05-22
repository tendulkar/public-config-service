// service/type_service.go
package service

import (
	"context"
	"errors"

	"golang.org/x/exp/slog"
	"stellarsky.ai/platform/public-config-service/model"
	"stellarsky.ai/platform/public-config-service/repository"
)

type TypeService struct {
	repo   *repository.TypeRepository
	logger *slog.Logger
}

func NewTypeService(repo *repository.TypeRepository, logger *slog.Logger) *TypeService {
	return &TypeService{
		repo:   repo,
		logger: logger,
	}
}

func (s *TypeService) GetAllTypes() ([]model.Type, error) {
	types, err := s.repo.GetAll(context.Background())
	if err != nil {
		s.logger.Error("error getting all types", slog.Any("error", err))
		return nil, err
	}
	return types, nil
}

func (s *TypeService) GetType(id int64) (*model.Type, error) {
	t, err := s.repo.GetByID(context.Background(), id)
	if err != nil {
		s.logger.Error("error getting type by id", slog.Any("error", err))
		return nil, err
	}
	if t == nil {
		return nil, errors.New("type not found")
	}
	return t, nil
}

func (s *TypeService) CreateType(t *model.Type) error {
	if err := s.repo.Create(context.Background(), t); err != nil {
		s.logger.Error("error creating type", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *TypeService) UpdateType(t *model.Type) error {
	if err := s.repo.Update(context.Background(), t); err != nil {
		s.logger.Error("error updating type", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *TypeService) DeleteType(id int64) error {
	if err := s.repo.Delete(context.Background(), id); err != nil {
		s.logger.Error("error deleting type", slog.Any("error", err))
		return err
	}
	return nil
}
