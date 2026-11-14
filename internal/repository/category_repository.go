package repository

import (
	"context"
	"evermos-hasby/ecommerce/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	return r.DB.WithContext(ctx).Create(category).Error
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, id uint, name string) error {
	return r.DB.WithContext(ctx).Model(&models.Category{}).Where("id = ?", id).Update("category_name", name).Error
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&models.Category{}, id).Error
}
