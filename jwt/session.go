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
