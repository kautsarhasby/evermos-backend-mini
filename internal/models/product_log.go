package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductLog struct {
	gorm.Model
	ProductID     uint      `json:"product_id"`
	ProductName   string    `json:"product_name"`
	Slug          string    `json:"slug"`
	ResellerPrice float64   `json:"reseller_price"`
	ConsumerPrice float64   `json:"consumer_price"`
	Description   string    `json:"description"`
	Image         string    `json:"image"`
	StoreID       uint      `json:"store_id"`
	CategoryID    uint      `json:"category_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
