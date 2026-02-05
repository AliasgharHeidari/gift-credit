package postgres

import (
	"log"

	"github.com/AliasgharHeidari/gift-credit/config"
	"github.com/AliasgharHeidari/gift-credit/internal/model"
	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg config.DatabaseConfig) {

	dsn := cfg.DSN.String()

	var err error
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
