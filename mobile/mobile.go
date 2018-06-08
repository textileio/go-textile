package mobile

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/central/models"
	tcore "github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/net"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/wallet"
	"github.com/textileio/textile-go/wallet/thread"
	"gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
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
	LogFiles      bool
}

// NewNode is the mobile entry point for creating a node
// NOTE: logLevel is one of: CRITICAL ERROR WARNING NOTICE INFO DEBUG
func NewNode(config *NodeConfig, messenger Messenger) (*Wrapper, error) {
	var m Mobile
	return m.NewNode(config, messenger)
}

// Mobile is the name of the framework (must match package name)
type Mobile struct{}

// Blocks is a wrapper around a list of Blocks, which makes decoding json from a little cleaner
// on the mobile side
type Blocks struct {
	Items []repo.Block `json:"items"`
}

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
		LogFiles:      config.LogFiles,
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
	online, err := tcore.Node.StartWallet()
	if err != nil {
		if err == wallet.ErrStopped {
			return nil
		}
		return err
	}

	go func() {
		<-online
		// join existing threads
		for _, thrd := range tcore.Node.Wallet.Threads() {
			w.subscribe(thrd)
		}

		// notify UI we're ready
		w.messenger.Notify(newEvent("onOnline", map[string]interface{}{}))

		// publish
		tcore.Node.Wallet.PublishThreads()
	}()

	return nil
}

// Stop the mobile node
func (w *Wrapper) Stop() error {
	if err := tcore.Node.StopWallet(); err != nil && err != wallet.ErrStopped {
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
	return tcore.Node.Wallet.SignUp(reg)
}

// SignIn build credentials and calls core SignIn
func (w *Wrapper) SignIn(username string, password string) error {
	// build creds
	creds := &models.Credentials{
		Username: username,
		Password: password,
	}
	return tcore.Node.Wallet.SignIn(creds)
}

// SignOut calls core SignOut
func (w *Wrapper) SignOut() error {
	return tcore.Node.Wallet.SignOut()
}

// IsSignedIn calls core IsSignedIn
func (w *Wrapper) IsSignedIn() bool {
	si, _ := tcore.Node.Wallet.IsSignedIn()
	return si
}

// GetUsername calls core GetUsername
func (w *Wrapper) GetUsername() (string, error) {
	return tcore.Node.Wallet.GetUsername()
}

// GetAccessToken calls core GetAccessToken
func (w *Wrapper) GetAccessToken() (string, error) {
	return tcore.Node.Wallet.GetAccessToken()
}

// AddPhoto adds a photo by path and shares it to the default thread
func (w *Wrapper) AddPhoto(path string, threadName string) (*net.MultipartRequest, error) {
	thrd := tcore.Node.Wallet.GetThreadByName(threadName)
	if thrd == nil {
		return nil, errors.New(fmt.Sprintf("could not find thread: %s", threadName))
	}
	added, err := tcore.Node.Wallet.AddPhoto(path)
	if err != nil {
		return nil, err
	}
	// TODO, fire off shared request to remote cluster
	_, err = thrd.AddPhoto(added.Id, "", added.Key)
	if err != nil {
		return nil, err
	}
	return added.RemoteRequest, nil
}

// SharePhoto adds an existing photo to a new thread
func (w *Wrapper) SharePhoto(id string, threadName string, caption string) (string, error) {
	block, err := tcore.Node.Wallet.FindBlock(id)
	if err != nil {
		return "", err
	}
	fromThreadId := libp2pc.ConfigEncodeKey(block.ThreadPubKey)
	fromThread := tcore.Node.Wallet.GetThread(fromThreadId)
	if fromThread == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", fromThreadId))
	}
	toThread := tcore.Node.Wallet.GetThreadByName(threadName)
	if toThread == nil {
		return "", errors.New(fmt.Sprintf("could not find thread named %s", threadName))
	}
	key, err := fromThread.Decrypt(block.TargetKey)
	if err != nil {
		return "", err
	}

	// TODO: owner challenge
	// TODO: fire off shared request to remote cluster
	shared, err := toThread.AddPhoto(id, caption, key)
	if err != nil {
		return "", err
	}
	return shared.Id, nil
}

// Get Photos returns core GetPhotos with json encoding
func (w *Wrapper) GetPhotos(offsetId string, limit int, threadName string) (string, error) {
	thrd := tcore.Node.Wallet.GetThreadByName(threadName)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("thread not found: %s", threadName))
	}

	if tcore.Node.Wallet.Online() {
		go thrd.Publish()
	}

	blocks := &Blocks{thrd.Blocks(offsetId, limit)}

	// gomobile does not allow slices. so, convert to json
	jsonb, err := json.Marshal(blocks)
	if err != nil {
		log.Errorf("error marshaling json: %s", err)
		return "", err
	}

	return string(jsonb), nil
}

// GetFileBase64 call core GetFileBase64
func (w *Wrapper) GetFileBase64(path string, blockId string) (string, error) {
	return tcore.Node.Wallet.GetFileBase64(path, blockId)
}

// GetIPFSPeerID returns the wallet's ipfs peer id
func (w *Wrapper) GetIPFSPeerID() (string, error) {
	return tcore.Node.Wallet.GetIPFSPeerID()
}

// PairDesktop publishes this nodes default thread key to a desktop node
// which is listening at it's own peer id
func (w *Wrapper) PairDesktop(pkb64 string) (string, error) {
	if !tcore.Node.Wallet.Online() {
		return "", wallet.ErrOffline
	}
	log.Info("pairing with desktop...")

	pkb, err := libp2pc.ConfigDecodeKey(pkb64)
	if err != nil {
		return "", err
	}
	pk, err := libp2pc.UnmarshalPublicKey(pkb)
	if err != nil {
		return "", err
	}

	// we invite the desktop to _read and write_ to our default album
	defaultThread := tcore.Node.Wallet.GetThreadByName("default")
	if defaultThread == nil {
		err = errors.New("default thread not found")
		log.Error(err.Error())
		return "", err
	}
	// encypt thread secret key with the desktop's pub key
	secret, err := defaultThread.PrivKey.Bytes()
	if err != nil {
		return "", err
	}
	secretcypher, err := crypto.Encrypt(pk, secret)
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
	// TODO: connect first
	err = tcore.Node.Wallet.Ipfs.Floodsub.Publish(topic, secretcypher)
	if err != nil {
		return "", err
	}
	log.Infof("published key phrase to desktop: %s", topic)

	return topic, nil
}

// subscribe to thread and pass updates to messenger
func (w *Wrapper) subscribe(thrd *thread.Thread) {
	datac := make(chan thread.Update)
	go thrd.Subscribe(datac)
	go func() {
		for {
			select {
			case update, ok := <-datac:
				if !ok {
					return
				}
				w.messenger.Notify(newEvent("onThreadUpdate", map[string]interface{}{
					"id":        update.Id,
					"thread":    update.Thread,
					"thread_id": update.ThreadID,
				}))
			}
		}
	}()
}
