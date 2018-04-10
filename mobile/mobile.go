package mobile

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"

	"fmt"
	tcore "github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/net"

	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

type Wrapper struct {
	RepoPath string
	Cancel   context.CancelFunc
	node     *tcore.TextileNode
}

func NewNode(repoPath string) (*Wrapper, error) {
	var m Mobile
	return m.NewNode(repoPath)
}

type Mobile struct{}

// Create a gomobile compatible wrapper around TextileNode
func (m *Mobile) NewNode(repoPath string) (*Wrapper, error) {
	node, err := tcore.NewNode(repoPath, true)
	if err != nil {
		return nil, err
	}

	return &Wrapper{RepoPath: repoPath, Cancel: node.Cancel, node: node}, nil
}

func (w *Wrapper) Start() error {
	return w.node.Start()
}

func (w *Wrapper) ConfigureDatastore(mnemonic string) error {
	return w.node.ConfigureDatastore(mnemonic, "")
}

func (w *Wrapper) IsDatastoreConfigured() bool {
	_, err := w.GetRecoveryPhrase()
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			fmt.Printf("error checking if datastore is configured: %s", err)
			return false
		}
	}
	return true
}

func (w *Wrapper) Stop() error {
	return w.node.Stop()
}

func (w *Wrapper) AddPhoto(path string, thumb string) (*net.MultipartRequest, error) {
	return w.node.AddPhoto(path, thumb)
}

func (w *Wrapper) GetPhotos(offsetId string, limit int) (string, error) {
	list := w.node.GetPhotos(offsetId, limit)

	// gomobile does not allow slices. so, convert to json
	jsonb, err := json.Marshal(list)
	if err != nil {
		return "", err
	}

	return string(jsonb), nil
}

func (w *Wrapper) GetFileBase64(path string) (string, error) {
	b, err := w.node.GetFile(path)
	if err != nil {
		return "error", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// provides the user's recovery phrase for their private key
func (w *Wrapper) GetRecoveryPhrase() (string, error) {
	return w.node.Datastore.Config().GetMnemonic()
}

func (w *Wrapper) PairDesktop(pkb64 string) (string, error) {
	pkb, err := base64.StdEncoding.DecodeString(pkb64)
	if err != nil {
		return "", err
	}

	pk, err := libp2p.UnmarshalPublicKey(pkb)
	if err != nil {
		return "", err
	}

	// the phrase will be used by the desktop client to create
	// the private key needed to decrypt photos
	ph, err := w.GetRecoveryPhrase()
	if err != nil {
		return "", err
	}
	// encypt with the desktop's pub key
	cph, err := net.Encrypt(pk, []byte(ph))
	if err != nil {
		return "", err
	}

	// get the topic to pair with from the pub key
	peerID, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return "", err
	}
	topic := peerID.Pretty()

	// finally, publish the encrypted phrase
	err = w.node.IpfsNode.Floodsub.Publish(topic, cph)
	if err != nil {
		return "", err
	}

	return topic, nil
}
