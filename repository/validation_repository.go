// repository/validation_repository.go
package repository

import (
	"context"

	"stellarsky.ai/platform/public-config-service/model"

	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

type ValidationRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewValidationRepository(db *gorm.DB, logger *slog.Logger) *ValidationRepository {
	return &ValidationRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ValidationRepository) GetAll(ctx context.Context) ([]model.Validation, error) {
	var validations []model.Validation
	result := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&validations)
	if result.Error != nil {
		r.logger.Error("error querying all validations", slog.Any("error", result.Error))
		return nil, result.Error
	}
	return validations, nil
}

func (r *ValidationRepository) GetByID(ctx context.Context, id int64) (*model.Validation, error) {
	var v model.Validation
	result := r.db.WithContext(ctx).First(&v, "id = ? AND deleted_at IS NULL", id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		r.logger.Error("error querying validation by id", slog.Any("error", result.Error))
		return nil, result.Error
	}
	return &v, nil
}

func (r *ValidationRepository) Create(ctx context.Context, v *model.Validation) error {
	result := r.db.WithContext(ctx).Create(v)
	if result.Error != nil {
		r.logger.Error("error creating validation", slog.Any("error", result.Error))
		return result.Error
	}
	return nil
}

func (r *ValidationRepository) Update(ctx context.Context, v *model.Validation) error {
	result := r.db.WithContext(ctx).Model(v).Updates(map[string]interface{}{
		"namespace":         v.Namespace,
		"family":            v.Family,
		"name":              v.Name,
		"rule_name":         v.RuleName,
		"validation_params": v.ValidationParams,
		"updated_at":        gorm.Expr("CURRENT_TIMESTAMP"),
		"version":           gorm.Expr("version + 1"),
	})
	if result.Error != nil {
		r.logger.Error("error updating validation", slog.Any("error", result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ValidationRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&model.Validation{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP"))
	if result.Error != nil {
		r.logger.Error("error deleting validation", slog.Any("error", result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
