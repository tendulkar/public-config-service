// service/form_service.go
package service

import (
	"context"
	"errors"

	"golang.org/x/exp/slog"
	"stellarsky.ai/platform/public-config-service/model"
	"stellarsky.ai/platform/public-config-service/repository"
)

type FormService struct {
	repo   *repository.FormRepository
	logger *slog.Logger
}

func NewFormService(repo *repository.FormRepository, logger *slog.Logger) *FormService {
	return &FormService{
		repo:   repo,
		logger: logger,
	}
}

func (s *FormService) GetAllForms() ([]model.Form, error) {
	forms, err := s.repo.GetAll(context.Background())
	if err != nil {
		s.logger.Error("error getting all forms", slog.Any("error", err))
		return nil, err
	}
	return forms, nil
}

func (s *FormService) GetForm(id int64) (*model.Form, error) {
	f, err := s.repo.GetByID(context.Background(), id)
	if err != nil {
		s.logger.Error("error getting form by id", slog.Any("error", err))
		return nil, err
	}
	if f == nil {
		return nil, errors.New("form not found")
	}
	return f, nil
}

func (s *FormService) CreateForm(f *model.Form) error {
	if err := s.repo.Create(context.Background(), f); err != nil {
		s.logger.Error("error creating form", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *FormService) UpdateForm(f *model.Form) error {
	if err := s.repo.Update(context.Background(), f); err != nil {
		s.logger.Error("error updating form", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *FormService) DeleteForm(id int64) error {
	if err := s.repo.Delete(context.Background(), id); err != nil {
		s.logger.Error("error deleting form", slog.Any("error", err))
		return err
	}
	return nil
}
