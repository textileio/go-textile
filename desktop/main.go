package main

import (
	"flag"
	"fmt"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/op/go-logging"
	"github.com/pkg/errors"

	"github.com/textileio/textile-go/core"
)

var (
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	w       *astilectron.Window
)

var textile *core.TextileNode
var gateway string

func main() {
	AppName := "Textile"
	flag.Parse()
	astilog.FlagInit()

	// create a desktop textile node
	// TODO: on darwin, repo should live in Application Support
	var err error
	textile, err = core.NewNode("output/.ipfs", false, logging.DEBUG)
	if err != nil {
		astilog.Errorf("create desktop node failed: %s", err)
		return
	}

	// bring the node online and startup the gateway
	err = textile.Start()
	if err != nil {
		astilog.Errorf("start desktop node failed: %s", err)
		return
	}

	// save off the gateway address
	gatewayPort, err := textile.GatewayPort()
	if err != nil {
		astilog.Errorf("get gateway port failed: %s", err)
		return
	}
	gateway = fmt.Sprintf("https://localhost:%d", gatewayPort)

	// start garbage collection and gateway services
	errc, err := textile.StartServices()
	if err != nil {
		astilog.Errorf("start service error: %s", err)
		return
	}
	go func() {
		for {
			select {
			case err := <-errc:
				if err != nil {
					astilog.Errorf("service error: %s", err)
				}
			}
		}
	}()

	// run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset: Asset,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
			// TODO: Revisit: slightly dangerous because this will ignore _all_ certificate errors
			ElectronSwitches: []string{"ignore-certificate-errors", "true"},
		},

		Debug:          *debug,
		Homepage:       "index.html",
		MessageHandler: handleMessages,
		OnWait:         start,
		RestoreAssets:  RestoreAssets,
		WindowOptions: &astilectron.WindowOptions{
			BackgroundColor: astilectron.PtrStr("#333333"),
			Center:          astilectron.PtrBool(true),
			Height:          astilectron.PtrInt(633),
			Width:           astilectron.PtrInt(1024),
		},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}
