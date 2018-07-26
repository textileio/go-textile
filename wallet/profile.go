package wallet

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	cmodels "github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/core/cafe"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
	"github.com/textileio/textile-go/wallet/model"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/namesys/opts"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/path"
	uio "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/unixfs/io"
	"time"
)

var profileTTL = time.Hour * 24 * 7 * 4
var profileCacheTTL = time.Hour * 24 * 7

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

	// re-pub profile
	go func() {
		if _, err := w.PublishProfile(); err != nil {
			log.Errorf("error publishing profile: %s", err)
		}
	}()

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

	// re-pub profile
	go func() {
		if _, err := w.PublishProfile(); err != nil {
			log.Errorf("error publishing profile: %s", err)
		}
	}()

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

// GetProfile return a model representation of a peer profile
func (w *Wallet) GetProfile() (*model.Profile, error) {
	id, err := w.GetId()
	if err != nil {
		return nil, err
	}
	username, _ := w.GetUsername()

	return &model.Profile{
		Id:       id,
		Username: username,
		AvatarId: "", // TODO: what goes here? choose photo for avatar pic?
	}, nil
}

// ResolveProfile looks up a peer's profile on ipns
func (w *Wallet) ResolveProfile(name string) (path.Path, error) {
	name = fmt.Sprintf("/ipns/%s", name)
	var ropts []nsopts.ResolveOpt
	ropts = append(ropts, nsopts.Depth(1))
	ropts = append(ropts, nsopts.DhtRecordCount(4))
	ropts = append(ropts, nsopts.DhtTimeout(5))

	return w.ipfs.Namesys.Resolve(w.ipfs.Context(), name, ropts...)
}

// PublishProfile publishes the peer profile to ipns
func (w *Wallet) PublishProfile() (*util.IpnsEntry, error) {
	if !w.IsOnline() {
		return nil, ErrOffline
	}

	if w.ipfs.Mounts.Ipns != nil && w.ipfs.Mounts.Ipns.IsActive() {
		return nil, errors.New("cannot manually publish while IPNS is mounted")
	}

	// get current profile
	prof, err := w.GetProfile()
	if err != nil {
		return nil, err
	}

	// create a virtual directory for the photo
	dirb := uio.NewDirectory(w.ipfs.DAG)
	if err := util.AddFileToDirectory(w.ipfs, dirb, bytes.NewReader([]byte(prof.Id)), "id"); err != nil {
		return nil, err
	}
	if err := util.AddFileToDirectory(w.ipfs, dirb, bytes.NewReader([]byte(prof.Username)), "username"); err != nil {
		return nil, err
	}
	if err := util.AddFileToDirectory(w.ipfs, dirb, bytes.NewReader([]byte(prof.AvatarId)), "avatar_id"); err != nil {
		return nil, err
	}

	// pin the directory locally
	dir, err := dirb.GetNode()
	if err != nil {
		return nil, err
	}
	if err := util.PinDirectory(w.ipfs, dir, []string{}); err != nil {
		return nil, err
	}
	id := dir.Cid().Hash().B58String()

	// request cafe pin
	go func() {
		if err := w.putPinRequest(id); err != nil {
			// TODO: #202 (Properly handle database/sql errors)
			log.Warningf("pin request %s exists: %s", id)
		}
	}()

	// extract path
	pth, err := path.ParsePath(id)
	if err != nil {
		return nil, err
	}

	// load our private key
	sk, err := w.GetPrivKey()
	if err != nil {
		return nil, err
	}

	// finish
	popts := &util.PublishOpts{
		VerifyExists: true,
		PubValidTime: profileCacheTTL,
	}
	ctx := context.WithValue(w.ipfs.Context(), "ipns-publish-ttl", profileTTL)
	entry, err := util.Publish(ctx, w.ipfs, sk, pth, popts)
	if err != nil {
		return nil, err
	}

	log.Debugf("updated profile: %s -> %s", entry.Name, entry.Value)

	return entry, nil
}
