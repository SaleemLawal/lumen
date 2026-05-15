package main

import "net/http"

// createPlaidLinkTokenHandler godoc
//
//	@Summary		Create Plaid link token
//	@Description	Create a Plaid link token
//	@Tags			plaid
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"Plaid link token"
//	@Router			/api/v1/plaid/link-token [get]
func (app *application) createPlaidLinkTokenHandler(w http.ResponseWriter, r *http.Request) {
	linkToken, err := app.plaidClient.CreateLinkToken()
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, linkToken)
}
