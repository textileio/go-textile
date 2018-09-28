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
	"github.com/textileio/textile-go/gateway"
	"github.com/textileio/textile-go/keypair"
	rconfig "github.com/textileio/textile-go/repo/config"
	"gopkg.in/abiosoft/ishell.v2"
	icore "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
)

type InitOptions struct {
	AccountSeed string `required:"true" short:"s" long:"seed" description:"account seed (run 'wallet' command to generate new seeds)"`
	ServerMode  bool   `long:"server" description:"start in server mode"`
	SwarmPorts  string `long:"swarm-ports" description:"set the swarm ports (tcp,ws)" default:"random"`
}

type LogOptions struct {
	Level   string `short:"l" long:"log-level" description:"set the logging level [debug, info, notice, warning, error, critical]" default:"debug"`
	NoFiles bool   `short:"n" long:"no-log-files" description:"do not save logs on disk"`
}

type GatewayOptions struct {
	BindAddr string `short:"g" long:"gateway-bind-addr" description:"set the gateway address" default:"127.0.0.1:random"`
}

type CafeOptions struct {
	// client settings
	Addr string `short:"c" long:"cafe" description:"cafe host address"`

	// host settings
	BindAddr    string `long:"cafe-bind-addr" description:"set the cafe address"`
	DBHosts     string `long:"cafe-db-hosts" description:"set the cafe mongo db hosts uri"`
	DBName      string `long:"cafe-db-name" description:"set the cafe mongo db name"`
	DBUser      string `long:"cafe-db-user" description:"set the cafe mongo db user"`
	DBPassword  string `long:"cafe-db-password" description:"set the cafe mongo db user password"`
	DBTLS       bool   `long:"cafe-db-tls" description:"use TLS for the cafe mongo db connection"`
	TokenSecret string `long:"cafe-token-secret" description:"set the cafe token secret"`
	ReferralKey string `long:"cafe-referral-key" description:"set the cafe referral key"`
}

type Options struct {
	RepoPath string     `short:"r" long:"repo-dir" description:"specify a custom repository path"`
	Logs     LogOptions `group:"Log Options"`
}

type VersionCommand struct{}

type InitCommand struct {
	Init InitOptions `group:"Init Options"`
}

type DaemonCommand struct {
	Gateway GatewayOptions `group:"Gateway Options"`
	Cafe    CafeOptions    `group:"Cafe Options"`
}

type ShellCommand struct {
	Gateway GatewayOptions `group:"Gateway Options"`
	Cafe    CafeOptions    `group:"Cafe Options"`
}

var initCommand InitCommand
var versionCommand VersionCommand
var shellCommand ShellCommand
var daemonCommand DaemonCommand
var options Options
var parser = flags.NewParser(&options, flags.Default)
var shell *ishell.Shell

func init() {
	parser.AddCommand("version",
		"Print version and exit",
		"Print the current version and exit.",
		&versionCommand)
	parser.AddCommand("init",
		"Init the node repo",
		"Initialize the node repository and exit.",
		&initCommand)
	parser.AddCommand("shell",
		"Start a node shell",
		"Start an interactive node shell session.",
		&shellCommand)
	parser.AddCommand("daemon",
		"Start a node daemon",
		"Start a node daemon session (useful w/ a cafe).",
		&daemonCommand)
}

func main() {
	parser.Parse()
}

func (x *VersionCommand) Execute(args []string) error {
	fmt.Println(core.Version)
	return nil
}

