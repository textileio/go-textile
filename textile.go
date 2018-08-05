package main

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"
	"github.com/textileio/textile-go/cafe"
	"github.com/textileio/textile-go/cafe/dao"
	"github.com/textileio/textile-go/cmd"
	"github.com/textileio/textile-go/core"
	rconfig "github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/wallet"
	"github.com/textileio/textile-go/wallet/thread"
	"gopkg.in/abiosoft/ishell.v2"
	icore "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
)

type Opts struct {
	Version bool `short:"v" long:"version" description:"print the version number and exit"`

	// repo location
	DataDir string `short:"r" long:"data-dir" description:"specify the data directory to be used"`

	// logging options
	LogLevel   string `short:"l" long:"log-level" description:"set the logging level [debug, info, notice, warning, error, critical]" default:"debug"`
	NoLogFiles bool   `short:"n" long:"no-log-files" description:"do not save logs on disk"`

	// modes
	DaemonMode bool `short:"d" long:"daemon" description:"start in a non-interactive daemon mode"`
	ServerMode bool `short:"s" long:"server" description:"start in server mode"`

	// gateway settings
	GatewayBindAddr string `short:"g" long:"gateway-bind-addr" description:"set the gateway address" default:"127.0.0.1:random"`

	// swarm settings
	SwarmPorts string `long:"swarm-ports" description:"set the swarm ports (tcp,ws)" default:"random"`

	// cafe client settings
	CafeAddr string `short:"c" long:"cafe" description:"cafe host address"`

	// cafe host settings
	CafeBindAddr string `long:"cafe-bind-addr" description:"set the cafe address"`

	CafeDBHosts    string `long:"cafe-db-hosts" description:"set the cafe mongo db hosts uri"`
	CafeDBName     string `long:"cafe-db-name" description:"set the cafe mongo db name"`
	CafeDBUser     string `long:"cafe-db-user" description:"set the cafe mongo db user"`
	CafeDBPassword string `long:"cafe-db-password" description:"set the cafe mongo db user password"`
	CafeDBTLS      bool   `long:"cafe-db-tls" description:"use TLS for the cafe mongo db connection"`

	CafeTokenSecret string `long:"cafe-token-secret" description:"set the cafe token secret"`
	CafeReferralKey string `long:"cafe-referral-key" description:"set the cafe referral key"`
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

	// node setup
	config := core.NodeConfig{
		WalletConfig: wallet.Config{
			RepoPath:   dataDir,
			SwarmPorts: Options.SwarmPorts,
			IsMobile:   false,
			IsServer:   Options.ServerMode,
			CafeAddr:   Options.CafeAddr,
		},
		LogLevel: level,
		LogFiles: !Options.NoLogFiles,
	}

	// create a desktop node
	node, _, err := core.NewNode(config)
	if err != nil {
		fmt.Println(fmt.Errorf("create desktop node failed: %s", err))
		return
	}
	core.Node = node

	// check cafe mode
	if Options.CafeBindAddr != "" {
		cafe.Host = &cafe.Cafe{
			Ipfs: func() *icore.IpfsNode {
				return core.Node.Wallet.Ipfs()
			},
			Dao: &dao.DAO{
				Hosts:    Options.CafeDBHosts,
				Name:     Options.CafeDBName,
				User:     Options.CafeDBUser,
				Password: Options.CafeDBPassword,
				TLS:      Options.CafeDBTLS,
			},
			TokenSecret: Options.CafeTokenSecret,
			ReferralKey: Options.CafeReferralKey,
			NodeVersion: core.Version,
		}
	}

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
		fmt.Printf("shutting down...")
		if err := stop(); err != nil && err != wallet.ErrStopped {
			fmt.Println(err.Error())
		} else {
			fmt.Print("done\n")
		}
		os.Exit(1)

	} else {

		// create a new shell
		shell = ishell.New()
		shell.SetHomeHistoryPath(".ishell_history")

		// handle interrupt
		shell.Interrupt(func(c *ishell.Context, count int, input string) {
			if count == 1 {
				shell.Println("input Ctrl-C once more to exit")
				return
			}
			shell.Println("interrupted")
			shell.Printf("shutting down...")
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
				if !core.Node.Wallet.IsOnline() {
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
			cafeCmd := &ishell.Cmd{
				Name:     "cafe",
				Help:     "manage cafe session",
				LongHelp: "Mange your cafe user session.",
			}
			cafeCmd.AddCmd(&ishell.Cmd{
				Name: "add-referral",
				Help: "add cafe referrals",
				Func: cmd.CafeReferral,
			})
			cafeCmd.AddCmd(&ishell.Cmd{
				Name: "referrals",
				Help: "list cafe referrals",
				Func: cmd.ListCafeReferrals,
			})
			cafeCmd.AddCmd(&ishell.Cmd{
				Name: "register",
				Help: "show connected peers (same as `ipfs swarm peers`)",
				Func: cmd.CafeRegister,
			})
			cafeCmd.AddCmd(&ishell.Cmd{
				Name: "login",
				Help: "cafe login",
				Func: cmd.CafeLogin,
			})
			cafeCmd.AddCmd(&ishell.Cmd{
				Name: "status",
				Help: "cafe status",
				Func: cmd.CafeStatus,
			})
			cafeCmd.AddCmd(&ishell.Cmd{
				Name: "logout",
				Help: "cafe logout",
				Func: cmd.CafeLogout,
			})
			shell.AddCmd(cafeCmd)
		}
		{
			profileCmd := &ishell.Cmd{
				Name:     "profile",
				Help:     "manage cafe profiles",
				LongHelp: "Resolve other profiles, get and publish local profile.",
			}
			profileCmd.AddCmd(&ishell.Cmd{
				Name: "publish",
				Help: "publish local profile",
				Func: cmd.PublishProfile,
			})
			profileCmd.AddCmd(&ishell.Cmd{
				Name: "resolve",
				Help: "resolve profiles",
				Func: cmd.ResolveProfile,
			})
			profileCmd.AddCmd(&ishell.Cmd{
				Name: "get",
				Help: "get peer profiles",
				Func: cmd.GetProfile,
			})
			profileCmd.AddCmd(&ishell.Cmd{
				Name: "set-avatar",
				Help: "set local profile avatar",
				Func: cmd.SetAvatarId,
			})
			shell.AddCmd(profileCmd)
		}
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
				Help: "add a new photo",
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
				Help: "list photos from a thread",
				Func: cmd.ListPhotos,
			})
			photoCmd.AddCmd(&ishell.Cmd{
				Name: "ignore",
				Help: "ignore a photo in a thread (requires block id, not photo id)",
				Func: cmd.IgnorePhoto,
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
				Name: "peers",
				Help: "list peers",
				Func: cmd.ListThreadPeers,
			})
			threadCmd.AddCmd(&ishell.Cmd{
				Name: "invite",
				Help: "invite a peer to a thread",
				Func: cmd.AddThreadInvite,
			})
			threadCmd.AddCmd(&ishell.Cmd{
				Name: "accept",
				Help: "accept a thread invite",
				Func: cmd.AcceptThreadInvite,
			})
			threadCmd.AddCmd(&ishell.Cmd{
				Name: "invite-external",
				Help: "create an external invite link",
				Func: cmd.AddExternalThreadInvite,
			})
			threadCmd.AddCmd(&ishell.Cmd{
				Name: "accept-external",
				Help: "accept an external thread invite",
				Func: cmd.AcceptExternalThreadInvite,
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
	if err := core.Node.StartWallet(); err != nil {
		return err
	}
	<-core.Node.Wallet.Online()

	// subscribe to thread updates
	peerId, err := core.Node.Wallet.GetId()
	if err != nil {
		return err
	}
	for _, thrd := range core.Node.Wallet.Threads() {
		go func(t *thread.Thread) {
			cmd.Subscribe(t, peerId)
		}(thrd)
	}

	// subscribe to wallet updates
	go func() {
		for {
			select {
			case update, ok := <-core.Node.Wallet.Updates():
				if !ok {
					return
				}
				switch update.Type {
				case wallet.ThreadAdded:
					if _, thrd := core.Node.Wallet.GetThread(update.Id); thrd != nil {
						go cmd.Subscribe(thrd, peerId)
					}
				case wallet.ThreadRemoved:
					break
				case wallet.DeviceAdded:
					break
				case wallet.DeviceRemoved:
					break
				}
			}
		}
	}()

	// start the gateway
	core.Node.StartGateway(resolveAddress(Options.GatewayBindAddr))

	// start cafe server
	if Options.CafeBindAddr != "" {
		cafe.Host.Start(resolveAddress(Options.CafeBindAddr))
	}

	return nil
}

func stop() error {
	if err := core.Node.StopGateway(); err != nil {
		return err
	}
	if Options.CafeBindAddr != "" {
		if err := cafe.Host.Stop(); err != nil {
			return err
		}
	}
	return core.Node.StopWallet()
}

func printSplashScreen() {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()
	blue := color.New(color.FgHiBlue).SprintFunc()
	grey := color.New(color.FgHiBlack).SprintFunc()
	fmt.Println(cyan("Textile"))
	fmt.Println(grey("version: ") + blue(core.Version))
	fmt.Println(grey("repo: ") + blue(core.Node.Wallet.GetRepoPath()))
	fmt.Println(grey("gateway: ") + yellow(core.Node.GetGatewayAddr()))
	if Options.CafeBindAddr != "" {
		fmt.Println(grey("cafe: ") + yellow(Options.CafeBindAddr))
	}
	if Options.CafeAddr != "" {
		fmt.Println(grey("cafe api: ") + yellow(core.Node.Wallet.GetCafeApiAddr()))
	}
	if Options.ServerMode {
		fmt.Println(grey("server mode: ") + green("enabled"))
	}
	if Options.DaemonMode {
		fmt.Println(grey("daemon mode: ") + green("enabled"))
	} else {
		fmt.Println(grey("type 'help' for available commands"))
	}
}

func resolveAddress(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		log.Fatalf("invalid address: %s", addr)
	}
	port := parts[1]
	if port == "random" {
		port = strconv.Itoa(rconfig.GetRandomPort())
	}
	return fmt.Sprintf("%s:%s", parts[0], port)
}
