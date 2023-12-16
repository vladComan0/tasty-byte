package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodGet, "/ping", http.HandlerFunc(app.ping))

	// CRUD
	router.Handler(http.MethodPost, "/v1/recipes", http.HandlerFunc(app.createRecipe))
	router.Handler(http.MethodGet, "/v1/recipes/:id", http.HandlerFunc(app.getRecipe))

	return router
}
