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
	var input struct {
		Name            string `json:"name"`
		Description     string `json:"description"`
		Instructions    string `json:"instructions"`
		PreparationTime string `json:"preparation_time"`
		CookingTime     string `json:"cooking_time"`
		Portions        int    `json:"portions"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	recipe := &models.Recipe{
		Name:            input.Name,
		Description:     input.Description,
		Instructions:    input.Instructions,
		PreparationTime: input.PreparationTime,
		CookingTime:     input.CookingTime,
		Portions:        input.Portions,
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

	// Make the application aware of that new location -> add the headers to the right json helper function
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("v1/recipes/%d", id))

	if err = app.writeJSON(w, http.StatusCreated, envelope{"recipe": recipe}, headers); err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("Created new recipe with id: %d", id)
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

	if err = app.writeJSON(w, http.StatusOK, envelope{"recipe": recipe}, nil); err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("Retrieved recipe with id: %d", id)
}

func (app *application) updateRecipe(w http.ResponseWriter, r *http.Request) {
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

	var input struct {
		Name            *string `json:"name"`
		Description     *string `json:"description"`
		Instructions    *string `json:"instructions"`
		PreparationTime *string `json:"preparation_time"`
		CookingTime     *string `json:"cooking_time"`
		Portions        *int    `json:"portions"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if input.Name != nil {
		recipe.Name = *input.Name
	}

	if input.Description != nil {
		recipe.Description = *input.Description
	}

	if input.Instructions != nil {
		recipe.Instructions = *input.Instructions
	}

	if input.PreparationTime != nil {
		recipe.PreparationTime = *input.PreparationTime
	}

	if input.CookingTime != nil {
		recipe.CookingTime = *input.CookingTime
	}

	if input.Portions != nil {
		recipe.Portions = *input.Portions
	}

	if err := app.recipes.Update(recipe); err != nil {
		app.serverError(w, err)
		return
	}

	if err = app.writeJSON(w, http.StatusOK, envelope{"recipe": recipe}, nil); err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("Updated recipe with id: %d", id)
}

func (app *application) deleteRecipe(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if err := app.recipes.Delete(id); err != nil {
		switch {
		case errors.Is(err, models.ErrNoRecord):
			app.clientError(w, http.StatusNotFound)
		default:
			app.serverError(w, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "Recipe successfully deleted"}, nil); err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("Deleted recipe with id: %d", id)
}
