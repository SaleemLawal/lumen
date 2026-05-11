package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type application struct {
	logger *zap.SugaredLogger
	config config
}

type config struct {
	addr        string
	env         string
	frontendUrl string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json", "text/xml"))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 60))

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
