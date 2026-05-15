package plaid

import (
	"context"
	"fmt"
	"time"

	plaid "github.com/plaid/plaid-go/v42/plaid"
)

type PlaidClient struct {
	plaidApi *plaid.PlaidApiService
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
		plaidApi: plaid.NewAPIClient(configuration).PlaidApi,
	}
}

func (c *PlaidClient) CreateLinkToken() (string, error) {
	ctx := context.Background()

	// TODO: Client User ID
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: time.Now().String(),
	}

	request := plaid.NewLinkTokenCreateRequest(
		"Plaid Quickstart",
		"en",
		[]plaid.CountryCode{plaid.COUNTRYCODE_US},
	)
	request.SetUser(user)
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_TRANSACTIONS})

	response, _, err := c.plaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", err
	}

	fmt.Println(response)

	return response.GetLinkToken(), nil

}
