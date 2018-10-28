package mobile

import (
	"encoding/json"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/keypair"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/repo/fsrepo"
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
	Account  string
	PinCode  string
	RepoPath string
	LogLevel string
	LogFiles bool
}

// Mobile is the name of the framework (must match package name)
type Mobile struct {
	RepoPath  string
	messenger Messenger
}

// MigrateRepo calls core MigrateRepo
func MigrateRepo(config *NodeConfig) error {
	return core.MigrateRepo(core.MigrateConfig{
		RepoPath: config.RepoPath,
		PinCode:  config.PinCode,
	})
}

// Create a gomobile compatible wrapper around TextileNode
func NewNode(config *NodeConfig, messenger Messenger) (*Mobile, error) {
	// determine log level
	logLevel, err := logging.LogLevel(config.LogLevel)
	if err != nil {
		logLevel = logging.INFO
	}

	// run init if needed
	if !fsrepo.IsInitialized(config.RepoPath) {
		if config.Account == "" {
			return nil, core.ErrAccountRequired
		}
		kp, err := keypair.Parse(config.Account)
		if err != nil {
			return nil, err
		}
		accnt, ok := kp.(*keypair.Full)
		if !ok {
			return nil, keypair.ErrInvalidKey
		}
		initc := core.InitConfig{
			Account:  *accnt,
			PinCode:  config.PinCode,
			RepoPath: config.RepoPath,
			IsMobile: true,
			LogLevel: logLevel,
			LogFiles: config.LogFiles,
		}
		if err := core.InitRepo(initc); err != nil {
			return nil, err
		}
	}

	// build textile node
	runc := core.RunConfig{
		PinCode:  config.PinCode,
		RepoPath: config.RepoPath,
		LogLevel: logLevel,
		LogFiles: config.LogFiles,
	}
	node, err := core.NewTextile(runc)
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
				case update, ok := <-core.Node.Updates():
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
				case update, ok := <-core.Node.ThreadUpdates():
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
				case notification, ok := <-core.Node.Notifications():
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
