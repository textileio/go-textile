package mobile

import (
	"encoding/json"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	logging "gx/ipfs/QmZChCsSt8DctjceaL56Eibc29CVQq4dGKRXC5JRZ6Ppae/go-log"

	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/wallet"
)

var log = logging.Logger("tex-mobile")

// Messenger is used to inform the bridge layer of new data waiting to be queried
type Messenger interface {
	Notify(event []byte)
}

// Callback is used for asyc methods (payload is a protobuf)
type Callback interface {
	Call(payload []byte, err error)
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
func WalletAccountAt(phrase string, index int, password string) (*WalletAccount, error) {
	w := wallet.NewWalletFromRecoveryPhrase(phrase)
	accnt, err := w.AccountAt(index, password)
	if err != nil {
		return nil, err
	}
	return &WalletAccount{
		Seed:    accnt.Seed(),
		Address: accnt.Address(),
	}, nil
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
					data, err := proto.Marshal(update)
					if err != nil {
						log.Errorf("error marshaling event data: %s", err)
						continue
					}
					m.notify(&pb.MobileEvent{
						Name: pb.MobileEvent_WALLET_UPDATE.String(),
						Data: data,
					})
				}
			}
		}()

		// subscribe to thread updates
		go func() {
			for {
				select {
				case value, ok := <-m.listener.Ch:
					if !ok {
						return
					}
					if update, ok := value.(*pb.FeedItem); ok {
						data, err := proto.Marshal(update)
						if err != nil {
							log.Errorf("error marshaling event data: %s", err)
							continue
						}
						m.notify(&pb.MobileEvent{
							Name: pb.MobileEvent_THREAD_UPDATE.String(),
							Data: data,
						})
					}
				}
			}
		}()

		// subscribe to notifications
		go func() {
			for {
				select {
				case note, ok := <-m.node.NotificationCh():
					if !ok {
						return
					}
					data, err := proto.Marshal(note)
					if err != nil {
						log.Errorf("error marshaling event data: %s", err)
						continue
					}
					m.notify(&pb.MobileEvent{
						Name: pb.MobileEvent_NOTIFICATION.String(),
						Data: data,
					})
				}
			}
		}()

		// ready
		m.notify(&pb.MobileEvent{
			Name: pb.MobileEvent_NODE_ONLINE.String(),
		})
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

// OnlineCh returns core OnlineCh
func (m *Mobile) OnlineCh() <-chan struct{} {
	return m.node.OnlineCh()
}

// blockView returns marshaled view of a block
func (m *Mobile) blockView(hash mh.Multihash) ([]byte, error) {
	view, err := m.node.BlockView(hash.B58String())
	if err != nil {
		return nil, err
	}
	return proto.Marshal(view)
}

func (m *Mobile) notify(event *pb.MobileEvent) {
	payload, err := proto.Marshal(event)
	if err != nil {
		log.Errorf("error marshaling event: %s", err)
		return
	}

	m.messenger.Notify(payload)
}