func (x *InitCommand) Execute(args []string) error {
	// build keypair from provided seed
	kp, err := keypair.Parse(x.Init.AccountSeed)
	if err != nil {
		return errors.New(fmt.Sprintf("parse account seed failed: %s", err))
	}
	accnt, ok := kp.(*keypair.Full)
	if !ok {
		return keypair.ErrInvalidKey
	}

	// handle repo path
	repoPath, err := getRepoPath()
	if err != nil {
		return err
	}

	// build config
	config := core.InitConfig{
		Account:    *accnt,
		RepoPath:   repoPath,
		SwarmPorts: x.Init.SwarmPorts,
		IsMobile:   false,
		IsServer:   x.Init.ServerMode,
	}

	// initialize a node
	if err := core.InitRepo(config); err != nil {
		return errors.New(fmt.Sprintf("initialize node failed: %s", err))
	}

	fmt.Printf("Textile node initialized with public address: %s\n", accnt.Address())

	return nil
}

func (x *DaemonCommand) Execute(args []string) error {
	if err := buildNode(x.Cafe, x.Gateway); err != nil {
		return err
	}
	printSplashScreen(x.Cafe, true)

	// handle interrupt
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("interrupted")
	fmt.Printf("shutting down...")
	if err := stopNode(x.Cafe); err != nil && err != core.ErrStopped {
		fmt.Println(err.Error())
	} else {
		fmt.Print("done\n")
	}
	os.Exit(1)
	return nil
}

