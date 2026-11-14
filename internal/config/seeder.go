package config

import (
	"evermos-hasby/ecommerce/internal/models"
	"fmt"
)

func SeedData() {
	var count int64
	DB.Model(&models.Category{}).Count(&count)
	if count == 0 {
		categories := []models.Category{
			{CategoryName: "Fashion"},
			{CategoryName: "Accesory"},
			{CategoryName: "House Tools"},
			{CategoryName: "Electronic"},
		}
		if err := DB.Create(&categories).Error; err != nil {
			fmt.Println("Failed to seed:", err)
		} else {
			fmt.Println("Success seed category")
		}
	}
}
