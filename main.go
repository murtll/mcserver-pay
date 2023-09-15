package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/murtll/mcserver-pay/pkg/config"
	"github.com/murtll/mcserver-pay/pkg/entities"
	"github.com/murtll/mcserver-pay/pkg/repository"
	"github.com/murtll/mcserver-pay/pkg/routes"
	"github.com/murtll/mcserver-pay/pkg/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var migrate = flag.Bool("migrate", false, "Run migrations.")

func main() {
	// init db
	db, err := gorm.Open(postgres.Open(config.PostgresString))
	if err != nil {
		panic(err)
	}

	flag.Parse()

	// init repositories
	dr := repository.NewDonateRepository(db)
	if *migrate {
		err := dr.Migrate()
		if err != nil {
			panic(err)
		}
	}

	ir, err := repository.NewItemRepository(config.ApiUrl)
	if err != nil {
		panic(err)
	}

	// init services
	ds := service.NewDonateService(dr, ir)

	// init main chi router
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Get(config.HealthPath, func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, entities.StatusResponse{
			Status:  "ok",
			Version: config.Version,
		})
	})

	// mount routers
	router.Mount("/process", routes.NewPaymentRouter(ds))

	log.Default().Printf("Starting server v%s on %s", config.Version, config.ListenAddr)

	err = http.ListenAndServe(config.ListenAddr, router)
	if err != nil {
		panic(err)
	}
}
