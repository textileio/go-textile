package mobile

import (
	"encoding/json"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	logging "gx/ipfs/QmZChCsSt8DctjceaL56Eibc29CVQq4dGKRXC5JRZ6Ppae/go-log"

	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/wallet"
)

var log = logging.Logger("tex-mobile")

// Message is a generic go -> bridge message structure
type Event struct {
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

// Messenger is used to inform the bridge layer of new data waiting to be queried
type Messenger interface {
	Notify(event *Event)
}

// ProtoCallback is used for asyc methods (payload is a protobuf)
type ProtoCallback interface {
	Call(payload []byte, err error)
}

// StringCallback is used for asyc methods (payload is a string)
type StringCallback interface {
	Call(payload string, err error)
}

// NewWallet creates a brand new wallet and returns its recovery phrase
func NewWallet(wordCount int) (string, error) {
	wcount, err := wallet.NewWordCount(wordCount)
	if err != nil {
		return "", err
	}

	w, err := wallet.NewWallet(wcount.EntropySize())
	if err != nil {
		return "", err
	}

	return w.RecoveryPhrase, nil
}

// WalletAccount represents a derived account in a wallet
type WalletAccount struct {
	Seed    string `json:"seed"`
	Address string `json:"address"`
}

// WalletAccountAt derives the account at the given index
func WalletAccountAt(phrase string, index int, password string) (string, error) {
	w := wallet.NewWalletFromRecoveryPhrase(phrase)
	accnt, err := w.AccountAt(index, password)
	if err != nil {
		return "", err
	}
	return toJSON(WalletAccount{
		Seed:    accnt.Seed(),
		Address: accnt.Address(),
	})
}

// InitConfig is used to setup a textile node
type InitConfig struct {
	Seed      string
	RepoPath  string
	LogToDisk bool
	Debug     bool
}

// MigrateConfig is used to define options during a major migration
type MigrateConfig struct {
	RepoPath string
}

// RunConfig is used to define run options for a mobile node
type RunConfig struct {
	RepoPath string
	Debug    bool
}

// Mobile is the name of the framework (must match package name)
type Mobile struct {
	RepoPath  string
	node      *core.Textile
	messenger Messenger
	listener  *broadcast.Listener
}

// InitRepo calls core InitRepo
func InitRepo(config *InitConfig) error {
	if config.Seed == "" {
		return core.ErrAccountRequired
	}
	kp, err := keypair.Parse(config.Seed)
	if err != nil {
		return err
	}
	accnt, ok := kp.(*keypair.Full)
	if !ok {
		return keypair.ErrInvalidKey
	}

	return core.InitRepo(core.InitConfig{
		Account:   accnt,
		RepoPath:  config.RepoPath,
		IsMobile:  true,
		LogToDisk: config.LogToDisk,
		Debug:     config.Debug,
	})
}

// MigrateRepo calls core MigrateRepo
func MigrateRepo(config *MigrateConfig) error {
	return core.MigrateRepo(core.MigrateConfig{
		RepoPath: config.RepoPath,
	})
}

// Create a gomobile compatible wrapper around Textile
func NewTextile(config *RunConfig, messenger Messenger) (*Mobile, error) {
	node, err := core.NewTextile(core.RunConfig{
		RepoPath: config.RepoPath,
		Debug:    config.Debug,
	})
	if err != nil {
		return nil, err
	}

	return &Mobile{
		RepoPath:  config.RepoPath,
		node:      node,
		messenger: messenger,
		listener:  node.ThreadUpdateListener(),
	}, nil
}

// SetLogLevels provides access to the underlying node's setLogLevels method
func (m *Mobile) SetLogLevels(logLevelsString string) error {
	var logLevels map[string]string
	if logLevelsString != "" {
		err := json.Unmarshal([]byte(logLevelsString), &logLevels)
		if err != nil {
			return err
		}
	}
	return m.node.SetLogLevels(logLevels)
}

// Start the mobile node
func (m *Mobile) Start() error {
	if err := m.node.Start(); err != nil {
		if err == core.ErrStarted {
			return nil
		}
		return err
	}

	go func() {
		<-m.node.OnlineCh()

		// subscribe to wallet updates
		go func() {
			for {
				select {
				case update, ok := <-m.node.UpdateCh():
					if !ok {
						return
					}
					payload, err := toJSON(update)
					if err != nil {
						return
					}
					var name string
					switch update.Type {
					case core.ThreadAdded:
						name = "onThreadAdded"
					case core.ThreadRemoved:
						name = "onThreadRemoved"
					case core.AccountPeerAdded:
						name = "onAccountPeerAdded"
					case core.AccountPeerRemoved:
						name = "onAccountPeerRemoved"
					}
					m.messenger.Notify(&Event{Name: name, Payload: payload})
				}
			}
		}()

		// subscribe to thread updates
		go func() {
			for {
				select {
				case update, ok := <-m.listener.Ch:
					if !ok {
						return
					}
					payload, err := toJSON(update)
					if err == nil {
						m.messenger.Notify(&Event{Name: "onThreadUpdate", Payload: payload})
					}
				}
			}
		}()

		// subscribe to notifications
		go func() {
			for {
				select {
				case notification, ok := <-m.node.NotificationCh():
					if !ok {
						return
					}
					payload, err := toJSON(notification)
					if err == nil {
						m.messenger.Notify(&Event{Name: "onNotification", Payload: payload})
					}
				}
			}
		}()

		// notify UI we're ready
		m.messenger.Notify(&Event{Name: "onOnline", Payload: "{}"})
	}()

	return nil
}

// Stop the mobile node
func (m *Mobile) Stop() error {
	if err := m.node.Stop(); err != nil && err != core.ErrStopped {
		return err
	}
	return nil
}

// Version returns core Version
func (m *Mobile) Version() string {
	return core.Version
}

// OnlineCh returns core OnlineCh
func (m *Mobile) OnlineCh() <-chan struct{} {
	return m.node.OnlineCh()
}

// PeerId returns the ipfs peer id
func (m *Mobile) PeerId() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	pid, err := m.node.PeerId()
	if err != nil {
		return "", err
	}
	return pid.Pretty(), nil
}

// Overview calls core Overview
func (m *Mobile) Overview() (string, error) {
	if !m.node.Started() {
		return "", core.ErrStopped
	}

	stats, err := m.node.Overview()
	if err != nil {
		return "", err
	}
	return toJSON(stats)
}

// blockInfo returns json info view of a block
func (m *Mobile) blockInfo(hash mh.Multihash) (string, error) {
	info, err := m.node.BlockInfo(hash.B58String())
	if err != nil {
		return "", err
	}
	return toJSON(info)
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
