package repository

import (
	"context"
	"evermos-hasby/ecommerce/internal/models"
	"gorm.io/gorm"
)

type AddressRepository struct {
	DB *gorm.DB
}

func (r *AddressRepository) CreateAddress(ctx context.Context, address *models.Address) error {
	return r.DB.WithContext(ctx).Create(address).Error
}

func (r *AddressRepository) GetAddressesByUserID(ctx context.Context, userID uint, search string, limit, offset int) ([]models.Address, error) {
	var addresses []models.Address
	query := r.DB.WithContext(ctx).Where("user_id = ?", userID)
	if search != "" {
		query = query.Where("address_name LIKE ? OR receiver_name LIKE ? OR city LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&addresses).Error
	return addresses, err
}

func (r *AddressRepository) UpdateAddress(ctx context.Context, address *models.Address, userID uint) error {
	result := r.DB.WithContext(ctx).Model(&models.Address{}).Where("id = ? AND user_id = ?", address.ID, userID).Updates(map[string]interface{}{"address_name": address.AddressName, "receiver_name": address.ReceiverName, "phone": address.Phone, "province": address.Province, "city": address.City, "subdistrict": address.Subdistrict, "ward": address.Ward, "address_detail": address.AddressDetail})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *AddressRepository) DeleteAddress(ctx context.Context, id, userID uint) error {
	result := r.DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.Address{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
