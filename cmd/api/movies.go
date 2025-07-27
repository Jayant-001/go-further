package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.jayant.com/internal/data"
	"greenlight.jayant.com/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title   string       `'json:"title"`
		Year    int32        `'json:"year"`
		Runtime data.Runtime `'json:"runtime"`
		Genres  []string     `'json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the input struct to a new Movie struct.
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	// Use the Check() method to execute our validation checks. This will add the
	// provided key and error message to the errors map if the check does not evaluate
	// to true. For example, in the first line here we "check that the title is not
	// equal to the empty string". In the second, we "check that the length of the title
	// is less than or equal to 500 bytes" and so on.
	// v.Check(input.Title != "", "title", "must be provided")
	// v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")

	// v.Check(input.Year != 0, "year", "must be provided")
	// v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	// v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	// v.Check(input.Runtime != 0, "runtime", "must be provided")
	// v.Check(input.Runtime > 0, "runtime", "must be a positive integer")

	// v.Check(input.Genres != nil, "genres", "must be provided")
	// v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	// v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
	// // Note that we're using the Unique helper in the line below to check that all
	// // values in the input.Genres slice are unique.
	// v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")

	// // Use the Valid() method to see if any of the checks failed. If they did, then use
	// // the failedValidationResponse() helper to send a response to the client, passing
	// // in the v.Errors map.
	// if !v.Valid() {
	//     app.failedValidationResponse(w, r, v.Errors)
	//     return
	// }

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJson(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getMoviesHandler(w http.ResponseWriter, r *http.Request) {

	movie := data.Movie{
		ID:        10,
		Title:     "Avengers",
		Year:      2017,
		Runtime:   120,
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
		Runtime:   120,
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
