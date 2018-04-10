package main

import (
	"flag"
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
	textile, err = core.NewNode("output/.ipfs", false)

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

	// tmp: sub to own peer id for pairing setup
	// this should really only happen when you click "pair"
	// and then be closed
	var errc2 = make(chan error)
	go func() {
		errc2 <- textile.StartPairing()
		close(errc2)
	}()

	// tmp: if paired id is present, start listening
	pid, err := textile.Datastore.Config().GetPairedID()
	if err != nil {
		astilog.Errorf("get paired id failed: %s", err)
	}
	if pid != "" {
		var errc3 = make(chan error)
		go func() {
			errc3 <- textile.StartSync(pid)
			close(errc3)
		}()
	}

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
		RestoreAssets:  RestoreAssets,
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
