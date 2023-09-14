package service

import (
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
	if promo != "" {
		ds.ir.GetPromo(promo)
	}

	return nil
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
