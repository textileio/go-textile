package main

import (
	"flag"
	//"fmt"
	"github.com/asticode/go-astilectron"
	//"path/filepath"
	//"runtime"
	"encoding/base64"
	"fmt"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/skip2/go-qrcode"
	"time"
)

var (
	appName = "Textile"
	builtAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	window  *astilectron.Window
	//tray      *astilectron.Tray
	gatewayAddr string
	expanded    bool
	menuVisible = false
	onboarded   = false
	rootFolder  = false
	connected   = false
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

func main() {
	//run()
	flag.Parse()
	astilog.FlagInit()
	bootstrapApp()
}

func start(a *astilectron.Astilectron, w []*astilectron.Window, _ *astilectron.Menu, t *astilectron.Tray, _ *astilectron.Menu) error {

	// remove the dock icon
	var d = a.Dock()
	d.Hide()

	window = w[0]

	var optionMenu = t.NewMenu(rightClickMenu())

	// Doesn't currently work.
	// TODO: Debug right-click, seems to be intercepted by normal click listener
	t.On(astilectron.EventNameTrayEventRightClicked, func(e astilectron.Event) (deleteListener bool) {
		if menuVisible {
			optionMenu.Destroy()
			menuVisible = false
		} else {
			optionMenu.Create()
			menuVisible = true
		}
		return
	})

	t.On(astilectron.EventNameTrayEventClicked, func(e astilectron.Event) (deleteListener bool) {
		astilog.Info("Tray Clicked")
		if window.IsShown() {
			window.Hide()
		} else {
			window.Show()
		}
		return
	})

	if !onboarded {
		// Show the window on the first launch
		window.Show()
		onboarded = true
	}

	// TODO: ensure that the user has chosen a folder
	if window.IsShown() && !rootFolder {
		showFolderSelect()
	}

	return nil
}

func showFolderSelect() {
	sendMessage("setup")
}

func showQRCode() {
	// TODO, get real PID on start...
	//	pid, err := core.Node.PeerId()
	//	if err != nil {
	//		return "", "", err
	//	}
	//	// create a qr code
	//	url := fmt.Sprintf("https://www.textile.photos/invites/device#id=%s", pid.Pretty())
	var pid = "blah-blah-blah"
	url := fmt.Sprintf("textile://www.textile.photos/invites/device#id=%s", pid)
	png, err := qrcode.Encode(url, qrcode.Medium, QRCodeSize)
	if err != nil {
		astilog.Info("openFolderDialog: TODO store folder dir")
		return
	}

	sendData("pair", map[string]interface{}{
		"qr": base64.StdEncoding.EncodeToString(png),
		"pk": pid,
	})
}

func showWarmup() {
	sendMessage("preready")
}

func showMain() {
	sendMessage("ready")
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
	case "openFolderDialog":
		astilog.Info(fmt.Errorf("openFolderDialog, TODO store folder dir: %s", m.Payload))
		showQRCode()
		// TODO: can use this intermediate step while pairing
		time.AfterFunc(5*time.Second, showWarmup)
		// TODO: complete the pairing
		time.AfterFunc(10*time.Second, showMain)
		return map[string]interface{}{}, nil
	default:
		return map[string]interface{}{}, nil
	}
}

func Exit(e astilectron.Event) (deleteListener bool) {
	astilog.Info("Exit clicked")
	return
}

func rightClickMenu() []*astilectron.MenuItemOptions {
	return []*astilectron.MenuItemOptions{
		{
			Label:   astilectron.PtrStr("Exit"),
			OnClick: Exit,
		},
	}
}

//func start(a *astilectron.Astilectron, w []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {

