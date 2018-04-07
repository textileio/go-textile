package main

import (
	"encoding/json"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/asticode/go-astilog"
	"encoding/base64"
)
// Init exploration
type IpfsResponse struct {
	path string `json:"path"`
	data string `json:"data"`
}
// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "peer.qr":

		var png []byte
		// I've registered this URL so that apple will do an App Link from any url like it directly into our
		// app. Just need to do a PR in the app to receive it
		png, err := qrcode.Encode("https://www.textile.io/threads/pair=testing", qrcode.Medium, 256)
		if err != nil {
			astilog.Errorf("qr generation failed: %s", err)
			return err.Error(), err
		}
		payload = base64.StdEncoding.EncodeToString(png)

	case "ipfs.getPath":
		// Unmarshal payload
		var path string
		if err = json.Unmarshal(m.Payload, &path); err != nil {
			return err.Error(), err
		}

		photoBase, _ := textile.GetPhotoBase64String(path)
		if err != nil {
			return err.Error(), err
		} else {
			return IpfsResponse {
				path: path,
				data: photoBase,
			}, nil
		}
	}
	return
}