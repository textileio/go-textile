package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/repo/fsrepo"
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
	// window  *astilectron.Window
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

func start(app *astilectron.Astilectron, w []*astilectron.Window, _ *astilectron.Menu, t *astilectron.Tray, _ *astilectron.Menu) error {
	// remove the dock icon
	dock := app.Dock()
	dock.Hide()

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

	// temp create new wallet each time
	wcount, err := wallet.NewWordCount(12)
	if err != nil {
		return err
	}

	wallet, err := wallet.NewWallet(wcount.EntropySize())
	if err != nil {
		return err
	}
	fmt.Println(wallet.RecoveryPhrase)
	// show first account
	kp, err := wallet.AccountAt(0, "password")
	if err != nil {
		return err
	}
	fmt.Println(kp.Address())
	fmt.Println(kp.Seed())

	repoPath := filepath.Join(appDir, kp.Address())

	// run init if needed
	if !fsrepo.IsInitialized(repoPath) {
		accnt := keypair.Random()
		initc := core.InitConfig{
			Account:     accnt,
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
	node, err = core.NewTextile(core.RunConfig{RepoPath: repoPath})
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

	// pid, err := node.PeerId()
	// if err != nil {
	// 	astilog.Fatalf("get peer id failed: %s", err)
	// }

	return nil
}

// handleMessage handles incoming messages from Javascript/Electron
func handleMessage(_ *astilectron.Window, m bootstrap.MessageIn) (interface{}, error) {
	switch m.Name {
	default:
		return nil, nil
	}
}
