package wallet

import (
	"errors"
	"fmt"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/cafe"
	cmodels "github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/core/cafe"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"time"
)

// GetCafeAddr returns the cafe address if set
func (w *Wallet) GetCafeAddr() string {
	return w.cafeAddr
}

// GetCafeApiAddr returns the cafe address if set
func (w *Wallet) GetCafeApiAddr() string {
	if w.cafeAddr == "" {
		return ""
	}
	return fmt.Sprintf("%s/api/%s", w.cafeAddr, cafe.Version)
}

// getCafeChallenge requests a challenge from a cafe and signs it
func (w *Wallet) getCafeChallenge(key libp2pc.PrivKey) (*cmodels.SignedChallenge, error) {
	if w.cafeAddr == "" {
		return nil, ErrNoCafeHost
	}
	pks, err := util.EncodeKey(key.GetPublic())
	if err != nil {
		return nil, err
	}
	req := &cmodels.ChallengeRequest{Pk: pks}
	cres, err := client.ProfileChallenge(req, fmt.Sprintf("%s/profiles/challenge", w.GetCafeApiAddr()))
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
	sigb, err := key.Sign([]byte(*cres.Value + cnonce))
	if err != nil {
		return nil, err
	}
	return &cmodels.SignedChallenge{
		Pk:        pks,
		Value:     *cres.Value,
		Nonce:     cnonce,
		Signature: libp2pc.ConfigEncodeKey(sigb),
	}, nil
}

// CafeRegister registers a public key w/ a cafe, requests a session token, and saves it locally
func (w *Wallet) CafeRegister(referral string) error {
	if w.cafeAddr == "" {
		return ErrNoCafeHost
	}
	if err := w.touchDatastore(); err != nil {
		return err
	}

	// get a challenge from the cafe
	key, err := w.GetKey()
	if err != nil {
		return err
	}
	challenge, err := w.getCafeChallenge(key)
	if err != nil {
		return err
	}
	reg := &cmodels.ProfileRegistration{
		Challenge: *challenge,
		Referral:  referral,
	}

	log.Debugf("cafe register: %s %s %s", reg.Challenge.Pk, reg.Challenge.Signature, reg.Referral)

	// remote register
	res, err := client.RegisterProfile(reg, fmt.Sprintf("%s/profiles", w.GetCafeApiAddr()))
	if err != nil {
		log.Errorf("register error: %s", err)
		return err
	}
	if res.Error != nil {
		log.Errorf("register error from cafe: %s", *res.Error)
		return errors.New(*res.Error)
	}

	// local login
	tokens := &repo.CafeTokens{
		Access:  res.Session.AccessToken,
		Refresh: res.Session.RefreshToken,
		Expiry:  time.Unix(res.Session.ExpiresAt, 0),
	}
	if err := w.datastore.Profile().CafeLogin(tokens); err != nil {
		log.Errorf("local login error: %s", err)
		return err
	}

	// initial profile publish
	go func() {
		<-w.Online()
		if _, err := w.PublishProfile(nil); err != nil {
			log.Errorf("error publishing initial profile: %s", err)
		}
	}()

	return nil
}

// CafeLogin requests a session token from a cafe and saves it locally
func (w *Wallet) CafeLogin() error {
	if w.cafeAddr == "" {
		return ErrNoCafeHost
	}
	if err := w.touchDatastore(); err != nil {
		return err
	}

	// get a challenge from the cafe
	key, err := w.GetKey()
	if err != nil {
		return err
	}
	challenge, err := w.getCafeChallenge(key)
	if err != nil {
		return err
	}

	log.Debugf("login: %s %s", challenge.Pk, challenge.Signature)

	// remote login
	res, err := client.LoginProfile(challenge, fmt.Sprintf("%s/profile", w.GetCafeApiAddr()))
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
	if err := w.datastore.Profile().CafeLogin(tokens); err != nil {
		log.Errorf("local login error: %s", err)
		return err
	}

	return nil
}

// CreateCafeReferral requests a referral from a cafe via key
func (w *Wallet) CreateCafeReferral(req *cmodels.ReferralRequest) (*cmodels.ReferralResponse, error) {
	if w.cafeAddr == "" {
		return nil, ErrNoCafeHost
	}
	log.Debug("requesting a referral")

	// remote request
	res, err := client.CreateReferral(req, fmt.Sprintf("%s/referrals", w.GetCafeApiAddr()))
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
func (w *Wallet) ListCafeReferrals(key string) (*cmodels.ReferralResponse, error) {
	if w.cafeAddr == "" {
		return nil, ErrNoCafeHost
	}
	log.Debug("listing referrals")

	// remote request
	res, err := client.ListReferrals(key, fmt.Sprintf("%s/referrals", w.GetCafeApiAddr()))
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
func (w *Wallet) CafeLogout() error {
	if err := w.touchDatastore(); err != nil {
		return err
	}
	log.Debug("logging out...")

	// remote is stateless, so we just ditch the local token
	if err := w.datastore.Profile().CafeLogout(); err != nil {
		log.Errorf("local logout error: %s", err)
		return err
	}

	return nil
}

// CafeLoggedIn returns whether or not the profile is logged into a cafe
func (w *Wallet) CafeLoggedIn() (bool, error) {
	if err := w.touchDatastore(); err != nil {
		return false, err
	}
	tokens, err := w.datastore.Profile().GetCafeTokens()
	if err != nil {
		return false, err
	}
	return tokens != nil, nil
}

// GetCafeTokens returns cafe json web tokens, refreshing if needed or if forceRefresh is true
func (w *Wallet) GetCafeTokens(forceRefresh bool) (*repo.CafeTokens, error) {
	if err := w.touchDatastore(); err != nil {
		return nil, err
	}
	tokens, err := w.datastore.Profile().GetCafeTokens()
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
	url := fmt.Sprintf("%s/tokens", w.GetCafeApiAddr())
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
	if err := w.datastore.Profile().CafeLogin(tokens); err != nil {
		log.Errorf("update tokens error: %s", err)
		return nil, err
	}
	return tokens, nil
}
