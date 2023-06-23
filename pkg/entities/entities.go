package entities

import (
	"gorm.io/gorm"
)

type Donate struct {
	gorm.Model
	DonaterUsername string
	DonateItemID int
	Amount int `gorm:"default:1"`
	Date int
	PaymentID string `gorm:"unique"`
	PaymentPrice int
}

func (Donate) TableName() string {
    return "donates2"
}

