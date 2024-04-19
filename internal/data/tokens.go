package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
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
	// we are ultimatly ecnoding here to when send it to user then it can easly read it 
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	//we are making hash of size 256 bit of plaintext
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}
