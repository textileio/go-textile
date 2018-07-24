package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/cafe/models"
	"time"
)

type textileClaims struct {
	Scope scope `json:"scopes"`
	jwt.StandardClaims
}

type scope string

const (
	access  scope = "access"
	refresh scope = "refresh"
)

const (
	week  = time.Hour * 24 * 7
	month = week * 4
)

func NewSession(subject string, secret string, issuer string) (*models.Session, error) {
	id := ksuid.New().String()
	expiresAt := time.Now().Add(month * 3)
	access, err := NewToken(id, subject, expiresAt, access, secret, issuer)
	if err != nil {
		return nil, err
	}
	refreshExpiresAt := time.Now().Add(month * 6)
	refresh, err := NewToken("r"+id, subject, refreshExpiresAt, refresh, secret, issuer)
	if err != nil {
		return nil, err
	}
	return &models.Session{
		AccessToken:      access,
		ExpiresAt:        expiresAt.Unix(),
		RefreshToken:     refresh,
		RefreshExpiresAt: refreshExpiresAt.Unix(),
		SubjectId:        subject,
		TokenType:        "JWT",
	}, nil
}

func NewToken(id string, subject string, expiry time.Time, scope scope, secret string, issuer string) (string, error) {
	claims := &textileClaims{
		Scope: scope,
		StandardClaims: jwt.StandardClaims{
			Audience:  "/textile/app/1.0.0",
			ExpiresAt: expiry.Unix(),
			Id:        id,
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
			Subject:   subject,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}
