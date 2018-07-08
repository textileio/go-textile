package main

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/cmd"
	"github.com/textileio/textile-go/core"
	"github.com/textileio/textile-go/wallet"
	"github.com/textileio/textile-go/wallet/thread"
	"gopkg.in/abiosoft/ishell.v2"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

type Opts struct {
	DataDir    string `short:"d" long:"datadir" description:"specify the data directory to be used"`
	DaemonMode bool   `short:"m" long:"daemon" description:"start in a non-interactive daemon mode"`
	ServerMode bool   `short:"s" long:"server" description:"start in server mode"`
	LogLevel   string `short:"l" long:"loglevel" description:"set the logging level [debug, info, notice, warning, error, critical]" default:"debug"`
	NoLogFiles bool   `short:"f" long:"logfiles" description:"do not save logs on disk"`
	Version    bool   `short:"v" long:"version" description:"print the version number and exit"`
	ApiPort    string `short:"p" long:"apiport" description:"set the api port (daemon only)" default:"3000"`
}

var Options Opts
var parser = flags.NewParser(&Options, flags.Default)

var shell *ishell.Shell

func main() {
	// parse flags
	if _, err := parser.Parse(); err != nil {
		return
	}

	// handle version flag
	if Options.Version {
		fmt.Println(core.Version)
		return
	}

	// handle data dir
	dataDir := Options.DataDir
	if len(Options.DataDir) == 0 {
		// get homedir
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(fmt.Errorf("get homedir failed: %s", err.Error()))
			return
		}

		// ensure app folder is created
		appDir := filepath.Join(home, ".textile")
		if err := os.MkdirAll(appDir, 0755); err != nil {
			fmt.Println(fmt.Errorf("create repo directory failed: %s", err.Error()))
			return
		}

		dataDir = filepath.Join(appDir, "repo")
	}

	// determine log level
	level, err := logging.LogLevel(strings.ToUpper(Options.LogLevel))
	if err != nil {
		fmt.Println(fmt.Errorf("determine log level failed: %s", err))
		return
	}

	// create and start a desktop textile node
	config := core.NodeConfig{
		LogLevel: level,
		LogFiles: !Options.NoLogFiles,
		WalletConfig: wallet.Config{
			RepoPath:   dataDir,
			CentralAPI: "https://api.textile.io",
			IsMobile:   false,
			IsServer:   Options.ServerMode,
		},
	}
	node, _, err := core.NewNode(config)
	if err != nil {
		fmt.Println(fmt.Errorf("create desktop node failed: %s", err))
		return
	}
	core.Node = node

	// auto start it
	if err := start(); err != nil {
		fmt.Println(fmt.Errorf("start desktop node failed: %s", err))
	}

	// welcome
	printSplashScreen()

	// run it
	if Options.DaemonMode {

		// handle interrupt
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		fmt.Println("interrupted")
		fmt.Printf("textile node shutting down...")
		if err := stop(); err != nil && err != wallet.ErrStopped {
			fmt.Println(err.Error())
		} else {
			fmt.Print("done\n")
		}
		os.Exit(1)

	} else {

		// create a new shell
		shell = ishell.New()

		// handle interrupt
		shell.Interrupt(func(c *ishell.Context, count int, input string) {
			if count == 1 {
				shell.Println("input Ctrl-C once more to exit")
				return
			}
			shell.Println("interrupted")
			shell.Printf("textile node shutting down...")
			if err := stop(); err != nil && err != wallet.ErrStopped {
				c.Err(err)
			} else {
				shell.Printf("done\n")
			}
			os.Exit(1)
		})

		// add commands
		shell.AddCmd(&ishell.Cmd{
			Name: "start",
			Help: "start the node",
			Func: func(c *ishell.Context) {
				if core.Node.Wallet.Started() {
					c.Println("already started")
					return
				}
				if err := start(); err != nil {
					c.Println(fmt.Errorf("start desktop node failed: %s", err))
					return
				}
				c.Println("ok, started")
			},
		})
		shell.AddCmd(&ishell.Cmd{
			Name: "stop",
			Help: "stop the node",
			Func: func(c *ishell.Context) {
				if !core.Node.Wallet.Started() {
					c.Println("already stopped")
					return
				}
				if err := stop(); err != nil {
					c.Println(fmt.Errorf("stop desktop node failed: %s", err))
					return
				}
				c.Println("ok, stopped")
			},
		})
		shell.AddCmd(&ishell.Cmd{
			Name: "id",
			Help: "show node id",
			Func: cmd.ShowId,
		})
		shell.AddCmd(&ishell.Cmd{
			Name: "ping",
			Help: "ping another textile node",
			Func: func(c *ishell.Context) {
				if !core.Node.Wallet.Online() {
					c.Println("not online yet")
					return
				}
				if len(c.Args) == 0 {
					c.Err(errors.New("missing node id"))
					return
				}
				status, err := core.Node.Wallet.GetPeerStatus(c.Args[0])
				if err != nil {
					c.Println(fmt.Errorf("ping failed: %s", err))
					return
				}
				c.Println(status)
			},
		})
		{
			swarmCmd := &ishell.Cmd{
				Name:     "swarm",
				Help:     "same as ipfs swarm",
				LongHelp: "Inspect IPFS swarm peers.",
			}
			swarmCmd.AddCmd(&ishell.Cmd{
				Name: "peers",
				Help: "show connected peers (same as `ipfs swarm peers`)",
				Func: cmd.SwarmPeers,
			})
			swarmCmd.AddCmd(&ishell.Cmd{
				Name: "ping",
				Help: "ping a peer (same as `ipfs ping`)",
				Func: cmd.SwarmPing,
			})
			swarmCmd.AddCmd(&ishell.Cmd{
				Name: "connect",
				Help: "connect to a peer (same as `ipfs swarm connect`)",
				Func: cmd.SwarmConnect,
			})
			shell.AddCmd(swarmCmd)
		}
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
				Name: "key",
				Help: "decrypt and print the key for a photo",
				Func: cmd.GetPhotoKey,
			})
			photoCmd.AddCmd(&ishell.Cmd{
				Name: "meta",
				Help: "cat photo metadata",
				Func: cmd.CatPhotoMetadata,
			})
			photoCmd.AddCmd(&ishell.Cmd{
				Name: "ls",
				Help: "list photos from a thread (defaults to \"#default\")",
				Func: cmd.ListPhotos,
			})
			shell.AddCmd(photoCmd)
		}
		{
			threadCmd := &ishell.Cmd{
				Name:     "thread",
				Help:     "manage threads",
				LongHelp: "Add, remove, list, invite to, and get info about textile threads.",
			}
			threadCmd.AddCmd(&ishell.Cmd{
				Name: "add",
				Help: "add a new thread",
				Func: cmd.AddThread,
			})
			threadCmd.AddCmd(&ishell.Cmd{
				Name: "rm",
				Help: "remove a thread by name",
				Func: cmd.RemoveThread,
			})
			threadCmd.AddCmd(&ishell.Cmd{
				Name: "ls",
				Help: "list threads",
				Func: cmd.ListThreads,
			})
			threadCmd.AddCmd(&ishell.Cmd{
				Name: "publish",
				Help: "publish latest update",
				Func: cmd.PublishThread,
			})
			threadCmd.AddCmd(&ishell.Cmd{
				Name: "peers",
				Help: "list peers",
				Func: cmd.ListThreadPeers,
			})
			threadCmd.AddCmd(&ishell.Cmd{
				Name: "invite",
				Help: "invite a peer to a thread",
				Func: cmd.AddThreadInvite,
			})
			shell.AddCmd(threadCmd)
		}
		{
			deviceCmd := &ishell.Cmd{
				Name:     "device",
				Help:     "manage connected devices",
				LongHelp: "Add, remove, and list connected devices.",
			}
			deviceCmd.AddCmd(&ishell.Cmd{
				Name: "add",
				Help: "add a new device",
				Func: cmd.AddDevice,
			})
			deviceCmd.AddCmd(&ishell.Cmd{
				Name: "rm",
				Help: "remove a device by name",
				Func: cmd.RemoveDevice,
			})
			deviceCmd.AddCmd(&ishell.Cmd{
				Name: "ls",
				Help: "list devices",
				Func: cmd.ListDevices,
			})
			shell.AddCmd(deviceCmd)
		}

		shell.Run()
	}
}

func start() error {
	online, err := core.Node.StartWallet()
	if err != nil {
		return err
	}
	<-online

	// subscribe to thread updates
	for _, thrd := range core.Node.Wallet.Threads() {
		go func(t *thread.Thread) {
			cmd.Subscribe(t)
		}(thrd)
	}

	// start the server
	core.Node.StartServer()

	return nil
}

func stop() error {
	err := core.Node.StopServer()
	if err != nil {
		return err
	}
	return core.Node.StopWallet()
}

func printSplashScreen() {
	cyan := color.New(color.FgCyan).SprintFunc()
	banner := "*** textile node ***"
	fmt.Println(cyan(banner))
	fmt.Println("version: " + core.Version)
	fmt.Printf("repo: %s\n", core.Node.Wallet.GetRepoPath())
	if Options.ServerMode {
		fmt.Println("server mode: enabled")
	}
	if Options.DaemonMode {
		fmt.Println("daemon mode: enabled")
	} else {
		fmt.Println("type 'help' for available commands")
	}
}
