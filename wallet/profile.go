package wallet

import (
	"errors"
	"fmt"
	ccafe "github.com/textileio/textile-go/cafe"
	cmodels "github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/core/cafe"
)

// SignUp requests a new username and token from the central api and saves them locally
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
		log.Errorf("signup error from central: %s", *res.Error)
		return errors.New(*res.Error)
	}

	// local signin
	if err := w.datastore.Profile().SignIn(
		reg.Username,
		res.Session.AccessToken, res.Session.RefreshToken,
	); err != nil {
		log.Errorf("local signin error: %s", err)
		return err
	}
	return nil
}

// SignIn requests a token with a username from the central api and saves them locally
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
		log.Errorf("signin error from central: %s", *res.Error)
		return errors.New(*res.Error)
	}

	// local signin
	if err := w.datastore.Profile().SignIn(
		creds.Username,
		res.Session.AccessToken, res.Session.RefreshToken,
	); err != nil {
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

// GetAccessToken returns the current access_token (jwt) for central
func (w *Wallet) GetAccessToken() (string, error) {
	if w.cafeAddr == "" {
		return "", ErrNoCafeHost
	}
	if err := w.touchDatastore(); err != nil {
		return "", err
	}
	at, _, err := w.datastore.Profile().GetTokens()
	if err != nil {
		return "", err
	}
	return at, nil
}

// GetCafeAddr returns the cafe address is set
func (w *Wallet) GetCafeAddr() string {
	if w.cafeAddr == "" {
		return ""
	}
	return fmt.Sprintf("%s/api/%s", w.cafeAddr, ccafe.Version)
}