//window = w[0]
//window.Show()
//
//// get homedir
//home, err := homedir.Dir()
//if err != nil {
//	astilog.Fatal(fmt.Errorf("get homedir failed: %s", err))
//}
//
//// ensure app support folder is created
//var appDir string
//if runtime.GOOS == "darwin" {
//	appDir = filepath.Join(home, "Library/Application Support/Textile")
//} else {
//	appDir = filepath.Join(home, ".textile")
//}
//if err := os.MkdirAll(appDir, 0755); err != nil {
//	return err
//}
//repoPath := filepath.Join(appDir, "repo")
//
//// run init if needed
//if !fsrepo.IsInitialized(repoPath) {
//	accnt := keypair.Random()
//	initc := core.InitConfig{
//		Account:   accnt,
//		RepoPath:  repoPath,
//		LogLevel:  logger.ERROR,
//		LogToDisk: true,
//	}
//	if err := core.InitRepo(initc); err != nil {
//		return err
//	}
//}
//
//// build textile node
//core.Node, err = core.NewTextile(core.RunConfig{RepoPath: repoPath})
//if err != nil {
//	return err
//}
//
//// bring the node online and startup the gateway
//if err := core.Node.Start(); err != nil {
//	return err
//}
//<-core.Node.OnlineCh()
//
//// subscribe to wallet updates
//go func() {
//	for {
//		select {
//		case update, ok := <-core.Node.UpdateCh():
//			if !ok {
//				return
//			}
//			payload := map[string]interface{}{
//				"update": update,
//			}
//			switch update.Type {
//			case core.ThreadAdded:
//				if expanded {
//					sendData("wallet.update", payload)
//				} else {
//					sendPreReady()
//					window.Hide()
//					expandWindow()
//					sendData("wallet.update", payload)
//					window.Show()
//					window.Focus()
//				}
//			default:
//				sendData("wallet.update", payload)
//			}
//		}
//	}
//}()
//
//// subscribe to thread updates
//go func() {
//	for {
//		select {
//		case update, ok := <-core.Node.ThreadUpdateCh():
//			if !ok {
//				return
//			}
//			sendData("thread.update", map[string]interface{}{
//				"update": update,
//			})
//		}
//	}
//}()
//
//// subscribe to notifications
//go func() {
//	for {
//		select {
//		case note, ok := <-core.Node.NotificationCh():
//			if !ok {
//				return
//			}
//			username := core.Node.ContactUsername(note.ActorId)
//			var uinote = a.NewNotification(&astilectron.NotificationOptions{
//				Title: note.Subject,
//				Body:  fmt.Sprintf("%s %s.", username, note.Body),
//				Icon:  "/resources/icon.png",
//			})
//
//			// tmp auto-accept thread invites
//			if note.Type == repo.ReceivedInviteNotification {
//				go func(tid string) {
//					if _, err := core.Node.AcceptThreadInvite(tid); err != nil {
//						astilog.Error(err)
//					}
//				}(note.BlockId)
//			}
//
//			// show notification
//			go func(n *astilectron.Notification) {
//				if err := n.Create(); err != nil {
//					astilog.Error(err)
//					return
//				}
//				if err := n.Show(); err != nil {
//					astilog.Error(err)
//					return
//				}
//			}(uinote)
//		}
//	}
//}()
//
//// start the gateway
//gateway.Host.Start(fmt.Sprintf("127.0.0.1:%s", core.GetRandomPort()))
//
//// save off the server address
//gatewayAddr = fmt.Sprintf("http://%s", gateway.Host.Addr())
//
//// sleep for a bit on the landing screen, it feels better
//time.Sleep(SleepOnLoad)
//
//// send cookie info to front-end
//sendData("login", map[string]interface{}{
//	"name":    "SessionId",
//	"value":   "not used",
//	"gateway": gatewayAddr,
//})
//
//// check if we're configured yet
//threads := core.Node.Threads()
//if len(threads) > 0 {
//	// load threads for UI
//	var threadsJSON []map[string]interface{}
//	for _, thrd := range threads {
//		threadsJSON = append(threadsJSON, map[string]interface{}{
//			"id":   thrd.Id,
//			"name": thrd.Name,
//		})
//	}
//
//	// reveal
//	sendPreReady()
//	window.Hide()
//	expandWindow()
//	sendData("ready", map[string]interface{}{
//		"threads": threadsJSON,
//	})
//	window.Show()
//	window.Focus()
//
//} else {
//	// get qr code for setup
//	qr, pk, err := getQRCode()
//	if err != nil {
//		astilog.Error(err)
//		return err
//	}
//	sendData("setup", map[string]interface{}{
//		"qr": qr,
//		"pk": pk,
//	})
//}
//
//return nil
//}

