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
	"github.com/textileio/textile-go/wallet"
)

var (
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
)

var textile *core.TextileNode
var gateway string

func main() {
	AppName := "Textile"
	flag.Parse()
	astilog.FlagInit()

	// create a desktop textile node
	// TODO: on darwin, repo should live in Application Support
	config := core.NodeConfig{
		LogLevel: logging.DEBUG,
		LogFiles: true,
		WalletConfig: wallet.Config{
			RepoPath:   "output/.ipfs",
			CentralAPI: "https://api.textile.io",
			IsMobile:   false,
		},
	}
	var err error
	textile, _, err = core.NewNode(config)
	if err != nil {
		astilog.Errorf("create desktop node failed: %s", err)
		return
	}

	// bring the node online and startup the gateway
	online, err := textile.StartWallet()
	if err != nil {
		astilog.Errorf("start desktop node failed: %s", err)
		return
	}
	<-online

	err = textile.StartGateway()
	if err != nil {
		astilog.Errorf("start gateway failed: %s", err)
		return
	}

	// save off the gateway address
	gateway = fmt.Sprintf("http://localhost%s", textile.GetGatewayAddress())

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
