package services

import (
	"context"
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/models"
	"evermos-hasby/ecommerce/utils"

	"github.com/pkg/errors"
)

type AddressService struct {
	AddressRepository interfaces.IAddressRepository
}

func (s *AddressService) CreateAddress(ctx context.Context, request models.Address, userID uint) (models.Address, error) {
	var (
		log = utils.Logger
	)
	if request.AddressName == "" || request.ReceiverName == "" || request.Phone == "" || request.AddressDetail == "" {
		return models.Address{}, errors.New("Required address fields are missing")
	}

	addressModel := models.Address{
		UserID:        userID,
		AddressName:   request.AddressName,
		ReceiverName:  request.ReceiverName,
		Phone:         request.Phone,
		Province:      request.Province,
		City:          request.City,
		Subdistrict:   request.Subdistrict,
		Ward:          request.Ward,
		AddressDetail: request.AddressDetail,
	}

	err := s.AddressRepository.CreateAddress(ctx, &addressModel)
	if err != nil {
		log.WithError(err).Error("Failed to create address")
		return models.Address{}, errors.Wrap(err, "Failed to create address")
	}

	return addressModel, nil
}

func (s *AddressService) GetAddresses(ctx context.Context, userID uint, search string, limit, offset int) ([]models.Address, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.AddressRepository.GetAddressesByUserID(ctx, userID, search, limit, offset)
}
func (s *AddressService) UpdateAddress(ctx context.Context, request models.Address, id, userID uint) error {
	if request.AddressName == "" || request.ReceiverName == "" || request.Phone == "" || request.AddressDetail == "" {
		return errors.New("Required address fields are missing")
	}
	request.ID = id
	return s.AddressRepository.UpdateAddress(ctx, &request, userID)
}
func (s *AddressService) DeleteAddress(ctx context.Context, id, userID uint) error {
	return s.AddressRepository.DeleteAddress(ctx, id, userID)
}
