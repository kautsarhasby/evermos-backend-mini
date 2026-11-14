package repository

import (
	"context"
	"evermos-hasby/ecommerce/internal/models"

	"gorm.io/gorm"
)

type StoreRepository struct {
	DB *gorm.DB
}

func (r *StoreRepository) CreateStore(ctx context.Context, store *models.Store) error {
	return r.DB.WithContext(ctx).Create(store).Error
}

func (r *StoreRepository) GetStoreByUserID(ctx context.Context, userID uint) (*models.Store, error) {
	var store models.Store
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).First(&store).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &store, nil
}

func (r *StoreRepository) UpdateStore(ctx context.Context, store *models.Store, userID uint) error {
	result := r.DB.WithContext(ctx).Model(&models.Store{}).Where("user_id = ?", userID).Updates(map[string]interface{}{"name": store.Name, "photo": store.Photo})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
