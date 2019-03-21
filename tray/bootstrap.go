package main

import (
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// bootstrapApp runs bootstrap. Moved to own file so we don't have to see Asset and RestoreAsset highlighed as errors :)
func bootstrapApp() {
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            appName,
			AppIconDarwinPath:  "resources/assets/icon.icns",
			AppIconDefaultPath: "resources/assets/icon.png",
		},
		Debug:         *debug,
		OnWait:        start,
		RestoreAssets: RestoreAssets,
		TrayOptions:   &astilectron.TrayOptions{
			Image:   astilectron.PtrStr("resources/assets/tray.png"),
			Tooltip: astilectron.PtrStr("Textile"),
		},
		Windows: []*bootstrap.Window{{
			Homepage:       "../index.html",
			MessageHandler: handleMessage,
			Options: &astilectron.WindowOptions{
				// AlwaysOnTop:   astilectron.PtrBool(true),
				// Resizable:     astilectron.PtrBool(false),
				// Closable:      astilectron.PtrBool(false),
				// Maximizable:   astilectron.PtrBool(false),
				// Minimizable:   astilectron.PtrBool(false),
				// SkipTaskbar:   astilectron.PtrBool(false),
				// TitleBarStyle: astilectron.TitleBarStyleHiddenInset,
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "bootstrap failed"))
	}
}
