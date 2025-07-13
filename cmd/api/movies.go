package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.jayant.com/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title   string       `'json:"title"`
		Year    int32        `'json:"year"`
		Runtime data.Runtime `'json:"runtime"`
		Genres  []string     `'json:"genres"`
	}

	// err := json.NewDecoder(r.Body).Decode(&input)
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Fprintln(w, "Create a new movie", input)
}

func (app *application) getMoviesHandler(w http.ResponseWriter, r *http.Request) {

	movie := data.Movie{
		ID:        10,
		Title:     "Avengers",
		Year:      2017,
		RunTime:   120,
		Genres:    []string{"action", "adventure"},
		Version:   2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := app.writeJson(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		Title:     "Avengers",
		Year:      2017,
		RunTime:   120,
		Genres:    []string{"action", "adventure"},
		Version:   2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = app.writeJson(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
