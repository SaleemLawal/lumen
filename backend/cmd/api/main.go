package main

import (
	"github.com/saleemlawal/lumen/internal"
	"go.uber.org/zap"
)

func main() {
	appEnv := internal.GetEnvString("APP_ENV", "development")
	var logger *zap.SugaredLogger

	if appEnv == "production" {
		logger = zap.Must(zap.NewProduction()).Sugar()
	} else {
		logger = zap.Must(zap.NewDevelopment()).Sugar()
	}

	defer logger.Sync()

	app := &application{
		logger: logger,
		config: config{
			addr:        internal.GetEnvString("PORT", ":8080"),
			env:         appEnv,
			frontendUrl: internal.GetEnvString("FRONTEND_URL", "http://localhost:5173"),
		},
	}

	chiRouter := app.mount()

	err := app.run(chiRouter)

	logger.Fatalw("Server has stopped", "error", err)
}
