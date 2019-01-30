package core

import (
	"errors"
	"strings"
	"time"

	"github.com/mr-tron/base58/base58"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/repo"
	"golang.org/x/crypto/bcrypt"
)

// CreateCafeToken creates a random developer access token, returns a base58 encoded version,
// and stores a bcrypt hashed version for later comparison
func (t *Textile) CreateCafeToken() (string, error) {
	key, err := crypto.GenerateAESKey()
	if err != nil {
		return "", err
	}

	id := ksuid.New().String()
	created := time.Now()
	rawToken := key[:32]

	safeToken, err := bcrypt.GenerateFromPassword(rawToken, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	err = t.datastore.CafeDevTokens().Add(
		&repo.CafeDevToken{
			Id:      id,
			Token:   safeToken,
			Created: created,
		})
	if err != nil {
		return "", err
	}

	return id + "+" + base58.FastBase58Encoding(rawToken), nil
}

// CafeDevTokens lists all stored (bcrypt encrypted) dev tokens
func (t *Textile) CafeDevTokens() ([]string, error) {
	tokens := t.datastore.CafeDevTokens().List()
	strings := make([]string, len(tokens))
	for i, token := range tokens {
		strings[i] = token.Id + "+" + base58.FastBase58Encoding(token.Token)
	}
	return strings, nil
}

// CompareCafeDevToken checks whether a supplied base58 encoded dev token matches the stored
// bcrypt hashed equivalent
func (t *Textile) CompareCafeDevToken(token string) (bool, error) {
	// dev tokens are actually ksuid+base58(token)
	s := strings.Split(token, "+")
	if len(s) != 2 {
		return false, errors.New("invalid token format")
	}
	id, token := s[0], s[1]

	plainBytes, err := base58.FastBase58Decoding(token)
	if err != nil {
		return false, err
	}

	encodedToken := t.datastore.CafeDevTokens().Get(id)
	if encodedToken == nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(encodedToken.Token, plainBytes)
	if err != nil {
		return false, err
	}

	return true, nil
}

// RemoveCafeDevToken removes a given cafe dev token by id
func (t *Textile) RemoveCafeDevToken(token string) error {
	// dev tokens are actually ksuid+base58(token)
	s := strings.Split(token, "+")
	if len(s) != 2 {
		return errors.New("invalid token format")
	}
	return t.datastore.CafeDevTokens().Delete(s[0])
}
