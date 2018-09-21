package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/textileio/textile-go/ipfs"
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
func (t *Textile) GetId() (*peer.ID, error) {
	if err := t.touchDatastore(); err != nil {
		return nil, err
	}
	id, err := t.datastore.Config().GetId()
	if err != nil {
		return nil, err
	}
	if id == nil {
		return nil, ErrProfileNotFound
	}
	pid, err := peer.IDFromString(*id)
	if err != nil {
		return nil, err
	}
	return &pid, nil
}

// GetKey returns profile master secret key
func (t *Textile) GetKey() (libp2pc.PrivKey, error) {
	if err := t.touchDatastore(); err != nil {
		return nil, err
	}
	return t.datastore.Config().GetKey()
}

// GetUsername returns profile username
func (t *Textile) GetUsername() (*string, error) {
	if err := t.touchDatastore(); err != nil {
		return nil, err
	}
	return t.datastore.Profile().GetUsername()
}

// SetUsername updates profile with a new username
func (t *Textile) SetUsername(username string) error {
	if err := t.touchDatastore(); err != nil {
		return err
	}

	// update
	if err := t.datastore.Profile().SetUsername(username); err != nil {
		return err
	}

	go func() {
		<-t.Online()

		// publish
		pid, err := t.GetId()
		if err != nil {
			log.Errorf("error getting id (set username): %s", err)
			return
		}
		prof, err := t.GetProfile(pid.Pretty())
		if err != nil {
			log.Errorf("error getting profile (set username): %s", err)
			return
		}
		if _, err := t.PublishProfile(prof); err != nil {
			log.Errorf("error publishing profile (set username): %s", err)
			return
		}
	}()
	return nil
}

// GetAvatarId returns profile avatar id
func (t *Textile) GetAvatarId() (*string, error) {
	if err := t.touchDatastore(); err != nil {
		return nil, err
	}
	return t.datastore.Profile().GetAvatarId()
}

// SetAvatarId updates profile with a new avatar
func (t *Textile) SetAvatarId(id string) error {
	if err := t.touchDatastore(); err != nil {
		return err
	}

	// get the public key for this photo
	key, err := t.GetPhotoKey(id)
	if err != nil {
		return err
	}

	// use the cafe address w/ public url
	link := fmt.Sprintf("/ipfs/%s/thumb?key=%s", id, key)

	// update
	if err := t.datastore.Profile().SetAvatarId(link); err != nil {
		return err
	}

	go func() {
		<-t.Online()

		// publish
		pid, err := t.GetId()
		if err != nil {
			log.Errorf("error getting id (set avatar): %s", err)
			return
		}
		prof, err := t.GetProfile(pid.Pretty())
		if err != nil {
			log.Errorf("error getting profile (set avatar): %s", err)
			return
		}
		if _, err := t.PublishProfile(prof); err != nil {
			log.Errorf("error publishing profile (set avatar): %s", err)
			return
		}
	}()
	return nil
}

// GetProfile return a model representation of an ipns profile
func (t *Textile) GetProfile(peerId string) (*Profile, error) {
	profile := &Profile{Id: peerId}

	// if peer id is local, return profile from db
	pid, err := t.GetId()
	if err != nil {
		return nil, err
	}
	if pid.Pretty() == peerId {
		username, err := t.GetUsername()
		if err != nil {
			return nil, err
		}
		if username != nil {
			profile.Username = *username
		}
		avatarId, err := t.GetAvatarId()
		if err != nil {
			return nil, err
		}
		if avatarId != nil {
			profile.AvatarId = *avatarId
		}
		return profile, nil
	}

	// resolve profile at peer id
	entry, err := t.ResolveProfile(peerId)
	if err != nil {
		return nil, err
	}
	root := entry.String()

	// get components from entry
	var usernameb, avatarIdb []byte
	usernameb, _ = ipfs.GetDataAtPath(t.ipfs, fmt.Sprintf("%s/%s", root, "username"))
	if usernameb != nil {
		profile.Username = string(usernameb)
	}
	avatarIdb, _ = ipfs.GetDataAtPath(t.ipfs, fmt.Sprintf("%s/%s", root, "avatar_id"))
	if avatarIdb != nil {
		profile.AvatarId = string(avatarIdb)
	}
	return profile, nil
}

// ResolveProfile looks up a profile on ipns
func (t *Textile) ResolveProfile(name string) (*path.Path, error) {
	if !t.IsOnline() {
		return nil, ErrOffline
	}

	// setup query
	name = fmt.Sprintf("/ipns/%s", name)
	var ropts []nsopts.ResolveOpt
	ropts = append(ropts, nsopts.Depth(1))
	ropts = append(ropts, nsopts.DhtRecordCount(4))
	ropts = append(ropts, nsopts.DhtTimeout(5))

	pth, err := t.ipfs.Namesys.Resolve(t.ipfs.Context(), name, ropts...)
	if err != nil {
		return nil, err
	}
	return &pth, nil
}

// PublishProfile publishes profile to ipns
func (t *Textile) PublishProfile(prof *Profile) (*ipfs.IpnsEntry, error) {
	if !t.IsOnline() {
		return nil, ErrOffline
	}
	if t.ipfs.Mounts.Ipns != nil && t.ipfs.Mounts.Ipns.IsActive() {
		return nil, errors.New("cannot manually publish while IPNS is mounted")
	}

	// if nil profile, use current
	if prof == nil {
		pid, err := t.GetId()
		if err != nil {
			return nil, err
		}
		prof, err = t.GetProfile(pid.Pretty())
		if err != nil {
			return nil, err
		}
	}

	// create a virtual directory for the photo
	dirb := uio.NewDirectory(t.ipfs.DAG)
	if err := ipfs.AddFileToDirectory(t.ipfs, dirb, bytes.NewReader([]byte(prof.Id)), "id"); err != nil {
		return nil, err
	}
	if err := ipfs.AddFileToDirectory(t.ipfs, dirb, bytes.NewReader([]byte(prof.Username)), "username"); err != nil {
		return nil, err
	}
	if err := ipfs.AddFileToDirectory(t.ipfs, dirb, bytes.NewReader([]byte(prof.AvatarId)), "avatar_id"); err != nil {
		return nil, err
	}

	// pin the directory locally
	dir, err := dirb.GetNode()
	if err != nil {
		return nil, err
	}
	if err := ipfs.PinDirectory(t.ipfs, dir, []string{}); err != nil {
		return nil, err
	}
	id := dir.Cid().Hash().B58String()

	// request cafe pin
	go func() {
		if err := t.putPinRequest(id); err != nil {
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
	sk, err := t.GetKey()
	if err != nil {
		return nil, err
	}

	// finish
	popts := &ipfs.PublishOpts{
		VerifyExists: true,
		PubValidTime: profileCacheTTL,
	}
	ctx := context.WithValue(t.ipfs.Context(), "ipns-publish-ttl", profileTTL)
	entry, err := ipfs.Publish(ctx, t.ipfs, sk, pth, popts)
	if err != nil {
		return nil, err
	}

	log.Debugf("updated profile: %s -> %s", entry.Name, entry.Value)

	return entry, nil
}
