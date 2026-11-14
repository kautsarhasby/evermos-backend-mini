package interfaces

import (
	"context"
	"evermos-hasby/ecommerce/internal/models"
)

type IProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	Get(ctx context.Context, id, userID uint) (*models.Product, error)
	List(ctx context.Context, userID uint, search, sort string, limit, offset int) ([]models.Product, error)
	Update(ctx context.Context, product *models.Product, userID uint) error
	Delete(ctx context.Context, id, userID uint) error
}

type IProductService interface {
	Create(ctx context.Context, product *models.Product, userID uint) error
	Get(ctx context.Context, id, userID uint) (*models.Product, error)
	List(ctx context.Context, userID uint, search, sort string, limit, offset int) ([]models.Product, error)
	Update(ctx context.Context, product *models.Product, userID uint) error
	Delete(ctx context.Context, id, userID uint) error
}
