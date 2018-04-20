package main

import (
	"fmt"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
)

var gateway = "http://localhost:9182"

func start(_ *astilectron.Astilectron, iw *astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
	astilog.Info("TEXTILE STARTED")

	// check if we're configured yet
	if textile.IsDatastoreConfigured() {
		// can join room
		astilog.Info("ALREADY CONFIGURED")

		// tell app we're ready and send initial html
		sendData(iw, "sync.ready", map[string]interface{}{
			"html": getPhotosHTML(),
		})

	} else {
		// otherwise, start onboaring
		astilog.Info("NOT CONFIGURED")
		astilog.Info("STARTING PAIRING")

		go func() {
			// sub to own peer id for pairing setup and wait
			textile.WaitForRoom()

			if !textile.IsDatastoreConfigured() {
				astilog.Error("failed to join room")
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

	datac := make(chan string)
	go textile.JoinRoom(datac)
	for {
		select {
		case hash, ok := <-datac:
			sendData(iw, "sync.data", map[string]interface{}{
				"hash": hash,
			})
			if !ok {
				return nil
			}
		}
	}
}

func getPhotosHTML() string {
	var html string
	for _, photo := range textile.Datastore.Photos().GetPhotos("", -1) {
		ph := fmt.Sprintf("%s/ipfs/%s/photo", gateway, photo.Cid)
		th := fmt.Sprintf("%s/ipfs/%s/thumb", gateway, photo.Cid)
		md := fmt.Sprintf("%s/ipfs/%s/meta", gateway, photo.Cid)
		img := fmt.Sprintf("<img src=\"%s\" />", th)
		html += fmt.Sprintf("<div class=\"grid-item\" data-url=\"%s\" data-meta=\"%s\">%s</div>", ph, md, img)
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
