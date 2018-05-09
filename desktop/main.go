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
	// TODO: make api url configurable somehow
	config := core.NodeConfig{
		RepoPath:      "output/.ipfs",
		CentralApiURL: "https://api.textile.io",
		IsMobile:      false,
		LogLevel:      logging.DEBUG,
		LogFiles:      true,
	}
	var err error
	textile, err = core.NewNode(config)
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
	gateway = fmt.Sprintf("http://localhost%s", textile.GatewayProxy.Addr)

	// start garbage collection and gateway services
	// TODO: see method todo before enabling
	//errc, err := textile.StartGarbageCollection()
	//if err != nil {
	//	astilog.Errorf("auto gc error: %s", err)
	//	return
	//}
	//go func() {
	//	for {
	//		select {
	//		case err, ok := <-errc:
	//			if err != nil {
	//				astilog.Errorf("auto gc error: %s", err)
	//			}
	//			if !ok {
	//				astilog.Info("auto gc stopped")
	//				return
	//			}
	//		}
	//	}
	//}()

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
