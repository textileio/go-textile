package mobile

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	tcore "github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/net"
)

type Wrapper struct {
	node   *tcore.TextileNode
	cancel context.CancelFunc
}
type Mobile struct{}

func NewTextile(repoPath string) *Wrapper {
	var m Mobile
	node, err := m.NewNode(repoPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return node
}

// Create a gomobile compatible wrapper around TextileNode
func (m *Mobile) NewNode(repoPath string) (*Wrapper, error) {
	node, err := tcore.NewNode(repoPath, true)
	if err != nil {
		return nil, err
	}

	return &Wrapper{node: node, cancel: node.Cancel}, nil
}

func (w *Wrapper) Start() error {
	return w.node.Start()
}

func (w *Wrapper) ConfigureDatastore(mnemonic string) error {
	return w.node.ConfigureDatastore(mnemonic)
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
