package core

import (
	"bytes"
	"crypto/rand"
	"fmt"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/path"
	uio "gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/unixfs/io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/textileio/textile-go/mill"

	"github.com/textileio/textile-go/schema/textile"

	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
)

// Profile is an account-wide public profile
// NOTE: any account peer can publish profile entries to the same IPNS key
type Profile struct {
	Address   string   `json:"address"`
	Inboxes   []string `json:"inboxes,omitempty"`
	Username  string   `json:"username,omitempty"`
	AvatarUri string   `json:"avatar_uri,omitempty"`
}

// profileLifetime is the duration the ipns profile record will be considered valid
var profileLifetime = time.Hour * 24 * 7

// profileTTL is the duration the ipns profile record will be locally cached
var profileTTL = time.Hour

// Username returns profile username
func (t *Textile) Username() (*string, error) {
	return t.datastore.Profile().GetUsername()
}

// SetUsername updates profile with a new username
func (t *Textile) SetUsername(username string) error {
	if err := t.datastore.Profile().SetUsername(username); err != nil {
		return err
	}

	for _, thrd := range t.threads {
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
	thrd := t.ThreadByKey("avatar")
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

	thumb, err := t.AddFile(&mill.ImageResize{
		Opts: mill.ImageResizeOpts{
			Width:   thrd.Schema.Links["thumb"].Opts["width"],
			Quality: thrd.Schema.Links["thumb"].Opts["quality"],
		},
	}, AddFileConfig{
		Input:     input,
		Media:     file.Media,
		Plaintext: thrd.Schema.Links["thumb"].Plaintext,
	})
	if err != nil {
		return err
	}
	dir := Directory{"small": *small, "thumb": *thumb}

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
			profile.Inboxes = append(profile.Inboxes, ses.CafeId)
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
		profile.Inboxes = strings.Split(string(inboxesb), ",")
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
		entry, err := t.publishProfile(*prof)
		if err != nil {
			log.Errorf("error publishing profile: %s", err)
			return
		}
		log.Debugf("published: %s -> %s", entry.Name, entry.Value)
	}()
	return nil
}

// ResolveProfile looks up a profile on ipns
func (t *Textile) ResolveProfile(name peer.ID) (*path.Path, error) {
	return ipfs.Resolve(t.node, name)
}

// publishProfile publishes profile to ipns
func (t *Textile) publishProfile(prof Profile) (*ipfs.IpnsEntry, error) {
	dir := uio.NewDirectory(t.node.DAG)

	addressId, err := ipfs.AddDataToDirectory(t.node, dir, "address", bytes.NewReader([]byte(prof.Address)))
	if err != nil {
		return nil, err
	}

	var inboxesId *cid.Cid
	sessions := t.datastore.CafeSessions().List()
	if len(sessions) > 0 {
		var inboxes []string
		for _, ses := range t.datastore.CafeSessions().List() {
			inboxes = append(inboxes, ses.CafeId)
		}
		inboxesStr := strings.Join(inboxes, ",")
		inboxesId, err = ipfs.AddDataToDirectory(t.node, dir, "inboxes", bytes.NewReader([]byte(inboxesStr)))
		if err != nil {
			return nil, err
		}
	}

	var usernameId *cid.Cid
	if prof.Username != "" {
		usernameId, err = ipfs.AddDataToDirectory(t.node, dir, "username", bytes.NewReader([]byte(prof.Username)))
		if err != nil {
			return nil, err
		}
	}

	var avatarId *cid.Cid
	if prof.AvatarUri != "" {
		avatarId, err = ipfs.AddDataToDirectory(t.node, dir, "avatar_uri", bytes.NewReader([]byte(prof.AvatarUri)))
		if err != nil {
			return nil, err
		}
	}

	node, err := dir.GetNode()
	if err != nil {
		return nil, err
	}
	if err := ipfs.PinNode(t.node, node, false); err != nil {
		return nil, err
	}

	t.cafeOutbox.Add(addressId.Hash().B58String(), repo.CafeStoreRequest)
	if inboxesId != nil {
		t.cafeOutbox.Add(inboxesId.Hash().B58String(), repo.CafeStoreRequest)
	}
	if usernameId != nil {
		t.cafeOutbox.Add(usernameId.Hash().B58String(), repo.CafeStoreRequest)
	}
	if avatarId != nil {
		t.cafeOutbox.Add(avatarId.Hash().B58String(), repo.CafeStoreRequest)
	}
	t.cafeOutbox.Add(node.Cid().Hash().B58String(), repo.CafeStoreRequest)
	go t.cafeOutbox.Flush()

	value := node.Cid().Hash().B58String()
	return ipfs.Publish(t.node, t.node.PrivateKey, value, profileLifetime, profileTTL)
}
