package core

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/cafe"
	"github.com/textileio/textile-go/cafe/client"
	cmodels "github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"time"
)

// GetCafeAddr returns the cafe address if set
func (t *Textile) GetCafeAddr() string {
	return t.cafeAddr
}

// GetCafeApiAddr returns the cafe address if set
func (t *Textile) GetCafeApiAddr() string {
	if t.cafeAddr == "" {
		return ""
	}
	return fmt.Sprintf("%s/api/%s", t.cafeAddr, cafe.Version)
}

// getCafeChallenge requests a challenge from a cafe and signs it
func (t *Textile) getCafeChallenge(accnt *keypair.Full) (*cmodels.SignedChallenge, error) {
	if t.cafeAddr == "" {
		return nil, ErrNoCafeHost
	}
	req := &cmodels.ChallengeRequest{Address: accnt.Address()}
	cres, err := client.ProfileChallenge(req, fmt.Sprintf("%s/profiles/challenge", t.GetCafeApiAddr()))
	if err != nil {
		log.Errorf("get challenge error: %s", err)
		return nil, err
	}
	if cres.Error != nil {
		log.Errorf("get challenge error from cafe: %s", *cres.Error)
		return nil, errors.New(*cres.Error)
	}
	if cres.Value == nil {
		return nil, errors.New("cafe returned nil challenge")
	}
	cnonce := ksuid.New().String()
	sigb, err := accnt.Sign([]byte(*cres.Value + cnonce))
	if err != nil {
		return nil, err
	}
	return &cmodels.SignedChallenge{
		Address:   accnt.Address(),
		Value:     *cres.Value,
		Nonce:     cnonce,
		Signature: libp2pc.ConfigEncodeKey(sigb),
	}, nil
}

// CafeRegister registers a public key w/ a cafe, requests a session token, and saves it locally
func (t *Textile) CafeRegister(peerId string) error {
	if err := t.touchDatastore(); err != nil {
		return err
	}
	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return err
	}

	// get a challenge from the cafe
	accnt, err := t.Account()
	if err != nil {
		return err
	}
	res, err := t.cafeService.RequestChallenge(accnt, pid)
	if err != nil {
		return err
	}

	// complete the challenge
	cnonce := ksuid.New().String()
	sig, err := accnt.Sign([]byte(res.Value + cnonce))
	if err != nil {
		return err
	}
	reg := &pb.CafeRegistration{
		Address: accnt.Address(),
		Value:   res.Value,
		Nonce:   cnonce,
		Sig:     sig,
	}
	session, err := t.cafeService.Register(reg, pid)
	if err != nil {
		return err
	}

	// local login
	exp, err := ptypes.Timestamp(session.Exp)
	if err != nil {
		return err
	}
	tokens := &repo.CafeTokens{
		Access:  session.Access,
		Refresh: session.Refresh,
		Expiry:  exp,
	}
	if err := t.datastore.Profile().CafeLogin(tokens); err != nil {
		log.Errorf("local login error: %s", err)
		return err
	}
	return nil
}

// CafeLogin requests a session token from a cafe and saves it locally
func (t *Textile) CafeLogin() error {
	if t.cafeAddr == "" {
		return ErrNoCafeHost
	}
	if err := t.touchDatastore(); err != nil {
		return err
	}

	// get a challenge from the cafe
	accnt, err := t.Account()
	if err != nil {
		return err
	}
	challenge, err := t.getCafeChallenge(accnt)
	if err != nil {
		return err
	}

	log.Debugf("login: %s %s", challenge.Address, challenge.Signature)

	// remote login
	res, err := client.LoginProfile(challenge, fmt.Sprintf("%s/profiles", t.GetCafeApiAddr()))
	if err != nil {
		log.Errorf("login error: %s", err)
		return err
	}
	if res.Error != nil {
		log.Errorf("login error from cafe: %s", *res.Error)
		return errors.New(*res.Error)
	}

	// local login
	tokens := &repo.CafeTokens{
		Access:  res.Session.AccessToken,
		Refresh: res.Session.RefreshToken,
		Expiry:  time.Unix(res.Session.ExpiresAt, 0),
	}
	if err := t.datastore.Profile().CafeLogin(tokens); err != nil {
		log.Errorf("local login error: %s", err)
		return err
	}

	return nil
}

