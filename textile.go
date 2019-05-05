package main

import (
	"bytes"
	_ "expvar"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	logging "github.com/ipfs/go-log"
	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/go-homedir"
	"github.com/textileio/go-textile/cmd"
	"github.com/textileio/go-textile/common"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/gateway"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
	"github.com/textileio/go-textile/wallet"
)

type ipfsOptions struct {
	ServerMode bool   `long:"server" description:"Apply IPFS server profile."`
	SwarmPorts string `long:"swarm-ports" description:"Set the swarm ports (TCP,WS). A random TCP port is chosen by default."`
}

type logOptions struct {
	NoFiles bool `short:"n" long:"no-log-files" description:"Write logs to stdout instead of rolling files."`
	Debug   bool `short:"d" long:"debug" description:"Set the logging level to debug."`
}

type addressOptions struct {
	ApiBindAddr       string `short:"a" long:"api-bind-addr" description:"Set the local API address." default:"127.0.0.1:40600"`
	CafeApiBindAddr   string `short:"c" long:"cafe-bind-addr" description:"Set the cafe REST API address." default:"0.0.0.0:40601"`
	GatewayBindAddr   string `short:"g" long:"gateway-bind-addr" description:"Set the IPFS gateway address." default:"127.0.0.1:5050"`
	ProfilingBindAddr string `long:"profile-bind-addr" description:"Set the profiling address." default:"127.0.0.1:6060"`
}

type cafeOptions struct {
	Open        bool   `long:"cafe-open" description:"Open the p2p Cafe Service for other peers."`
	URL         string `long:"cafe-url" description:"Specify a custom URL of this cafe, e.g., https://mycafe.com"`
	NeighborURL string `long:"cafe-neighbor-url" description:"Specify the URL of a secondary cafe. Must return cafe info, e.g., via a Gateway: https://my-gateway.yolo.com/cafe, or a Cafe API: https://my-cafe.yolo.com"`
}

type options struct{}

type walletCmd struct {
	Init     walletInitCmd     `command:"init" description:"Initialize a new wallet"`
	Accounts walletAccountsCmd `command:"accounts" description:"Show derived accounts"`
}

type walletInitCmd struct {
	WordCount int    `short:"w" long:"word-count" description:"Number of mnemonic recovery phrase words: 12,15,18,21,24." default:"12"`
	Password  string `short:"p" long:"password" description:"Mnemonic recovery phrase password (omit if none)."`
}

func (x *walletInitCmd) Usage() string {
	return `

Initializes a new account wallet backed by a mnemonic recovery phrase.`
}

type walletAccountsCmd struct {
	Password string `short:"p" long:"password" description:"Mnemonic recovery phrase password (omit if none)."`
	Depth    int    `short:"d" long:"depth" description:"Number of accounts to show." default:"1"`
	Offset   int    `short:"o" long:"offset" description:"Account depth to start from." default:"0"`
}

func (x *walletAccountsCmd) Usage() string {
	return `

Shows the derived accounts (address/seed pairs) in a wallet.`
}

type versionCmd struct {
	Git bool `short:"g" long:"git" description:"Show full git version summary."`
}

type initCmd struct {
	AccountSeed string         `required:"true" short:"s" long:"seed" description:"Account seed (run 'wallet' command to generate new seeds)."`
	PinCode     string         `short:"p" long:"pin-code" description:"Specify a pin code for datastore encryption."`
	RepoPath    string         `short:"r" long:"repo-dir" description:"Specify a custom repository path."`
	Addresses   addressOptions `group:"Address Options"`
	CafeOptions cafeOptions    `group:"Cafe Options"`
	IPFS        ipfsOptions    `group:"IPFS Options"`
	Logs        logOptions     `group:"Log Options"`
}

type migrateCmd struct {
	RepoPath string `short:"r" long:"repo-dir" description:"Specify a custom repository path."`
	PinCode  string `short:"p" long:"pin-code" description:"Specify the pin code for datastore encryption (omit of none was used during init)."`
}

type daemonCmd struct {
	PinCode  string `short:"p" long:"pin-code" description:"Specify the pin code for datastore encryption (omit of none was used during init)."`
	RepoPath string `short:"r" long:"repo-dir" description:"Specify a custom repository path."`
	Debug    bool   `short:"d" long:"debug" description:"Set the logging level to debug."`
	Docs     bool   `short:"s" long:"serve-docs" description:"Whether to serve the local REST API docs."`
}

type commandsCmd struct {
}

type docsCmd struct {
}

var node *core.Textile
var log = logging.Logger("tex-main")
var parser = flags.NewParser(&options{}, flags.Default)

