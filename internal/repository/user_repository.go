package repository

import (
	"context"
	"errors"
	"evermos-hasby/ecommerce/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) ListUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.DB.WithContext(ctx).Omit("Password").Order("id ASC").Find(&users).Error
	return users, err
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	err := r.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUserProfile(ctx context.Context, user *models.User) error {
	return r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{"username": user.Username, "email": user.Email, "phone": user.Phone, "date_of_birth": user.DateOfBirth, "gender": user.Gender, "occupation": user.Occupation, "province_id": user.ProvinceID, "city_id": user.CityID}).Error
}

func (r *UserRepository) UpdateManagedUser(ctx context.Context, user *models.User) error {
	return r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"username": user.Username, "email": user.Email, "phone": user.Phone,
		"date_of_birth": user.DateOfBirth, "gender": user.Gender, "occupation": user.Occupation,
		"role": user.Role, "province_id": user.ProvinceID, "city_id": user.CityID,
	}).Error
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}
