package core

import (
	"encoding/hex"
	"errors"
	"time"

	"github.com/mr-tron/base58/base58"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/repo"
	"golang.org/x/crypto/bcrypt"
)

// CreateCafeToken creates a random access token, returns a base58 encoded version,
// and stores a bcrypt hashed version for later comparison
func (t *Textile) CreateCafeToken() (string, error) {
	key, err := crypto.GenerateAESKey()
	if err != nil {
		return "", err
	}

	date := time.Now()
	id := hex.EncodeToString(key[:12])
	rawToken := key[12:]

	safeToken, err := bcrypt.GenerateFromPassword(rawToken, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	err = t.datastore.CafeTokens().Add(
		&repo.CafeToken{
			Id:    id,
			Token: safeToken,
			Date:  date,
		})
	if err != nil {
		return "", err
	}
	return base58.FastBase58Encoding(key), nil
}

// CafeTokens lists all stored (bcrypt hashed) tokens
func (t *Textile) CafeTokens() ([]string, error) {
	tokens := t.datastore.CafeTokens().List()
	strings := make([]string, len(tokens))
	for i, token := range tokens {
		id, err := hex.DecodeString(token.Id)
		if err != nil {
			return []string{}, err
		}
		strings[i] = base58.FastBase58Encoding(append(id, token.Token...))
	}
	return strings, nil
}

// ValidateCafeToken checks whether a supplied base58 encoded token matches the stored
// bcrypt hashed equivalent
func (t *Textile) ValidateCafeToken(token string) (bool, error) {
	// dev tokens are actually base58(id+token)
	plainBytes, err := base58.FastBase58Decoding(token)
	if err != nil {
		return false, err
	}
	if len(plainBytes) < 44 {
		return false, errors.New("invalid token format")
	}
	encodedToken := t.datastore.CafeTokens().Get(hex.EncodeToString(plainBytes[:12]))
	if encodedToken == nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword(encodedToken.Token, plainBytes[12:])
	if err != nil {
		return false, err
	}
	return true, nil
}

// RemoveCafeToken removes a given cafe token
func (t *Textile) RemoveCafeToken(token string) error {
	// dev tokens are actually base58(id+token)
	plainBytes, err := base58.FastBase58Decoding(token)
	if err != nil {
		return err
	}
	if len(plainBytes) < 44 {
		return errors.New("invalid token format")
	}
	return t.datastore.CafeTokens().Delete(hex.EncodeToString(plainBytes[:12]))
}
