package main

import (
	_ "github.com/saleemlawal/lumen/docs"
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
	cfg := loadConfig()
	logger := newLogger(cfg.env)
	defer func() { _ = logger.Sync() }()

	app, cleanup, err := newApplication(logger, cfg)
	if err != nil {
		logger.Fatalw("failed to connect to database", "error", err)
	}
	defer cleanup()

	if err := app.run(app.mount()); err != nil {
		logger.Fatalw("server stopped", "error", err)
	}
}
