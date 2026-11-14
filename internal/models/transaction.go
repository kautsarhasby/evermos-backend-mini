package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ReceiverAddress uint    `json:"receiver_address"`
	Total           float64 `json:"total"`
	Invoice         string  `json:"invoice"`
	PaymentMethod   string  `json:"payment_method"`
	Status          string  `json:"status"`
	UserID          uint    `json:"user_id"`

	User              User                `json:"-"`
	TransactionDetail []TransactionDetail `json:"transaction_detail" gorm:"foreignKey:TransactionID"`
}

type TransactionItemRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type CreateTransactionRequest struct {
	AddressID     uint                     `json:"address_id"`
	PaymentMethod string                   `json:"payment_method"`
	Items         []TransactionItemRequest `json:"items"`
}
