package models

import (
	"gorm.io/gorm"
)

type TransactionDetail struct {
	gorm.Model
	ProductLogID  uint    `json:"product_log_id"`
	TransactionID uint    `json:"transaction_id"`
	Quantity      int     `json:"quantity"`
	Total         float64 `json:"total"`

	ProductLog  ProductLog  `json:"product_log" gorm:"foreignKey:ProductLogID"`
	Transaction Transaction `json:"-" gorm:"foreignKey:TransactionID"`
}
