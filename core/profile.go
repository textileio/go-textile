package core

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"gx/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7/go-cid"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	libp2pc "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"
	"gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/core/coreapi/interface"
	uio "gx/ipfs/QmfB3oNXGGq9S4B2a9YeCajoATms3Zw2VvDm8fK7VeLSV8/go-unixfs/io"
	"io/ioutil"

	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/mill"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/schema/textile"
)

// Profile is an account-wide public profile
// NOTE: any account peer can publish profile entries to the same IPNS key
type Profile struct {
	Address   string      `json:"address"`
	Inboxes   []repo.Cafe `json:"inboxes,omitempty"`
	Username  string      `json:"username,omitempty"`
	AvatarUri string      `json:"avatar_uri,omitempty"`
}

// Username returns profile username
func (t *Textile) Username() (*string, error) {
	return t.datastore.Profile().GetUsername()
}

// SetUsername updates profile with a new username
func (t *Textile) SetUsername(username string) error {
	if err := t.datastore.Profile().SetUsername(username); err != nil {
		return err
	}

	for _, thrd := range t.loadedThreads {
		if _, err := thrd.annouce(); err != nil {
			return err
		}
	}

	return t.PublishProfile()
}

// Avatar returns profile avatar
func (t *Textile) Avatar() (*string, error) {
	return t.datastore.Profile().GetAvatar()
}

// SetAvatar updates profile with a new avatar at the given file hash.
func (t *Textile) SetAvatar(hash string) error {
	data, file, err := t.FileData(hash)
	if err != nil {
		return err
	}
	input, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	// create a plaintext files thread for tracking avatars
	thrd := t.ThreadByKey("avatars")
	if thrd == nil {
		sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
		if err != nil {
			return err
		}

		sf, err := t.AddSchema(textile.Avatars, "avatars")
		if err != nil {
			return err
		}
		shash, err := mh.FromB58String(sf.Hash)
		if err != nil {
			return err
		}

		thrd, err = t.AddThread(sk, AddThreadConfig{
			Key:       "avatars",
			Name:      "avatars",
			Schema:    shash,
			Initiator: t.account.Address(),
			Type:      repo.PrivateThread,
			Join:      true,
		})
		if err != nil {
			return err
		}
	}

	large, err := t.AddFile(&mill.ImageResize{
		Opts: mill.ImageResizeOpts{
			Width:   thrd.Schema.Links["large"].Opts["width"],
			Quality: thrd.Schema.Links["large"].Opts["quality"],
		},
	}, AddFileConfig{
		Input:     input,
		Media:     file.Media,
		Plaintext: thrd.Schema.Links["large"].Plaintext,
	})
	if err != nil {
		return err
	}

	small, err := t.AddFile(&mill.ImageResize{
		Opts: mill.ImageResizeOpts{
			Width:   thrd.Schema.Links["small"].Opts["width"],
			Quality: thrd.Schema.Links["small"].Opts["quality"],
		},
	}, AddFileConfig{
		Input:     input,
		Media:     file.Media,
		Plaintext: thrd.Schema.Links["small"].Plaintext,
	})
	if err != nil {
		return err
	}

	dir := Directory{"large": *large, "small": *small}
	node, keys, err := t.AddNodeFromDirs([]Directory{dir})
	if err != nil {
		return err
	}

	if _, err := thrd.AddFiles(node, "", keys); err != nil {
		return err
	}

	uri := fmt.Sprintf("/ipfs/%s", node.Cid().Hash().B58String())
	if err := t.datastore.Profile().SetAvatar(uri); err != nil {
		return err
	}

	return t.PublishProfile()
}

