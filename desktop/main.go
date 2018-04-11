package main

import (
	"flag"
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
	"github.com/textileio/textile-go/core"
)

var (
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	w       *astilectron.Window
)

var textile *core.TextileNode

func main() {
	// Init
	AppName := "Textile"
	flag.Parse()
	astilog.FlagInit()

	// Create a desktop textile node
	// TODO: on darwin, I think repo should live in Application Support
	var err error
	textile, _ = core.NewNode("output/.ipfs", false)

	// Bring the node online
	err = textile.Start()
	if err != nil {
		astilog.Errorf("start mobile node failed: %s", err)
	}

	// Start garbage collection and gateway services
	// NOTE: on desktop, gateway runs on 8081
	var errc = make(chan error)
	go func() {
		errc <- textile.StartServices()
		close(errc)
	}()

	pairId, _ := textile.Datastore.Config().GetPairedID()
	//if err != nil {
	//	astilog.Errorf("get paired id failed: %s", err)
	//}

	// Run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset: Asset,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug:          *debug,
		Homepage:       "index.html",
		MessageHandler: handleMessages,
		OnWait: func(_ *astilectron.Astilectron, iw *astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			//w = iw
			if pairId != "" {
				// tmp: if paired id is present, start listening
				var errc3 = make(chan error)
				go func() {
					errc3 <- textile.StartSync(pairId)
					close(errc3)
				}()

				iw.SendMessage(map[string]string{"name": "ready", "gateway": "http://localhost:9192"})

				iw.SendMessage(map[string]string{
					"name":      "sync.new",
					"hash":      "QmbZPvPb9kmsRJDC5Zi8Jw5M6stT7J6i7L8FkPJYbAfE3F",
					"timestamp": "today",
				})
				// TODO:
				// Each time a new hash is delivered, issue
				// iw.SendMessage(map[string]string{"name": "sync.new", "hash": hash})

				// Grab all the current hashes in the datastore and populate the gallery
				for _, photo := range textile.Datastore.Photos().GetPhotos("-1", 10) {
					fmt.Printf("Photo: %s", photo.Cid)
					iw.SendMessage(map[string]string{
						"name":      "sync.new",
						"hash":      photo.Cid,
						"timestamp": photo.Timestamp.String(),
					})
				}
			} else {
				// tmp: sub to own peer id for pairing setup
				// this should really only happen when you click "pair"
				// and then be closed
				var errc2 = make(chan error)
				go func() {
					errc2 <- textile.StartPairing()
					// TODO: When this completes we should send the onboard.complete message
					// iw.SendMessage(map[string]string{"name": "onboard.complete"})
					close(errc2)
				}()
				/*
					If the desktop client hasn't already paired with a mobile device, tell the UI to init onboarding
				*/
				iw.SendMessage(map[string]string{"name": "onboard"})

				// TODO:
				// TODO: When StartPairing is complete, we should issue the message
				// iw.SendMessage(map[string]string{"name": "pairing.complete"})
			}
			return nil
		},
		RestoreAssets: RestoreAssets,
		WindowOptions: &astilectron.WindowOptions{
			BackgroundColor: astilectron.PtrStr("#333"),
			Center:          astilectron.PtrBool(true),
			Height:          astilectron.PtrInt(700),
			Width:           astilectron.PtrInt(700),
		},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}
