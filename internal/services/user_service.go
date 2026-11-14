package services

import (
	"context"
	"errors"
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountService struct {
	Repository interfaces.IUserRepository
	DB         *gorm.DB
}

func (s *AccountService) List(ctx context.Context, adminID uint) ([]models.User, error) {
	return s.Repository.ListUsers(ctx)
}

func (s *AccountService) Get(ctx context.Context, adminID, id uint) (*models.User, error) {

	user, err := s.Repository.GetUserById(ctx, id)
	if user != nil {
		user.Password = ""
	}
	return user, err
}

func (s *AccountService) Create(ctx context.Context, adminID uint, user *models.User) error {

	if user.Username == "" || user.Email == "" || user.Phone == "" || len(user.Password) < 6 {
		return errors.New("username, email, phone, and password (at least 6 characters) are required")
	}
	if user.Role == "" {
		user.Role = "user"
	}
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(password)
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return tx.Create(&models.Store{Name: user.Username + "Store", UserID: user.ID}).Error
	})
}

func (s *AccountService) Update(ctx context.Context, adminID, id uint, request *models.User) (*models.User, error) {

	if adminID == id {
		return nil, errors.New("admin cannot update own managed account")
	}
	if request.Username == "" || request.Email == "" || request.Phone == "" || request.Role == "" {
		return nil, errors.New("username, email, phone, and role are required")
	}
	user, err := s.Repository.GetUserById(ctx, id)
	if err != nil || user == nil {
		return nil, gorm.ErrRecordNotFound
	}
	user.Username, user.Email, user.Phone, user.Role = request.Username, request.Email, request.Phone, request.Role
	user.DateOfBirth, user.Gender, user.Occupation = request.DateOfBirth, request.Gender, request.Occupation
	user.ProvinceID, user.CityID = request.ProvinceID, request.CityID
	if err := s.Repository.UpdateManagedUser(ctx, user); err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *AccountService) Delete(ctx context.Context, adminID, id uint) error {
	if adminID == id {
		return errors.New("admin cannot delete own account")
	}
	return s.Repository.DeleteUser(ctx, id)
}
