package main

import "net/http"

// getTransactions godoc
//
//	@Summary		Get transactions
//	@Description	Get all transactions
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]domain.Transaction
//	@Failure		400	{object}	main.ErrorResponse
//	@Failure		500	{object}	main.ErrorResponse
//	@Router			/api/v1/transactions [get]
func (app *application) getTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	transactions, err := app.storage.Transactions.GetAll(r.Context(), nil, nil)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, transactions); err != nil {
		app.internalServerError(w, r, err)
	}
}
