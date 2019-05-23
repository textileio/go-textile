package mobile

import (
	"github.com/golang/protobuf/proto"
	logging "github.com/ipfs/go-log"
	mh "github.com/multiformats/go-multihash"
	"github.com/textileio/go-textile/broadcast"
	"github.com/textileio/go-textile/common"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/wallet"
)

var log = logging.Logger("tex-mobile")

// Messenger is a push mechanism to the bridge
type Messenger interface {
	Notify(event *Event)
}

// Event is sent by Messenger to the bridge (data is a protobuf,
// name is the string value of a pb.MobileEvent_Type)
type Event struct {
	Name string
	Type int32
	Data []byte
}

// Callback is used for asyc methods (data is a protobuf)
type Callback interface {
	Call(data []byte, err error)
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

// WalletAccountAt derives the account at the given index
func WalletAccountAt(phrase string, index int, password string) ([]byte, error) {
	w := wallet.NewWalletFromRecoveryPhrase(phrase)
	accnt, err := w.AccountAt(index, password)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(&pb.MobileWalletAccount{
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
	RepoPath          string
	Debug             bool
	CafeOutboxHandler core.CafeOutboxHandler
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
		RepoPath:          config.RepoPath,
		CafeOutboxHandler: config.CafeOutboxHandler,
		Debug:             config.Debug,
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

		// subscribe to account updates
		go func() {
			for {
				select {
				case update, ok := <-m.node.UpdateCh():
					if !ok {
						return
					}
					m.notify(pb.MobileEventType_WALLET_UPDATE, update)
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
						m.notify(pb.MobileEventType_THREAD_UPDATE, update)
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
					m.notify(pb.MobileEventType_NOTIFICATION, note)
				}
			}
		}()

		m.notify(pb.MobileEventType_NODE_ONLINE, nil)
	}()

	m.notify(pb.MobileEventType_NODE_START, nil)
	return nil
}

// Stop the mobile node
func (m *Mobile) Stop() error {
	if err := m.node.Stop(); err != nil && err != core.ErrStopped {
		return err
	}
	m.notify(pb.MobileEventType_NODE_STOP, nil)
	return nil
}

// OnlineCh returns core OnlineCh
func (m *Mobile) OnlineCh() <-chan struct{} {
	return m.node.OnlineCh()
}

// Version returns common Version
func (m *Mobile) Version() string {
	return "v" + common.Version
}

// GitSummary returns common GitSummary
func (m *Mobile) GitSummary() string {
	return common.GitSummary
}

// Summary calls core Summary
func (m *Mobile) Summary() ([]byte, error) {
	if !m.node.Started() {
		return nil, core.ErrStopped
	}

	return proto.Marshal(m.node.Summary())
}

// blockView returns marshaled view of a block
func (m *Mobile) blockView(hash mh.Multihash) ([]byte, error) {
	view, err := m.node.BlockView(hash.B58String())
	if err != nil {
		return nil, err
	}
	return proto.Marshal(view)
}

func (m *Mobile) notify(etype pb.MobileEventType, msg proto.Message) {
	var data []byte
	if msg != nil {
		var err error
		data, err = proto.Marshal(msg)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
	m.messenger.Notify(&Event{
		Name: etype.String(),
		Type: int32(etype),
		Data: data,
	})
}
