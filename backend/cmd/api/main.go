package main

import (
	_ "github.com/saleemlawal/lumen/docs"
	"github.com/saleemlawal/lumen/internal/env"

	"go.uber.org/zap"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	appEnv := env.GetEnvString("APP_ENV", "development")
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
			addr:        env.GetEnvString("PORT", ":8080"),
			env:         appEnv,
			frontendUrl: env.GetEnvString("FRONTEND_URL", "http://localhost:5173"),
		},
	}

	chiRouter := app.mount()

	err := app.run(chiRouter)

	logger.Fatalw("Server has stopped", "error", err)
}
