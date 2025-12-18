package model

import (
	"time"

	"gorm.io/gorm"
)

type GiftCode struct {
	gorm.Model
	Code      string
	UsedCount int
	MaxUsage  int
	IsActive  bool
}

type GiftCodeUsage struct {
	GiftCodeID  int
	PhoneNumber int
	UsedAt      time.Time
}

type Input struct {
	Phone int    `json:"phone"`
	Code  string `json:"code"`
}
