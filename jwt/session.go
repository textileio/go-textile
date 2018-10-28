package jwt

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/ptypes"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/pb"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"time"
)

var ErrClaimsInvalid = errors.New("claims invalid")
var ErrNoToken = errors.New("no token found")
var ErrExpired = errors.New("token expired")
var ErrInvalid = errors.New("token invalid")

type TextileClaims struct {
	Scope Scope `json:"scopes"`
	jwt.StandardClaims
}

type Scope string

const (
	Access  Scope = "access"
	Refresh Scope = "refresh"
)

func NewSession(sk libp2pc.PrivKey, pid peer.ID, proto protocol.ID, duration time.Duration, httpAddr string) (*pb.CafeSession, error) {
	issuer, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}
	id := ksuid.New().String()

	// build access token
	now := time.Now()
	exp := now.Add(duration)
	claims := &TextileClaims{
		Scope: Access,
		StandardClaims: jwt.StandardClaims{
			Audience:  string(proto),
			ExpiresAt: exp.Unix(),
			Id:        id,
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer.Pretty(),
			Subject:   pid.Pretty(),
		},
	}
	access, err := jwt.NewWithClaims(SigningMethodEd25519i, claims).SignedString(sk)
	if err != nil {
		return nil, err
	}

	// build refresh token
	rexp := now.Add(duration * 2)
	rclaims := &TextileClaims{
		Scope: Refresh,
		StandardClaims: jwt.StandardClaims{
			Audience:  string(proto),
			ExpiresAt: rexp.Unix(),
			Id:        "r" + id,
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer.Pretty(),
			Subject:   pid.Pretty(),
		},
	}
	refresh, err := jwt.NewWithClaims(SigningMethodEd25519i, rclaims).SignedString(sk)
	if err != nil {
		return nil, err
	}

	// build session
	pexp, err := ptypes.TimestampProto(exp)
	if err != nil {
		return nil, err
	}
	prexp, err := ptypes.TimestampProto(rexp)
	if err != nil {
		return nil, err
	}
	return &pb.CafeSession{
		Access:   access,
		Exp:      pexp,
		Refresh:  refresh,
		Rexp:     prexp,
		Subject:  pid.Pretty(),
		Type:     "JWT",
		HttpAddr: httpAddr,
	}, nil
}

func ParseClaims(claims jwt.Claims) (*TextileClaims, error) {
	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrClaimsInvalid
	}
	claimsb, err := json.Marshal(mapClaims)
	if err != nil {
		return nil, ErrClaimsInvalid
	}
	var tclaims *TextileClaims
	if err := json.Unmarshal(claimsb, &tclaims); err != nil {
		return nil, ErrClaimsInvalid
	}
	return tclaims, nil
}

func Validate(tokenString string, keyfunc jwt.Keyfunc, refreshing bool, audience string, subject *string) error {
	// parse it
	token, pErr := jwt.Parse(tokenString, keyfunc)
	if token == nil {
		return ErrNoToken
	}

	// pull out claims
	claims, err := ParseClaims(token.Claims)
	if err != nil {
		return ErrInvalid
	}

	// check valid
	if pErr != nil {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			return ErrExpired
		}
		return ErrInvalid
	}

	// check scope
	switch claims.Scope {
	case Access:
		if refreshing {
			return ErrInvalid
		}
	case Refresh:
		if !refreshing {
			return ErrInvalid
		}
	default:
		return ErrInvalid
	}

	// verify owner
	if subject != nil && *subject != claims.Subject {
		return ErrInvalid
	}

	// verify protocol
	if !claims.VerifyAudience(audience, true) {
		return ErrInvalid
	}
	return nil
}
