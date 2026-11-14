package repository

import (
	"context"
	"errors"
	"evermos-hasby/ecommerce/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepository struct{ DB *gorm.DB }

func (r *TransactionRepository) Create(ctx context.Context, userID uint, request models.CreateTransactionRequest) (*models.Transaction, error) {
	txModel := &models.Transaction{UserID: userID, ReceiverAddress: request.AddressID, PaymentMethod: request.PaymentMethod, Status: "pending", Invoice: uuid.NewString()}
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(txModel).Error; err != nil {
			return err
		}
		for _, item := range request.Items {
			var product models.Product
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", item.ProductID).First(&product).Error; err != nil {
				return err
			}
			if product.Stock < item.Quantity {
				return errors.New("Insufficient product stock")
			}
			if err := tx.Model(&product).Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return err
			}
			log := models.ProductLog{ProductID: product.ID, ProductName: product.ProductName, Slug: product.Slug, ResellerPrice: product.ResellerPrice, ConsumerPrice: product.ConsumerPrice, Description: product.Description, Image: product.Image, StoreID: product.StoreID, CategoryID: product.CategoryID}
			if err := tx.Create(&log).Error; err != nil {
				return err
			}
			lineTotal := product.ConsumerPrice * float64(item.Quantity)
			detail := models.TransactionDetail{ProductLogID: log.ID, TransactionID: txModel.ID, Quantity: item.Quantity, Total: lineTotal}
			if err := tx.Create(&detail).Error; err != nil {
				return err
			}
			txModel.Total += lineTotal
		}
		return tx.Model(txModel).Update("total", txModel.Total).Error
	})
	if err != nil {
		return nil, err
	}
	return txModel, nil
}

func (r *TransactionRepository) List(ctx context.Context, userID uint, status string, limit, offset int) ([]models.Transaction, error) {
	var result []models.Transaction
	q := r.DB.WithContext(ctx).Preload("TransactionDetail.ProductLog").Where("user_id = ?", userID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&result).Error
	return result, err
}
