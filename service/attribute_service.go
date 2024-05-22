// service/attribute_service.go
package service

import (
	"context"
	"errors"

	"golang.org/x/exp/slog"
	"stellarsky.ai/platform/public-config-service/model"
	"stellarsky.ai/platform/public-config-service/repository"
)

type AttributeService struct {
	repo   *repository.AttributeRepository
	logger *slog.Logger
}

func NewAttributeService(repo *repository.AttributeRepository, logger *slog.Logger) *AttributeService {
	return &AttributeService{
		repo:   repo,
		logger: logger,
	}
}

func (s *AttributeService) GetAllAttributes() ([]model.Attribute, error) {
	attributes, err := s.repo.GetAll(context.Background())
	if err != nil {
		s.logger.Error("error getting all attributes", slog.Any("error", err))
		return nil, err
	}
	return attributes, nil
}

func (s *AttributeService) GetAttribute(id int64) (*model.Attribute, error) {
	a, err := s.repo.GetByID(context.Background(), id)
	if err != nil {
		s.logger.Error("error getting attribute by id", slog.Any("error", err))
		return nil, err
	}
	if a == nil {
		return nil, errors.New("attribute not found")
	}
	return a, nil
}

func (s *AttributeService) CreateAttribute(a *model.Attribute) error {
	if err := s.repo.Create(context.Background(), a); err != nil {
		s.logger.Error("error creating attribute", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *AttributeService) UpdateAttribute(a *model.Attribute) error {
	if err := s.repo.Update(context.Background(), a); err != nil {
		s.logger.Error("error updating attribute", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *AttributeService) DeleteAttribute(id int64) error {
	if err := s.repo.Delete(context.Background(), id); err != nil {
		s.logger.Error("error deleting attribute", slog.Any("error", err))
		return err
	}
	return nil
}