func init() {
	// Add the commands that are located in this file
	_, _ = parser.AddCommand("version",
		"Print version and exit",
		"Print the current version and exit.",
		&versionCmd{})
	_, _ = parser.AddCommand("wallet",
		"Manage or create an account wallet",
		"Initialize a new wallet, or view accounts from an existing wallet.",
		&walletCmd{})
	_, _ = parser.AddCommand("init",
		"Init the node repo and exit",
		"Initialize the node repository and exit.",
		&initCmd{})
	_, _ = parser.AddCommand("migrate",
		"Migrate the node repo and exit",
		"Migrate the node repository and exit.",
		&migrateCmd{})
	_, _ = parser.AddCommand("daemon",
		"Start the daemon",
		"Start a node daemon session.",
		&daemonCmd{})
	_, _ = parser.AddCommand("commands",
		"List available commands",
		"List all available textile commands.",
		&commandsCmd{})
	_, _ = parser.AddCommand("docs",
		"Print docs",
		"Prints markdown docs for the command-line client.",
		&docsCmd{})

	// Add the commands that are located in the cmd directory/package
	for _, c := range cmd.Cmds() {
		_, _ = parser.AddCommand(c.Name(), c.Short(), c.Long(), c)
	}
}

func main() {
	_, _ = parser.Parse()
}

// Output the available commands to the user
func (x *commandsCmd) Execute(args []string) error {
	for _, c := range parser.Commands() {
		if len(c.Commands()) == 0 {
			fmt.Println(fmt.Sprintf("textile %s", c.Name))
		}
		for _, sub := range c.Commands() {
			fmt.Println(fmt.Sprintf("textile %s %s", c.Name, sub.Name))
		}
	}
	return nil
}

func (x *docsCmd) Execute(args []string) error {
	doc, err := getCommandDoc("textile")
	if err != nil {
		return err
	}
	fmt.Println(doc)

	var list []string
	set := make(map[string]*flags.Command)
	for _, c := range parser.Commands() {
		list = append(list, c.Name)
		set[c.Name] = c
	}

	sort.Strings(list)
	for _, n := range list {
		c := set[n]
		doc, err := getCommandDoc(c.Name, c.Name)
		if err != nil {
			return err
		}
		fmt.Println(doc)
		for _, sub := range c.Commands() {
			doc, err := getCommandDoc(sub.Name, c.Name, sub.Name)
			if err != nil {
				return err
			}
			fmt.Println(doc)
		}
	}
	return nil
}

func getCommandDoc(name string, args ...string) (string, error) {
	args = append(args, "--help")
	e := exec.Command("textile", args...)
	buf := new(bytes.Buffer)
	e.Stdout = buf
	if err := e.Run(); err != nil {
		return "", err
	}

	pre := strings.Repeat("#", len(args)+1)

	doc := pre + " `" + name + "`\n\n"
	doc += "```\n"
	doc += buf.String()
	doc += "```\n"

	return doc, nil
}

func (x *walletInitCmd) Execute(args []string) error {
	wcount, err := wallet.NewWordCount(x.WordCount)
	if err != nil {
		return err
	}

	w, err := wallet.NewWallet(wcount.EntropySize())
	if err != nil {
		return err
	}

	// Print the recovery phrase surrounded by a box of dashes
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
	fmt.Println(kp.Address())
	fmt.Println(kp.Seed())

	return nil
}

func (x *walletAccountsCmd) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing recovery phrase")
	}

	if x.Depth < 1 || x.Depth > 100 {
		return fmt.Errorf("depth must be greater than 0 and less than 100")
	}
	if x.Offset < 0 || x.Offset > x.Depth {
		return fmt.Errorf("offset must be greater than 0 and less than depth")
	}

	wall := wallet.NewWalletFromRecoveryPhrase(args[0])

	for i := x.Offset; i < x.Offset+x.Depth; i++ {
		kp, err := wall.AccountAt(i, x.Password)
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("--- ACCOUNT %d ---", i))
		fmt.Println(kp.Address())
		fmt.Println(kp.Seed())
	}
	return nil
}

func (x *versionCmd) Execute(args []string) error {
	if x.Git {
		fmt.Println("go-textile version " + common.GitSummary)
	} else {
		fmt.Println("go-textile version v" + common.Version)
	}
	return nil
}

