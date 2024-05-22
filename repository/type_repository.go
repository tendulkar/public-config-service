// repository/type_repository.go
package repository

import (
	"context"

	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"stellarsky.ai/platform/public-config-service/model"
)

type TypeRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewTypeRepository(db *gorm.DB, logger *slog.Logger) *TypeRepository {
	return &TypeRepository{
		db:     db,
		logger: logger,
	}
}

func (r *TypeRepository) GetAll(ctx context.Context) ([]model.Type, error) {
	var types []model.Type
	result := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&types)
	if result.Error != nil {
		r.logger.Error("error querying all types", slog.Any("error", result.Error))
		return nil, result.Error
	}
	return types, nil
}

func (r *TypeRepository) GetByID(ctx context.Context, id int64) (*model.Type, error) {
	var t model.Type
	result := r.db.WithContext(ctx).First(&t, "id = ? AND deleted_at IS NULL", id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		r.logger.Error("error querying type by id", slog.Any("error", result.Error))
		return nil, result.Error
	}
	return &t, nil
}

func (r *TypeRepository) Create(ctx context.Context, t *model.Type) error {
	result := r.db.WithContext(ctx).Create(t)
	if result.Error != nil {
		r.logger.Error("error creating type", slog.Any("error", result.Error))
		return result.Error
	}
	return nil
}

func (r *TypeRepository) Update(ctx context.Context, t *model.Type) error {
	result := r.db.WithContext(ctx).Model(t).Updates(map[string]interface{}{
		"namespace":    t.Namespace,
		"family":       t.Family,
		"name":         t.Name,
		"element_type": t.ElementType,
		"widget_type":  t.WidgetType,
		"updated_at":   gorm.Expr("CURRENT_TIMESTAMP"),
		"version":      gorm.Expr("version + 1"),
	})
	if result.Error != nil {
		r.logger.Error("error updating type", slog.Any("error", result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *TypeRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&model.Type{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP"))
	if result.Error != nil {
		r.logger.Error("error deleting type", slog.Any("error", result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
