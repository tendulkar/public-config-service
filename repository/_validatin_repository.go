package repository

// import (
// 	"context"

// 	"config-service/model"
// 	"golang.org/x/exp/slog"
// 	"gorm.io/gorm"
// )

// type AttributeRepository struct {
// 	db     *gorm.DB
// 	logger *slog.Logger
// }

// func NewAttributeRepository(db *gorm.DB, logger *slog.Logger) *AttributeRepository {
// 	return &AttributeRepository{
// 		db:     db,
// 		logger: logger,
// 	}
// }

// func (r *AttributeRepository) GetAll(ctx context.Context) ([]model.Attribute, error) {
// 	var attributes []model.Attribute
// 	result := r.db.WithContext(ctx).Preload("Type").Preload("Validations").Preload("Forms").Find(&attributes)
// 	if result.Error != nil {
// 		r.logger.Error("error querying all attributes", slog.Any("error", result.Error))
// 		return nil, result.Error
// 	}
// 	return attributes, nil
// }

// func (r *AttributeRepository) GetByID(ctx context.Context, id uint) (*model.Attribute, error) {
// 	var a model.Attribute
// 	result := r.db.WithContext(ctx).Preload("Type").Preload("Validations").Preload("Forms").First(&a, id)
// 	if result.Error == gorm.ErrRecordNotFound {
// 		return nil, nil
// 	}
// 	if result.Error != nil {
// 		r.logger.Error("error querying attribute by id", slog.Any("error", result.Error))
// 		return nil, result.Error
// 	}
// 	return &a, nil
// }

// func (r *AttributeRepository) Create(ctx context.Context, a *model.Attribute) error {
// 	result := r.db.WithContext(ctx).Create(a)
// 	if result.Error != nil {
// 		r.logger.Error("error creating attribute", slog.Any("error", result.Error))
// 		return result.Error
// 	}
// 	return nil
// }

// func (r *AttributeRepository) Update(ctx context.Context, a *model.Attribute) error {
// 	result := r.db.WithContext(ctx).Save(a)
// 	if result.Error != nil {
// 		r.logger.Error("error updating attribute", slog.Any("error", result.Error))
// 		return result.Error
// 	}
// 	return nil
// }

// func (r *AttributeRepository) Delete(ctx context.Context, id uint) error {
// 	result := r.db.WithContext(ctx).Delete(&model.Attribute{}, id)
// 	if result.Error != nil {
// 		r.logger.Error("error deleting attribute", slog.Any("error", result.Error))
// 		return result.Error
// 	}
// 	return nil
// }
