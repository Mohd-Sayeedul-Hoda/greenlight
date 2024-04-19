package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"greenlight/internal/validator"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int64 `json:"id"`
	CreateAt time.Time `json:"created_at"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password password `json:"-"`
	Activated bool `json:"activated"`
	Version int `json:"-"`
}

type password struct {
	plaintext *string
	hash []byte
}

var(
	ErrDuplicateEmail = errors.New("duplicate emial")
)

type UserModel struct{
	DB *sql.DB
}

func (p *password) Set(plainTextPassword string) error{
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil{
		return err
	}
	p.plaintext = &plainTextPassword
	p.hash = hash
	
	return nil
}

func (p *password) Matches (plaintextPassword string) (bool, error){
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil{
		switch{
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// validation on user
func ValidateEmail(v *validator.Validator, email string){
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EMailRX), "email", "email must be valid address")
}

func ValidatePasswordPlainText(v *validator.Validator, password string){
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User){
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil{
		ValidatePasswordPlainText(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil{
		panic("missing password hash for user")
	}
}

// crud on user

func (m UserModel) Insert(user *User) error{
	query := `INSERT INTO users (name, email, password_hash, activated) VALUES ($1, $2, $3, $4) RETURNING id, created_at, version`

	args := []any{user.Name, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()


	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreateAt, &user.Version)
	if err != nil{
		switch{
			case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
				return ErrDuplicateEmail
			default:
				return err
			}
		}

	return nil
	
}

func (m UserModel) GetByEmail(email string)(*User, error){
	query := `SELECT id, created_at, name, email, password_hash, activated, version FROM users WHERE email = $1`

	var user User
	
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreateAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil{
		switch{
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error{
	query := `UPDATE users
	SET name = $1, email = $2, password_hash = $3, activated = $4, version = version = version + 1 WHERE id = $5 AND version = $6 RETURNING version`

	args := []any{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch{
		case err.Error() == `pq: dubplicate key value violates unique constraint "users_email_key"`:
		return ErrDuplicateEmail
		}
	}

	return nil
}