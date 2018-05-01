package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/central/models"
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
	week = time.Hour * 24 * 7
)

func NewSession(subject string) (*models.Session, error) {
	id := ksuid.New().String()
	ae := time.Now().Add(week)
	at, err := NewToken(id, subject, ae, access)
	if err != nil {
		return nil, err
	}
	re := time.Now().Add(week * 4)
	rt, err := NewToken("r"+id, subject, re, refresh)
	if err != nil {
		return nil, err
	}
	return &models.Session{
		AccessToken:      at,
		ExpiresAt:        ae.Unix(),
		RefreshToken:     rt,
		RefreshExpiresAt: re.Unix(),
		SubjectID:        subject,
		TokenType:        "JWT",
	}, nil
}

func NewToken(id string, subject string, expiry time.Time, scope scope) (string, error) {
	claims := &textileClaims{
		scope,
		jwt.StandardClaims{
			Audience:  "textile",
			ExpiresAt: expiry.Unix(),
			Id:        id,
			IssuedAt:  time.Now().Unix(),
			Issuer:    os.Getenv("TOKEN_ISSUER"),
			Subject:   subject,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}
