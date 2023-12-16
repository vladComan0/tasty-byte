package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/vladComan0/tasty-byte/internal/models"
)

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	if err := app.recipes.DB.Ping(); err != nil {
		app.errorLog.Printf("Unable to establish connection with database: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("pong"))
}

func (app *application) createRecipe(w http.ResponseWriter, r *http.Request) {
	// mock the recipe
	recipe := &models.Recipe{
		Name:            "Pancakes",
		Description:     "Delicious pancakes",
		Instructions:    "Mix the ingredients and fry them",
		PreparationTime: "10 minutes",
		CookingTime:     "10 minutes",
		Portions:        "2",
	}

	id, err := app.recipes.Insert(
		recipe.Name,
		recipe.Description,
		recipe.Instructions,
		recipe.PreparationTime,
		recipe.CookingTime,
		recipe.Portions,
	)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// send the recipe as JSON
	fmt.Fprintf(w, "The recipe with id %d was created successfully!", id)
}

func (app *application) getRecipe(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	recipe, err := app.recipes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoRecord):
			app.clientError(w, http.StatusNotFound)
		default:
			app.serverError(w, err)
		}
		return
	}

	// send the recipe as JSON
	_ = recipe
	fmt.Fprint(w, "Nailed it!")
}
