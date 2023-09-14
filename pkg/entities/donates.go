package entities

import (
	"database/sql/driver"
	"strconv"
	"time"
)

type OrmModel struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

type Donate struct {
	OrmModel
	DonaterUsername string   `json:"donaterUsername"`
	DonateItemID    int      `json:"itemId"`
	Amount          int      `json:"amount" gorm:"default:1"`
	Date            JSONTime `json:"date" gorm:"type:timestamp;default:current_timestamp"`
	PaymentID       int      `json:"-" gorm:"unique"`
	PaymentPrice    int      `json:"price"`
}

type Donateable interface {
	ToDonate() *Donate
}

type Donater struct {
	DonaterUsername string
	MoneySpent      int
}

func (Donate) TableName() string {
	return "donates"
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).UnixMilli(), 10)), nil
}

func (t JSONTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}
