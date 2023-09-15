package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/murtll/mcserver-pay/pkg/config"
	"github.com/murtll/mcserver-pay/pkg/crypto"
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

	log.Default().Print(*request)

	err := p.ds.ProcessDonate(request,
		request.Promo,
		crypto.CheckSignFk,
		request.Sign,
		config.FkMerchantID,
		request.Amount,
		config.FkSigningKey,
		request.MerchantOrderID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, entities.NewErrorResponse(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, entities.NewOkResponse())
}

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
