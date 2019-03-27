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
		TrayMenuOptions: []*astilectron.MenuItemOptions{
			{
				Type: astilectron.MenuItemTypeSeparator,
			},
		},
		Windows: []*bootstrap.Window{{
			Homepage: "../index.html",
			MessageHandler: handleMessage,
			Options: &astilectron.WindowOptions{
				Show: astilectron.PtrBool(true),
				Movable: astilectron.PtrBool(true),
				Width: astilectron.PtrInt(500),
				Height: astilectron.PtrInt(500),
				Focusable: astilectron.PtrBool(true),
				SkipTaskbar:   astilectron.PtrBool(true),
				TitleBarStyle: astilectron.TitleBarStyleHidden,
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "bootstrap failed"))
	}
}
