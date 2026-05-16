package main

import "net/http"

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusInternalServerError, "The server encountered a problem and could not process your request")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, msg string) {
	app.logger.Warnw("Conflict", "method", r.Method, "path", r.URL.Path, "reason", msg)
	writeJSONError(w, http.StatusConflict, msg)
}
