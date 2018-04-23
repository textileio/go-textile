package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/skip2/go-qrcode"
)

// handleMessages handles messages
func handleMessages(iw *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "pair.start":
		astilog.Info("PAIRING STARTED")
		astilog.Info("GENERATING QR CODE")

		// create a random confirmation code
		code := fmt.Sprintf("%04d", rand.Int63n(1e4))

		// get our own rsa public key
		pk, err := textile.GetPublicPeerKeyString()
		if err != nil {
			astilog.Errorf("public key generation failed: %s", err)
			return nil, err
		}

		// create a qr code
		url := fmt.Sprintf("https://www.textile.io/clients?code=%s&key=%s", code, pk)
		png, err := qrcode.Encode(url, qrcode.Medium, 256)
		if err != nil {
			astilog.Errorf("qr generation failed: %s", err)
			return nil, err
		}

		// pass the qr code and info back to app
		return map[string]interface{}{
			"png":  base64.StdEncoding.EncodeToString(png),
			"code": code,
			"url":  url,
			"key":  pk,
		}, nil

	case "sync.start":
		astilog.Info("GOT START SYNC MESSAGE")

		// finally, start syncing
		go joinRoom(iw)

		// return empty response
		return map[string]interface{}{}, nil

	case "login.request":
		astilog.Info("GOT LOGIN REQUEST MESSAGE")
		// TODO: Make this more secure with salt, SAH hashing, and additional randomness?
		return map[string]interface{}{
			"name":  "SessionId",
			"value": textile.Password,
		}, nil
	}

	return
}
