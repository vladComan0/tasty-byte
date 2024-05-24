package main

import (
	"github.com/vladComan0/tasty-byte/internal/mocks"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

type testServer struct {
	*httptest.Server
}

func newTestApplication() *application {
	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
		recipes: &mocks.MockRecipeModelInterface{
			IngredientModel:       &mocks.MockIngredientModelInterface{},
			TagModel:              &mocks.MockTagModelInterface{},
			RecipeIngredientModel: &mocks.MockRecipeIngredientModelInterface{},
			RecipeTagModel:        &mocks.MockRecipeTagModelInterface{},
		},
	}
}

func newTestServer(h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}
