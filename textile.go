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
		dataDir = filepath.Join(hd, ".textile")
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
		Name: "id",
		Help: "show peer id",
		Func: cmd.ShowId,
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "peers",
		Help: "show connected peers (same as `ipfs swarm peers`)",
		Func: cmd.SwarmPeers,
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "ping",
		Help: "ping a peer (same as `ipfs ping`)",
		Func: cmd.SwarmPing,
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "connect",
		Help: "connect to a peer (same as `ipfs swarm connect`)",
		Func: cmd.SwarmConnect,
	})
	{
		photoCmd := &ishell.Cmd{
			Name:     "photo",
			Help:     "manage photos",
			LongHelp: "Add, list, and get info about photos.",
		}
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "add",
			Help: "add a new photo (default thread is \"#default\")",
			Func: cmd.AddPhoto,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "share",
			Help: "share a photo to a different thread",
			Func: cmd.SharePhoto,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "get",
			Help: "save a photo to a local file",
			Func: cmd.GetPhoto,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "ls",
			Help: "list photos from a thread (defaults to \"#default\")",
			Func: cmd.ListPhotos,
		})
		shell.AddCmd(photoCmd)
	}
	{
		albumsCmd := &ishell.Cmd{
			Name:     "thread",
			Help:     "manage photo threads",
			LongHelp: "Add, list, enable, disable, and get info about photo threads.",
		}
		albumsCmd.AddCmd(&ishell.Cmd{
			Name: "add",
			Help: "add a new thread",
			Func: cmd.CreateAlbum,
		})
		albumsCmd.AddCmd(&ishell.Cmd{
			Name: "ls",
			Help: "list threads",
			Func: cmd.ListAlbums,
		})
		albumsCmd.AddCmd(&ishell.Cmd{
			Name: "enable",
			Help: "enable a thread",
			Func: cmd.EnableAlbum,
		})
		albumsCmd.AddCmd(&ishell.Cmd{
			Name: "disable",
			Help: "disable a thread",
			Func: cmd.DisableAlbum,
		})
		albumsCmd.AddCmd(&ishell.Cmd{
			Name: "mnemonic",
			Help: "show mnemonic phrase",
			Func: cmd.AlbumMnemonic,
		})
		albumsCmd.AddCmd(&ishell.Cmd{
			Name: "publish",
			Help: "publish latest update",
			Func: cmd.RepublishAlbum,
		})
		albumsCmd.AddCmd(&ishell.Cmd{
			Name: "peers",
			Help: "list peers",
			Func: cmd.ListAlbumPeers,
		})
		shell.AddCmd(albumsCmd)
	}

	// create a desktop textile node
	// TODO: darwin should use App. Support dir, not home dir
	node, err := core.NewNode(dataDir, false, logging.DEBUG)
	if err != nil {
		shell.Println(fmt.Errorf("create desktop node failed: %s", err))
		return
	}
	core.Node = node
	// start node and https gateway server
	if err = core.Node.Start(); err != nil {
		shell.Println(fmt.Errorf("start desktop node failed: %s", err))
		return
	}

	// start garbage collection
	go startServices()

	// join existing rooms
	albums := core.Node.Datastore.Albums().GetAlbums("")
	for _, a := range albums {
		go core.Node.JoinRoom(a.Id, make(chan string))
	}

	// run shell
	shell.Run()
}

// Start garbage collection
func startServices() {
	errc, err := core.Node.StartServices()
	if err != nil {
		shell.Println(fmt.Errorf("start service error: %s", err))
		return
	}

	// TODO: if we want the 'raw' gateway server when running from cmd line
	// TODO: that can be added here directly (would need to export method)

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
	shell.Println("textile node v" + core.VERSION)
	shell.Println("type `help` for available commands")
}
