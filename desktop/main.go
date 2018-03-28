package main

import (
	"flag"
	//"time"
	//
	//"encoding/json"

	textilego "github.com/textileio/textile-go/mobile"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Constants
const htmlAbout = `Welcome on <b>Textile</b> demo!<br>
This is using the bootstrap and the bundler.`

// Vars
var (
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	w       *astilectron.Window
)

var textile *textilego.Node

func main() {
	// Init
	AppName := "Textile"
	flag.Parse()
	astilog.FlagInit()


	textile = textilego.NewTextile("output/.ipfs", "https://ipfs.textile.io")

	err := textile.Start()
	if err != nil {
		astilog.Errorf("start mobile node failed: %s", err)
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
		Debug:    *debug,
		Homepage: "index.html",
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
