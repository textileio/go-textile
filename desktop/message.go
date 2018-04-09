package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/skip2/go-qrcode"
)

// Init exploration
type IpfsResponse struct {
	path string `json:"path"`
	data []byte `json:"data"`
}
type QRCodeResponse struct {
	png  string `json:"png"`
	code string `json:"code"`
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "peer.qr":

		// create a random confirmation code
		code := fmt.Sprintf("%04d", rand.Int63n(1e4))

		pubKey, err := textile.GetPublicKey()
		if err != nil {
			astilog.Errorf("public key generation failed: %s", err)
			return nil, err
		}
		pubKeyS := string(base64.StdEncoding.EncodeToString(pubKey))

		var png []byte
		// I've registered this URL so that apple will do an App Link from any url like it directly into our
		// app. Just need to do a PR in the app to receive it
		png, err = qrcode.Encode(fmt.Sprintf("https://www.textile.io/clients?code=%s&key=%s", code, pubKeyS), qrcode.Medium, 256)
		if err != nil {
			astilog.Errorf("qr generation failed: %s", err)
			return nil, err
		}
		res := map[string]interface{}{
			"png":  string(base64.StdEncoding.EncodeToString(png)),
			"code": code,
			"url":  fmt.Sprintf("https://www.textile.io/clients?code=%s&key=%s", code, pubKeyS),
			"key":  pubKeyS,
		}
		return res, nil

	case "ipfs.getPath":
		// Unmarshal payload
		var path string
		if err = json.Unmarshal(m.Payload, &path); err != nil {
			return err.Error(), err
		}

		file, err := textile.GetFile(path)
		if err != nil {
			return err.Error(), err
		}
		return IpfsResponse{path: path, data: file}, nil
	}

	return
}
