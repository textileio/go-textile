package mobile

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/cafe/models"
	tcore "github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
	"github.com/textileio/textile-go/wallet"
	"github.com/textileio/textile-go/wallet/model"
	"github.com/textileio/textile-go/wallet/thread"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"time"
)

var log = logging.MustGetLogger("mobile")

// Message is a generic go -> bridge message structure
type Event struct {
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

// Messenger is used to inform the bridge layer of new data waiting to be queried
type Messenger interface {
	Notify(event *Event)
}

// NodeConfig is used to configure the mobile node
// NOTE: logLevel is one of: CRITICAL ERROR WARNING NOTICE INFO DEBUG
type NodeConfig struct {
	RepoPath string
	CafeAddr string
	LogLevel string
	LogFiles bool
}

// Mobile is the name of the framework (must match package name)
type Mobile struct {
	RepoPath  string
	Mnemonic  string
	messenger Messenger
}

// Thread is a simple meta data wrapper around a Thread
type Thread struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Peers int    `json:"peers"`
}

// Threads is a wrapper around a list of Threads
type Threads struct {
	Items []Thread `json:"items"`
}

// Device is a simple meta data wrapper around a Device
type Device struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Devices is a wrapper around a list of Devices
type Devices struct {
	Items []Device `json:"items"`
}

// Photo is a simple meta data wrapper around a photo block
type Photo struct {
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	AuthorId string    `json:"author_id"`
	Caption  string    `json:"caption"`
}

// Photos is a wrapper around a list of photos
type Photos struct {
	Items []Photo `json:"items"`
}

// ImageData is a wrapper around an image data url and meta data
type ImageData struct {
	Url      string               `json:"url"`
	Metadata *model.PhotoMetadata `json:"metadata"`
}

// ExternalInvite is a wrapper around an invite id and key
type ExternalInvite struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	Inviter string `json:"inviter"`
}

// Create a gomobile compatible wrapper around TextileNode
func NewNode(config *NodeConfig, messenger Messenger) (*Mobile, error) {
	ll, err := logging.LogLevel(config.LogLevel)
	if err != nil {
		ll = logging.INFO
	}
	cconfig := tcore.NodeConfig{
		LogLevel: ll,
		LogFiles: config.LogFiles,
		WalletConfig: wallet.Config{
			RepoPath: config.RepoPath,
			IsMobile: true,
			CafeAddr: config.CafeAddr,
		},
	}
	node, mnemonic, err := tcore.NewNode(cconfig)
	if err != nil {
		return nil, err
	}
	tcore.Node = node

	return &Mobile{RepoPath: config.RepoPath, Mnemonic: mnemonic, messenger: messenger}, nil
}

// Start the mobile node
func (m *Mobile) Start() error {
	if err := tcore.Node.StartWallet(); err != nil {
		if err == wallet.ErrStarted {
			return nil
		}
		return err
	}

	go func() {
		<-tcore.Node.Wallet.Online()
		// subscribe to thread updates
		for _, thrd := range tcore.Node.Wallet.Threads() {
			go func(t *thread.Thread) {
				m.subscribe(t)
			}(thrd)
		}

		// subscribe to wallet updates
		go func() {
			for {
				select {
				case update, ok := <-tcore.Node.Wallet.Updates():
					if !ok {
						return
					}
					payload, err := toJSON(update)
					if err != nil {
						return
					}
					var name string
					switch update.Type {
					case wallet.ThreadAdded:
						// subscribe to updates
						if _, thrd := tcore.Node.Wallet.GetThread(update.Id); thrd != nil {
							go m.subscribe(thrd)
						}
						name = "onThreadAdded"

					case wallet.ThreadRemoved:
						name = "onThreadRemoved"

					case wallet.DeviceAdded:
						name = "onDeviceAdded"

					case wallet.DeviceRemoved:
						name = "onDeviceRemoved"
					}
					m.messenger.Notify(&Event{Name: name, Payload: payload})
				}
			}
		}()

		// notify UI we're ready
		m.messenger.Notify(&Event{Name: "onOnline", Payload: "{}"})

		// check for new messages
		m.RefreshMessages()

		// run the pinner
		tcore.Node.Wallet.RunPinner()
	}()

	return nil
}

// Stop the mobile node
func (m *Mobile) Stop() error {
	if err := tcore.Node.StopWallet(); err != nil && err != wallet.ErrStopped {
		return err
	}
	return nil
}

