package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/markponce/greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create movie here...")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablance",
		Year:      2001,
		RunTime:   102,
		Genres:    []string{},
		Vesion:    1,
	}

	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "he server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
