package service

import (
	"fmt"
	"math"

	"github.com/murtll/mcserver-pay/pkg/entities"
	"github.com/murtll/mcserver-pay/pkg/repository"
)

type DonateService struct {
	dr *repository.DonateRepository
	ir *repository.ItemRepository
}

func NewDonateService(dr *repository.DonateRepository,
	ir *repository.ItemRepository) *DonateService {
	return &DonateService{
		dr: dr,
		ir: ir,
	}
}

func (ds *DonateService) ProcessDonate(d entities.Donateable, promo string) error {
	donate := d.ToDonate()

	if ok, err := ds.dr.PaymentExist(donate.ID); err != nil || ok {
		if ok {
			return fmt.Errorf("payment with id '%d' already exist", donate.ID)
		} else {
			return err
		}
	}

	multiplier := 1.0
	if promo != "" {
		fullPromo, err := ds.ir.GetPromo(promo)
		if err != nil {
			return err
		}
		multiplier = fullPromo.Multiplier
	}

	item, err := ds.ir.GetItem(donate.DonateItemID)
	if err != nil {
		return err
	}

	return nil
}

func calculatePrice(price, amount int, multiplier float64) int {
	if amount > 1 {
		return int(multiplier * float64(amount) * math.Round(float64(price) * (((100.0 - float64(calculateSale(amount))) / 100.0))))
	} else {
		return int(float64(price) * multiplier)
	}
}

func calculateSale(amount int) int {
	return int(math.Round(50 / (math.Pow(math.E, 3 - (float64(amount) / math.Pow(math.Pi, 2))) + 1)))
}


//     if (await db.checkIfPaymentExists(info.intid)) return res.status(400).json({ error: 'Payment already exists.' })

//     let promo = null
//     if (req.body.us_promo)
//     	promo = await db.getPromo(req.body.us_promo)

//     try {
//         const item = await db.getItemById(Number(info.us_donate))
// 		if (Math.round((promo ? promo.multiplier : 1) * calculatePrice(item.price, info.us_number)) !== Number(info.AMOUNT)) {
// 			console.log('Invalid amount: ' + info.AMOUNT + ' Real price: ' + calculatePrice(item.price, info.us_number))
// 			return res.status(400).json({ error: 'Invalid amount' })
// 		}

// 		if (item.command) {
// 	    	var command = item.command.replaceAll('%user%', info.us_username)
//             if (info.us_number)
//     		   command = command.replaceAll('%number%', info.us_number)

//             console.log(`sending "${command}" to server`)
// 		    await rcon.connect()
//             await rcon.send(command)
//             rcon.end()
// 		} else {
// 		    console.log('No command for item')
// 		}

//         await db.addDonateInfo(Number(info.us_donate), info.us_username, Number(info.us_number), Date.now(), info.MERCHANT_ORDER_ID, Number(info.AMOUNT))

//         res.json({ok: true})

//     } catch(error) {
//         console.log(error);
//         res.status(400).json({ error: error.toString() })
//     }
