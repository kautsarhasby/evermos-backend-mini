package services

import (
	"context"
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/models"
	"evermos-hasby/ecommerce/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type StoreService struct {
	StoreRepository interfaces.IStoreRepository
	DB              *gorm.DB
}

func (s *StoreService) CreateStore(ctx context.Context, request models.Store, userID uint) (models.Store, error) {
	var (
		log = utils.Logger
	)

	storeModel := models.Store{
		Name:   request.Name,
		UserID: userID,
	}

	err := s.StoreRepository.CreateStore(ctx, &storeModel)
	if err != nil {
		log.WithError(err).Error("Failed to create store")
		return models.Store{}, errors.Wrap(err, "Failed to create store")
	}

	return storeModel, nil
}

func (s *StoreService) GetStore(ctx context.Context, userID uint) (*models.Store, error) {
	return s.StoreRepository.GetStoreByUserID(ctx, userID)
}
func (s *StoreService) UpdateStore(ctx context.Context, request models.Store, userID uint) (*models.Store, error) {
	if request.Name == "" {
		return nil, errors.New("Store name is required")
	}
	if err := s.StoreRepository.UpdateStore(ctx, &request, userID); err != nil {
		return nil, err
	}
	return s.StoreRepository.GetStoreByUserID(ctx, userID)
}
