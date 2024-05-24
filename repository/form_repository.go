// repository/form_repository.go
package repository

import (
	"context"

	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"stellarsky.ai/platform/public-config-service/model"
)

type FormRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewFormRepository(db *gorm.DB, logger *slog.Logger) *FormRepository {
	return &FormRepository{
		db:     db,
		logger: logger,
	}
}

func (r *FormRepository) GetAll(ctx context.Context) ([]model.Form, error) {
	var forms []model.Form
	result := r.db.WithContext(ctx).
		Joins("LEFT JOIN form_attributes fas ON fas.form_id = forms.id").
		Joins("LEFT JOIN attributes ON fas.attribute_id = attributes.id AND attributes.deleted_at is NULL").
		Joins("LEFT JOIN attribute_validations avs ON attributes.id = avs.attribute_id").
		Joins("LEFT JOIN validations ON avs.validation_id = validations.id AND validations.deleted_at is NULL").
		Joins("LEFT JOIN types ON attributes.type_id = types.id AND types.deleted_at is NULL").
		Preload("Attributes").
		Preload("Attributes.Type").
		Preload("Attributes.Validations").
		Find(&forms, "forms.deleted_at IS NULL")
	if result.Error != nil {
		r.logger.Error("error querying all forms", slog.Any("error", result.Error))
		return nil, result.Error
	}
	return forms, nil
}

func (r *FormRepository) GetByID(ctx context.Context, id int64) (*model.Form, error) {
	var f model.Form
	result := r.db.WithContext(ctx).
		// Joins("LEFT JOIN form_attributes fas ON fas.form_id = forms.id").
		// Joins("LEFT JOIN attributes ON fas.attribute_id = attributes.id AND attributes.deleted_at is NULL").
		// Joins("LEFT JOIN attribute_validations avs ON attributes.id = avs.attribute_id").
		// Joins("LEFT JOIN validations ON avs.validation_id = validations.id AND validations.deleted_at is NULL").
		// Joins("LEFT JOIN types ON attributes.type_id = types.id AND types.deleted_at is NULL").
		Preload("Attributes").
		Preload("Attributes.Type").
		Preload("Attributes.Validations").
		First(&f, "forms.id = ? AND forms.deleted_at IS NULL", id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		r.logger.Error("error querying form by id", slog.Any("error", result.Error))
		return nil, result.Error
	}
	return &f, nil
}

func (r *FormRepository) Create(ctx context.Context, f *model.Form) error {
	r.logger.Info("Form to be created", slog.Any("form", f))
	result := r.db.WithContext(ctx).Create(f)
	if result.Error != nil {
		r.logger.Error("error creating form", slog.Any("error", result.Error))
		return result.Error
	}
	return nil
}

func (r *FormRepository) Update(ctx context.Context, f *model.Form) error {
	result := r.db.WithContext(ctx).Model(f).Updates(map[string]interface{}{
		"namespace":   f.Namespace,
		"family":      f.Family,
		"name":        f.Name,
		"action_name": f.ActionName,
		"updated_at":  gorm.Expr("CURRENT_TIMESTAMP"),
		"version":     gorm.Expr("version + 1"),
	})
	if result.Error != nil {
		r.logger.Error("error updating form", slog.Any("error", result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *FormRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&model.Form{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP"))
	if result.Error != nil {
		r.logger.Error("error deleting form", slog.Any("error", result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
