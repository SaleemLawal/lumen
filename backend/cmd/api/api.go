package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/saleemlawal/lumen/internal/plaid"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

const version = "1.0.0"

type application struct {
	logger      *zap.SugaredLogger
	config      config
	plaidClient *plaid.PlaidClient
}

type config struct {
	addr        string
	env         string
	frontendUrl string
}

type plaidConfig struct {
	plaidClientId string
	plaidSecret   string
	plaidEnv      string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 60))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json", "text/xml"))
		r.Get("/health", app.healthcheckHandler)

		r.Get("/plaid/link-token", app.createPlaidLinkTokenHandler)
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 10,
	}

	app.logger.Infow("Server has started", "addr", app.config.addr)

	return srv.ListenAndServe()
}
