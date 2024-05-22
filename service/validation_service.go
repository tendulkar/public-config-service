// service/validation_service.go
package service

import (
	"context"
	"errors"

	"golang.org/x/exp/slog"
	"stellarsky.ai/platform/public-config-service/model"
	"stellarsky.ai/platform/public-config-service/repository"
)

type ValidationService struct {
	repo   *repository.ValidationRepository
	logger *slog.Logger
}

func NewValidationService(repo *repository.ValidationRepository, logger *slog.Logger) *ValidationService {
	return &ValidationService{
		repo:   repo,
		logger: logger,
	}
}

func (s *ValidationService) GetAllValidations() ([]model.Validation, error) {
	validations, err := s.repo.GetAll(context.Background())
	if err != nil {
		s.logger.Error("error getting all validations", slog.Any("error", err))
		return nil, err
	}
	return validations, nil
}

func (s *ValidationService) GetValidation(id int64) (*model.Validation, error) {
	v, err := s.repo.GetByID(context.Background(), id)
	if err != nil {
		s.logger.Error("error getting validation by id", slog.Any("error", err))
		return nil, err
	}
	if v == nil {
		return nil, errors.New("validation not found")
	}
	return v, nil
}

func (s *ValidationService) CreateValidation(v *model.Validation) error {
	if err := s.repo.Create(context.Background(), v); err != nil {
		s.logger.Error("error creating validation", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *ValidationService) UpdateValidation(v *model.Validation) error {
	if err := s.repo.Update(context.Background(), v); err != nil {
		s.logger.Error("error updating validation", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *ValidationService) DeleteValidation(id int64) error {
	if err := s.repo.Delete(context.Background(), id); err != nil {
		s.logger.Error("error deleting validation", slog.Any("error", err))
		return err
	}
	return nil
}
