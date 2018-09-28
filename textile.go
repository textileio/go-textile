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
	"github.com/textileio/textile-go/wallet"
	"gopkg.in/abiosoft/ishell.v2"
	icore "gx/ipfs/Qmb8jW1F6ZVyYPW1epc2GFRipmd3S8tJ48pZKBVPzVqj9T/go-ipfs/core"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var wordsRegexp = regexp.MustCompile(`^[a-z]+$`)

type IPFSOptions struct {
	ServerMode  bool   `long:"server" description:"apply IPFS server profile"`
	SwarmPorts  string `long:"swarm-ports" description:"set the swarm ports (tcp,ws)" default:"random"`
}

type LogOptions struct {
	Level   string `short:"l" long:"log-level" description:"set the logging level [debug, info, notice, warning, error, critical]" default:"info"`
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

type WalletCommand struct{
	Init WalletInitCommand `command:"init"`
	Accounts WalletAccountsCommand `command:"accounts"`
}

type WalletInitCommand struct {
	WordCount int `short:"w" long:"word-count" description:"number of mnemonic recovery phrase words: 12,15,18,21,24" default:"12"`
	Password  string `short:"p" long:"password" description:"mnemonic recovery phrase password (omit if none)"`
}

type WalletAccountsCommand struct {
	Password string `short:"p" long:"password" description:"mnemonic recovery phrase password (omit if none)"`
	Depth    int    `short:"d" long:"depth" description:"number of accounts to show" default:"1"`
	Offset   int    `short:"o" long:"offset" description:"account depth to start from" default:"0"`
}

type VersionCommand struct{}

type InitCommand struct {
	AccountSeed string `required:"true" short:"s" long:"seed" description:"account seed (run 'wallet' command to generate new seeds)"`
	IPFS        IPFSOptions `group:"IPFS Options"`
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
var walletCommand WalletCommand
var shellCommand ShellCommand
var daemonCommand DaemonCommand
var options Options
var parser = flags.NewParser(&options, flags.Default)

func init() {
	parser.AddCommand("version",
		"Print version and exit",
		"Print the current version and exit.",
		&versionCommand)
	parser.AddCommand("wallet",
		"Manage a wallet of accounts",
		"Initialize a new wallet, or view accounts from an existing wallet.",
		&walletCommand)
	parser.AddCommand("init",
		"Init the node repo and exit",
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

func (x *WalletInitCommand) Execute(args []string) error {
	// determine word count
	wcount, err := wallet.NewWordCount(x.WordCount)
	if err != nil {
		return err
	}

	// create a new wallet
	w, err := wallet.NewWallet(wcount.EntropySize())
	if err != nil {
		return err
	}

	// show info
	fmt.Println(strings.Repeat("-", len(w.RecoveryPhrase)+4))
	fmt.Println("| " + w.RecoveryPhrase + " |")
	fmt.Println(strings.Repeat("-", len(w.RecoveryPhrase)+4))
	fmt.Println("WARNING! Store these words above in a safe place!")
	fmt.Println("WARNING! If you lose your words, you will lose access to data in all derived accounts!")
	fmt.Println("WARNING! Anyone who has access to these words can access your wallet accounts!")
	fmt.Println("")
	fmt.Println("Use: `wallet accounts` command to inspect more accounts.")
	fmt.Println("")

	// show first account
	kp, err := w.AccountAt(0, x.Password)
	if err != nil {
		return err
	}
	fmt.Println("--- ACCOUNT 0 ---")
	fmt.Println(fmt.Sprintf("PUBLIC KEY: %s", kp.Address()))
	fmt.Println(fmt.Sprintf("SECRET KEY: %s", kp.Seed()))

	return nil
}

func (x *WalletAccountsCommand) Execute(args []string) error {
	if x.Depth < 1 || x.Depth > 100 {
		return errors.New("depth must be greater than 0 and less than 100")
	}
	if x.Offset < 0 || x.Offset > x.Depth {
		return errors.New("offset must be greater than 0 and less than depth")
	}

	// create a shell for reading input
	shell := ishell.New()

	// determine word count
	count := shell.MultiChoice([]string{
		"12",
		"15",
		"18",
		"21",
		"24",
	}, "How many words are in your mnemonic recovery phrase?")
	var wcount wallet.WordCount
	switch count {
	case 0:
		wcount = wallet.TwelveWords
	case 1:
		wcount = wallet.FifteenWords
	case 2:
		wcount = wallet.EighteenWords
	case 3:
		wcount = wallet.TwentyOneWords
	case 4:
		wcount = wallet.TwentyFourWords
	default:
		return wallet.ErrInvalidWordCount
	}

	// read input
	words := make([]string, int(wcount))
	for i := 0; i < int(wcount); i++ {
		shell.Print(fmt.Sprintf("Enter word #%d: ", i+1))
		words[i] = shell.ReadLine()
		if !wordsRegexp.MatchString(words[i]) {
			shell.Println("Invalid word, try again.")
			i--
		}
	}
	wall := wallet.NewWalletFromRecoveryPhrase(strings.Join(words, " "))

	// show info
	for i := x.Offset; i < x.Offset + x.Depth; i++ {
		kp, err := wall.AccountAt(i, x.Password)
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("--- ACCOUNT %d ---", i))
		fmt.Println(fmt.Sprintf("PUBLIC KEY: %s", kp.Address()))
		fmt.Println(fmt.Sprintf("SECRET KEY: %s", kp.Seed()))
	}
	return nil
}

func (x *VersionCommand) Execute(args []string) error {
	fmt.Println(core.Version)
	return nil
}

func (x *InitCommand) Execute(args []string) error {
	// build keypair from provided seed
	kp, err := keypair.Parse(x.AccountSeed)
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

	// determine log level
	level, err := logging.LogLevel(strings.ToUpper(options.Logs.Level))
	if err != nil {
		return errors.New(fmt.Sprintf("determine log level failed: %s", err))
	}

	// build config
	config := core.InitConfig{
		Account:    *accnt,
		RepoPath:   repoPath,
		SwarmPorts: x.IPFS.SwarmPorts,
		IsMobile:   false,
		IsServer:   x.IPFS.ServerMode,
		LogLevel:   level,
		LogFiles:   !options.Logs.NoFiles,
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

	// run the shell
	cmd.RunShell(func() error {
		return startNode(x.Cafe, x.Gateway)
	}, func() error {
		return stopNode(x.Cafe)
	})
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
		CafeAddr: cafeOpts.Addr,
		LogLevel: level,
		LogFiles: !options.Logs.NoFiles,
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
