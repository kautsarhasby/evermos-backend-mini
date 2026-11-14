package config

import (
	"evermos-hasby/ecommerce/internal/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Store{},
		&models.Address{},
		&models.Category{},
		&models.Product{},
		&models.Transaction{},
		&models.ProductLog{},
		&models.TransactionDetail{},
	)
	if err != nil {
		log.Fatal("Gagal migrasi database:", err)
	}

	log.Println("Success connected to database")
}
