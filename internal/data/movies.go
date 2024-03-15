package data

import (
	"context"
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

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// we are using int64 on uint because error
func(m MovieModel) Get(id int64) (*Movie, error){
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, created_at, title, year,
	runtime, genres, version FROM movies where id 
	= $1`

	var movie Movie

	//using ctx we are declaring context or meta data for we generally 
	// specifiy for the operation we are sending here
	// to db driver that is executing the function and
	// it interally do the query and listen for the ctx from the channel
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// we are clearing all the resource that context is taken at last
	defer cancel()

	//when listen for the signal then it terminate the query and return
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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
	query := `UPDATE movies set title = $1, year = 
		$2, runtime = $3, genres = $4, version = version + 1 WHERE id = $5 AND version = $6 RETURNING version`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil{
		switch{
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func(m MovieModel) Delete(id int64) error{

	if id < 1{
		return ErrRecordNotFound
	}
	query := `DELETE from movies WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil{
		return nil
	}

	rowAffected, err := result.RowsAffected()
	if err != nil{
		return err
	}

	if rowAffected == 0{
		return ErrRecordNotFound
	}

	return nil
}

func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error){

	// @> say if contian pq array 
	query := `SELECT id, created_at, title, year, runtime, genres, version from movies WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') AND (genres @> $2 OR $2 = '{}') ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, title, pq.Array(genres))
	if err != nil{
		return nil, err
	}

	defer rows.Close()

	movies := []*Movie{}

	for rows.Next(){
		
		var movie Movie

		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)

		if err != nil{
			return nil, err
		}

		movies = append(movies, &movie)
	}
	if err = rows.Err(); err != nil{
		return nil, err
	}
	return movies, nil
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

func (m MockMovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error){
	return nil, nil
}
