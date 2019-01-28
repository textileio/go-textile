package core

import (
	"fmt"
	"time"

	"github.com/mr-tron/base58/base58"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/repo"
	"golang.org/x/crypto/bcrypt"
)

// CreateCafeToken creates a random developer access token, returns a base58 encoded version,
// and stores a bcrypt hashed version for later comparison
func (t *Textile) CreateCafeToken() (*repo.CafeDevToken, error) {
	key, err := crypto.GenerateAESKey()
	if err != nil {
		return nil, err
	}

	id := ksuid.New().String()
	created := time.Now()
	rawToken := key[:32]

	safeToken, err := bcrypt.GenerateFromPassword(rawToken, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	fmt.Println(t.datastore)

	err = t.datastore.CafeDevTokens().Add(
		&repo.CafeDevToken{
			Id:      id,
			Token:   string(safeToken),
			Created: created,
		})
	if err != nil {
		return nil, err
	}

	return &repo.CafeDevToken{
		Id:      id,
		Token:   base58.FastBase58Encoding(rawToken),
		Created: created,
	}, nil
}

// CafeDevTokens lists all stored (bcrypt encrypted) dev tokens
func (t *Textile) CafeDevTokens() ([]repo.CafeDevToken, error) {
	return t.datastore.CafeDevTokens().List(), nil
}

// CheckCafeDevToken checks whether a supplied base58 encoded dev token matches the stored
// bcrypt hashed equivalent
func (t *Textile) CompareCafeDevToken(id string, token string) (bool, error) {
	plainBytes, err := base58.FastBase58Decoding(token)
	if err != nil {
		return false, err
	}

	hashedToken := t.datastore.CafeDevTokens().Get(id)
	if hashedToken == nil {
		return false, err
	}

	hashBytes := []byte(hashedToken.Token)
	err = bcrypt.CompareHashAndPassword(hashBytes, plainBytes)
	if err != nil {
		return false, err
	}

	return true, nil
}

// RemoveDevToken removes a given cafe dev token by id
func (t *Textile) RemoveCafeDevToken(id string) error {
	return t.datastore.CafeDevTokens().Delete(id)
}