//
//func handleMessage(_ *astilectron.Window, m bootstrap.MessageIn) (interface{}, error) {
//	switch m.Name {
//	case "refresh":
//		if err := core.Node.CheckCafeMail(); err != nil {
//			return nil, err
//		}
//		return map[string]interface{}{}, nil
//	case "thread.load":
//		var threadId string
//		if err := json.Unmarshal(m.Payload, &threadId); err != nil {
//			return nil, err
//		}
//		html, err := getThreadPhotos(threadId)
//		if err != nil {
//			return nil, err
//		}
//		return map[string]interface{}{
//			"html": html,
//		}, nil
//	default:
//		return map[string]interface{}{}, nil
//	}
//}
//
//func getQRCode() (string, string, error) {
//	// get our own peer id for receiving an account key
//	pid, err := core.Node.PeerId()
//	if err != nil {
//		return "", "", err
//	}
//
//	// create a qr code
//	url := fmt.Sprintf("https://www.textile.photos/invites/device#id=%s", pid.Pretty())
//	png, err := qrcode.Encode(url, qrcode.Medium, QRCodeSize)
//	if err != nil {
//		return "", "", err
//	}
//
//	return base64.StdEncoding.EncodeToString(png), pid.Pretty(), nil
//}
//
//func getThreadPhotos(id string) (string, error) {
//	thrd := core.Node.Thread(id)
//	if thrd == nil {
//		return "", errors.New("thread not found")
//	}
//	var html string
//	query := fmt.Sprintf("threadId='%s' and type=%d", thrd.Id, repo.FileBlock)
//	for _, block := range core.Node.Blocks("", -1, query) {
//		photo := fmt.Sprintf("%s/ipfs/%s/photo?block=%s", gatewayAddr, block.DataId, block.Id)
//		small := fmt.Sprintf("%s/ipfs/%s/small?block=%s", gatewayAddr, block.DataId, block.Id)
//		meta := fmt.Sprintf("%s/ipfs/%s/meta?block=%s", gatewayAddr, block.DataId, block.Id)
//		img := fmt.Sprintf("<img src=\"%s\" />", small)
//		html += fmt.Sprintf(
//			"<div id=\"%s\" class=\"grid-item\" ondragstart=\"imageDragStart(event);\" draggable=\"true\" data-url=\"%s\" data-meta=\"%s\">%s</div>",
//			block.Id, photo, meta, img)
//	}
//	return html, nil
//}
//
//func sendPreReady() {
//	sendMessage("preready")
//	time.Sleep(SleepOnPreReady)
//}
//
//func expandWindow() {
//	expanded = true
//	go window.Resize(InitialWidth, InitialHeight)
//	go window.Center()
//	time.Sleep(SleepOnExpand)
//}

//var p = os.Getenv("GOPATH") + "/src/github.com/textileio/textile-go/desktop"
//
//// Create astilectron
//var a *astilectron.Astilectron
//if a, err = astilectron.New(astilectron.Options{
//	AppName:            "Textile",
//	AppIconDefaultPath: p + "/resources/icon.png",
//	AppIconDarwinPath:  p + "/resources/icon.icns",
//	BaseDirectoryPath:  p,
//}); err != nil {
//	//astilog.Fatal(errors.Wrap(err, "creating new astilectron failed"))
//	astilog.Fatal(fmt.Errorf("creating new astilectron failed: %s", err))
//}
//defer a.Close()
//a.HandleSignals()

// Start
//if err = a.Start(); err != nil {
//	//astilog.Fatal(errors.Wrap(err, "starting failed"))
//	astilog.Fatal(fmt.Errorf("starting failed: %s", err))
//}
//
//// New tray
//var t = a.NewTray(&astilectron.TrayOptions{
//	Image:   astilectron.PtrStr(p + "/gopher.png"),
//	Tooltip: astilectron.PtrStr("Tray's tooltip"),
//})

// New tray menu
//var m = t.NewMenu([]*astilectron.MenuItemOptions{
//	{
//		Label: astilectron.PtrStr("Root 1"),
//		SubMenu: []*astilectron.MenuItemOptions{
//			{Label: astilectron.PtrStr("Item 1")},
//			{Label: astilectron.PtrStr("Item 2")},
//			{Type: astilectron.MenuItemTypeSeparator},
//			{Label: astilectron.PtrStr("Item 3")},
//		},
//	},
//	{
//		Label: astilectron.PtrStr("Root 2"),
//		SubMenu: []*astilectron.MenuItemOptions{
//			{Label: astilectron.PtrStr("Item 1")},
//			{Label: astilectron.PtrStr("Item 2")},
//		},
//	},
//})

//// Create the menu
//if err = m.Create(); err != nil {
//	astilog.Fatal(fmt.Errorf("creating tray menu failed: %s", err))
//	//astilog.Fatal(errors.Wrap(err, "creating tray menu failed"))
//}
//
//// Create tray
//if err = t.Create(); err != nil {
//	astilog.Fatal(fmt.Errorf("creating tray failed: %s", err))
//	//astilog.Fatal(errors.Wrap(err, "creating tray failed"))
//}

// Blocking pattern
//a.Wait()
