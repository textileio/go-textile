package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	logger "gx/ipfs/QmQvJiADDe7JR4m968MwXobTCCzUqQkP87aRHe29MEBGHV/go-logging"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/repo/fsrepo"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/mitchellh/go-homedir"
	"github.com/skip2/go-qrcode"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/gateway"
	"github.com/textileio/textile-go/keypair"
	"github.com/textileio/textile-go/repo"
)

var (
	appName     = "Textile"
	builtAt     string
	debug       = flag.Bool("d", false, "enables the debug mode")
	window      *astilectron.Window
	gatewayAddr string
	expanded    bool
)

const (
	SetupSize       = 384
	QRCodeSize      = 256
	InitialWidth    = 1024
	InitialHeight   = 633
	SleepOnLoad     = time.Second * 1
	SleepOnPreReady = time.Millisecond * 200
	SleepOnExpand   = time.Millisecond * 200
)

var node *core.Textile

func main() {
	flag.Parse()
	astilog.FlagInit()
	bootstrapApp()
}

func start(a *astilectron.Astilectron, w []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
	window = w[0]
	window.Show()

	// get homedir
	home, err := homedir.Dir()
	if err != nil {
		astilog.Fatal(fmt.Errorf("get homedir failed: %s", err))
	}

	// ensure app support folder is created
	var appDir string
	if runtime.GOOS == "darwin" {
		appDir = filepath.Join(home, "Library/Application Support/Textile")
	} else {
		appDir = filepath.Join(home, ".textile")
	}
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return err
	}
	repoPath := filepath.Join(appDir, "repo")

	// run init if needed
	if !fsrepo.IsInitialized(repoPath) {
		accnt := keypair.Random()
		initc := core.InitConfig{
			Account:   accnt,
			RepoPath:  repoPath,
			LogLevel:  logger.ERROR,
			LogToDisk: true,
		}
		if err := core.InitRepo(initc); err != nil {
			return err
		}
	}

	// build textile node
	node, err = core.NewTextile(core.RunConfig{RepoPath: repoPath})
	if err != nil {
		return err
	}

	// bring the node online and startup the gateway
	if err := node.Start(); err != nil {
		return err
	}
	<-node.OnlineCh()

	// subscribe to wallet updates
	go func() {
		for {
			select {
			case update, ok := <-node.UpdateCh():
				if !ok {
					return
				}
				payload := map[string]interface{}{
					"update": update,
				}
				switch update.Type {
				case core.ThreadAdded:
					if expanded {
						sendData("wallet.update", payload)
					} else {
						sendPreReady()
						window.Hide()
						expandWindow()
						sendData("wallet.update", payload)
						window.Show()
						window.Focus()
					}
				default:
					sendData("wallet.update", payload)
				}
			}
		}
	}()

	// subscribe to thread updates
	go func() {
		for {
			select {
			case update, ok := <-node.GetTheadUpdateListener():
				if !ok {
					return
				}
				sendData("thread.update", map[string]interface{}{
					"update": update,
				})
			}
		}
	}()

	// subscribe to notifications
	go func() {
		for {
			select {
			case note, ok := <-node.NotificationCh():
				if !ok {
					return
				}
				username := node.ContactUsername(note.ActorId)
				var uinote = a.NewNotification(&astilectron.NotificationOptions{
					Title: note.Subject,
					Body:  fmt.Sprintf("%s %s.", username, note.Body),
					Icon:  "/resources/icon.png",
				})

				// tmp auto-accept thread invites
				if note.Type == repo.InviteReceivedNotification {
					go func(tid string) {
						if _, err := node.AcceptThreadInvite(tid); err != nil {
							astilog.Error(err)
						}
					}(note.BlockId)
				}

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

	// start the gateway
	gateway.Host.Start(fmt.Sprintf("127.0.0.1:%s", core.GetRandomPort()))

	// save off the server address
	gatewayAddr = fmt.Sprintf("http://%s", gateway.Host.Addr())

	// sleep for a bit on the landing screen, it feels better
	time.Sleep(SleepOnLoad)

	// send cookie info to front-end
	sendData("login", map[string]interface{}{
		"name":    "SessionId",
		"value":   "not used",
		"gateway": gatewayAddr,
	})

	// check if we're configured yet
	threads := node.Threads()
	if len(threads) > 0 {
		// load threads for UI
		var threadsJSON []map[string]interface{}
		for _, thrd := range threads {
			threadsJSON = append(threadsJSON, map[string]interface{}{
				"id":   thrd.Id,
				"name": thrd.Name,
			})
		}

		// reveal
		sendPreReady()
		window.Hide()
		expandWindow()
		sendData("ready", map[string]interface{}{
			"threads": threadsJSON,
		})
		window.Show()
		window.Focus()

	} else {
		// get qr code for setup
		qr, pk, err := getQRCode()
		if err != nil {
			astilog.Error(err)
			return err
		}
		sendData("setup", map[string]interface{}{
			"qr": qr,
			"pk": pk,
		})
	}

	return nil
}

func sendMessage(name string) {
	window.SendMessage(map[string]string{"name": name})
}

func sendData(name string, data map[string]interface{}) {
	data["name"] = name
	window.SendMessage(data)
}

func handleMessage(_ *astilectron.Window, m bootstrap.MessageIn) (interface{}, error) {
	switch m.Name {
	case "refresh":
		if err := node.CheckCafeMail(); err != nil {
			return nil, err
		}
		return map[string]interface{}{}, nil
	case "thread.load":
		var threadId string
		if err := json.Unmarshal(m.Payload, &threadId); err != nil {
			return nil, err
		}
		html, err := getThreadPhotos(threadId)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"html": html,
		}, nil
	default:
		return map[string]interface{}{}, nil
	}
}

func getQRCode() (string, string, error) {
	// get our own peer id for receiving an account key
	pid, err := node.PeerId()
	if err != nil {
		return "", "", err
	}

	// create a qr code
	url := fmt.Sprintf("https://www.textile.photos/invites/device#id=%s", pid.Pretty())
	png, err := qrcode.Encode(url, qrcode.Medium, QRCodeSize)
	if err != nil {
		return "", "", err
	}

	return base64.StdEncoding.EncodeToString(png), pid.Pretty(), nil
}

func getThreadPhotos(id string) (string, error) {
	thrd := node.Thread(id)
	if thrd == nil {
		return "", core.ErrThreadNotFound
	}
	var html string
	query := fmt.Sprintf("threadId='%s' and type=%d", thrd.Id, repo.FilesBlock)
	for range node.Blocks("", -1, query) {
		//photo := fmt.Sprintf("%s/ipfs/%s/photo?block=%s", gatewayAddr, block.DataId, block.Id)
		//small := fmt.Sprintf("%s/ipfs/%s/small?block=%s", gatewayAddr, block.DataId, block.Id)
		//meta := fmt.Sprintf("%s/ipfs/%s/meta?block=%s", gatewayAddr, block.DataId, block.Id)
		//img := fmt.Sprintf("<img src=\"%s\" />", small)
		//html += fmt.Sprintf(
		//	"<div id=\"%s\" class=\"grid-item\" ondragstart=\"imageDragStart(event);\" draggable=\"true\" data-url=\"%s\" data-meta=\"%s\">%s</div>",
		//	block.Id, photo, meta, img)
	}
	return html, nil
}

func sendPreReady() {
	sendMessage("preready")
	time.Sleep(SleepOnPreReady)
}

func expandWindow() {
	expanded = true
	go window.Resize(InitialWidth, InitialHeight)
	go window.Center()
	time.Sleep(SleepOnExpand)
}
