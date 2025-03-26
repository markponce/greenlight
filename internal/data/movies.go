package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/markponce/greenlight/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitzero"`
	Runtime   Runtime   `json:"runtime,omitzero"`
	Genres    []string  `json:"genres,omitempty"`
	Vesion    int32     `json:"version"`
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

func (m *MovieModel) List(*[]Movie) error {
	return nil
}

func (m *MovieModel) Insert(movie *Movie) error {
	// query statement
	query := `

		INSERT INTO movies(title, year, runtime, genres)
		VALUES ($1,$2,$3,$4)
		RETURNING id, created_at, version
	`
	// create slice type any for the arguments
	args := []any{
		movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres),
	}

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute statement in the db. convert args using variadics and reference to update the movie id, createdAt, version
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Vesion)
}

func (m *MovieModel) Get(id int64) (*Movie, error) {
	// check if less than 1 return ErrRecordNotfound
	if id < 1 {
		return nil, ErrRecordNotfound
	}

	// query statement
	stmt := `
		SELECT id, created_at, title, year, runtime, genres, version
		from movies 
		where id=$1
	`

	// movie variable
	var movie Movie

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute query to the db.
	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Vesion,
	)

	// check for result
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			// if error is no rows return standard error ErrRecordNotfound
			return nil, ErrRecordNotfound
		default:
			// return whatever error
			return nil, err
		}
	}

	return &movie, nil
}

func (m *MovieModel) Update(movie *Movie) error {
	// update query statement
	// avoid race condition where version
	stmt := `
		UPDATE movies
		SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		where id = $5 and version = $6
		returning version
		`
	// create slice any for the arguments
	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Vesion,
	}

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&movie.Vesion)
	if err != nil {
		switch {
		// if no updated record it means that the record has been updated already (data race condition)
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}

	}

	return nil
}

func (m *MovieModel) Delete(id int64) error {
	// if less than 1 return standard app error ErrRecordNotfound
	if id < 1 {
		return ErrRecordNotfound
	}

	// query statement
	query := `
		DELETE FROM movies
		where id = $1
	`

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	//check for row result
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// if there is no result, return standard app error ErrRecordNotfound
	if rowsAffected == 0 {
		return ErrRecordNotfound
	}

	// deletion is successfull
	return nil
}

func (m *MovieModel) GetAll(title string, genres []string, filter Filters) ([]*Movie, error) {
	// Construct the SQL query to retrieve all movie records.
	query := `

		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		ORDER by id
	`

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()

	// Initialize an empty slice to hold the movie data.
	movies := []*Movie{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		var movie Movie

		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Vesion,
		)

		if err != nil {
			return nil, err
		}

		// Add the Movie struct to the slice.
		movies = append(movies, &movie)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK, then return the slice of movies.
	return movies, nil

}
