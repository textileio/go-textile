package main

import (
	"errors"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
)

// bootstrapApp runs bootstrap. Moved to own file so we don't have to see Asset and RestoreAsset highlighed as errors :)
func bootstrapApp() {
	astilog.Debugf("Running app built at %s", builtAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            appName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug:         *debug,
		OnWait:        start,
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessage,
			Options: &astilectron.WindowOptions{
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(SetupSize),
				Width:           astilectron.PtrInt(SetupSize),
				BackgroundColor: astilectron.PtrStr("#ffffff"),
				TitleBarStyle:   astilectron.TitleBarStyleHiddenInset,
				Show:            astilectron.PtrBool(false),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "bootstrap failed"))
	}
}
