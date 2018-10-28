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
	"github.com/textileio/textile-go/gateway"
	"github.com/textileio/textile-go/keypair"
	rconfig "github.com/textileio/textile-go/repo/config"
	"github.com/textileio/textile-go/wallet"
	"gopkg.in/abiosoft/ishell.v2"
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
	ServerMode bool   `long:"server" description:"apply IPFS server profile"`
	SwarmPorts string `long:"swarm-ports" description:"set the swarm ports (tcp,ws)" default:"random"`
}

type LogOptions struct {
	Level   string `short:"l" long:"log-level" description:"set the logging level [debug, info, notice, warning, error, critical]" default:"info"`
	NoFiles bool   `short:"n" long:"no-log-files" description:"do not save logs on disk"`
}

type ApiOptions struct {
	BindAddr string `short:"a" long:"api-bind-addr" description:"set the rest api address" default:"127.0.0.1:random"`
}

type GatewayOptions struct {
	BindAddr string `short:"g" long:"gateway-bind-addr" description:"set the gateway address" default:"127.0.0.1:random"`
}

type CafeApiOptions struct {
	Open     bool   `short:"c" long:"open-cafe" description:"opens the cafe service for other peers"`
	BindAddr string `long:"cafe-bind-addr" description:"set the cafe rest api address" default:"127.0.0.1:random"`
}

type Options struct{}

type WalletCommand struct {
	Init     WalletInitCommand     `command:"init"`
	Accounts WalletAccountsCommand `command:"accounts"`
}

type WalletInitCommand struct {
	WordCount int    `short:"w" long:"word-count" description:"number of mnemonic recovery phrase words: 12,15,18,21,24" default:"12"`
	Password  string `short:"p" long:"password" description:"mnemonic recovery phrase password (omit if none)"`
}

type WalletAccountsCommand struct {
	Password string `short:"p" long:"password" description:"mnemonic recovery phrase password (omit if none)"`
	Depth    int    `short:"d" long:"depth" description:"number of accounts to show" default:"1"`
	Offset   int    `short:"o" long:"offset" description:"account depth to start from" default:"0"`
}

type VersionCommand struct{}

type InitCommand struct {
	AccountSeed string      `required:"true" short:"s" long:"seed" description:"account seed (run 'wallet' command to generate new seeds)"`
	RepoPath    string      `short:"r" long:"repo-dir" description:"specify a custom repository path"`
	Logs        LogOptions  `group:"Log Options"`
	IPFS        IPFSOptions `group:"IPFS Options"`
}

type MigrateCommand struct {
	RepoPath string `short:"r" long:"repo-dir" description:"specify a custom repository path"`
}

type DaemonCommand struct {
	RepoPath string         `short:"r" long:"repo-dir" description:"specify a custom repository path"`
	Logs     LogOptions     `group:"Log Options"`
	Api      ApiOptions     `group:"API Options"`
	Gateway  GatewayOptions `group:"Gateway Options"`
	CafeApi  CafeApiOptions `group:"Cafe API Options"`
}

type ShellCommand struct {
	RepoPath string         `short:"r" long:"repo-dir" description:"specify a custom repository path"`
	Logs     LogOptions     `group:"Log Options"`
	Api      ApiOptions     `group:"API Options"`
	Gateway  GatewayOptions `group:"Gateway Options"`
	CafeApi  CafeApiOptions `group:"Cafe API Options"`
}

// main commands
var initCommand InitCommand
var migrateCommand MigrateCommand
var versionCommand VersionCommand
var walletCommand WalletCommand
var shellCommand ShellCommand
var daemonCommand DaemonCommand

var options Options
var parser = flags.NewParser(&options, flags.Default)

func init() {
	// add main commands
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
	parser.AddCommand("migrate",
		"Migrate the node repo and exit",
		"Migrate the node repository and exit.",
		&migrateCommand)
	parser.AddCommand("shell",
		"Start a node shell",
		"Start an interactive node shell session.",
		&shellCommand)
	parser.AddCommand("daemon",
		"Start a node daemon",
		"Start a node daemon session.",
		&daemonCommand)

	// add cmd commands
	for _, c := range cmd.Cmds() {
		parser.AddCommand(c.Name(), c.Short(), c.Long(), c)
	}
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
	fmt.Println(fmt.Sprintf("PUBLIC ADDR: %s", kp.Address()))
	fmt.Println(fmt.Sprintf("SECRET SEED: %s", kp.Seed()))

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
	count := shell.MultiChoice([]string{"12", "15", "18", "21", "24"}, "How many words are in your mnemonic recovery phrase?")
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
	for i := x.Offset; i < x.Offset+x.Depth; i++ {
		kp, err := wall.AccountAt(i, x.Password)
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("--- ACCOUNT %d ---", i))
		fmt.Println(fmt.Sprintf("PUBLIC ADDR: %s", kp.Address()))
		fmt.Println(fmt.Sprintf("SECRET SEED: %s", kp.Seed()))
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
	repoPath, err := getRepoPath(x.RepoPath)
	if err != nil {
		return err
	}

	// determine log level
	level, err := logging.LogLevel(strings.ToUpper(x.Logs.Level))
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
		LogFiles:   !x.Logs.NoFiles,
	}

	// initialize a node
	if err := core.InitRepo(config); err != nil {
		return errors.New(fmt.Sprintf("initialize failed: %s", err))
	}
	fmt.Printf("ok, address: %s\n", accnt.Address())
	return nil
}

