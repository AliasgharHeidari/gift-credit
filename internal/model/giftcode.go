package model

import (
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
	gorm.Model
	MobileNumber int
}

type GiftCodeStatus struct {
	GiftCode string
}

type Input struct {
	Phone int    `json:"phone"`
	Code  string `json:"code"`
}

type NewBalance struct {
	Balance float64 `json:"balance"`
} 