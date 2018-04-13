// TODO: use lumberjack logger, not stdout, see #33
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"
	"gopkg.in/abiosoft/ishell.v2"

	"github.com/textileio/textile-go/cmd"
	"github.com/textileio/textile-go/core"
)

var log = logging.MustGetLogger("main")

type Opts struct {
	Version bool   `short:"v" long:"version" description:"print the version number and exit"`
	DataDir string `short:"d" long:"datadir" description:"specify the data directory to be used"`
}

var Options Opts
var parser = flags.NewParser(&Options, flags.Default)

func main() {
	// create a new shell
	shell := ishell.New()

	// handle version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		shell.Println(core.VERSION)
		return
	}

	// handle data dir
	var dataDir string
	if len(os.Args) > 1 && (os.Args[1] == "--datadir" || os.Args[1] == "-d") {
		if len(os.Args) < 3 {
			shell.Println(errors.New("datadir option provided but missing value"))
			return
		}
		dataDir = os.Args[2]
	} else {
		hd, err := homedir.Dir()
		if err != nil {
			shell.Println(errors.New("could not determine home directory"))
			return
		}
		dataDir = filepath.Join(hd, ".ipfs")
	}

	// parse flags
	if _, err := parser.Parse(); err != nil {
		return
	}

	// handle interrupt
	// TODO: shutdown on 'exit' command too
	shell.Interrupt(func(c *ishell.Context, count int, input string) {
		if count == 1 {
			shell.Println("input Ctrl-C once more to exit")
			return
		}
		shell.Println("interrupted")
		shell.Printf("textile server shutting down...")
		if core.Node.IpfsNode != nil {
			core.Node.Stop()
		}
		shell.Printf("done\n")
		os.Exit(1)
	})

	// welcome
	printSplashScreen(shell)

	// add commands
	shell.AddCmd(&ishell.Cmd{
		Name: "init",
		Help: "initialize a new datastore",
		Func: func(c *ishell.Context) {
			err := core.Node.ConfigureDatastore("", "")
			if err != nil {
				c.Err(fmt.Errorf("configure node datastore failed: %s", err))
				return
			}
		},
	})

	{
		photosCmd := &ishell.Cmd{
			Name:     "photos",
			Help:     "manage wallet photos",
			LongHelp: `Manage your textile wallet photos.`,
		}
		photosCmd.AddCmd(&ishell.Cmd{
			Name: "add",
			Help: "add a new photo",
			Func: cmd.AddPhoto,
		})

		shell.AddCmd(photosCmd)
	}

	// create a desktop textile node
	// TODO: darwin should use App. Support dir, not home dir
	node, err := core.NewNode(dataDir, false)
	if err != nil {
		shell.Println(fmt.Errorf("create desktop node failed: %s", err))
		return
	}
	core.Node = node
	if err = core.Node.Start(); err != nil {
		shell.Println(fmt.Errorf("create desktop node failed: %s", err))
		return
	}

	// Start garbage collection and gateway services
	// NOTE: on desktop, gateway runs on 8182, decrypting file gateway on 9192
	var servErrc = make(chan error)
	go func() {
		servErrc <- core.Node.StartServices()
		close(servErrc)
	}()
	go func() {
		for {
			select {
			case err := <-servErrc:
				shell.Println(fmt.Errorf("server error: %s", err))
			}
		}
	}()

	// run shell
	shell.Run()
}

func printSplashScreen(shell *ishell.Shell) {
	blue := color.New(color.FgBlue).SprintFunc()
	banner :=
		`
  __                   __  .__.__          
_/  |_  ____ ___  ____/  |_|__|  |   ____  
\   __\/ __ \\  \/  /\   __\  |  | _/ __ \ 
 |  | \  ___/ >    <  |  | |  |  |_\  ___/ 
 |__|  \___  >__/\_ \ |__| |__|____/\___  >
           \/      \/                   \/ 
`
	shell.Println(blue(banner))
	shell.Println("")
	shell.Println("textile server v" + core.VERSION)
	shell.Println("")
}
