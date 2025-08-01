package data

import (
	"database/sql"
	"errors"
	"time"

	pg "github.com/lib/pq"
	"greenlight.jayant.com/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitzero"`
	Runtime   Runtime   `json:"runtime,omitzero"` // string directive would convert value to string
	Genres    []string  `json:"genres"`
	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	query := `
        INSERT INTO movies (title, year, runtime, genres)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pg.Array(movie.Genres)}

	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, title, year, runtime, genres, version, created_at
        FROM movies
        WHERE id = $1`

	var movie Movie

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pg.Array(&movie.Genres),
		&movie.Version,
		&movie.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {

    query := `
        UPDATE movies
        SET title = $1, year = $2, runtime = $3, genres = $4, version = version+1
        WHERE id = $5
        RETURNING version`

    args := []any {
        movie.Title,
        movie.Year, 
        movie.Runtime, 
        pg.Array(movie.Genres), 
        movie.ID,
    }

    return m.DB.QueryRow(query, args...).Scan(&movie.Version)
}

func (m MovieModel) Delete(id int64) error {

    if id < 1 {
        return ErrRecordNotFound
    }

    query := `
        DELETE FROM movies
        WHERE id = $1`

    result, err := m.DB.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return ErrRecordNotFound
    }
    
	return nil
}

func (m MovieModel) GetAll() ([]Movie, error) {

    query := `
        SELECT id, title, year, runtime, genres, version, created_at
        FROM movies`

    rows, err := m.DB.Query(query)
    if err != nil {
        return nil, err
    }

    defer rows.Close()
    var movies []Movie
    for rows.Next() {
        var movie Movie
        err = rows.Scan(
            &movie.ID, 
            &movie.Title, 
            &movie.Year, 
            &movie.Runtime, 
            pg.Array(&movie.Genres), 
            &movie.Version, 
            &movie.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        movies = append(movies, movie)
    }
    
    if err = rows.Err(); err != nil {
        return movies, err
    }

    return movies, nil
}