func (x *ShellCommand) Execute(args []string) error {
	if err := buildNode(x.Cafe, x.Gateway); err != nil {
		return err
	}
	printSplashScreen(x.Cafe, false)

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
		if err := stopNode(x.Cafe); err != nil && err != core.ErrStopped {
			c.Err(err)
		} else {
			shell.Printf("done\n")
		}
		os.Exit(1)
	})

	// add interactive commands
	shell.AddCmd(&ishell.Cmd{
		Name: "start",
		Help: "start the node",
		Func: func(c *ishell.Context) {
			if core.Node.Started() {
				c.Println("already started")
				return
			}
			if err := startNode(x.Cafe, x.Gateway); err != nil {
				c.Println(fmt.Errorf("start node failed: %s", err))
				return
			}
			c.Println("ok, started")
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "stop",
		Help: "stop the node",
		Func: func(c *ishell.Context) {
			if !core.Node.Started() {
				c.Println("already stopped")
				return
			}
			if err := stopNode(x.Cafe); err != nil {
				c.Println(fmt.Errorf("stop node failed: %s", err))
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
		Help: "ping another peer",
		Func: func(c *ishell.Context) {
			if !core.Node.IsOnline() {
				c.Println("not online yet")
				return
			}
			if len(c.Args) == 0 {
				c.Err(errors.New("missing peer id"))
				return
			}
			status, err := core.Node.GetPeerStatus(c.Args[0])
			if err != nil {
				c.Println(fmt.Errorf("ping failed: %s", err))
				return
			}
			c.Println(status)
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "fetch-messages",
		Help: "fetch offline messages from the DHT",
		Func: func(c *ishell.Context) {
			if !core.Node.IsOnline() {
				c.Println("not online yet")
				return
			}
			if err := core.Node.FetchMessages(); err != nil {
				c.Println(fmt.Errorf("fetch messages failed: %s", err))
				return
			}
			c.Println("ok, fetching")
		},
	})
	{
		walletCmd := &ishell.Cmd{
			Name:     "wallet",
			Help:     "manage wallets",
			LongHelp: "Create and manage your textile wallet.",
		}
		walletCmd.AddCmd(&ishell.Cmd{
			Name: "create",
			Help: "create a new textile wallet",
			Func: cmd.CreateWallet,
		})
		walletCmd.AddCmd(&ishell.Cmd{
			Name: "accounts",
			Help: "view wallet account keys",
			Func: cmd.WalletAccounts,
		})
		shell.AddCmd(walletCmd)
	}
	{
		cafeCmd := &ishell.Cmd{
			Name:     "cafe",
			Help:     "manage cafe session",
			LongHelp: "Manage your cafe user session.",
		}
		cafeCmd.AddCmd(&ishell.Cmd{
			Name: "add-referral",
			Help: "add cafe referrals",
			Func: cmd.CafeAddReferral,
		})
		cafeCmd.AddCmd(&ishell.Cmd{
			Name: "referrals",
			Help: "list cafe referrals",
			Func: cmd.ListCafeReferrals,
		})
		cafeCmd.AddCmd(&ishell.Cmd{
			Name: "register",
			Help: "cafe register",
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
			Name: "tokens",
			Help: "cafe tokens",
			Func: cmd.CafeTokens,
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
			Help: "get photo metadata",
			Func: cmd.GetPhotoMetadata,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "ls",
			Help: "list photos from a thread",
			Func: cmd.ListPhotos,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "comment",
			Help: "comment on a photo (terminate input w/ ';'",
			Func: cmd.AddPhotoComment,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "like",
			Help: "like a photo",
			Func: cmd.AddPhotoLike,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "comments",
			Help: "list photo comments",
			Func: cmd.ListPhotoComments,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "likes",
			Help: "list photo likes",
			Func: cmd.ListPhotoLikes,
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
			Name: "blocks",
			Help: "list blocks",
			Func: cmd.ListThreadBlocks,
		})
		threadCmd.AddCmd(&ishell.Cmd{
			Name: "head",
			Help: "show current HEAD",
			Func: cmd.GetThreadHead,
		})
		threadCmd.AddCmd(&ishell.Cmd{
			Name: "ignore",
			Help: "ignore a block",
			Func: cmd.IgnoreBlock,
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
	{
		notificationCmd := &ishell.Cmd{
			Name:     "notification",
			Help:     "manage notifications",
			LongHelp: "List and read notifications.",
		}
		notificationCmd.AddCmd(&ishell.Cmd{
			Name: "read",
			Help: "mark a notification as read",
			Func: cmd.ReadNotification,
		})
		notificationCmd.AddCmd(&ishell.Cmd{
			Name: "readall",
			Help: "mark all notifications as read",
			Func: cmd.ReadAllNotifications,
		})
		notificationCmd.AddCmd(&ishell.Cmd{
			Name: "ls",
			Help: "list notifications",
			Func: cmd.ListNotifications,
		})
		notificationCmd.AddCmd(&ishell.Cmd{
			Name: "accept",
			Help: "accept an invite via notification",
			Func: cmd.AcceptThreadInviteViaNotification,
		})
		shell.AddCmd(notificationCmd)
	}

	shell.Run()
	return nil
}

func getRepoPath() (string, error) {
	// handle repo path
	repoPath := options.RepoPath
	if len(repoPath) == 0 {
		// get homedir
		home, err := homedir.Dir()
		if err != nil {
			return "", errors.New(fmt.Sprintf("get homedir failed: %s", err))
		}

		// ensure app folder is created
		appDir := filepath.Join(home, ".textile")
		if err := os.MkdirAll(appDir, 0755); err != nil {
			return "", errors.New(fmt.Sprintf("create repo directory failed: %s", err))
		}
		repoPath = filepath.Join(appDir, "repo")
	}
	return repoPath, nil
}

func buildNode(cafeOpts CafeOptions, gatewayOpts GatewayOptions) error {
	// handle repo path
	repoPathf, err := getRepoPath()
	if err != nil {
		return err
	}

	// determine log level
	level, err := logging.LogLevel(strings.ToUpper(options.Logs.Level))
	if err != nil {
		return errors.New(fmt.Sprintf("determine log level failed: %s", err))
	}

	// node setup
	config := core.RunConfig{
		RepoPath: repoPathf,
		LogLevel: level,
		LogFiles: !options.Logs.NoFiles,
		CafeAddr: cafeOpts.Addr,
	}

	// create a node
	node, err := core.NewTextile(config)
	if err != nil {
		return errors.New(fmt.Sprintf("create node failed: %s", err))
	}
	core.Node = node

	// create the gateway
	gateway.Host = &gateway.Gateway{}

	// check cafe mode
	if cafeOpts.BindAddr != "" {
		cafe.Host = &cafe.Cafe{
			Ipfs: func() *icore.IpfsNode {
				return core.Node.Ipfs()
			},
			Dao: &dao.DAO{
				Hosts:    cafeOpts.DBHosts,
				Name:     cafeOpts.DBName,
				User:     cafeOpts.DBUser,
				Password: cafeOpts.DBPassword,
				TLS:      cafeOpts.DBTLS,
			},
			TokenSecret: cafeOpts.TokenSecret,
			ReferralKey: cafeOpts.ReferralKey,
			NodeVersion: core.Version,
		}
	}

	// auto start it
	if err := startNode(cafeOpts, gatewayOpts); err != nil {
		fmt.Println(fmt.Errorf("start node failed: %s", err))
	}

	return nil
}

func startNode(cafeOpts CafeOptions, gatewayOpts GatewayOptions) error {
	if err := core.Node.Start(); err != nil {
		return err
	}
	<-core.Node.Online()

	// subscribe to wallet updates
	go func() {
		for {
			select {
			case update, ok := <-core.Node.Updates():
				if !ok {
					return
				}
				switch update.Type {
				case core.ThreadAdded:
					break
				case core.ThreadRemoved:
					break
				case core.DeviceAdded:
					break
				case core.DeviceRemoved:
					break
				}
			}
		}
	}()

	// subscribe to thread updates
	go func() {
		green := color.New(color.FgHiGreen).SprintFunc()
		for {
			select {
			case update, ok := <-core.Node.ThreadUpdates():
				if !ok {
					return
				}
				msg := fmt.Sprintf("new %s block in thread '%s'", update.Block.Type.Description(), update.ThreadName)
				fmt.Println(green(msg))
			}
		}
	}()

	// subscribe to notifications
	go func() {
		yellow := color.New(color.FgHiYellow).SprintFunc()
		for {
			select {
			case notification, ok := <-core.Node.Notifications():
				if !ok {
					return
				}
				var username string
				if notification.ActorUsername != "" {
					username = notification.ActorUsername
				} else {
					username = notification.ActorId
				}
				note := fmt.Sprintf("#%s: %s %s.", notification.Subject, username, notification.Body)
				fmt.Println(yellow(note))
			}
		}
	}()

	// start the gateway
	gateway.Host.Start(resolveAddress(gatewayOpts.BindAddr))

	// start cafe server
	if cafeOpts.BindAddr != "" {
		cafe.Host.Start(resolveAddress(cafeOpts.BindAddr))
	}

	return nil
}

func stopNode(cafeOpts CafeOptions) error {
	if err := gateway.Host.Stop(); err != nil {
		return err
	}
	if cafeOpts.BindAddr != "" {
		if err := cafe.Host.Stop(); err != nil {
			return err
		}
	}
	return core.Node.Stop()
}

func printSplashScreen(cafeOpts CafeOptions, daemon bool) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()
	blue := color.New(color.FgHiBlue).SprintFunc()
	grey := color.New(color.FgHiBlack).SprintFunc()
	addr, err := core.Node.Address()
	if err != nil {
		log.Fatalf("get address failed: %s", err)
	}
	if daemon {
		fmt.Println(cyan("Textile Daemon"))
	} else {
		fmt.Println(cyan("Textile Shell"))
	}
	fmt.Println(grey("address: ") + green(addr))
	fmt.Println(grey("version: ") + blue(core.Version))
	fmt.Println(grey("repo: ") + blue(core.Node.GetRepoPath()))
	fmt.Println(grey("gateway: ") + yellow(gateway.Host.Addr()))
	if cafeOpts.BindAddr != "" {
		fmt.Println(grey("cafe: ") + yellow(cafeOpts.BindAddr))
	}
	if cafeOpts.Addr != "" {
		fmt.Println(grey("cafe api: ") + yellow(core.Node.GetCafeApiAddr()))
	}
	if !daemon {
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