// CreateCafeReferral requests a referral from a cafe via key
func (t *Textile) CreateCafeReferral(req *cmodels.ReferralRequest) (*cmodels.ReferralResponse, error) {
	if t.cafeAddr == "" {
		return nil, ErrNoCafeHost
	}
	log.Debug("requesting a referral")

	// remote request
	res, err := client.CreateReferral(req, fmt.Sprintf("%s/referrals", t.GetCafeApiAddr()))
	if err != nil {
		log.Errorf("create referral error: %s", err)
		return nil, err
	}
	if res.Error != nil {
		log.Errorf("create referral error from cafe: %s", *res.Error)
		return nil, errors.New(*res.Error)
	}
	return res, nil
}

// ListCafeReferrals lists existing referrals from a cafe via key
func (t *Textile) ListCafeReferrals(key string) (*cmodels.ReferralResponse, error) {
	if t.cafeAddr == "" {
		return nil, ErrNoCafeHost
	}
	log.Debug("listing referrals")

	// remote request
	res, err := client.ListReferrals(key, fmt.Sprintf("%s/referrals", t.GetCafeApiAddr()))
	if err != nil {
		log.Errorf("list referrals error: %s", err)
		return nil, err
	}
	if res.Error != nil {
		log.Errorf("list referrals error from cafe: %s", *res.Error)
		return nil, errors.New(*res.Error)
	}
	return res, nil
}

// CafeLogout deletes the locally saved profile key and cafe session if present
func (t *Textile) CafeLogout() error {
	if err := t.touchDatastore(); err != nil {
		return err
	}
	log.Debug("logging out...")

	// remote is stateless, so we just ditch the local token
	if err := t.datastore.Profile().CafeLogout(); err != nil {
		log.Errorf("local logout error: %s", err)
		return err
	}

	return nil
}

// CafeLoggedIn returns whether or not the profile is logged into a cafe
func (t *Textile) CafeLoggedIn() (bool, error) {
	if err := t.touchDatastore(); err != nil {
		return false, err
	}
	tokens, err := t.datastore.Profile().GetCafeTokens()
	if err != nil {
		return false, err
	}
	return tokens != nil, nil
}

// GetCafeTokens returns cafe json web tokens, refreshing if needed or if forceRefresh is true
func (t *Textile) GetCafeTokens(forceRefresh bool) (*repo.CafeTokens, error) {
	if err := t.touchDatastore(); err != nil {
		return nil, err
	}
	tokens, err := t.datastore.Profile().GetCafeTokens()
	if err != nil {
		return nil, err
	}
	if tokens == nil {
		return nil, nil
	}

	// check expiry
	if tokens.Expiry.After(time.Now()) && !forceRefresh {
		return tokens, nil
	}

	// remote refresh
	url := fmt.Sprintf("%s/tokens", t.GetCafeApiAddr())
	res, err := client.RefreshSession(tokens.Access, tokens.Refresh, url)
	if err != nil {
		log.Errorf("get tokens error: %s", err)
		return nil, err
	}
	if res.Error != nil {
		log.Errorf("get tokens error from cafe: %s", *res.Error)
		return nil, errors.New(*res.Error)
	}

	// update tokens
	tokens = &repo.CafeTokens{
		Access:  res.Session.AccessToken,
		Refresh: res.Session.RefreshToken,
		Expiry:  time.Unix(res.Session.ExpiresAt, 0),
	}
	if err := t.datastore.Profile().CafeLogin(tokens); err != nil {
		log.Errorf("update tokens error: %s", err)
		return nil, err
	}
	return tokens, nil
}
