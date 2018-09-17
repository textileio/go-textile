package wallet

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/textileio/textile-go/util"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/namesys/opts"
	"gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/path"
	uio "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/unixfs/io"
	"time"
)

type Profile struct {
	Id       string `json:"id"`
	Username string `json:"username,omitempty"`
	AvatarId string `json:"avatar_id,omitempty"`
}

var profileTTL = time.Hour * 24 * 7 * 4
var profileCacheTTL = time.Hour * 24 * 7

// GetId returns profile id
func (w *Wallet) GetId() (*peer.ID, error) {
	if err := w.touchDatastore(); err != nil {
		return nil, err
	}
	key, err := w.datastore.Profile().GetKey()
	if err != nil {
		return nil, err
	}
	id, err := peer.IDFromPrivateKey(key)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// GetKey returns profile master secret key
func (w *Wallet) GetKey() (libp2pc.PrivKey, error) {
	if err := w.touchDatastore(); err != nil {
		return nil, err
	}
	return w.datastore.Profile().GetKey()
}

// GetUsername returns profile username
func (w *Wallet) GetUsername() (*string, error) {
	if err := w.touchDatastore(); err != nil {
		return nil, err
	}
	return w.datastore.Profile().GetUsername()
}

// SetUsername updates profile with a new username
func (w *Wallet) SetUsername(username string) error {
	if err := w.touchDatastore(); err != nil {
		return err
	}

	// update
	if err := w.datastore.Profile().SetUsername(username); err != nil {
		return err
	}

	go func() {
		<-w.Online()

		// publish
		pid, err := w.GetId()
		if err != nil {
			log.Errorf("error getting id (set username): %s", err)
			return
		}
		prof, err := w.GetProfile(pid.Pretty())
		if err != nil {
			log.Errorf("error getting profile (set username): %s", err)
			return
		}
		if _, err := w.PublishProfile(prof); err != nil {
			log.Errorf("error publishing profile (set username): %s", err)
			return
		}
	}()
	return nil
}

// GetAvatarId returns profile avatar id
func (w *Wallet) GetAvatarId() (*string, error) {
	if err := w.touchDatastore(); err != nil {
		return nil, err
	}
	return w.datastore.Profile().GetAvatarId()
}

// SetAvatarId updates profile with a new avatar
func (w *Wallet) SetAvatarId(id string) error {
	if err := w.touchDatastore(); err != nil {
		return err
	}

	// get the public key for this photo
	key, err := w.GetPhotoKey(id)
	if err != nil {
		return err
	}

	// use the cafe address w/ public url
	link := fmt.Sprintf("/ipfs/%s/thumb?key=%s", id, key)

	// update
	if err := w.datastore.Profile().SetAvatarId(link); err != nil {
		return err
	}

	go func() {
		<-w.Online()

		// publish
		pid, err := w.GetId()
		if err != nil {
			log.Errorf("error getting id (set avatar): %s", err)
			return
		}
		prof, err := w.GetProfile(pid.Pretty())
		if err != nil {
			log.Errorf("error getting profile (set avatar): %s", err)
			return
		}
		if _, err := w.PublishProfile(prof); err != nil {
			log.Errorf("error publishing profile (set avatar): %s", err)
			return
		}
	}()
	return nil
}

// GetProfile return a model representation of an ipns profile
func (w *Wallet) GetProfile(peerId string) (*Profile, error) {
	profile := &Profile{Id: peerId}

	// if peer id is local, return profile from db
	pid, err := w.GetId()
	if err != nil {
		return nil, err
	}
	if pid.Pretty() == peerId {
		username, err := w.GetUsername()
		if err != nil {
			return nil, err
		}
		if username != nil {
			profile.Username = *username
		}
		avatarId, err := w.GetAvatarId()
		if err != nil {
			return nil, err
		}
		if avatarId != nil {
			profile.AvatarId = *avatarId
		}
		return profile, nil
	}

	// resolve profile at peer id
	entry, err := w.ResolveProfile(peerId)
	if err != nil {
		return nil, err
	}
	root := entry.String()

	// get components from entry
	var usernameb, avatarIdb []byte
	usernameb, _ = util.GetDataAtPath(w.ipfs, fmt.Sprintf("%s/%s", root, "username"))
	if usernameb != nil {
		profile.Username = string(usernameb)
	}
	avatarIdb, _ = util.GetDataAtPath(w.ipfs, fmt.Sprintf("%s/%s", root, "avatar_id"))
	if avatarIdb != nil {
		profile.AvatarId = string(avatarIdb)
	}
	return profile, nil
}

// ResolveProfile looks up a profile on ipns
func (w *Wallet) ResolveProfile(name string) (*path.Path, error) {
	if !w.IsOnline() {
		return nil, ErrOffline
	}

	// setup query
	name = fmt.Sprintf("/ipns/%s", name)
	var ropts []nsopts.ResolveOpt
	ropts = append(ropts, nsopts.Depth(1))
	ropts = append(ropts, nsopts.DhtRecordCount(4))
	ropts = append(ropts, nsopts.DhtTimeout(5))

	pth, err := w.ipfs.Namesys.Resolve(w.ipfs.Context(), name, ropts...)
	if err != nil {
		return nil, err
	}
	return &pth, nil
}

// PublishProfile publishes profile to ipns
func (w *Wallet) PublishProfile(prof *Profile) (*util.IpnsEntry, error) {
	if !w.IsOnline() {
		return nil, ErrOffline
	}
	if w.ipfs.Mounts.Ipns != nil && w.ipfs.Mounts.Ipns.IsActive() {
		return nil, errors.New("cannot manually publish while IPNS is mounted")
	}

	// if nil profile, use current
	if prof == nil {
		pid, err := w.GetId()
		if err != nil {
			return nil, err
		}
		prof, err = w.GetProfile(pid.Pretty())
		if err != nil {
			return nil, err
		}
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
			log.Warningf("pin request exists: %s", id)
		}
	}()

	// extract path
	pth, err := path.ParsePath(id)
	if err != nil {
		return nil, err
	}

	// load our private key
	sk, err := w.GetKey()
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
