package main

import "net/http"

type ExchangePublicTokenRequest struct {
	PublicToken string `json:"public_token" example:"public-sandbox-abc123"`
}

// PlaidExchangePublicTokenResponse is returned after a successful exchange; access_token and item_id stay server-side only.
type PlaidExchangePublicTokenResponse struct {
	Status string `json:"status" example:"linked"`
}

// createPlaidLinkTokenHandler godoc
//
//	@Summary		Create Plaid link token
//	@Description	Create a Plaid link token
//	@Tags			plaid
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"plaid-sandbox-1234567890"
//	@Failure		415	{object}	github.com/saleemlawal/lumen/internal/apiresponse.ErrorResponse	"Unsupported Content-Type (allowed: application/json, text/xml)"
//	@Failure		500	{object}	github.com/saleemlawal/lumen/internal/apiresponse.ErrorResponse	"Plaid or server error"
//	@Router			/api/v1/plaid/link-token [get]
func (app *application) createPlaidLinkTokenHandler(w http.ResponseWriter, r *http.Request) {
	linkToken, err := app.plaidClient.CreateLinkToken()
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, linkToken)
}

// exchangePublicTokenHandler godoc
//
//	@Summary		Exchange Plaid public token
//	@Description	Exchanges a Link public_token server-side.
//	@Tags			plaid
//	@Accept			json
//	@Produce		json
//	@Param			body	body		ExchangePublicTokenRequest	true	"public_token from Link onSuccess"
//	@Success		200		{object}	PlaidExchangePublicTokenResponse
//	@Failure		400	{object}	github.com/saleemlawal/lumen/internal/apiresponse.ErrorResponse	"Invalid JSON body or unknown fields"
//	@Failure		415	{object}	github.com/saleemlawal/lumen/internal/apiresponse.ErrorResponse	"Unsupported Content-Type (allowed: application/json, text/xml)"
//	@Failure		500	{object}	github.com/saleemlawal/lumen/internal/apiresponse.ErrorResponse	"Plaid or server error"
//	@Router			/api/v1/plaid/exchange-public-token [post]
func (app *application) exchangePublicTokenHandler(w http.ResponseWriter, r *http.Request) {
	var request ExchangePublicTokenRequest

	if err := readJSON(w, r, &request); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// TODO: Persist result.AccessToken (encrypted) and result.ItemID for the authenticated user.

	_, err := app.plaidClient.ExchangePublicToken(request.PublicToken)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, PlaidExchangePublicTokenResponse{Status: "linked"})
}
