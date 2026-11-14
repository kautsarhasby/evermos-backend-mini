package services

import (
	"context"
	"errors"
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/models"
	"evermos-hasby/ecommerce/utils"

	"gorm.io/gorm"
)

type CategoryService struct {
	CategoryRepository interfaces.ICategoryRepository
	DB                 *gorm.DB
}

func (s *CategoryService) CreateCategory(ctx context.Context, categoryName string, userID uint) (models.Category, error) {
	log := utils.Logger

	var user models.User
	if err := s.DB.WithContext(ctx).First(&user, userID).Error; err != nil {
		log.WithError(err).Warn("User not found during category creation")
		return models.Category{}, errors.New("User not found")
	}

	if user.Role != "admin" {
		log.WithField("user_id", userID).Warn("Non-admin tried to create category")
		return models.Category{}, errors.New("Forbidden: Hanya Admin yang dapat mengelola kategori")
	}

	categoryModel := models.Category{
		CategoryName: categoryName,
	}

	err := s.CategoryRepository.CreateCategory(ctx, &categoryModel)
	if err != nil {
		log.WithError(err).Error("Failed to create category")
		return models.Category{}, errors.New("Gagal membuat kategori baru")
	}

	return categoryModel, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id uint, name string, userID uint) error {
	return s.CategoryRepository.UpdateCategory(ctx, id, name)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id, userID uint) error {

	return s.CategoryRepository.DeleteCategory(ctx, id)
}
