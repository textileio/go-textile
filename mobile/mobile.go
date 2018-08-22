package mobile

import (
	"encoding/json"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/wallet"
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

// Create a gomobile compatible wrapper around TextileNode
func NewNode(config *NodeConfig, messenger Messenger) (*Mobile, error) {
	ll, err := logging.LogLevel(config.LogLevel)
	if err != nil {
		ll = logging.INFO
	}
	cconfig := core.NodeConfig{
		LogLevel: ll,
		LogFiles: config.LogFiles,
		WalletConfig: wallet.Config{
			RepoPath: config.RepoPath,
			IsMobile: true,
			CafeAddr: config.CafeAddr,
		},
	}
	node, mnemonic, err := core.NewNode(cconfig)
	if err != nil {
		return nil, err
	}
	core.Node = node

	return &Mobile{RepoPath: config.RepoPath, Mnemonic: mnemonic, messenger: messenger}, nil
}

// Start the mobile node
func (m *Mobile) Start() error {
	if err := core.Node.StartWallet(); err != nil {
		if err == wallet.ErrStarted {
			return nil
		}
		return err
	}

	go func() {
		<-core.Node.Wallet.Online()

		// subscribe to wallet updates
		go func() {
			for {
				select {
				case update, ok := <-core.Node.Wallet.Updates():
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

		// subscribe to thread updates
		go func() {
			for {
				select {
				case update, ok := <-core.Node.Wallet.ThreadUpdates():
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
				case notification, ok := <-core.Node.Wallet.Notifications():
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

		// run the pinner
		core.Node.Wallet.RunPinner()
	}()

	return nil
}

// Stop the mobile node
func (m *Mobile) Stop() error {
	if err := core.Node.StopWallet(); err != nil && err != wallet.ErrStopped {
		return err
	}
	return nil
}

// RefreshMessages run the message retriever and repointer jobs
func (m *Mobile) RefreshMessages() error {
	return core.Node.Wallet.RefreshMessages()
}

// Overview calls core Overview
func (m *Mobile) Overview() (string, error) {
	stats, err := core.Node.Wallet.Overview()
	if err != nil {
		return "", err
	}
	return toJSON(stats)
}

// waitForOnline waits up to 5 seconds for the node to go online
func (m *Mobile) waitForOnline() {
	if core.Node.Wallet.IsOnline() {
		return
	}
	deadline := time.Now().Add(time.Second * 5)
	tick := time.NewTicker(time.Millisecond * 10)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if core.Node.Wallet.IsOnline() || time.Now().After(deadline) {
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
