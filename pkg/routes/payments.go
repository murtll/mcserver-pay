package routes

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/murtll/mcserver-pay/pkg/config"
	"github.com/murtll/mcserver-pay/pkg/entities"
	"github.com/murtll/mcserver-pay/pkg/service"
)

var trustedIps = map[string][]string{
	"fk": config.FkTrustedIps,
}

type PaymentRouter struct {
	chi.Mux
	ds *service.DonateService
}

func NewPaymentRouter(ds *service.DonateService) *PaymentRouter {
	p := &PaymentRouter{
		Mux: *chi.NewRouter(),
		ds:  ds,
	}

	p.Use(p.checkIpMiddleware)

	// process freekassa payment
	p.Post("/fk", p.processFkPayment)

	return p
}

func (p *PaymentRouter) checkIpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pieces := strings.Split(r.URL.Path, "/")
		paymentProvider := pieces[len(pieces)-1]

		realIp := r.Header.Get("x-forwarded-for")

		if ips, ok := trustedIps[paymentProvider]; ok {
			for _, ip := range ips {
				if ip == realIp {
					next.ServeHTTP(w, r)
					return
				}
			}
		} else {
			next.ServeHTTP(w, r)
			return
		}

		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, entities.ErrorResponse{
			Error: "invalid ip",
		})
	})
}

func (p *PaymentRouter) processFkPayment(w http.ResponseWriter, r *http.Request) {
	request := &entities.FkPaymentRequest{}
	if err := render.DefaultDecoder(r, request); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, entities.NewErrorResponse(err))
		return
	}

	p.ds.ProcessDonate(request, request.Promo)
}

// router.post('/process-payment-fk', async (req, res) => {
//     const trustedIps = ['168.119.157.136', '168.119.60.227', '138.201.88.124', '178.154.197.79']

//     const ip = req.headers['x-forwarded-for'] || req.socket.remoteAddress

//     console.log('from: ' + ip)

//     if (!trustedIps.includes(ip)) {
//         return res.status(400).json({ error: 'Bad IP' })
//     }

//     const info = req.body
//     console.log(info)

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
// })
