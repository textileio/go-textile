package main

import (
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/wallet/thread"
)

var mobileThread *thread.Thread

func start(_ *astilectron.Astilectron, iw *astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
	astilog.Info("TEXTILE STARTED")

	astilog.Info("SENDING COOKIE INFO")
	sendData(iw, "login.cookie", map[string]interface{}{
		"name":    "SessionId",
		"value":   "not used",
		"gateway": gateway,
	})

	// check if we're configured yet
	mobileThread := textile.Wallet.GetThreadByName("mobile")
	if mobileThread != nil {
		astilog.Info("FOUND MOBILE THREAD")

		// tell app we're ready and send initial html
		sendData(iw, "sync.ready", map[string]interface{}{
			"html": getPhotosHTML(),
		})

	} else {
		// otherwise, start onboaring
		astilog.Info("COULD NOT FIND MOBILE THREAD")
		astilog.Info("STARTING PAIRING")

		go func() {
			// sub to owr peer id for pairing setup and wait
			//textile.Wallet.WaitForInvite()

			mobileThread := textile.Wallet.GetThreadByName("mobile")
			if mobileThread == nil {
				astilog.Error("failed to create mobile album")
				return
			}

			// let the app know we're done pairing
			sendMessage(iw, "onboard.complete")

			// and that we're ready to go
			sendData(iw, "sync.ready", map[string]interface{}{
				"html": getPhotosHTML(),
			})
		}()
		sendMessage(iw, "onboard.start")
	}

	return nil
}

func joinRoom(iw *astilectron.Window) error {
	astilog.Info("STARTING SYNC")

	datac := make(chan thread.Update)
	go mobileThread.Subscribe(datac)
	for {
		select {
		case update, ok := <-datac:
			if !ok {
				return nil
			}
			sendData(iw, "sync.data", map[string]interface{}{
				"update":  update,
				"gateway": gateway,
			})
		}
	}
}

func getPhotosHTML() string {
	var html string
	for _, block := range mobileThread.Blocks("", -1, repo.PhotoBlock) {
		ph := fmt.Sprintf("%s/ipfs/%s/photo?block=%s", gateway, block.Target, block.Id)
		th := fmt.Sprintf("%s/ipfs/%s/thumb?block=%s", gateway, block.Target, block.Id)
		md := fmt.Sprintf("%s/ipfs/%s/meta?block=%s", gateway, block.Target, block.Id)
		img := fmt.Sprintf("<img src=\"%s\" />", th)
		html += fmt.Sprintf("<div id=\"%s\" class=\"grid-item\" ondragstart=\"imageDragStart(event);\" draggable=\"true\" data-url=\"%s\" data-meta=\"%s\">%s</div>", block.Id, ph, md, img)
	}
	return html
}

func sendMessage(iw *astilectron.Window, name string) {
	iw.SendMessage(map[string]string{"name": name})
}

func sendData(iw *astilectron.Window, name string, data map[string]interface{}) {
	data["name"] = name
	iw.SendMessage(data)
}
