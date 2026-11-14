package interfaces

import (
	"context"
	"evermos-hasby/ecommerce/internal/models"
)

type ITransactionRepository interface {
	Create(ctx context.Context, userID uint, request models.CreateTransactionRequest) (*models.Transaction, error)
	List(ctx context.Context, userID uint, status string, limit, offset int) ([]models.Transaction, error)
}

type ITransactionService interface {
	Create(ctx context.Context, userID uint, request models.CreateTransactionRequest) (*models.Transaction, error)
	List(ctx context.Context, userID uint, status string, limit, offset int) ([]models.Transaction, error)
}
