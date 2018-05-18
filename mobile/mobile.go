package mobile

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/op/go-logging"

	tcore "github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/net"

	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"

	"github.com/textileio/textile-go/central/models"
)

var log = logging.MustGetLogger("mobile")

// Message is a generic go -> bridge message structure
type Event struct {
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

// newEvent transforms an event name and structured data in Event
func newEvent(name string, payload map[string]interface{}) *Event {
	event := &Event{Name: name}
	jsonb, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("error creating event data json: %s", err)
	}
	event.Payload = string(jsonb)
	return event
}

// Messenger is used to inform the bridge layer of new data waiting to be queried
type Messenger interface {
	Notify(event *Event)
}

// Wrapper is the object exposed in the frameworks
type Wrapper struct {
	RepoPath  string
	messenger Messenger
}

// NodeConfig is used to configure the mobile node
type NodeConfig struct {
	RepoPath      string
	CentralApiURL string
	LogLevel      string
}

// NewNode is the mobile entry point for creating a node
// NOTE: logLevel is one of: CRITICAL ERROR WARNING NOTICE INFO DEBUG
func NewNode(config *NodeConfig, messenger Messenger) (*Wrapper, error) {
	var m Mobile
	return m.NewNode(config, messenger)
}

// Mobile is the name of the framework (must match package name)
type Mobile struct{}

// Create a gomobile compatible wrapper around TextileNode
func (m *Mobile) NewNode(config *NodeConfig, messenger Messenger) (*Wrapper, error) {
	ll, err := logging.LogLevel(config.LogLevel)
	if err != nil {
		ll = logging.INFO
	}
	cconfig := tcore.NodeConfig{
		RepoPath:      config.RepoPath,
		CentralApiURL: config.CentralApiURL,
		IsMobile:      true,
		LogLevel:      ll,
		LogFiles:      true,
	}
	node, err := tcore.NewNode(cconfig)
	if err != nil {
		return nil, err
	}
	tcore.Node = node

	return &Wrapper{RepoPath: config.RepoPath, messenger: messenger}, nil
}

// Start the mobile node
func (w *Wrapper) Start() error {
	if err := tcore.Node.Start(); err != nil {
		if err == tcore.ErrNodeRunning {
			return nil
		}
		return err
	}

	// join existing rooms
	for _, album := range tcore.Node.Datastore.Albums().GetAlbums("") {
		w.joinRoom(album.Id)
	}

	return nil
}

// Stop the mobile node
func (w *Wrapper) Stop() error {
	if err := tcore.Node.Stop(); err != nil && err != tcore.ErrNodeNotRunning {
		return err
	}
	return nil
}

// SignUpWithEmail creates an email based registration and calls core signup
func (w *Wrapper) SignUpWithEmail(username string, password string, email string, referral string) error {
	// build registration
	reg := &models.Registration{
		Username: username,
		Password: password,
		Identity: &models.Identity{
			Type:  models.EmailAddress,
			Value: email,
		},
		Referral: referral,
	}

	// signup
	return tcore.Node.SignUp(reg)
}

// SignIn build credentials and calls core SignIn
func (w *Wrapper) SignIn(username string, password string) error {
	// build creds
	creds := &models.Credentials{
		Username: username,
		Password: password,
	}

	// signin
	return tcore.Node.SignIn(creds)
}

// SignOut calls core SignOut
func (w *Wrapper) SignOut() error {
	return tcore.Node.SignOut()
}

// IsSignedIn calls core IsSignedIn
func (w *Wrapper) IsSignedIn() bool {
	si, _ := tcore.Node.IsSignedIn()
	return si
}

// Update thread allows the mobile client to choose the 'all' thread to subscribe
func (w *Wrapper) UpdateThread(mnemonic string, name string) error {
	if mnemonic == "" {
		return errors.New("mnemonic must not be empty")
	}

	ba := tcore.Node.Datastore.Albums().GetAlbumByName(name)
	if ba == nil {
		log.Debugf("creating album: %s", name)
		err := tcore.Node.CreateAlbum(mnemonic, name)
		if err != nil {
			log.Errorf("error creating album %s: %s", name, err)
			return err
		}
	} else {
		log.Debugf("removing album: %s", name)
		err := tcore.Node.DeleteAlbum(ba.Id)
		if err != nil {
			log.Errorf("error deleting album %s: %s", name, err)
			return err
		}
		return tcore.Node.CreateAlbum(mnemonic, name)
	}
	return nil
}

// GetUsername calls core GetUsername
func (w *Wrapper) GetUsername() (string, error) {
	return tcore.Node.GetUsername()
}

