package plaid

import (
	plaid "github.com/plaid/plaid-go/v42/plaid"
)

type PlaidClient struct {
	api *plaid.APIClient
}

var environments = map[string]plaid.Environment{
	"sandbox":    plaid.Sandbox,
	"production": plaid.Production,
}

func NewPlaidClient(clientId, secret, env string) *PlaidClient {
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", clientId)
	configuration.AddDefaultHeader("PLAID-SECRET", secret)
	configuration.UseEnvironment(environments[env])
	return &PlaidClient{
		api: plaid.NewAPIClient(configuration),
	}
}
