package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gx/ipfs/QmUf5i9YncsDbikKC5wWBmPeLVxz35yKSQwbp11REBGFGi/go-ipfs/repo/fsrepo"
	"github.com/pkg/browser"
	"github.com/atotto/clipboard"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/mitchellh/go-homedir"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/gateway"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
)

var (
	appName = "Textile"
	debug   = flag.Bool("d", false, "enables debug mode")
	app     *astilectron.Astilectron
	menu    *astilectron.Menu
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
				user := node.User(note.Actor)
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

func start(app *astilectron.Astilectron, _ []*astilectron.Window, _ *astilectron.Menu, t *astilectron.Tray, m *astilectron.Menu) error {
	// remove the dock icon
	var d = app.Dock()
	d.Hide()

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
	repoPath := filepath.Join(appDir, "repo")

	// run init if needed
	if !fsrepo.IsInitialized(repoPath) {
		accnt := keypair.Random()
		initc := core.InitConfig{
			Account:   accnt,
			RepoPath:  repoPath,
			LogToDisk: true,
			GatewayAddr: fmt.Sprintf("127.0.0.1:5052"),
			ApiAddr: fmt.Sprintf("127.0.0.1:40602"),

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

	pid, err := node.PeerId()
	if err != nil {
		astilog.Fatalf("get peer id failed: %s", err)
	}

	items := []*astilectron.MenuItemOptions{
		{
			Label: astilectron.PtrStr("Online/Offline"),
			OnClick: func(e astilectron.Event) (deleteListener bool) {
				if *e.MenuItemOptions.Checked {
					startNode()
				} else {
					stopNode()
				}
				return
			},
			Type: astilectron.MenuItemTypeCheckbox,
			Checked: astilectron.PtrBool(true),
		},
		{
			Label: astilectron.PtrStr("Check Messages"),
			OnClick: func(e astilectron.Event) (deleteListener bool) {
				node.CheckCafeMessages()
				return
			},
		},
		{
			Type: astilectron.MenuItemTypeSeparator,
		},
		{
			Label: astilectron.PtrStr("Peer"),
			SubMenu: []*astilectron.MenuItemOptions{
        {
					Label:   astilectron.PtrStr("Copy Peer ID"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						clipboard.WriteAll(pid.Pretty())
						return
					},
				},
				{
					Label:   astilectron.PtrStr("Copy Peer Address"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						clipboard.WriteAll(node.Account().Address())
						return
					},
				},
			},
		},
		{
			Label: astilectron.PtrStr("API"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label:   astilectron.PtrStr(fmt.Sprintf("Copy URL (%s)", node.ApiAddr())),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						clipboard.WriteAll(fmt.Sprintf("http://%s/api/v0", node.ApiAddr()))
						return
					},
				},
				{
					Label:   astilectron.PtrStr(fmt.Sprintf("Copy gateway (%s)", gateway.Host.Addr())),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						clipboard.WriteAll(fmt.Sprintf("http://%s", gateway.Host.Addr()))
						return
					},
				},
				{
					Label: astilectron.PtrStr("View docs"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						browser.OpenURL(fmt.Sprintf("http://%s/docs/index.html", node.ApiAddr()))
						return
					},
				},
			},
		},
		{
			Label:   astilectron.PtrStr("Repo"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astilectron.PtrStr("View/edit config file"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						browser.OpenFile(filepath.Join(repoPath, "textile"))
						return
					},
				},
				{
					Label:   astilectron.PtrStr("Open repo folder"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						browser.OpenFile(repoPath)
						return
					},
				},
			},
		},
		{
			Type: astilectron.MenuItemTypeSeparator,
		},
		{
			Label: astilectron.PtrStr("Quit"),
			OnClick: func(e astilectron.Event) (deleteListener bool) {
				stopNode()
				app.Quit()
				return
			},
		},
	}

	for _, item := range items {
		var i = m.NewItem(item)
		m.Append(i)
	}

	return nil
}

// handleMessage handles incoming messages from Javascript/Electron
func handleMessage(_ *astilectron.Window, m bootstrap.MessageIn) (interface{}, error) {
	switch m.Name {
	default:
		return nil, nil
	}
}
