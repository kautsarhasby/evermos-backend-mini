package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductName   string  `json:"product_name"`
	Slug          string  `json:"slug"`
	ResellerPrice float64 `json:"reseller_price"`
	ConsumerPrice float64 `json:"consumer_price"`
	Stock         int     `json:"stock"`
	Description   string  `json:"description"`
	Image         string  `json:"image"`
	StoreID       uint    `json:"store_id"`
	CategoryID    uint    `json:"category_id"`

	Store    Store    `json:"store"`
	Category Category `json:"category"`
}