// SignUpWithEmail creates an email based registration and calls core signup
func (m *Mobile) SignUpWithEmail(email string, username string, password string, referral string) error {
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
func (m *Mobile) SignIn(username string, password string) error {
	// build creds
	creds := &models.Credentials{
		Username: username,
		Password: password,
	}
	return tcore.Node.Wallet.SignIn(creds)
}

// SignOut calls core SignOut
func (m *Mobile) SignOut() error {
	return tcore.Node.Wallet.SignOut()
}

// IsSignedIn calls core IsSignedIn
func (m *Mobile) IsSignedIn() bool {
	si, _ := tcore.Node.Wallet.IsSignedIn()
	return si
}

// GetId calls core GetId
func (m *Mobile) GetId() (string, error) {
	return tcore.Node.Wallet.GetId()
}

// GetPubKey calls core GetPubKeyString
func (m *Mobile) GetPubKey() (string, error) {
	return tcore.Node.Wallet.GetPubKeyString()
}

// GetUsername calls core GetUsername
func (m *Mobile) GetUsername() (string, error) {
	return tcore.Node.Wallet.GetUsername()
}

// GetTokens calls core GetTokens
func (m *Mobile) GetTokens() (string, error) {
	tokens, err := tcore.Node.Wallet.GetTokens()
	if err != nil {
		return "", err
	}
	return toJSON(tokens)
}

// SetAvatarId calls core SetAvatarId
func (m *Mobile) SetAvatarId(id string) error {
	return tcore.Node.Wallet.SetAvatarId(id)
}

// GetProfile returns this peer's profile
func (m *Mobile) GetProfile() (string, error) {
	id, err := tcore.Node.Wallet.GetId()
	if err != nil {
		log.Errorf("error getting id %s: %s", id, err)
		return "", err
	}
	prof, err := tcore.Node.Wallet.GetProfile(id)
	if err != nil {
		log.Errorf("error getting profile %s: %s", id, err)
		return "", err
	}
	return toJSON(prof)
}

// GetPeerProfile uses a peer id to look up a profile
func (m *Mobile) GetPeerProfile(peerId string) (string, error) {
	prof, err := tcore.Node.Wallet.GetProfile(peerId)
	if err != nil {
		log.Errorf("error getting profile %s: %s", peerId, err)
		return "", err
	}
	return toJSON(prof)
}

// RefreshMessages run the message retriever and repointer jobs
func (m *Mobile) RefreshMessages() error {
	return tcore.Node.Wallet.RefreshMessages()
}

// Threads lists all threads
func (m *Mobile) Threads() (string, error) {
	threads := Threads{Items: make([]Thread, 0)}
	for _, thrd := range tcore.Node.Wallet.Threads() {
		peers := thrd.Peers()
		item := Thread{Id: thrd.Id, Name: thrd.Name, Peers: len(peers)}
		threads.Items = append(threads.Items, item)
	}
	return toJSON(threads)
}

// PhotoThreads call core PhotoThreads
func (m *Mobile) PhotoThreads(id string) (string, error) {
	threads := Threads{Items: make([]Thread, 0)}
	for _, thrd := range tcore.Node.Wallet.PhotoThreads(id) {
		peers := thrd.Peers()
		item := Thread{Id: thrd.Id, Name: thrd.Name, Peers: len(peers)}
		threads.Items = append(threads.Items, item)
	}
	return toJSON(threads)
}

// AddThread adds a new thread with the given name
func (m *Mobile) AddThread(name string, mnemonic string) (string, error) {
	var mnem *string
	if mnemonic != "" {
		mnem = &mnemonic
	}
	thrd, _, err := tcore.Node.Wallet.AddThreadWithMnemonic(name, mnem)
	if err != nil {
		return "", err
	}

	// build json
	peers := thrd.Peers()
	item := Thread{
		Id:    thrd.Id,
		Name:  thrd.Name,
		Peers: len(peers),
	}
	return toJSON(item)
}

// AddThreadInvite adds a new invite to a thread
func (m *Mobile) AddThreadInvite(threadId string, inviteePk string) (string, error) {
	_, thrd := tcore.Node.Wallet.GetThread(threadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread: %s", threadId))
	}

	// decode pubkey
	ikb, err := libp2pc.ConfigDecodeKey(inviteePk)
	if err != nil {
		return "", err
	}
	ipk, err := libp2pc.UnmarshalPublicKey(ikb)
	if err != nil {
		return "", err
	}

	// add it
	addr, err := thrd.AddInvite(ipk)
	if err != nil {
		return "", err
	}

	return addr.B58String(), nil
}

// AddExternalThreadInvite generates a new external invite link to a thread
func (m *Mobile) AddExternalThreadInvite(threadId string) (string, error) {
	_, thrd := tcore.Node.Wallet.GetThread(threadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread: %s", threadId))
	}

	// add it
	addr, key, err := thrd.AddExternalInvite()
	if err != nil {
		return "", err
	}

	// create a structured invite
	username, _ := m.GetUsername()
	invite := ExternalInvite{
		Id:      addr.B58String(),
		Key:     string(key),
		Inviter: username,
	}

	return toJSON(invite)
}

// AcceptExternalThreadInvite notifies the thread of a join
func (m *Mobile) AcceptExternalThreadInvite(id string, key string) (string, error) {
	m.waitForOnline()
	addr, err := tcore.Node.Wallet.AcceptExternalThreadInvite(id, []byte(key))
	if err != nil {
		return "", err
	}
	return addr.B58String(), nil
}

// RemoveThread call core RemoveDevice
func (m *Mobile) RemoveThread(id string) (string, error) {
	addr, err := tcore.Node.Wallet.RemoveThread(id)
	if err != nil {
		return "", err
	}
	return addr.B58String(), err
}

// Devices lists all devices
func (m *Mobile) Devices() (string, error) {
	devices := Devices{Items: make([]Device, 0)}
	for _, dev := range tcore.Node.Wallet.Devices() {
		item := Device{Id: dev.Id, Name: dev.Name}
		devices.Items = append(devices.Items, item)
	}
	return toJSON(devices)
}

// AddDevice calls core AddDevice
func (m *Mobile) AddDevice(name string, pubKey string) error {
	m.waitForOnline()
	pkb, err := libp2pc.ConfigDecodeKey(pubKey)
	if err != nil {
		return err
	}
	pk, err := libp2pc.UnmarshalPublicKey(pkb)
	if err != nil {
		return err
	}
	return tcore.Node.Wallet.AddDevice(name, pk)
}

// RemoveDevice call core RemoveDevice
func (m *Mobile) RemoveDevice(id string) error {
	return tcore.Node.Wallet.RemoveDevice(id)
}

// AddPhoto adds a photo by path
func (m *Mobile) AddPhoto(path string) (string, error) {
	added, err := tcore.Node.Wallet.AddPhoto(path)
	if err != nil {
		return "", err
	}
	return toJSON(added)
}

// SharePhoto adds an existing photo to a new thread
func (m *Mobile) AddPhotoToThread(dataId string, key string, threadId string, caption string) (string, error) {
	_, thrd := tcore.Node.Wallet.GetThread(threadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", threadId))
	}

	addr, err := thrd.AddPhoto(dataId, caption, []byte(key))
	if err != nil {
		return "", err
	}

	return addr.B58String(), nil
}

// SharePhoto adds an existing photo to a new thread
func (m *Mobile) SharePhotoToThread(dataId string, threadId string, caption string) (string, error) {
	block, err := tcore.Node.Wallet.GetBlockByDataId(dataId)
	if err != nil {
		return "", err
	}
	_, fromThread := tcore.Node.Wallet.GetThread(block.ThreadId)
	if fromThread == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId))
	}
	_, toThread := tcore.Node.Wallet.GetThread(threadId)
	if toThread == nil {
		return "", errors.New(fmt.Sprintf("could not find thread %s", threadId))
	}
	key, err := fromThread.Decrypt(block.DataKeyCipher)
	if err != nil {
		return "", err
	}

	// TODO: owner challenge
	addr, err := toThread.AddPhoto(dataId, caption, key)
	if err != nil {
		return "", err
	}

	return addr.B58String(), nil
}

