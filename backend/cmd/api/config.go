package main

import (
	"github.com/saleemlawal/lumen/internal/env"
	"go.uber.org/zap"
)

func loadConfig() config {
	envName := env.GetEnvString("APP_ENV", "development")
	return config{
		env:         envName,
		addr:        env.GetEnvString("PORT", ":8080"),
		frontendUrl: env.GetEnvString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:        env.GetEnvString("DB_URL", ""),
			maxOpen:     env.GetEnvInt("DB_MAX_OPEN", 10),
			maxIdle:     env.GetEnvInt("DB_MAX_IDLE", 10),
			maxIdleTime: env.GetEnvString("DB_MAX_IDLE_TIME", "10m"),
		},
		plaid: plaidConfig{
			plaidClientId: env.GetEnvString("PLAID_CLIENT_ID", ""),
			plaidSecret:   env.GetEnvString("PLAID_SECRET", ""),
			plaidEnv:      env.GetEnvString("PLAID_ENV", "sandbox"),
		},
	}
}

func newLogger(envName string) *zap.SugaredLogger {
	if envName == "production" {
		return zap.Must(zap.NewProduction()).Sugar()
	}
	return zap.Must(zap.NewDevelopment()).Sugar()
}
