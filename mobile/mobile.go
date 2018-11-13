package mobile

import (
	"encoding/json"
	logger "gx/ipfs/QmQvJiADDe7JR4m968MwXobTCCzUqQkP87aRHe29MEBGHV/go-logging"
	logging "gx/ipfs/QmcVVHfdyv15GVPk7NrxdWjh2hLVccXnoD8j2tyQShiXJb/go-log"
	"strings"
	"time"

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

// NewWallet creates a brand new wallet and returns its recovery phrase
func NewWallet(wordCount int) (string, error) {
	// determine word count
	wcount, err := wallet.NewWordCount(wordCount)
	if err != nil {
		return "", err
	}

	// create a new wallet
	w, err := wallet.NewWallet(wcount.EntropySize())
	if err != nil {
		return "", err
	}

	// return the new recovery phrase
	return w.RecoveryPhrase, nil
}

// WalletAccount represents a derived account in a wallet
type WalletAccount struct {
	Seed    string
	Address string
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
	LogLevel  string
	LogToDisk bool
}

// MigrateConfig is used to define options during a major migration
type MigrateConfig struct {
	RepoPath string
}

// RunConfig is used to define run options for a mobile node
type RunConfig struct {
	RepoPath string
}

// Mobile is the name of the framework (must match package name)
type Mobile struct {
	RepoPath  string
	messenger Messenger
}

// InitRepo calls core InitRepo
func InitRepo(config *InitConfig) error {
	// convert seed string to full account keypair
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

	// logLevel is one of: critical error warning notice info debug
	logLevel, err := logger.LogLevel(strings.ToUpper(config.LogLevel))
	if err != nil {
		logLevel = logger.ERROR
	}

	// ready to call core
	return core.InitRepo(core.InitConfig{
		Account:   accnt,
		RepoPath:  config.RepoPath,
		IsMobile:  true,
		LogLevel:  logLevel,
		LogToDisk: config.LogToDisk,
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
	// build textile node
	node, err := core.NewTextile(core.RunConfig{
		RepoPath: config.RepoPath,
	})
	if err != nil {
		return nil, err
	}
	core.Node = node

	return &Mobile{RepoPath: config.RepoPath, messenger: messenger}, nil
}

// Start the mobile node
func (m *Mobile) Start() error {
	if err := core.Node.Start(); err != nil {
		if err == core.ErrStarted {
			return nil
		}
		return err
	}

	go func() {
		<-core.Node.OnlineCh()

		// subscribe to wallet updates
		go func() {
			for {
				select {
				case update, ok := <-core.Node.UpdateCh():
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
						name = "onDeviceAdded"
					case core.AccountPeerRemoved:
						name = "onDeviceRemoved"
					}
					m.messenger.Notify(&Event{Name: name, Payload: payload})
				}
			}
		}()

		// subscribe to thread updates
		go func() {
			for {
				select {
				case update, ok := <-core.Node.ThreadUpdateCh():
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
				case notification, ok := <-core.Node.NotificationCh():
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
	if err := core.Node.Stop(); err != nil && err != core.ErrStopped {
		return err
	}
	return nil
}

// Overview calls core Overview
func (m *Mobile) Overview() (string, error) {
	stats, err := core.Node.Overview()
	if err != nil {
		return "", err
	}
	return toJSON(stats)
}

// Version returns core Version
func (m *Mobile) Version() string {
	return core.Version
}

// PeerId returns the ipfs peer id
func (m *Mobile) PeerId() (string, error) {
	pid, err := core.Node.PeerId()
	if err != nil {
		return "", err
	}
	return pid.Pretty(), nil
}

// waitForOnline waits up to 5 seconds for the node to go online
func (m *Mobile) waitForOnline() {
	if core.Node.Online() {
		return
	}
	deadline := time.Now().Add(time.Second * 5)
	tick := time.NewTicker(time.Millisecond * 10)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if core.Node.Online() || time.Now().After(deadline) {
				return
			}
		}
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
