package wallet

import (
	"errors"
	"fmt"
	cmodels "github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/core/cafe"
	"github.com/textileio/textile-go/repo"
)

// CreateReferral requests a referral from a cafe via key
func (w *Wallet) CreateReferral(req *cmodels.ReferralRequest) (*cmodels.ReferralResponse, error) {
	if w.cafeAddr == "" {
		return nil, ErrNoCafeHost
	}
	log.Debug("requesting a referral")

	// remote request
	res, err := client.CreateReferral(req, fmt.Sprintf("%s/referrals", w.GetCafeAddr()))
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

// ListReferrals lists existing referrals from a cafe via key
func (w *Wallet) ListReferrals(key string) (*cmodels.ReferralResponse, error) {
	if w.cafeAddr == "" {
		return nil, ErrNoCafeHost
	}
	log.Debug("listing referrals")

	// remote request
	res, err := client.ListReferrals(key, fmt.Sprintf("%s/referrals", w.GetCafeAddr()))
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

// SignUp requests a new username and token from a cafe and saves them locally
func (w *Wallet) SignUp(reg *cmodels.Registration) error {
	if w.cafeAddr == "" {
		return ErrNoCafeHost
	}
	if err := w.touchDatastore(); err != nil {
		return err
	}
	log.Debugf("signup: %s %s %s %s %s", reg.Username, "xxxxxx", reg.Identity.Type, reg.Identity.Value, reg.Referral)

	// remote signup
	res, err := client.SignUp(reg, fmt.Sprintf("%s/users", w.GetCafeAddr()))
	if err != nil {
		log.Errorf("signup error: %s", err)
		return err
	}
	if res.Error != nil {
		log.Errorf("signup error from cafe: %s", *res.Error)
		return errors.New(*res.Error)
	}

	// local signin
	tokens := &repo.CafeTokens{
		Access:  res.Session.AccessToken,
		Refresh: res.Session.RefreshToken,
	}
	if err := w.datastore.Profile().SignIn(reg.Username, tokens); err != nil {
		log.Errorf("local signin error: %s", err)
		return err
	}
	return nil
}

// SignIn requests a token with a username from a cafe and saves them locally
func (w *Wallet) SignIn(creds *cmodels.Credentials) error {
	if w.cafeAddr == "" {
		return ErrNoCafeHost
	}
	if err := w.touchDatastore(); err != nil {
		return err
	}
	log.Debugf("signin: %s %s", creds.Username, "xxxxxx")

	// remote signin
	res, err := client.SignIn(creds, fmt.Sprintf("%s/users", w.GetCafeAddr()))
	if err != nil {
		log.Errorf("signin error: %s", err)
		return err
	}
	if res.Error != nil {
		log.Errorf("signin error from cafe: %s", *res.Error)
		return errors.New(*res.Error)
	}

	// local signin
	tokens := &repo.CafeTokens{
		Access:  res.Session.AccessToken,
		Refresh: res.Session.RefreshToken,
	}
	if err := w.datastore.Profile().SignIn(creds.Username, tokens); err != nil {
		log.Errorf("local signin error: %s", err)
		return err
	}
	return nil
}

// SignOut deletes the locally saved user info (username and tokens)
func (w *Wallet) SignOut() error {
	if w.cafeAddr == "" {
		return ErrNoCafeHost
	}
	if err := w.touchDatastore(); err != nil {
		return err
	}
	log.Debug("signing out...")

	// remote is stateless, so we just ditch the local token
	if err := w.datastore.Profile().SignOut(); err != nil {
		log.Errorf("local signout error: %s", err)
		return err
	}
	return nil
}

// IsSignedIn returns whether or not a user is signed in
func (w *Wallet) IsSignedIn() (bool, error) {
	if w.cafeAddr == "" {
		return false, ErrNoCafeHost
	}
	if err := w.touchDatastore(); err != nil {
		return false, err
	}
	_, err := w.datastore.Profile().GetUsername()
	return err == nil, nil
}

// GetUsername returns the current user's username
func (w *Wallet) GetUsername() (string, error) {
	if w.cafeAddr == "" {
		return "", ErrNoCafeHost
	}
	if err := w.touchDatastore(); err != nil {
		return "", err
	}
	return w.datastore.Profile().GetUsername()
}

// GetAccessToken returns the current access_token (jwt) for a cafe
func (w *Wallet) GetTokens() (*repo.CafeTokens, error) {
	if w.cafeAddr == "" {
		return nil, ErrNoCafeHost
	}
	if err := w.touchDatastore(); err != nil {
		return nil, err
	}
	return w.datastore.Profile().GetTokens()
}
