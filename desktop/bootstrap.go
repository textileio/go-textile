package main

import (
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// bootstrapApp runs bootstrap. Moved to own file so we don't have to see the error highlighting :)
func bootstrapApp() {
	astilog.Debugf("Running app built at %s", builtAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset: Asset,
		AstilectronOptions: astilectron.Options{
			AppName:            appName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
			// TODO: Revisit: slightly dangerous because this will ignore _all_ certificate errors
			ElectronSwitches: []string{"ignore-certificate-errors", "true"},
		},

		Debug:          *debug,
		Homepage:       "index.html",
		MessageHandler: handleMessage,
		OnWait:         start,
		RestoreAssets:  RestoreAssets,
		WindowOptions: &astilectron.WindowOptions{
			Center:          astilectron.PtrBool(true),
			Height:          astilectron.PtrInt(SetupSize),
			Width:           astilectron.PtrInt(SetupSize),
			BackgroundColor: astilectron.PtrStr("#ffffff"),
			TitleBarStyle:   astilectron.TitleBarStyleHiddenInset,
		},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "bootstrap failed"))
	}
}