// GetAccessToken calls core GetAccessToken
func (w *Wrapper) GetAccessToken() (string, error) {
	return tcore.Node.GetAccessToken()
}

// GetGatewayPassword returns the current cookie value expected by the gateway
func (w *Wrapper) GetGatewayPassword() string {
	return tcore.Node.GatewayPassword
}

// AddPhoto calls core AddPhoto
func (w *Wrapper) AddPhoto(path string, thumb string, thread string) (*net.MultipartRequest, error) {
	return tcore.Node.AddPhoto(path, thumb, thread, "")
}

// SharePhoto calls core SharePhoto
func (w *Wrapper) SharePhoto(hash string, thread string, caption string) (*net.MultipartRequest, error) {
	return tcore.Node.SharePhoto(hash, thread, caption)
}

// GetHashRequest calls core GetHashRequest
func (w *Wrapper) GetHashRequest(hash string) (string, error) {
	request := tcore.Node.GetHashRequest(hash)

	// gomobile does not allow slices. so, convert to json
	jsonb, err := json.Marshal(request)
	if err != nil {
		log.Errorf("error marshaling json: %s", err)
		return "", err
	}

	return string(jsonb), nil
}

// Get Photos returns core GetPhotos with json encoding
func (w *Wrapper) GetPhotos(offsetId string, limit int, thread string) (string, error) {
	list := tcore.Node.GetPhotos(offsetId, limit, thread)
	if list == nil {
		list = &tcore.PhotoList{
			Hashes: make([]string, 0),
		}
	}

	// gomobile does not allow slices. so, convert to json
	jsonb, err := json.Marshal(list)
	if err != nil {
		log.Errorf("error marshaling json: %s", err)
		return "", err
	}

	return string(jsonb), nil
}

// GetFileBase64 call core GetFileBase64
func (w *Wrapper) GetFileBase64(path string) (string, error) {
	return tcore.Node.GetFileBase64(path)
}

// GetPeerID returns our peer id
func (w *Wrapper) GetPeerID() (string, error) {
	if !tcore.Node.Online() {
		return "", tcore.ErrNodeNotRunning
	}
	return tcore.Node.IpfsNode.Identity.Pretty(), nil
}

// PairDesktop publishes this nodes default album keys to a desktop node
// which is listening at it's own peer id
func (w *Wrapper) PairDesktop(pkb64 string) (string, error) {
	if !tcore.Node.Online() {
		return "", tcore.ErrNodeNotRunning
	}
	log.Info("pairing with desktop...")

	pkb, err := base64.StdEncoding.DecodeString(pkb64)
	if err != nil {
		log.Errorf("error decoding string: %s: %s", pkb64, err)
		return "", err
	}

	pk, err := libp2p.UnmarshalPublicKey(pkb)
	if err != nil {
		log.Errorf("error unmarshaling pub key: %s", err)
		return "", err
	}

	// the phrase will be used by the desktop client to create
	// the private key needed to decrypt photos
	// we invite the desktop to _read and write_ to our default album
	da := tcore.Node.Datastore.Albums().GetAlbumByName("default")
	if da == nil {
		err = errors.New("default album not found")
		log.Error(err.Error())
		return "", err
	}
	// encypt with the desktop's pub key
	cph, err := net.Encrypt(pk, []byte(da.Mnemonic))
	if err != nil {
		log.Errorf("encrypt failed: %s", err)
		return "", err
	}

	// get the topic to pair with from the pub key
	peerID, err := peer.IDFromPublicKey(pk)
	if err != nil {
		log.Errorf("id from public key failed: %s", err)
		return "", err
	}
	topic := peerID.Pretty()

	// finally, publish the encrypted phrase
	err = tcore.Node.IpfsNode.Floodsub.Publish(topic, cph)
	if err != nil {
		log.Errorf("publish %s failed: %s", topic, err)
		return "", err
	}
	log.Infof("published key phrase to desktop: %s", topic)

	// try a ping
	go func() {
		err = tcore.Node.PingPeer(topic, 1, make(chan string))
		if err != nil {
			log.Errorf("ping %s failed: %s", topic, err)
		}
	}()

	return topic, nil
}

// joinRoom and pass updates to messenger
func (w *Wrapper) joinRoom(id string) {
	datac := make(chan tcore.ThreadUpdate)
	go tcore.Node.JoinRoom(id, datac)
	go func() {
		for {
			select {
			case update, ok := <-datac:
				if !ok {
					return
				}
				w.messenger.Notify(newEvent("onThreadUpdate", map[string]interface{}{
					"cid":       update.Cid,
					"thread":    update.Thread,
					"thread_id": update.ThreadID,
				}))
			}
		}
	}()
}
