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

type Opts struct {
	Version bool   `short:"v" long:"version" description:"print the version number and exit"`
	DataDir string `short:"d" long:"datadir" description:"specify the data directory to be used"`
}

var Options Opts
var parser = flags.NewParser(&Options, flags.Default)

var shell *ishell.Shell

func main() {
	// create a new shell
	shell = ishell.New()

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
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)

			// get existing mnemonic
			c.Print("mnemonic phrase (optional): ")
			mnemonic := c.ReadLine()

			// configure
			err := core.Node.ConfigureDatastore(mnemonic)
			if err != nil {
				c.Err(fmt.Errorf("configure node datastore failed: %s", err))
				return
			}

			// start services
			if !core.Node.ServicesUp {
				go startServices()

				// leave old wallet room
				// TODO: need a diff way to determine if we prev. had a room subscription
				// TODO: this is ugly relying on ServiceUp
				core.Node.LeaveRoom()
				<-core.Node.LeftRoomCh
			}

			// join new room
			go core.Node.JoinRoom()
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "id",
		Help: "show node ids",
		Func: cmd.GetIds,
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
	node, err := core.NewNode(dataDir, false, logging.DEBUG)
	if err != nil {
		shell.Println(fmt.Errorf("create desktop node failed: %s", err))
		return
	}
	core.Node = node
	if err = core.Node.Start(); err != nil {
		shell.Println(fmt.Errorf("start desktop node failed: %s", err))
		return
	}

	// if datastore is configured, start services
	if core.Node.IsDatastoreConfigured() {
		go startServices()
		go core.Node.JoinRoom()
	}

	// run shell
	shell.Run()
}

// Start garbage collection and gateway services
// NOTE: on desktop, gateway runs on 8182, decrypting file gateway on 9192
// TODO: make this cancelable
func startServices() {
	errc, err := core.Node.StartServices()
	if err != nil {
		shell.Println(fmt.Errorf("start service error: %s", err))
		return
	}

	for {
		select {
		case err := <-errc:
			if err != nil {
				shell.Println(fmt.Errorf("service error: %s", err))
			}
		}
	}
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
	shell.Println("type `help` for available commands")
}