// Initialise the textile user data store at the repo path
func (x *initCmd) Execute(args []string) error {
	kp, err := keypair.Parse(x.AccountSeed)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("parse account seed failed: %s", err))
	}
	accnt, ok := kp.(*keypair.Full)
	if !ok {
		return keypair.ErrInvalidKey
	}

	repoPath, err := getRepoPath(x.RepoPath)
	if err != nil {
		return err
	}

	config := core.InitConfig{
		Account:         accnt,
		PinCode:         x.PinCode,
		RepoPath:        repoPath,
		SwarmPorts:      x.IPFS.SwarmPorts,
		ApiAddr:         x.Addresses.ApiBindAddr,
		CafeApiAddr:     x.Addresses.CafeApiBindAddr,
		GatewayAddr:     x.Addresses.GatewayBindAddr,
		ProfilingAddr:   x.Addresses.ProfilingBindAddr,
		IsMobile:        false,
		IsServer:        x.IPFS.ServerMode,
		LogToDisk:       !x.Logs.NoFiles,
		Debug:           x.Logs.Debug,
		CafeOpen:        x.CafeOptions.Open,
		CafeURL:         envOrFlag("CAFE_HOST_URL", x.CafeOptions.URL),
		CafeNeighborURL: envOrFlag("CAFE_HOST_NEIGHBOR_URL", x.CafeOptions.NeighborURL),
		// ^ @todo why do we prefer env over flag for these? shouldn't flag override the env?
	}

	if err := core.InitRepo(config); err != nil {
		return fmt.Errorf(fmt.Sprintf("initialize failed: %s", err))
	}
	fmt.Printf("Initialized account with address %s\n", accnt.Address())
	return nil
}

// Grab the repo path and migrate it to the latest version, passing the decryption pincode
func (x *migrateCmd) Execute(args []string) error {
	repoPath, err := getRepoPath(x.RepoPath)
	if err != nil {
		return err
	}

	if err := core.MigrateRepo(core.MigrateConfig{
		PinCode:  x.PinCode,
		RepoPath: repoPath,
	}); err != nil {
		return fmt.Errorf(fmt.Sprintf("migrate repo: %s", err))
	}
	fmt.Println("Repo was successfully migrated")
	return nil
}

// Start the daemon against the user repository
func (x *daemonCmd) Execute(args []string) error {
	repoPathf, err := getRepoPath(x.RepoPath)
	if err != nil {
		return err
	}

	node, err = core.NewTextile(core.RunConfig{
		PinCode:  x.PinCode,
		RepoPath: repoPathf,
		Debug:    x.Debug,
	})
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("create node failed: %s", err))
	}

	gateway.Host = &gateway.Gateway{
		Node: node,
	}

	if err := startNode(x.Docs); err != nil {
		return fmt.Errorf(fmt.Sprintf("start node failed: %s", err))
	}
	printSplash()

	// Shutdown gracefully if an SIGINT was received
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Interrupted")
	fmt.Printf("Shutting down...")
	if err := stopNode(); err != nil && err != core.ErrStopped {
		fmt.Println(err.Error())
	} else {
		fmt.Print("done\n")
	}
	os.Exit(1)
	return nil
}

// Get the repo path for the user, will create it if missing
// Unless provided, it defaults to ~/.textile/repo
func getRepoPath(repoPath string) (string, error) {
	if len(repoPath) == 0 {
		// get homedir
		home, err := homedir.Dir()
		if err != nil {
			return "", fmt.Errorf(fmt.Sprintf("get homedir failed: %s", err))
		}

		// ensure app folder is created
		appDir := filepath.Join(home, ".textile")
		if err := os.MkdirAll(appDir, 0755); err != nil {
			return "", fmt.Errorf(fmt.Sprintf("create repo directory failed: %s", err))
		}
		repoPath = filepath.Join(appDir, "repo")
	}
	return repoPath, nil
}

