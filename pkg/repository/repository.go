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

func (dr *DonateRepository) GetLastDonates(count int) ([]*entities.Donate, error) {
	var donates []*entities.Donate
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

// export const checkIfPaymentExists = (paymentId) => {
//     return db('donates').first().where({ payment_id: paymentId })
// }

// export const getTopDonaters = () => {
//     return db('donates').select('donater_username', 'sum(payment_price)').groupBy('donater_username').orderBy('sum(payment_price)', 'desc')
// }


func (dr *DonateRepository) Migrate() {
	dr.db.AutoMigrate(&entities.Donate{})
}