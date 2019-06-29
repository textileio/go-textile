package jwt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/ptypes"
	libp2pc "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/pb"
)

var ErrClaimsInvalid = fmt.Errorf("claims invalid")
var ErrNoToken = fmt.Errorf("no token found")
var ErrExpired = fmt.Errorf("token expired")
var ErrInvalid = fmt.Errorf("token invalid")

type TextileClaims struct {
	Scope Scope `json:"scopes"`
	jwt.StandardClaims
}

type Scope string

const (
	Access  Scope = "access"
	Refresh Scope = "refresh"
)

func NewSession(sk libp2pc.PrivKey, pid peer.ID, proto protocol.ID, duration time.Duration, cafe *pb.Cafe) (*pb.CafeSession, error) {
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
		Id:      issuer.Pretty(),
		Access:  access,
		Exp:     pexp,
		Refresh: refresh,
		Rexp:    prexp,
		Subject: pid.Pretty(),
		Type:    "JWT",
		Cafe:    cafe,
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

func Validate(tokenString string, keyfunc jwt.Keyfunc, refreshing bool, audience string, subject *string) (*TextileClaims, error) {
	token, pErr := jwt.Parse(tokenString, keyfunc)
	if token == nil {
		return nil, ErrNoToken
	}

	claims, err := ParseClaims(token.Claims)
	if err != nil {
		return nil, ErrInvalid
	}

	if pErr != nil {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			return nil, ErrExpired
		}
		return nil, ErrInvalid
	}

	switch claims.Scope {
	case Access:
		if refreshing {
			return nil, ErrInvalid
		}
	case Refresh:
		if !refreshing {
			return nil, ErrInvalid
		}
	default:
		return nil, ErrInvalid
	}

	// verify owner
	if subject != nil && *subject != claims.Subject {
		return nil, ErrInvalid
	}

	// verify protocol
	if !claims.VerifyAudience(audience, true) {
		return nil, ErrInvalid
	}
	return claims, nil
}