// Start the node, the API, and the Gateway
// And subsribe to updates of the wallet, thread, and notifications
func startNode(serveDocs bool) error {
	listener := node.ThreadUpdateListener()

	if err := node.Start(); err != nil {
		return err
	}

	// subscribe to wallet updates
	go func() {
		for {
			select {
			case update, ok := <-node.UpdateCh():
				if !ok {
					return
				}
				switch update.Type {
				case pb.WalletUpdate_THREAD_ADDED:
					break
				case pb.WalletUpdate_THREAD_REMOVED:
					break
				case pb.WalletUpdate_ACCOUNT_PEER_ADDED:
					break
				case pb.WalletUpdate_ACCOUNT_PEER_REMOVED:
					break
				}
			}
		}
	}()

	// subscribe to thread updates
	go func() {
		for {
			select {
			case value, ok := <-listener.Ch:
				if !ok {
					return
				}
				if update, ok := value.(*pb.FeedItem); ok {
					thrd := update.Thread[len(update.Thread)-8:]

					btype, err := core.FeedItemType(update)
					if err != nil {
						log.Error(err.Error())
						continue
					}

					payload, err := core.GetFeedItemPayload(update)
					if err != nil {
						log.Error(err.Error())
						continue
					}
					user := payload.GetUser()
					date := payload.GetDate()

					var txt string
					txt += time.Unix(0, util.ProtoNanos(date)).Format(time.RFC822)
					txt += "  "

					if user != nil {
						var name string
						if user.Name != "" {
							name = user.Name
						} else {
							if len(user.Address) >= 7 {
								name = user.Address[:7]
							} else {
								name = user.Address
							}
						}
						txt += name + " "
					}
					txt += "added "

					msg := cmd.Grey(txt) + cmd.Green(btype.String()) + cmd.Grey(" update to "+thrd)
					fmt.Println(msg)
				}
			}
		}
	}()

	// subscribe to notifications
	go func() {
		for {
			select {
			case note, ok := <-node.NotificationCh():
				if !ok {
					return
				}

				date := util.ProtoTime(note.Date).Format(time.RFC822)
				var subject string
				if len(note.Subject) >= 7 {
					subject = note.Subject[len(note.Subject)-7:]
				}

				msg := cmd.Grey(date+"  "+note.User.Name+" ") + cmd.Cyan(note.Body) +
					cmd.Grey(" "+subject)
				fmt.Println(msg)
			}
		}
	}()

	// start apis
	node.StartApi(node.Config().Addresses.API, serveDocs)
	gateway.Host.Start(node.Config().Addresses.Gateway)

	// start profiling api
	go func() {
		writeHeapDump("/debug/write-heap-dump/")
		freeOSMemory("/debug/free-os-memory/")
		mutexFractionOption("/debug/pprof-mutex/")
		if err := http.ListenAndServe(node.Config().Addresses.Profiling, http.DefaultServeMux); err != nil {
			log.Errorf("error staring profile listener: %s", err)
		}
	}()

	// Wait concurrently here until the node comes online
	// that is to say, until the online channel opens
	<-node.OnlineCh()

	// Textile is now online, continue
	return nil
}

// Stop the api, then the gateway, then the node, then if possible, the channels
// If a former fails, do not continue with the latter
func stopNode() error {
	if err := node.StopApi(); err != nil {
		return err
	}
	if err := gateway.Host.Stop(); err != nil {
		return err
	}
	if err := node.Stop(); err != nil {
		return err
	}

	node.CloseChns()
	return nil
}

// Output the instance environment for the daemon command
func printSplash() {
	pid, err := node.PeerId()
	if err != nil {
		log.Fatalf("get peer id failed: %s", err)
	}
	fmt.Println(cmd.Grey("go-textile version: " + common.GitSummary))
	fmt.Println(cmd.Grey("Repo version: ") + cmd.Grey(repo.Repover))
	fmt.Println(cmd.Grey("Repo path: ") + cmd.Grey(node.RepoPath()))
	fmt.Println(cmd.Grey("API address: ") + cmd.Grey(node.ApiAddr()))
	fmt.Println(cmd.Grey("Gateway address: ") + cmd.Grey(gateway.Host.Addr()))
	if node.CafeApiAddr() != "" {
		fmt.Println(cmd.Grey("Cafe address: ") + cmd.Grey(node.CafeApiAddr()))
	}
	fmt.Println(cmd.Grey("System version: ") + cmd.Grey(runtime.GOARCH+"/"+runtime.GOOS))
	fmt.Println(cmd.Grey("Golang version: ") + cmd.Grey(runtime.Version()))
	fmt.Println(cmd.Grey("PeerID:  ") + cmd.Green(pid.Pretty()))
	fmt.Println(cmd.Grey("Account: ") + cmd.Cyan(node.Account().Address()))
}

// If the env var value exists and is not empty, then use that, otherwise use the passed flag
func envOrFlag(env string, flag string) string {
	if os.Getenv(env) != "" {
		return os.Getenv(env)
	}
	return flag
}

// mutexFractionOption allows to set runtime.SetMutexProfileFraction via HTTP
// using POST request with parameter 'fraction'.
func mutexFractionOption(path string) {
	http.DefaultServeMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		asfr := r.Form.Get("fraction")
		if len(asfr) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fr, err := strconv.Atoi(asfr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		log.Infof("Setting MutexProfileFraction to %d", fr)
		runtime.SetMutexProfileFraction(fr)
	})
}

// writeHeapDump writes a description of the heap and the objects in
// it to the given file descriptor. (used here for debugging)
func writeHeapDump(path string) {
	http.DefaultServeMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		log.Infof("Writing heap dump")
		f, err := os.Create("heapdump")
		if err != nil {
			return
		}
		debug.WriteHeapDump(f.Fd())
	})
}

// freeOSMemory forces a garbage collection followed by an
// attempt to return as much memory to the operating system
// as possible. (used here for debugging)
func freeOSMemory(path string) {
	http.DefaultServeMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		log.Infof("Freeing OS memory")
		debug.FreeOSMemory()
	})
}
