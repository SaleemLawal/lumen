package plaid

import (
	"context"
	"time"

	plaid "github.com/plaid/plaid-go/v42/plaid"
)

type PlaidClient struct {
	plaidApi *plaid.PlaidApiService
}

type PublicTokenExchangeResult struct {
	AccessToken string
	ItemID      string
}

type SyncTransactionsResult struct {
	Added      []plaid.Transaction
	Modified   []plaid.Transaction
	Removed    []plaid.RemovedTransaction
	NextCursor *string
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

	return response.GetLinkToken(), nil
}

func (c *PlaidClient) ExchangePublicToken(publicToken string) (PublicTokenExchangeResult, error) {
	ctx := context.Background()

	exchangePublicTokenReq := plaid.NewItemPublicTokenExchangeRequest(publicToken)

	response, _, err := c.plaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(*exchangePublicTokenReq).Execute()
	if err != nil {
		return PublicTokenExchangeResult{}, err
	}

	return PublicTokenExchangeResult{
		AccessToken: response.GetAccessToken(),
		ItemID:      response.GetItemId(),
	}, nil
}

func (c *PlaidClient) FetchAccounts(accessToken string) ([]plaid.AccountBase, error) {
	ctx := context.Background()

	response, _, err := c.plaidApi.AccountsGet(ctx).AccountsGetRequest(*plaid.NewAccountsGetRequest(accessToken)).Execute()
	if err != nil {
		return nil, err
	}

	return response.GetAccounts(), nil
}

func (c *PlaidClient) SyncTransactions(accessToken, cursor *string) (*SyncTransactionsResult, error) {
	ctx := context.Background()

	var added []plaid.Transaction
	var modified []plaid.Transaction
	var removed []plaid.RemovedTransaction

	hasMore := true

	for hasMore {
		request := plaid.NewTransactionsSyncRequest(*accessToken)

		if cursor != nil {
			request.SetCursor(*cursor)
		}

		response, _, err := c.plaidApi.TransactionsSync(ctx).TransactionsSyncRequest(*request).Execute()
		if err != nil {
			return nil, err
		}

		nextCursor := response.GetNextCursor()
		cursor = &nextCursor

		if *cursor == "" {
			time.Sleep(1 * time.Second)
			continue
		}

		added = append(added, response.GetAdded()...)
		modified = append(modified, response.GetModified()...)
		removed = append(removed, response.GetRemoved()...)
	}

	return &SyncTransactionsResult{
		Added:      added,
		Modified:   modified,
		Removed:    removed,
		NextCursor: cursor,
	}, nil
}
