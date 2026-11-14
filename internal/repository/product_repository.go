package repository

import (
	"context"
	"evermos-hasby/ecommerce/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct{ DB *gorm.DB }

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	return r.DB.WithContext(ctx).Create(product).Error
}

func (r *ProductRepository) Get(ctx context.Context, id, userID uint) (*models.Product, error) {
	var product models.Product
	err := r.DB.WithContext(ctx).Joins("JOIN stores ON stores.id = products.store_id").Where("products.id = ? AND stores.user_id = ?", id, userID).First(&product).Error
	return &product, err
}

func (r *ProductRepository) List(ctx context.Context, userID uint, search, sort string, limit, offset int) ([]models.Product, error) {
	var products []models.Product
	query := r.DB.WithContext(ctx).Joins("JOIN stores ON stores.id = products.store_id").Where("stores.user_id = ?", userID)
	if search != "" {
		query = query.Where("(products.product_name LIKE ? OR products.description LIKE ?)", "%"+search+"%", "%"+search+"%")
	}
	if sort == "price_asc" {
		query = query.Order("products.consumer_price ASC")
	} else if sort == "price_desc" {
		query = query.Order("products.consumer_price DESC")
	} else {
		query = query.Order("products.created_at DESC")
	}
	err := query.Limit(limit).Offset(offset).Find(&products).Error
	return products, err
}

func (r *ProductRepository) Update(ctx context.Context, product *models.Product, userID uint) error {
	return r.DB.WithContext(ctx).Model(&models.Product{}).Where("id = ? AND store_id IN (?)", product.ID, r.DB.Model(&models.Store{}).Select("id").Where("user_id = ?", userID)).Updates(map[string]interface{}{"product_name": product.ProductName, "slug": product.Slug, "reseller_price": product.ResellerPrice, "consumer_price": product.ConsumerPrice, "stock": product.Stock, "description": product.Description, "image": product.Image, "category_id": product.CategoryID}).Error
}

func (r *ProductRepository) Delete(ctx context.Context, id, userID uint) error {
	return r.DB.WithContext(ctx).Where("store_id IN (?)", r.DB.Model(&models.Store{}).Select("id").Where("user_id = ?", userID)).Delete(&models.Product{}, id).Error
}