// GetPhotos returns thread photo blocks with json encoding
func (m *Mobile) GetPhotos(offsetId string, limit int, threadId string) (string, error) {
	_, thrd := tcore.Node.Wallet.GetThread(threadId)
	if thrd == nil {
		return "", errors.New(fmt.Sprintf("thread not found: %s", threadId))
	}

	// build json
	photos := &Photos{Items: make([]Photo, 0)}
	for _, b := range thrd.Blocks(offsetId, limit, repo.PhotoBlock) {
		var caption string
		if b.DataCaptionCipher != nil {
			captionb, err := thrd.Decrypt(b.DataCaptionCipher)
			if err != nil {
				return "", err
			}
			caption = string(captionb)
		}
		authorId, err := util.IdFromEncodedPublicKey(b.AuthorPk)
		if err != nil {
			return "", err
		}
		photos.Items = append(photos.Items, Photo{
			Id:       b.DataId,
			Date:     b.Date,
			Caption:  string(caption),
			AuthorId: authorId.Pretty(),
		})
	}

	// check for offline messages
	m.RefreshMessages()

	return toJSON(photos)
}

// GetPhotoData returns a data url for a photo
func (m *Mobile) GetPhotoData(id string) (string, error) {
	return m.getImageData(id, "photo", false)
}

