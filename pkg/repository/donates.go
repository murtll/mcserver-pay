package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/murtll/mcserver-pay/pkg/entities"
)

type DonateRepository struct {
	db *gorm.DB
}

func NewDonateRepository(db *gorm.DB) *DonateRepository {
	return &DonateRepository{
		db: db,
	}
}

func (dr *DonateRepository) AddDonate(donate *entities.Donate) (*entities.Donate, error) {
	result := dr.db.Create(donate)
	return donate, result.Error
}

func (dr *DonateRepository) GetLastDonates(count int) ([]entities.Donate, error) {
	var donates []entities.Donate
	result := dr.db.Order("date DESC").Limit(count).Find(&donates)
	if result.Error != nil {
		return nil, result.Error
	}
	return donates, nil
}

func (dr *DonateRepository) PaymentExist(paymentID string) (bool, error) {
	var donate entities.Donate
	result := dr.db.Where("payment_id = ?", paymentID).Take(&donate)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, result.Error
	}
	return true, nil
}

func (dr *DonateRepository) GetTopDonaters(count int) ([]entities.Donater, error) {
	var donaters []entities.Donater

	result := dr.db.Model(&entities.Donate{}).
		Select("donater_username, sum(payment_price) as money_spent").
		Group("donater_username").
		Order("money_spent DESC").
		Limit(count).
		Find(&donaters)

	if result.Error != nil {
		return nil, result.Error
	}

	return donaters, nil
}

func (dr *DonateRepository) Migrate() error {
	return dr.db.AutoMigrate(&entities.Donate{})
}
