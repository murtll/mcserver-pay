package routes

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/murtll/mcserver-pay/pkg/entities"
	"github.com/murtll/mcserver-pay/pkg/service"
)

type DonateRouter struct {
	chi.Mux
	ds *service.DonateService
}

func NewDonateRouter(ds *service.DonateService) *DonateRouter {
	p := &DonateRouter{
		Mux: *chi.NewRouter(),
		ds:  ds,
	}

	// process freekassa payment
	p.Get("/last", p.getLastDonates)

	return p
}

func (p *DonateRouter) getLastDonates(w http.ResponseWriter, r *http.Request) {

	count := 10
	countStr := r.URL.Query().Get("count")
	if countStr != "" {
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, entities.NewErrorResponse(err))	
		}
	}

	donates, err := p.ds.GetLastDonates(count)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, entities.NewErrorResponse(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, donates)
}
