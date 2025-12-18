package postgres

import (
	"log"
	"os"

	"github.com/AliasgharHeidari/gift-credit/internal/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}

	dsn := os.Getenv("DSN")

	if DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		log.Print("failed to connect to database, error:", err)
	}
}

func AutoMigrate() {
	err := DB.AutoMigrate(&model.GiftCode{}, &model.GiftCodeUsage{})
	if err != nil {
		log.Print("failed to automigrate, error:", err)
	}
}

func GetDB() *gorm.DB {
	return DB
}
