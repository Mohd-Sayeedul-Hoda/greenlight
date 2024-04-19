package data

import (
	"context"
	"database/sql"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"greenlight/internal/validator"
)

const (
	ScopeActivation = "activation"
)

type Token struct{
	Plaintext string
	Hash []byte
	UserID int64
	Expiry time.Time
	Scope string
}

func genrateToken(userID int64, ttl time.Duration, scope string)(*Token, error){

	//Creating token
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope: scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil{
		return nil, err
	}

		// so here we first getting struct for encoding that withNoPadding
	// then we are econding radombytes here
	// we are ultimatly encoding here to when send it to user then it can easly read it 
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	//we are making hash of size 256 bit of plaintext
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string){
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")

}

// Define the tokenModel type.
type TokenModel struct{
	DB *sql.DB
}

// The New() method is shortcut which creates a new token sturt and thes inserts the
// data in the token tabel.

func (m TokenModel) New(userID int64, ttl time.Duration, scope string)(*Token, error){
	token, err := genrateToken(userID, ttl, scope)
	if err != nil{
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

func (m TokenModel) Insert(token *Token) error{
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) VALUES($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)

	return err
}

func (m TokenModel) DeletAllForUser(scope string, userID int64) error{
	query := `DELETE FROM tokens 
	where scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}