// GetPhotoData returns a data url for a photo
func (m *Mobile) GetThumbData(id string) (string, error) {
	return m.getImageData(id, "thumb", true)
}

// GetPhotoMetadata returns a meta data object for a photo
func (m *Mobile) GetPhotoMetadata(id string) (string, error) {
	block, err := tcore.Node.Wallet.GetBlockByDataId(id)
	if err != nil {
		log.Errorf("could not find block for data id %s: %s", id, err)
		return "", err
	}
	_, thrd := tcore.Node.Wallet.GetThread(block.ThreadId)
	if thrd == nil {
		err := errors.New(fmt.Sprintf("could not find thread: %s", block.ThreadId))
		log.Error(err.Error())
		return "", err
	}
	meta, err := thrd.GetPhotoMetaData(id, block)
	if err != nil {
		log.Warningf("get photo meta data failed %s: %s", id, err)
		meta = &model.PhotoMetadata{}
	}
	return toJSON(meta)
}

// GetPhotoKey calls core GetPhotoKey
func (m *Mobile) GetPhotoKey(id string) (string, error) {
	return tcore.Node.Wallet.GetPhotoKey(id)
}

// getImageData returns a data url for an image under a path
func (m *Mobile) getImageData(id string, path string, isThumb bool) (string, error) {
	block, err := tcore.Node.Wallet.GetBlockByDataId(id)
	if err != nil {
		log.Errorf("could not find block for data id %s: %s", id, err)
		return "", err
	}
	_, thrd := tcore.Node.Wallet.GetThread(block.ThreadId)
	if thrd == nil {
		err := errors.New(fmt.Sprintf("could not find thread: %s", block.ThreadId))
		log.Error(err.Error())
		return "", err
	}
	url, err := thrd.GetBlockDataBase64(fmt.Sprintf("%s/%s", id, path), block)
	if err != nil {
		log.Errorf("get block data base64 failed %s: %s", id, err)
		return "", err
	}

	// get meta data for url type
	meta, err := thrd.GetPhotoMetaData(id, block)
	if err != nil {
		log.Warningf("get photo meta data failed %s: %s", id, err)
		meta = &model.PhotoMetadata{}
	}
	if isThumb {
		url = getThumbDataURLPrefix(meta) + url
	} else {
		url = getPhotoDataURLPrefix(meta) + url
	}
	data := &ImageData{Url: url, Metadata: meta}

	return toJSON(data)
}

// subscribe to thread and pass updates to messenger
func (m *Mobile) subscribe(thrd *thread.Thread) {
	for {
		select {
		case update, ok := <-thrd.Updates():
			if !ok {
				return
			}
			payload, err := toJSON(update)
			if err == nil {
				m.messenger.Notify(&Event{Name: "onThreadUpdate", Payload: payload})
			}
		}
	}
}

// waitForOnline waits up to 5 seconds for the node to go online
func (m *Mobile) waitForOnline() {
	if tcore.Node.Wallet.IsOnline() {
		return
	}
	deadline := time.Now().Add(time.Second * 5)
	tick := time.NewTicker(time.Millisecond * 10)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if tcore.Node.Wallet.IsOnline() || time.Now().After(deadline) {
				return
			}
		}
	}
}

// getDataURLPrefix adds the correct data url prefix to a data url
func getPhotoDataURLPrefix(meta *model.PhotoMetadata) string {
	switch util.Format(meta.Format) {
	case util.PNG:
		return "data:image/png;base64,"
	case util.GIF:
		return "data:image/gif;base64,"
	default:
		return "data:image/jpeg;base64,"
	}
}

func getThumbDataURLPrefix(meta *model.PhotoMetadata) string {
	switch util.Format(meta.ThumbnailFormat) {
	case util.GIF:
		return "data:image/gif;base64,"
	default:
		return "data:image/jpeg;base64,"
	}
}

// toJSON returns a json string and logs errors
func toJSON(any interface{}) (string, error) {
	jsonb, err := json.Marshal(any)
	if err != nil {
		log.Errorf("error marshaling json: %s", err)
		return "", err
	}
	return string(jsonb), nil
}
