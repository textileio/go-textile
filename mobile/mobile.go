package mobile

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/op/go-logging"

	tcore "github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/net"

	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

var log = logging.MustGetLogger("mobile")

type Wrapper struct {
	RepoPath       string
	Cancel         context.CancelFunc
	node           *tcore.TextileNode
	gatewayRunning bool
}

func NewNode(repoPath string) (*Wrapper, error) {
	var m Mobile
	return m.NewNode(repoPath)
}

type Mobile struct{}

// Create a gomobile compatible wrapper around TextileNode
func (m *Mobile) NewNode(repoPath string) (*Wrapper, error) {
	node, err := tcore.NewNode(repoPath, true, logging.DEBUG)
	if err != nil {
		return nil, err
	}

	return &Wrapper{RepoPath: repoPath, Cancel: node.Cancel, node: node}, nil
}

func (w *Wrapper) Start() error {
	if err := w.node.Start(); err != nil {
		if err == tcore.ErrNodeRunning {
			return nil
		}
		return err
	}

	// join existing rooms
	albums := w.node.Datastore.Albums().GetAlbums("")
	datacs := make([]chan string, len(albums))
	for i, a := range albums {
		go w.node.JoinRoom(a.Id, datacs[i])
	}

	return nil
}

func (w *Wrapper) StartGateway() error {
	if w.gatewayRunning {
		return nil
	}
	if _, err := tcore.ServeHTTPGatewayProxy(w.node); err != nil {
		return err
	}
	w.gatewayRunning = true
	return nil
}

func (w *Wrapper) Stop() error {
	return w.node.Stop()
}

func (w *Wrapper) AddPhoto(path string, thumb string, thread string) (*net.MultipartRequest, error) {
	return w.node.AddPhoto(path, thumb, thread)
}

func (w *Wrapper) SharePhoto(hash string, thread string) (*net.MultipartRequest, error) {
	return w.node.SharePhoto(hash, thread)
}

func (w *Wrapper) GetPhotos(offsetId string, limit int, thread string) (string, error) {
	list := w.node.GetPhotos(offsetId, limit, thread)

	// gomobile does not allow slices. so, convert to json
	jsonb, err := json.Marshal(list)
	if err != nil {
		return "", err
	}

	return string(jsonb), nil
}

func (w *Wrapper) GetFileBase64(path string) (string, error) {
	b, err := w.node.GetFile(path, nil)
	if err != nil {
		return "error", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (w *Wrapper) GetPeerID() (string, error) {
	if w.node.IpfsNode == nil {
		return "", errors.New("node not started")
	}
	return w.node.IpfsNode.Identity.Pretty(), nil
}

func (w *Wrapper) PairDesktop(pkb64 string) (string, error) {
	log.Info("pairing with desktop...")
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
	// we invite the desktop to _read and write_ to our default album
	da := w.node.Datastore.Albums().GetAlbumByName("default")
	if da == nil {
		err = errors.New("default album not found")
		log.Error(err.Error())
		return "", err
	}
	// encypt with the desktop's pub key
	cph, err := net.Encrypt(pk, []byte(da.Mnemonic))
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
	log.Infof("published key phrase to desktop: %s", topic)

	// try a ping
	err = w.node.PingPeer(topic, 1, make(chan string))
	if err != nil {
		log.Errorf("ping %s failed: %s", topic, err)
	}

	return topic, nil
}

func (w *Wrapper) SignUpWithEmail(username string, password string, email string, referral string) ([]byte, error) {
	res, err := w.node.SignUpWithEmail(username, password, email, referral)
	if err != nil {
		log.Errorf("signup error: %s", err)
		return nil, err
	}
	jsn, err := json.Marshal(res)
	if err != nil {
		log.Errorf("signup data conversion failed: %s", err)
		return nil, err
	}
	return jsn, nil
}

func (w *Wrapper) SignIn(username string, password string) ([]byte, error) {
	res, err := w.node.SignIn(username, password)
	if err != nil {
		log.Errorf("signin failed: %s", err)
		return nil, err
	}
	jsn, err := json.Marshal(res)
	if err != nil {
		log.Errorf("signin data conversion failed: %s", err)
		return nil, err
	}
	return jsn, nil
}
