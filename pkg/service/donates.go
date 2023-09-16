package service

import (
	"fmt"
	"log"
	"math"

	"github.com/murtll/mcserver-pay/pkg/entities"
	"github.com/murtll/mcserver-pay/pkg/pb"
	"github.com/murtll/mcserver-pay/pkg/repository"
)

type DonateService struct {
	dr *repository.DonateRepository
	ir *repository.ItemRepository
	mr *repository.MessageRepository
}

func NewDonateService(dr *repository.DonateRepository,
	ir *repository.ItemRepository,
	mr *repository.MessageRepository,
) *DonateService {
	return &DonateService{
		dr: dr,
		ir: ir,
		mr: mr,
	}
}

func (ds *DonateService) GetLastDonates(count int) ([]entities.Donate, error) {
	return ds.dr.GetLastDonates(count)
}

func (ds *DonateService) ProcessDonate(d entities.Donateable,
	promo string,
	checkSignFunc func(...interface{}) (bool, error),
	checkSignArgs ...interface{}) error {
	donate := d.ToDonate()

	log.Default().Print(*donate)

	if ok, err := ds.dr.PaymentExist(donate.PaymentID); err != nil || ok {
		if ok {
			return fmt.Errorf("payment with id '%d' already exist", donate.PaymentID)
		} else {
			return err
		}
	}

	mul := 1.0
	if promo != "" {
		fullPromo, err := ds.ir.GetPromo(promo)
		if err != nil {
			return err
		}
		mul = fullPromo.Multiplier
	}

	item, err := ds.ir.GetItem(donate.DonateItemID)
	if err != nil {
		return err
	}

	expectedPrice := calculatePrice(item.Price, donate.Amount, mul)

	if expectedPrice != donate.PaymentPrice {
		return fmt.Errorf("expected price %d and actual price %d does not match", expectedPrice, donate.PaymentPrice)
	}

	if checkSignFunc != nil {
		if ok, err := checkSignFunc(checkSignArgs); err != nil || !ok {
			return fmt.Errorf("can't verify request sign")
		}
	}

	_, err = ds.dr.AddDonate(donate)
	if err != nil {
		return err
	}

	err = ds.mr.PublishDonate(&pb.DonateMessage{
		DonateItemId: int32(donate.DonateItemID),
		DonaterUsername: donate.DonaterUsername,
		Amount: int32(donate.Amount),
	})
	if err != nil {
		return err
	}

	return nil
}

func calculatePrice(price, amount int, multiplier float64) int {
	if amount > 1 {
		return int(multiplier * float64(amount) * math.Round(float64(price)*((100.0-float64(calculateSale(amount)))/100.0)))
	} else {
		return int(float64(price) * multiplier)
	}
}

func calculateSale(amount int) int {
	return int(math.Round(50 / (math.Pow(math.E, 3-(float64(amount)/math.Pow(math.Pi, 2))) + 1)))
}
