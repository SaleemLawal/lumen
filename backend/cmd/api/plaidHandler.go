package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/plaid/plaid-go/v42/plaid"
	externalplaid "github.com/saleemlawal/lumen/external/plaid"
	"github.com/saleemlawal/lumen/internal/domain"
)

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
//	@Failure		415	{object}	ErrorResponse	"Unsupported Content-Type (allowed: application/json, text/xml)"
//	@Failure		500	{object}	ErrorResponse	"Plaid or server error"
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
//	@Description	Exchanges a Link public_token server-side, stores the encrypted access token and item, then fetches and persists the user's accounts from Plaid.
//	@Tags			plaid
//	@Accept			json
//	@Produce		json
//	@Param			body	body		ExchangePublicTokenRequest	true	"public_token from Link onSuccess"
//	@Success		200		{object}	PlaidExchangePublicTokenResponse
//	@Failure		400		{object}	ErrorResponse	"Invalid JSON body or unknown fields"
//	@Failure		415		{object}	ErrorResponse	"Unsupported Content-Type (allowed: application/json, text/xml)"
//	@Failure		500		{object}	ErrorResponse	"Plaid or server error"
//	@Router			/api/v1/plaid/exchange-public-token [post]
func (app *application) exchangePublicTokenHandler(w http.ResponseWriter, r *http.Request) {
	var request ExchangePublicTokenRequest

	if err := readJSON(w, r, &request); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	response, err := app.plaidClient.ExchangePublicToken(request.PublicToken)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	encryptedToken, err := app.encryptor.Encrypt(response.AccessToken)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	var accounts []plaid.AccountBase
	var syncResult *externalplaid.SyncTransactionsResult
	var institutionID string

	g, _ := errgroup.WithContext(r.Context())
	g.Go(func() error {
		var err error
		accounts, err = app.plaidClient.FetchAccounts(response.AccessToken)
		return err
	})
	g.Go(func() error {
		var err error
		syncResult, err = app.plaidClient.SyncTransactions(&response.AccessToken, nil)
		return err
	})
	g.Go(func() error {
		var err error
		institutionID, err = app.plaidClient.FetchInstitutionID(response.AccessToken)
		return err
	})
	if err := g.Wait(); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	linked, err := app.storage.Plaid.InstitutionLinked(r.Context(), institutionID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if linked {
		app.conflictError(w, r, "this institution is already connected")
		return
	}

	err = app.storage.StoreLinkSync(r.Context(),
		&domain.PlaidItem{AccessToken: encryptedToken, ItemID: response.ItemID, InstitutionID: institutionID},
		accounts,
		syncResult.Added,
		syncResult.Removed,
		*syncResult.NextCursor,
	)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, PlaidExchangePublicTokenResponse{Status: "linked"})
}

// getPlaidItemsHandler godoc
//
//	@Summary		List linked institutions
//	@Description	Returns all linked Plaid items with their synced accounts.
//	@Tags			plaid
//	@Produce		json
//	@Success		200	{array}		domain.PlaidItemSummary
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/plaid/items [get]
func (app *application) getPlaidItemsHandler(w http.ResponseWriter, r *http.Request) {
	items, err := app.storage.Plaid.GetAllItems(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// getUpdateLinkTokenHandler godoc
//
//	@Summary		Create update-mode link token
//	@Description	Returns a Plaid link token configured for account selection update mode for an existing item.
//	@Tags			plaid
//	@Produce		json
//	@Param			id	path		string	true	"Plaid item UUID"
//	@Success		200	{string}	string	"link token"
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/plaid/items/{id}/link-token [get]
func (app *application) getUpdateLinkTokenHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	item, err := app.storage.Plaid.GetItemByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, http.StatusNotFound, "item not found")
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	accessToken, err := app.encryptor.Decrypt(item.AccessToken)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	linkToken, err := app.plaidClient.CreateUpdateLinkToken(accessToken)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, linkToken)
}

// syncItemAccountsHandler godoc
//
//	@Summary		Sync accounts for an item
//	@Description	Re-fetches accounts from Plaid for an existing item and upserts them (soft sync — adds new, updates existing, never deletes).
//	@Tags			plaid
//	@Produce		json
//	@Param			id	path		string	true	"Plaid item UUID"
//	@Success		200	{object}	PlaidExchangePublicTokenResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/plaid/items/{id}/sync-accounts [post]
func (app *application) syncItemAccountsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	item, err := app.storage.Plaid.GetItemByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, http.StatusNotFound, "item not found")
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	accessToken, err := app.encryptor.Decrypt(item.AccessToken)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	accounts, err := app.plaidClient.FetchAccounts(accessToken)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.storage.SyncItemAccounts(r.Context(), item.ItemID, accounts); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, PlaidExchangePublicTokenResponse{Status: "synced"})
}
