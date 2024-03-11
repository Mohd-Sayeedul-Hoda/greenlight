package data

import (
	"database/sql"
	"errors"
	"greenlight/internal/validator"
	"time"

	"github.com/lib/pq"
)

type MockMovieModel struct{}

type MovieModel struct{
	DB *sql.DB
}

type Movie struct{
	ID int64 `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title string	`json:"title"`
	Year int32	`json:"year,omitempty"`
	Runtime Runtime `json:"runtime,omitempty"`// movie lenght
	Genres []string `json:"genres,omitempty"`
	Version int32 `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie){
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must no be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a postive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at leas 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")

	v.Check(validator.Unique(movie.Genres), "genres", "most not contain duplicate values")

}

func(m MovieModel) Insert(movie *Movie) error{
	query := `
		INSERT INTO movies (title, year, runtime, genres
		) VALUES ($1, $2, $3, $4) RETURNING id, created_at, version`
	
	// pq array is changing Array go array type into 
	// psql type array
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	

	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version,)
}

func(m MovieModel) Get(id int64) (*Movie, error){
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, created_at, title, year,
	runtime, genres, version FROM movies where id 
	= $1`

	var movie Movie

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
		)

	if err != nil{
		switch{
			case errors.Is(err, sql.ErrNoRows):
				return nil, ErrRecordNotFound
			default:
				return nil, err
		}
	}
	return &movie, nil
}

func(m MovieModel) Update(movie *Movie) error{
	return nil
}

func(m MovieModel) Delete(id int64) error{
	return nil
}

// Mock start here

func(m MockMovieModel) Insert(movie *Movie) error{
	return nil
}

func(m MockMovieModel) Get(id int64) (*Movie, error){
	return nil, nil
}

func(m MockMovieModel) Update(movie *Movie) error{
	return nil
}

func(m MockMovieModel) Delete(id int64) error{
	return nil
}
