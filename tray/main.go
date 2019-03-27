package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/repo/fsrepo"
	// "github.com/atotto/clipboard"
	"github.com/pkg/browser"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/mitchellh/go-homedir"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/gateway"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/wallet"
)

var (
	appName = "Textile"
	debug   = flag.Bool("d", false, "enables debug mode")
	app     *astilectron.Astilectron
	window  *astilectron.Window
)

var node *core.Textile

func main() {
	flag.Parse()
	astilog.FlagInit()
	bootstrapApp()
}

func startNode() error {
	if err := node.Start(); err != nil {
		astilog.Error(err)
		if err == core.ErrStarted {
			return nil
		}
		return err
	}

	// subscribe to notifications
	go func() {
		for {
			select {
			case note, ok := <-node.NotificationCh():
				if !ok {
					return
				}
				user := node.PeerUser(note.Actor)
				var uinote = app.NewNotification(&astilectron.NotificationOptions{
					Title: note.Subject,
					Body:  fmt.Sprintf("%s: %s.", user.Name, note.Body),
					Icon:  fmt.Sprintf("%s/ipfs/%s/0/small/d", gateway.Host.Addr(), user.Avatar),
				})

				// tmp auto-accept thread invites
				if note.Type == pb.Notification_INVITE_RECEIVED {
					go func(tid string) {
						if _, err := node.AcceptInvite(tid); err != nil {
							astilog.Error(err)
						}
					}(note.Block)
				}

				fmt.Println(fmt.Sprintf("%s: %s.", user.Name, note.Body))

				// show notification
				go func(n *astilectron.Notification) {
					if err := n.Create(); err != nil {
						astilog.Error(err)
						return
					}
					if err := n.Show(); err != nil {
						astilog.Error(err)
						return
					}
				}(uinote)
			}
		}
	}()

	// setup and start the apis
	gateway.Host = &gateway.Gateway{
		Node: node,
	}
	node.StartApi(node.Config().Addresses.API, true)
	gateway.Host.Start(node.Config().Addresses.Gateway)

	// wait for node to come online
	<-node.OnlineCh()

	return nil
}

func stopNode() error {
	if err := node.Stop(); err != nil {
		astilog.Error(err)
		if err == core.ErrStopped {
			return nil
		}
		return err
	}
	if err := node.StopApi(); err != nil {
		return err
	}
	if err := gateway.Host.Stop(); err != nil {
		return err
	}

	return nil
}

func initAndStartTextile(mnemonic string, password string) error {
	// get homedir
	home, err := homedir.Dir()
	if err != nil {
		astilog.Fatal(fmt.Errorf("get homedir failed: %s", err))
	}

	// ensure app support folder is created
	var appDir string
	if runtime.GOOS == "darwin" {
		appDir = filepath.Join(home, "Library", "Application Support", "Textile")
	} else {
		appDir = filepath.Join(home, ".textile")
	}
	if err := os.MkdirAll(appDir, 0755); err != nil {
		astilog.Fatal(fmt.Errorf("create app dir failed: %s", err))
	}

	wallet := wallet.NewWalletFromRecoveryPhrase(mnemonic)
	// show first account
	// TODO: Ask the user for a new password on new account init
	kp, err := wallet.AccountAt(0, "")
	if err != nil {
		return err
	}

	repoPath := filepath.Join(appDir, kp.Address())

	// run init if needed
	if !fsrepo.IsInitialized(repoPath) {
		accnt := keypair.Random()
		initc := core.InitConfig{
			Account:     accnt,
			PinCode:     password,
			RepoPath:    repoPath,
			LogToDisk:   true,
			GatewayAddr: fmt.Sprintf("127.0.0.1:5052"),
			ApiAddr:     fmt.Sprintf("127.0.0.1:40602"),
		}
		if err := core.InitRepo(initc); err != nil {
			astilog.Fatal(fmt.Errorf("create repo failed: %s", err))
		}
	}

	// build textile node
	node, err = core.NewTextile(core.RunConfig{
		PinCode:  password,
		RepoPath: repoPath,
	})
	if err != nil {
		astilog.Error(err)
		return err
	}

	// bring the node online
	err = startNode()
	if err != nil {
		astilog.Error(err)
		return err
	}
	return nil
}

func start(
	app *astilectron.Astilectron,
	windows []*astilectron.Window,
	_ *astilectron.Menu,
	t *astilectron.Tray,
	m *astilectron.Menu) error {
	// remove the dock icon
	dock := app.Dock()
	dock.Hide()

	window = windows[0]

	var i = m.NewItem(&astilectron.MenuItemOptions{
		Label: astilectron.PtrStr("Quit"),
		OnClick: func(e astilectron.Event) (deleteListener bool) {
			stopNode()
			app.Quit()
			return
		},
	})
	m.Append(i)
	return nil
}

func sendData(name string, data map[string]interface{}) {
	data["name"] = name
	window.SendMessage(data)
}

// handleMessage handles incoming messages from Javascript/Electron
func handleMessage(w *astilectron.Window, m bootstrap.MessageIn) (interface{}, error) {
	switch m.Name {
	case "init":
		type init struct {
			Mnemonic string `json:"mnemonic"`
			Password string `json:"password,omitempty"`
		}
		var payload init
		if err := json.Unmarshal(m.Payload, &payload); err != nil {
			return nil, err
		}
		err := initAndStartTextile(payload.Mnemonic, payload.Password)
		if err != nil {
			return nil, err
		}
		return true, nil
	case "hide":
		w.Hide()
		return true, nil
	case "open":
		type open struct {
			File string `json:"file,omitempty"`
			URL string `json:"url,omitempty"`
		}
		var payload open
		if err := json.Unmarshal(m.Payload, &payload); err != nil {
			return nil, err
		}
		if payload.URL != "" {
			browser.OpenURL(payload.URL)
		}
		if payload.File != "" {
			browser.OpenFile(payload.File)
		}
		return true, nil
	default:
		return nil, nil
	}
}
