package entities

import (
	"strconv"
	"time"
)

type StatusResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type ErrorResponse struct {
	Error string `json:"message"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error: err.Error(),
	}
}

type FkPaymentRequest struct {
	ID              int    `json:"intid" form:"intid"`
	DonateItemID    string `json:"us_donate" form:"us_donate"`
	DonaterUsername string `json:"us_username" form:"us_username"`
	Amount          string `json:"us_number" form:"us_number"`
	Promo           string `json:"us_promo" form:"us_promo"`
	PaymentPrice    string `json:"AMOUNT" form:"AMOUNT"`
}

func (fk *FkPaymentRequest) ToDonate() *Donate {
	itemId, _ := strconv.Atoi(fk.DonateItemID)
	paymentPrice, _ := strconv.Atoi(fk.PaymentPrice)
	amount, _ := strconv.Atoi(fk.Amount)

	return &Donate{
		DonaterUsername: fk.DonaterUsername,
		DonateItemID:    itemId,
		Amount:          amount,
		Date:            JSONTime(time.Now()),
		PaymentID:       fk.ID,
		PaymentPrice:    paymentPrice,
	}
}
