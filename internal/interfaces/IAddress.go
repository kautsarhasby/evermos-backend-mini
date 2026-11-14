package interfaces

import (
	"context"
	"evermos-hasby/ecommerce/internal/models"
)

type IAddressRepository interface {
	CreateAddress(ctx context.Context, address *models.Address) error
	GetAddressesByUserID(ctx context.Context, userID uint, search string, limit, offset int) ([]models.Address, error)
	UpdateAddress(ctx context.Context, address *models.Address, userID uint) error
	DeleteAddress(ctx context.Context, id, userID uint) error
}

type IAddressService interface {
	CreateAddress(ctx context.Context, request models.Address, userID uint) (models.Address, error)
	GetAddresses(ctx context.Context, userID uint, search string, limit, offset int) ([]models.Address, error)
	UpdateAddress(ctx context.Context, request models.Address, id, userID uint) error
	DeleteAddress(ctx context.Context, id, userID uint) error
}
