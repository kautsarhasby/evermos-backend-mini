package interfaces

import (
	"context"
	"evermos-hasby/ecommerce/internal/models"
)

type ICategoryRepository interface {
	CreateCategory(ctx context.Context, category *models.Category) error
	UpdateCategory(ctx context.Context, id uint, name string) error
	DeleteCategory(ctx context.Context, id uint) error
}

type ICategoryService interface {
	CreateCategory(ctx context.Context, categoryName string, userID uint) (models.Category, error)
	UpdateCategory(ctx context.Context, id uint, name string, userID uint) error
	DeleteCategory(ctx context.Context, id, userID uint) error
}
