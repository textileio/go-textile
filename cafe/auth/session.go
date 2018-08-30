package auth

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/cafe/models"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"time"
)

var ErrInvalidClaims = errors.New("invalid claims")

type TextileClaims struct {
	Scope Scope `json:"scopes"`
	jwt.StandardClaims
}

type Scope string

const (
	Access  Scope = "access"
	Refresh Scope = "refresh"
)

func NewSession(subject string, secret string, issuer string, audience protocol.ID, duration time.Duration) (*models.Session, error) {
	id := ksuid.New().String()
	now := time.Now()
	expiresAt := now.Add(duration)
	accessToken, err := NewToken(id, subject, expiresAt, Access, secret, issuer, string(audience))
	if err != nil {
		return nil, err
	}
	refreshExpiresAt := now.Add(duration * 2)
	refreshToken, err := NewToken("r"+id, subject, refreshExpiresAt, Refresh, secret, issuer, string(audience))
	if err != nil {
		return nil, err
	}
	return &models.Session{
		AccessToken:      accessToken,
		ExpiresAt:        expiresAt.Unix(),
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpiresAt.Unix(),
		SubjectId:        subject,
		TokenType:        "JWT",
	}, nil
}

func NewToken(id string, subject string, expiry time.Time, scope Scope, secret string, issuer string, audience string) (string, error) {
	claims := &TextileClaims{
		Scope: scope,
		StandardClaims: jwt.StandardClaims{
			Audience:  audience,
			ExpiresAt: expiry.Unix(),
			Id:        id,
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
			Subject:   subject,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func ParseClaims(claims jwt.Claims) (*TextileClaims, error) {
	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}
	claimsb, err := json.Marshal(mapClaims)
	if err != nil {
		return nil, ErrInvalidClaims
	}
	var tclaims *TextileClaims
	if err := json.Unmarshal(claimsb, &tclaims); err != nil {
		return nil, ErrInvalidClaims
	}
	return tclaims, nil
}
