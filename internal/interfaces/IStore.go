package interfaces

import (
	"context"
	"evermos-hasby/ecommerce/internal/models"
)

type IStoreRepository interface {
	CreateStore(ctx context.Context, store *models.Store) error
	GetStoreByUserID(ctx context.Context, userID uint) (*models.Store, error)
	UpdateStore(ctx context.Context, store *models.Store, userID uint) error
}

type IStoreService interface {
	CreateStore(ctx context.Context, request models.Store, userID uint) (models.Store, error)
	GetStore(ctx context.Context, userID uint) (*models.Store, error)
	UpdateStore(ctx context.Context, request models.Store, userID uint) (*models.Store, error)
}
