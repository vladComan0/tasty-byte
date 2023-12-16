package main

import "net/http"

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	if err := app.recipes.DB.Ping(); err != nil {
		app.errorLog.Printf("Unable to establish connection with database: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("pong"))
}