func (x *MigrateCommand) Execute(args []string) error {
	// handle repo path
	repoPath, err := getRepoPath(x.RepoPath)
	if err != nil {
		return err
	}

	// build config
	config := core.MigrateConfig{
		RepoPath: repoPath,
	}

	// run migrate
	if err := core.MigrateRepo(config); err != nil {
		return errors.New(fmt.Sprintf("migrate repo: %s", err))
	}
	fmt.Println("repo successfully migrated")
	return nil
}

func (x *DaemonCommand) Execute(args []string) error {
	if err := buildNode(x.RepoPath, x.Api, x.Gateway, x.CafeApi, x.Logs); err != nil {
		return err
	}
	printSplashScreen(true)

	// handle interrupt
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("interrupted")
	fmt.Printf("shutting down...")
	if err := stopNode(); err != nil && err != core.ErrStopped {
		fmt.Println(err.Error())
	} else {
		fmt.Print("done\n")
	}
	os.Exit(1)
	return nil
}

func (x *ShellCommand) Execute(args []string) error {
	if err := buildNode(x.RepoPath, x.Api, x.Gateway, x.CafeApi, x.Logs); err != nil {
		return err
	}
	printSplashScreen(false)

	// run the shell
	cmd.RunShell(func() error {
		return startNode(x.Api, x.Gateway)
	}, func() error {
		return stopNode()
	})
	return nil
}

func getRepoPath(repoPath string) (string, error) {
	// handle repo path
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

func buildNode(repoPath string, apiOpts ApiOptions, gatewayOpts GatewayOptions, cafeOpts CafeApiOptions, logOpts LogOptions) error {
	// handle repo path
	repoPathf, err := getRepoPath(repoPath)
	if err != nil {
		return err
	}

	// determine log level
	level, err := logging.LogLevel(strings.ToUpper(logOpts.Level))
	if err != nil {
		return errors.New(fmt.Sprintf("determine log level failed: %s", err))
	}

	// node setup
	config := core.RunConfig{
		RepoPath:     repoPathf,
		LogLevel:     level,
		LogFiles:     !logOpts.NoFiles,
		CafeOpen:     cafeOpts.Open,
		CafeBindAddr: resolveAddress(cafeOpts.BindAddr),
	}

	// create a node
	node, err := core.NewTextile(config)
	if err != nil {
		return errors.New(fmt.Sprintf("create node failed: %s", err))
	}
	core.Node = node

	// create the gateway
	gateway.Host = &gateway.Gateway{}

	// auto start it
	if err := startNode(apiOpts, gatewayOpts); err != nil {
		fmt.Println(fmt.Errorf("start node failed: %s", err))
	}
	return nil
}

func startNode(apiOpts ApiOptions, gatewayOpts GatewayOptions) error {
	if err := core.Node.Start(); err != nil {
		return err
	}

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
				case core.AccountPeerAdded:
					break
				case core.AccountPeerRemoved:
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

	// start api server
	core.Node.StartApi(resolveAddress(apiOpts.BindAddr))

	// start the gateway
	gateway.Host.Start(resolveAddress(gatewayOpts.BindAddr))

	// wait for the ipfs node to go online
	<-core.Node.OnlineCh()

	return nil
}

func stopNode() error {
	if err := core.Node.StopApi(); err != nil {
		return err
	}
	if err := gateway.Host.Stop(); err != nil {
		return err
	}
	return core.Node.Stop()
}

func printSplashScreen(daemon bool) {
	cyan := color.New(color.FgHiCyan).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()
	grey := color.New(color.FgHiBlack).SprintFunc()
	pid, err := core.Node.PeerId()
	if err != nil {
		log.Fatalf("get peer id failed: %s", err)
	}
	accnt, err := core.Node.Account()
	if err != nil {
		log.Fatalf("get account failed: %s", err)
	}
	if daemon {
		fmt.Println(grey("Textile daemon version v" + core.Version))
	} else {
		fmt.Println(grey("Textile shell version v" + core.Version))
	}
	fmt.Println(grey("repo:    ") + grey(core.Node.GetRepoPath()))
	fmt.Println(grey("api:     ") + grey(core.Node.ApiAddr()))
	fmt.Println(grey("gateway: ") + grey(gateway.Host.Addr()))
	if core.Node.CafeApiAddr() != "" {
		fmt.Println(grey("cafe:    ") + grey(core.Node.CafeApiAddr()))
	}
	fmt.Println(grey("peer:    ") + green(pid.Pretty()))
	fmt.Println(grey("account: ") + cyan(accnt.Address()))
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
