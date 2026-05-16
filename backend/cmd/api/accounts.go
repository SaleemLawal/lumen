package main

import (
	"net/http"

	"github.com/google/uuid"
)

// getAccountsHandler godoc
//
//	@Summary		Get accounts
//	@Description	Get accounts for a given item ID
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			item_id	query		string	false	"Item ID"	example:"7aa81d92-1fe6-419c-9253-23929031497c"
//	@Success		200		{array}	domain.Account
//	@Failure		400		{object}	ErrorResponse	"Invalid item ID"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/api/v1/accounts [get]
func (app *application) getAccountsHandler(w http.ResponseWriter, r *http.Request) {
	var itemID *uuid.UUID
	if raw := r.URL.Query().Get("item_id"); raw != "" {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}
		itemID = &parsed
	}

	accounts, err := app.storage.Accounts.GetAll(r.Context(), itemID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, accounts)
}
