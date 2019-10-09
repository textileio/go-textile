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

// Callback is used for asyc methods
type Callback interface {
	Call(err error)
}

// ProtoCallback is used for asyc methods that deliver a protobuf message
type ProtoCallback interface {
	Call(msg []byte, err error)
}

// DataCallback is used for asyc methods that deliver raw data
type DataCallback interface {
	Call(data []byte, media string, err error)
}

// NewWallet creates a brand new wallet and returns its recovery phrase
func NewWallet(wordCount int) (string, error) {
	w, err := wallet.WalletFromWordCount(wordCount)
	if err != nil {
		return "", err
	}
	return w.RecoveryPhrase, nil
}

// WalletAccountAt derives the account at the given index
func WalletAccountAt(mnemonic string, index int, passphrase string) ([]byte, error) {
	w := wallet.WalletFromMnemonic(mnemonic)
	account, err := w.AccountAt(index, passphrase)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(&pb.MobileWalletAccount{
		Seed:    account.Seed(),
		Address: account.Address(),
	})
}

// InitConfig is used to setup a textile node
type InitConfig struct {
	Seed         string
	RepoPath     string
	BaseRepoPath string
	LogToDisk    bool
	Debug        bool
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

// Repo returns the actual location of the configured repo
func (conf InitConfig) Repo() (string, error) {
	coreConf, err := conf.coreInitConfig()
	if err != nil {
		return "", err
	}

	repo, err := coreConf.Repo()
	if err != nil {
		return "", err
	}

	return repo, nil
}

// RepoExists return whether or not the configured repo already exists
func (conf InitConfig) RepoExists() (bool, error) {
	coreConf, err := conf.coreInitConfig()
	if err != nil {
		return false, err
	}

	exists, err := coreConf.RepoExists()
	if err != nil {
		return false, err
	}

	return exists, nil
}

// RepoExists return whether or not the repo at repoPath exists
func RepoExists(repoPath string) bool {
	return core.RepoExists(repoPath)
}

// AccountRepoExists return whether or not the repo at repoPath exists
func AccountRepoExists(baseRepoPath string, accountAddress string) bool {
	return core.AccountRepoExists(baseRepoPath, accountAddress)
}

func (conf InitConfig) coreInitConfig() (core.InitConfig, error) {
	var accnt *keypair.Full
	if len(conf.Seed) > 0 {
		var err error
		accnt, err = toAccount(conf.Seed)
		if err != nil {
			return core.InitConfig{}, err
		}
	}

	return core.InitConfig{
		Account:      accnt,
		RepoPath:     conf.RepoPath,
		BaseRepoPath: conf.BaseRepoPath,
		IsMobile:     true,
		LogToDisk:    conf.LogToDisk,
		Debug:        conf.Debug,
	}, nil
}

func toAccount(seed string) (*keypair.Full, error) {
	kp, err := keypair.Parse(seed)
	if err != nil {
		return nil, err
	}
	accnt, ok := kp.(*keypair.Full)
	if !ok {
		return nil, keypair.ErrInvalidKey
	}
	return accnt, nil
}

// InitRepo calls core InitRepo
func InitRepo(config *InitConfig) error {
	coreConf, err := config.coreInitConfig()
	if err != nil {
		return err
	}
	return core.InitRepo(coreConf)
}

// MigrateRepo calls core MigrateRepo
func MigrateRepo(config *MigrateConfig) error {
	return core.MigrateRepo(core.MigrateConfig{
		RepoPath: config.RepoPath,
	})
}

// Create a gomobile compatible wrapper around Textile
func NewTextile(config *RunConfig, messenger Messenger) (*Mobile, error) {
	mobile := &Mobile{
		RepoPath:  config.RepoPath,
		messenger: messenger,
	}

	node, err := core.NewTextile(core.RunConfig{
		RepoPath:          config.RepoPath,
		CafeOutboxHandler: config.CafeOutboxHandler,
		CheckMessages:     mobile.checkCafeMessages,
		Debug:             config.Debug,
	})
	if err != nil {
		return nil, err
	}

	mobile.node = node
	mobile.listener = node.ThreadUpdateListener()

	return mobile, nil
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
					m.notify(pb.MobileEventType_ACCOUNT_UPDATE, update)
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
func (m *Mobile) Stop(cb Callback) {
	go func() {
		cb.Call(m.stop())
	}()
}

// stop is the sync version of Stop
func (m *Mobile) stop() error {
	if err := m.node.Stop(); err != nil {
		if err == core.ErrStopped {
			return nil
		}
		return err
	}

	m.notify(pb.MobileEventType_NODE_STOP, nil)
	return nil
}

// Online returns core Online
func (m *Mobile) Online() bool {
	return m.node.Online()
}

// WaitAdd calls core WaitAdd
func (m *Mobile) WaitAdd(delta int, src string) {
	m.node.WaitAdd(delta, src)
}

// WaitDone marks a wait as done in the stop wait group
func (m *Mobile) WaitDone(src string) {
	m.node.WaitDone(src)
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

// onlineCh returns core OnlineCh
func (m *Mobile) onlineCh() <-chan struct{} {
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
