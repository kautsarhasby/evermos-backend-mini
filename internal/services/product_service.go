package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/models"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ProductService struct {
	Repository interfaces.IProductRepository
	DB         *gorm.DB
	Redis      *redis.Client
}

const productCacheTTL = 10 * time.Second

func (s *ProductService) listCacheKey(userID uint, search, sort string, limit, offset int) string {
	query := fmt.Sprintf("%d:%s:%s:%d:%d", userID, search, sort, limit, offset)
	sum := sha256.Sum256([]byte(query))
	return fmt.Sprintf("products:list:%d:%s", userID, hex.EncodeToString(sum[:]))
}

func (s *ProductService) invalidateProductCache(ctx context.Context, userID, productID uint) {
	if s.Redis == nil {
		return
	}
	keys, _, err := s.Redis.Scan(ctx, 0, fmt.Sprintf("products:list:%d:*", userID), 100).Result()
	if err == nil && len(keys) > 0 {
		_ = s.Redis.Del(ctx, keys...).Err()
	}
	_ = s.Redis.Del(ctx, fmt.Sprintf("products:item:%d:%d", userID, productID)).Err()
}

func (s *ProductService) Create(ctx context.Context, p *models.Product, userID uint) error {
	var store models.Store
	if err := s.DB.WithContext(ctx).Where("user_id = ?", userID).First(&store).Error; err != nil {
		return errors.New("Store not found")
	}
	var count int64
	if err := s.DB.WithContext(ctx).Model(&models.Category{}).Where("id = ?", p.CategoryID).Count(&count).Error; err != nil || count == 0 {
		return errors.New("Category not found")
	}
	if p.ProductName == "" || p.Stock < 0 || p.ConsumerPrice < 0 || p.ResellerPrice < 0 {
		return errors.New("Invalid product data")
	}
	p.StoreID = store.ID
	if p.Slug == "" {
		p.Slug = p.ProductName
	}
	if err := s.Repository.Create(ctx, p); err != nil {
		return err
	}
	s.invalidateProductCache(ctx, userID, p.ID)
	return nil
}
func (s *ProductService) Get(ctx context.Context, id, userID uint) (*models.Product, error) {
	key := fmt.Sprintf("products:item:%d:%d", userID, id)
	if s.Redis != nil {
		if cached, err := s.Redis.Get(ctx, key).Bytes(); err == nil {
			var product models.Product
			if json.Unmarshal(cached, &product) == nil {
				return &product, nil
			}
		}
	}
	product, err := s.Repository.Get(ctx, id, userID)
	if err == nil && s.Redis != nil {
		if cached, marshalErr := json.Marshal(product); marshalErr == nil {
			_ = s.Redis.Set(ctx, key, cached, productCacheTTL).Err()
		}
	}
	return product, err
}
func (s *ProductService) List(ctx context.Context, userID uint, search, sort string, limit, offset int) ([]models.Product, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	key := s.listCacheKey(userID, search, sort, limit, offset)
	if s.Redis != nil {
		if cached, err := s.Redis.Get(ctx, key).Bytes(); err == nil {
			var products []models.Product
			if json.Unmarshal(cached, &products) == nil {
				return products, nil
			}
		}
	}
	products, err := s.Repository.List(ctx, userID, search, sort, limit, offset)
	if err == nil && s.Redis != nil {
		if cached, marshalErr := json.Marshal(products); marshalErr == nil {
			_ = s.Redis.Set(ctx, key, cached, productCacheTTL).Err()
		}
	}
	return products, err
}
func (s *ProductService) Update(ctx context.Context, p *models.Product, userID uint) error {
	if p.ProductName == "" || p.Stock < 0 || p.ConsumerPrice < 0 || p.ResellerPrice < 0 {
		return errors.New("Invalid product data")
	}
	if p.Slug == "" {
		p.Slug = p.ProductName
	}
	if _, err := s.Repository.Get(ctx, p.ID, userID); err != nil {
		return errors.New("Product not found")
	}
	if err := s.Repository.Update(ctx, p, userID); err != nil {
		return err
	}
	s.invalidateProductCache(ctx, userID, p.ID)
	return nil
}
func (s *ProductService) Delete(ctx context.Context, id, userID uint) error {
	if _, err := s.Repository.Get(ctx, id, userID); err != nil {
		return errors.New("Product not found")
	}
	if err := s.Repository.Delete(ctx, id, userID); err != nil {
		return err
	}
	s.invalidateProductCache(ctx, userID, id)
	return nil
}
