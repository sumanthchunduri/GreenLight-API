package main

import (
	"GreenLight/internal/data"
	"GreenLight/internal/validator"
	"fmt"
	"net/http"
	"time"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(input.Runtime != 0, "runtime", "must be provided")
	v.Check(input.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "details of movie %d\n", id)

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "RRR",
		Runtime:   102,
		Genres:    []string{"Drama", "History"},
		Version:   1,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
		//http.Error(w, "the server encountered a problem", http.StatusInternalServerError)
	}
}