// Profile return a model representation of an ipns profile
func (t *Textile) Profile(pid peer.ID) (*Profile, error) {
	profile := &Profile{}

	// if peer id is local, return profile from db
	if t.node.Identity.Pretty() == pid.Pretty() {
		profile.Address = t.account.Address()
		for _, ses := range t.datastore.CafeSessions().List() {
			profile.Inboxes = append(profile.Inboxes, ses.Cafe)
		}
		username, err := t.Username()
		if err != nil {
			return nil, err
		}
		if username != nil {
			profile.Username = *username
		}
		avatar, err := t.Avatar()
		if err != nil {
			return nil, err
		}
		if avatar != nil {
			profile.AvatarUri = *avatar
		}
		return profile, nil
	}

	entry, err := t.ResolveProfile(pid)
	if err != nil {
		return nil, err
	}
	root := entry.String()

	var addrb, inboxesb, usernameb, avatarb []byte
	addrb, _ = ipfs.DataAtPath(t.node, fmt.Sprintf("%s/%s", root, "address"))
	if addrb != nil {
		profile.Address = string(addrb)
	}
	inboxesb, _ = ipfs.DataAtPath(t.node, fmt.Sprintf("%s/%s", root, "inboxes"))
	if inboxesb != nil && string(inboxesb) != "" {
		var list []repo.Cafe
		if err := json.Unmarshal(inboxesb, &list); err != nil {
			return nil, err
		}
		profile.Inboxes = list
	}
	usernameb, _ = ipfs.DataAtPath(t.node, fmt.Sprintf("%s/%s", root, "username"))
	if usernameb != nil {
		profile.Username = string(usernameb)
	}
	avatarb, _ = ipfs.DataAtPath(t.node, fmt.Sprintf("%s/%s", root, "avatar_uri"))
	if avatarb != nil {
		profile.AvatarUri = string(avatarb)
	}
	return profile, nil
}

// PublishProfile publishes the current profile
func (t *Textile) PublishProfile() error {
	prof, err := t.Profile(t.node.Identity)
	if err != nil {
		return err
	}
	if prof == nil {
		return nil
	}

	go func() {
		<-t.OnlineCh()
		if err := t.publishProfile(*prof); err != nil {
			log.Errorf("error publishing profile: %s", err)
			return
		}
	}()
	return nil
}

// ResolveProfile looks up a profile on ipns
func (t *Textile) ResolveProfile(name peer.ID) (iface.Path, error) {
	return ipfs.ResolveIPNS(t.node, name)
}

// publishProfile publishes profile to ipns
func (t *Textile) publishProfile(prof Profile) error {
	dir := uio.NewDirectory(t.node.DAG)

	addressId, err := ipfs.AddDataToDirectory(t.node, dir, "address", bytes.NewReader([]byte(prof.Address)))
	if err != nil {
		return err
	}

	var inboxesId *cid.Cid
	sessions := t.datastore.CafeSessions().List()
	if len(sessions) > 0 {
		var inboxes []repo.Cafe
		for _, ses := range t.datastore.CafeSessions().List() {
			inboxes = append(inboxes, ses.Cafe)
		}
		inboxesb, err := json.Marshal(inboxes)
		if err != nil {
			return nil, err
		}
		inboxesId, err = ipfs.AddDataToDirectory(t.node, dir, "inboxes", bytes.NewReader(inboxesb))
		if err != nil {
			return err
		}
	}

	var usernameId *cid.Cid
	if prof.Username != "" {
		usernameId, err = ipfs.AddDataToDirectory(t.node, dir, "username", bytes.NewReader([]byte(prof.Username)))
		if err != nil {
			return err
		}
	}

	var avatarId *cid.Cid
	if prof.AvatarUri != "" {
		avatarId, err = ipfs.AddDataToDirectory(t.node, dir, "avatar_uri", bytes.NewReader([]byte(prof.AvatarUri)))
		if err != nil {
			return err
		}
	}

	node, err := dir.GetNode()
	if err != nil {
		return err
	}
	if err := ipfs.PinNode(t.node, node, false); err != nil {
		return err
	}

	if err := t.cafeOutbox.Add(addressId.Hash().B58String(), repo.CafeStoreRequest); err != nil {
		return err
	}
	if inboxesId != nil {
		if err := t.cafeOutbox.Add(inboxesId.Hash().B58String(), repo.CafeStoreRequest); err != nil {
			return err
		}
	}
	if usernameId != nil {
		if err := t.cafeOutbox.Add(usernameId.Hash().B58String(), repo.CafeStoreRequest); err != nil {
			return err
		}
	}
	if avatarId != nil {
		if err := t.cafeOutbox.Add(avatarId.Hash().B58String(), repo.CafeStoreRequest); err != nil {
			return err
		}
	}
	if err := t.cafeOutbox.Add(node.Cid().Hash().B58String(), repo.CafeStoreRequest); err != nil {
		return err
	}
	go t.cafeOutbox.Flush()

	return t.cafeService.PublishContact(node.Cid().Hash().B58String())
}
