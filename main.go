package main

import (
	"fmt"

	"github.com/murtll/mcserver-pay/pkg/entities"
	"github.com/murtll/mcserver-pay/pkg/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, _ := gorm.Open(postgres.Open("host=localhost user=test password=test dbname=test port=54321 sslmode=disable"))
	dr := repository.NewDonateRepository(db)
	dr.Migrate()
	dr.AddDonate(&entities.Donate{
		DonaterUsername: "kek",
		DonateItemID: 25,
		Amount: 1,
		Date: 121212121,
		PaymentID: "oabsoabsoabs",
		PaymentPrice: 12,
	})
	fmt.Println(dr.PaymentExist("oabsoabsoabs"))
	fmt.Println(dr.PaymentExist("bebra"))
}