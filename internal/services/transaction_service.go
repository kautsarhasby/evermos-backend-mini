package services

import (
	"context"
	"errors"
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/models"
	"gorm.io/gorm"
)

type TransactionService struct {
	Repository interfaces.ITransactionRepository
	DB         *gorm.DB
}

func (s *TransactionService) Create(ctx context.Context, userID uint, request models.CreateTransactionRequest) (*models.Transaction, error) {
	if request.AddressID == 0 || len(request.Items) == 0 {
		return nil, errors.New("Address and transaction items are required")
	}
	var count int64
	if err := s.DB.WithContext(ctx).Model(&models.Address{}).Where("id = ? AND user_id = ?", request.AddressID, userID).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("Address not found")
	}
	for _, item := range request.Items {
		if item.ProductID == 0 || item.Quantity < 1 {
			return nil, errors.New("Invalid transaction item")
		}
	}
	return s.Repository.Create(ctx, userID, request)
}
func (s *TransactionService) List(ctx context.Context, userID uint, status string, limit, offset int) ([]models.Transaction, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.Repository.List(ctx, userID, status, limit, offset)
}
