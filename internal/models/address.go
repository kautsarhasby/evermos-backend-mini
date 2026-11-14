package models

import (
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	UserID        uint   `json:"user_id"`
	AddressName   string `json:"address_name"`
	ReceiverName  string `json:"receiver_name"`
	Phone         string `json:"phone"`
	Province      string `json:"province"`
	City          string `json:"city"`
	Subdistrict   string `json:"subdisctrict"`
	Ward          string `json:"ward"`
	AddressDetail string `json:"address_detail"`
}
