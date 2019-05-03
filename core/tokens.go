package core

import (
	"encoding/hex"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/mr-tron/base58/base58"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/pb"
	"golang.org/x/crypto/bcrypt"
)

// CafeTokens lists all locally-stored (bcrypt hashed) tokens
func (t *Textile) CafeTokens() ([]string, error) {
	tokens := t.datastore.CafeTokens().List()
	strings := make([]string, len(tokens))

	for i, token := range tokens {
		id, err := hex.DecodeString(token.Id)
		if err != nil {
			return []string{}, err
		}
		strings[i] = base58.FastBase58Encoding(append(id, token.Value...))
	}

	return strings, nil
}

// CreateCafeToken creates (or uses `token`) random access token, returns base58 encoded version,
// and stores (unless `store` is false) a bcrypt hashed version for later comparison
func (t *Textile) CreateCafeToken(token string, store bool) (string, error) {
	var key []byte
	var err error
	if token != "" {
		key, err = base58.FastBase58Decoding(token)
		if err != nil {
			return "", err
		}
		if len(key) != 44 {
			return "", fmt.Errorf("invalid token format")
		}
	} else {
		key, err = crypto.GenerateAESKey()
		if err != nil {
			return "", err
		}
	}

	rawToken := key[12:]
	safeToken, err := bcrypt.GenerateFromPassword(rawToken, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	if store {
		if err := t.datastore.CafeTokens().Add(&pb.CafeToken{
			Id:    hex.EncodeToString(key[:12]),
			Value: safeToken,
			Date:  ptypes.TimestampNow(),
		}); err != nil {
			return "", err
		}
	}

	return base58.FastBase58Encoding(key), nil
}

// ValidateCafeToken checks whether a supplied base58 encoded token matches the locally-stored
// bcrypt hashed equivalent
func (t *Textile) ValidateCafeToken(token string) (bool, error) {
	// dev tokens are actually base58(id+token)
	plainBytes, err := base58.FastBase58Decoding(token)
	if err != nil {
		return false, err
	}
	if len(plainBytes) < 44 {
		return false, fmt.Errorf("invalid token format")
	}

	encodedToken := t.datastore.CafeTokens().Get(hex.EncodeToString(plainBytes[:12]))
	if encodedToken == nil {
		return false, err
	}
	if err := bcrypt.CompareHashAndPassword(encodedToken.Value, plainBytes[12:]); err != nil {
		return false, err
	}
	return true, nil
}

// RemoveCafeToken removes a given cafe token from the local store
func (t *Textile) RemoveCafeToken(token string) error {
	// dev tokens are actually base58(id+token)
	plainBytes, err := base58.FastBase58Decoding(token)
	if err != nil {
		return err
	}
	if len(plainBytes) < 44 {
		return fmt.Errorf("invalid token format")
	}
	return t.datastore.CafeTokens().Delete(hex.EncodeToString(plainBytes[:12]))
}
