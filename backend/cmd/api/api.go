package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/saleemlawal/lumen/external/plaid"
	"github.com/saleemlawal/lumen/internal/db"
	"github.com/saleemlawal/lumen/internal/encryption"
	"github.com/saleemlawal/lumen/internal/store"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

const version = "1.0.0"

type application struct {
	logger      *zap.SugaredLogger
	config      config
	plaidClient *plaid.PlaidClient
	storage     *store.Storage
	encryptor   *encryption.AESEncryptor
}

type config struct {
	addr        string
	env         string
	frontendUrl string
	db          dbConfig
	plaid       plaidConfig
}

type plaidConfig struct {
	plaidClientId      string
	plaidSecret        string
	plaidEnv           string
	tokenEncryptionKey string
}

type dbConfig struct {
	addr        string
	maxOpen     int
	maxIdle     int
	maxIdleTime string
}

func newApplication(logger *zap.SugaredLogger, cfg config) (*application, func(), error) {
	sqlDB, err := db.New(cfg.db.addr,
		cfg.db.maxOpen,
		cfg.db.maxIdle,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { _ = sqlDB.Close() }

	keyBytes, err := base64.StdEncoding.DecodeString(cfg.plaid.tokenEncryptionKey)
	if err != nil {
		return nil, cleanup, fmt.Errorf("invalid PLAID_TOKEN_ENCRYPTION_KEY (must be base64): %w", err)
	}

	enc, err := encryption.NewAESEncryptor(keyBytes)
	if err != nil {
		return nil, cleanup, fmt.Errorf("failed to create encryptor: %w", err)
	}

	app := &application{
		logger: logger,
		config: cfg,
		plaidClient: plaid.NewPlaidClient(cfg.plaid.plaidClientId,
			cfg.plaid.plaidSecret,
			cfg.plaid.plaidEnv,
		),
		storage:   store.NewStorage(sqlDB),
		encryptor: enc,
	}
	return app, cleanup, nil
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

		r.Route("/plaid", func(r chi.Router) {
			r.Get("/link-token", app.createPlaidLinkTokenHandler)
			r.Post("/exchange-public-token", app.exchangePublicTokenHandler)
			r.Get("/items", app.getPlaidItemsHandler)
			r.Get("/items/{id}/link-token", app.getUpdateLinkTokenHandler)
			r.Post("/items/{id}/sync-accounts", app.syncItemAccountsHandler)
			r.Post("/items/{id}/sync-transactions", app.syncItemTransactionsHandler)
		})

		r.Route("/transactions", func(r chi.Router) {
			r.Get("/", app.getTransactionsHandler)
		})

		r.Route("/accounts", func(r chi.Router) {
			r.Get("/", app.getAccountsHandler)
		})
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
