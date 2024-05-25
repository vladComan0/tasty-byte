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
	ts := httptest.NewTLSServer(h)

	// Disable redirect-following for the test server client by setting a custom
	// CheckRedirect function. This function will be called whenever a 3xx response
	// status is received by the client and, by always returning a http.ErrUseLastResponse
	// error it forces the client to immediately return the received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}
