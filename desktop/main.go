package main

import (
	"flag"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
	"github.com/textileio/textile-go/mobile"
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

var textile *mobile.Wrapper

func main() {
	// Init
	AppName := "Textile"
	flag.Parse()
	astilog.FlagInit()

	textile = mobile.NewTextile("output/.ipfs")

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
