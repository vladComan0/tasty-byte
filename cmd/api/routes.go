package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodGet, "/ping", http.HandlerFunc(app.ping))

	// CRUD
	router.Handler(http.MethodPost, "/v1/recipes", http.HandlerFunc(app.createRecipe))
	router.Handler(http.MethodGet, "/v1/recipes/:id", http.HandlerFunc(app.getRecipe))
	router.Handler(http.MethodPut, "/v1/recipes/:id", http.HandlerFunc(app.updateRecipe))
	router.Handler(http.MethodDelete, "/v1/recipes/:id", http.HandlerFunc(app.deleteRecipe))

	standardChain := alice.New(app.recoverPanic, app.logRequests, app.enableCORS)

	return standardChain.Then(router)
}
