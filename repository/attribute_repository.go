// repository/attribute_repository.go
package repository

import (
	"context"

	"stellarsky.ai/platform/public-config-service/model"

	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

type AttributeRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewAttributeRepository(db *gorm.DB, logger *slog.Logger) *AttributeRepository {
	return &AttributeRepository{
		db:     db,
		logger: logger,
	}
}

func (r *AttributeRepository) GetAll(ctx context.Context) ([]model.Attribute, error) {
	var attributes []model.Attribute
	result := r.db.WithContext(ctx).
		Joins("LEFT JOIN attribute_validations avs ON attributes.id = avs.attribute_id").
		Joins("LEFT JOIN validations ON avs.validation_id = validations.id AND validations.deleted_at is NULL").
		Joins("LEFT JOIN types ON attributes.type_id = types.id AND types.deleted_at is NULL").
		Preload("Type").
		Preload("Validations").
		Find(&attributes, "attributes.deleted_at IS NULL")

	if result.Error != nil {
		r.logger.Error("error querying all attributes", slog.Any("error", result.Error))
		return nil, result.Error
	}
	return attributes, nil
}

func (r *AttributeRepository) GetByID(ctx context.Context, id int64) (*model.Attribute, error) {
	var a model.Attribute
	result := r.db.WithContext(ctx).
		Joins("LEFT JOIN attribute_validations avs ON attributes.id = avs.attribute_id AND avs.attribute_id = ?", id).
		Joins("LEFT JOIN validations ON avs.validation_id = validations.id AND validations.deleted_at is NULL").
		Joins("LEFT JOIN types ON attributes.type_id = types.id AND types.deleted_at is NULL").
		Preload("Type").
		Preload("Validations").
		First(&a, "attributes.id = ? AND attributes.deleted_at IS NULL", id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		r.logger.Error("error querying attribute by id", slog.Any("error", result.Error))
		return nil, result.Error
	}
	return &a, nil
}

func (r *AttributeRepository) Create(ctx context.Context, a *model.Attribute) error {
	result := r.db.WithContext(ctx).Create(a)
	if result.Error != nil {
		r.logger.Error("error creating attribute", slog.Any("error", result.Error))
		return result.Error
	}
	return nil
}

func (r *AttributeRepository) Update(ctx context.Context, a *model.Attribute) error {
	result := r.db.WithContext(ctx).Model(a).Updates(map[string]interface{}{
		"namespace":   a.Namespace,
		"family":      a.Family,
		"name":        a.Name,
		"label":       a.Label,
		"design_spec": a.DesignSpec,
		"updated_at":  gorm.Expr("CURRENT_TIMESTAMP"),
		"version":     gorm.Expr("version + 1"),
	})
	if result.Error != nil {
		r.logger.Error("error updating attribute", slog.Any("error", result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *AttributeRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&model.Attribute{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP"))
	if result.Error != nil {
		r.logger.Error("error deleting attribute", slog.Any("error", result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
