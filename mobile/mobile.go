package mobile

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	zxcvbn "github.com/nbutton23/zxcvbn-go"

	"github.com/op/go-logging"

	"github.com/textileio/textile-go/central/models"
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
	return w.node.Start()
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

// TODO: doesn't use a cleaned version of the phone number, if pn is supplied
func (w *Wrapper) CheckPassword(password string, identity string) (bool, error) {
	match := zxcvbn.PasswordStrength(password, []string{identity})
	if match.Score < 3 {
		return false, errors.New(fmt.Sprintf("weak password - crackable in %s", match.CrackTimeDisplay))
	}
	return true, nil
}

func (w *Wrapper) SignUpWithEmail(username string, password string, email string, referral string) (int, *models.Response, error) {
	apiURL := ""

	reg := models.Registration{
		Username: username,
		Password: password,
		Identity: &models.Identity{
			Type:  models.EmailAddress,
			Value: email,
		},
		Referral: referral,
	}

	url := fmt.Sprintf("%s/api/v1/users", apiURL)
	payload, err := json.Marshal(reg)
	if err != nil {
		return 0, nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	resp := &models.Response{}
	if err := resp.Read(res.Body); err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resp, nil
}

func (w *Wrapper) SignIn(username string, password string) (int, *models.Response, error) {
	apiURL := ""

	creds := models.Credentials{
		Username: username,
		Password: password,
	}

	url := fmt.Sprintf("%s/api/v1/users", apiURL)
	payload, err := json.Marshal(creds)
	if err != nil {
		return 0, nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	resp := &models.Response{}
	if err := resp.Read(res.Body); err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resp, nil
}
