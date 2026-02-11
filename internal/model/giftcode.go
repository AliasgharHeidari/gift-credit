package model

import (
	"gorm.io/gorm"
)

type GiftCode struct {
	gorm.Model
	Code      string `json:"Code"`
	UsedCount int    `json:"UseCount"`
	MaxUsage  int    `json:"MaxUsage"`
	IsActive  bool   `json:"IsActive"`
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
