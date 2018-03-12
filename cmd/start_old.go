package main
//package cmd
//
//import (
//	"context"
//	"fmt"
//	"time"
//	"net"
//	"os"
//	"path/filepath"
//	"sort"
//	"net/http"
//	"sync"
//	"syscall"
//
//	"github.com/fatih/color"
//	"github.com/op/go-logging"
//	"golang.org/x/crypto/ssh/terminal"
//
//	"github.com/textileio/textile-go/core"
//	"github.com/textileio/textile-go/repo"
//	"github.com/textileio/textile-go/repo/db"
//
//	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/commands"
//	ipfscore "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core"
//	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/corehttp"
//	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/config"
//	utilmain "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/cmd/ipfs/util"
//	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo"
//	lockfile "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo/lock"
//
//	//ipfslogging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
//	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
//	"gx/ipfs/QmRK2LxanhK2gZq6k6R7vk5ZoYZk8ULSSTB7FzDsMUX6CB/go-multiaddr-net"
//	mprome "gx/ipfs/QmSTf3wJXBQk2fxdmXtodvyczrCPgJaK1B1maY78qeebNX/go-metrics-prometheus"
//	"gx/ipfs/QmX3QZ5jHEPidwUrymXV1iSCSUhdGxj15sm2gP4jKMef7B/client_golang/prometheus"
//)
//
//var stdoutLogFormat = logging.MustStringFormatter(
//	`%{color:reset}%{color}%{time:15:04:05.000} [%{shortfunc}] [%{level}] %{message}`,
//)
//
//var fileLogFormat = logging.MustStringFormatter(
//	`%{time:15:04:05.000} [%{shortfunc}] [%{level}] %{message}`,
//)
//
//// defaultMux tells mux to serve path using the default muxer. This is
//// mostly useful to hook up things that register in the default muxer,
//// and don't provide a convenient http.Handler entry point, such as
//// expvar and http/pprof.
//func defaultMux(path string) corehttp.ServeOption {
//	return func(node *ipfscore.IpfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
//		mux.Handle(path, http.DefaultServeMux)
//		return mux, nil
//	}
//}
//
//type Start struct {
//	Password   string   `short:"p" long:"password" description:"the encryption password if the database is encrypted"`
//	LogLevel   string   `short:"l" long:"loglevel" description:"set the logging level [debug, info, notice, warning, error, critical]" defaut:"debug"`
//	NoLogFiles bool     `short:"f" long:"nologfiles" description:"save logs on disk"`
//	DataDir    string   `short:"d" long:"datadir" description:"specify the data directory to be used"`
//	Verbose    bool     `short:"v" long:"verbose" description:"print openbazaar logs to stdout"`
//}
//
//func (x *Start) Execute(args []string) error {
//	printSplashScreen(x.Verbose)
//
//	// Inject metrics before we do anything
//	err := mprome.Inject()
//	if err != nil {
//		log.Errorf("Injecting prometheus handler for metrics failed with message: %s\n", err.Error())
//	}
//
//	// let the user know we're going.
//	fmt.Printf("Initializing daemon...\n")
//
//	if err := utilmain.ManageFdLimit(); err != nil {
//		log.Errorf("setting file descriptor limit: %s", err)
//	}
//
//	// Set repo path
//	repoPath, err := repo.GetRepoPath()
//	if err != nil {
//		return err
//	}
//	if x.DataDir != "" {
//		repoPath = x.DataDir
//	}
//
//	repoLockFile := filepath.Join(repoPath, lockfile.LockFile)
//	os.Remove(repoLockFile)
//
//	// Initialize the repo (ipfs and sqlite)
//	sqliteDB, err := InitializeRepo(repoPath, x.Password, "", time.Now())
//	if err != nil && err != repo.ErrRepoExists {
//		return err
//	}
//
//	//// Logging
//	//w := &lumberjack.Logger{
//	//	Filename:   path.Join(repoPath, "logs", "textile.log"),
//	//	MaxSize:    10, // Megabytes
//	//	MaxBackups: 3,
//	//	MaxAge:     30, // Days
//	//}
//	//var backendStdoutFormatter logging.Backend
//	//if x.Verbose {
//	//	backendStdout := logging.NewLogBackend(os.Stdout, "", 0)
//	//	backendStdoutFormatter = logging.NewBackendFormatter(backendStdout, stdoutLogFormat)
//	//	logging.SetBackend(backendStdoutFormatter)
//	//}
//	//
//	//if !x.NoLogFiles {
//	//	backendFile := logging.NewLogBackend(w, "", 0)
//	//	backendFileFormatter := logging.NewBackendFormatter(backendFile, fileLogFormat)
//	//	if x.Verbose {
//	//		logging.SetBackend(backendFileFormatter, backendStdoutFormatter)
//	//	} else {
//	//		logging.SetBackend(backendFileFormatter)
//	//	}
//	//	ipfslogging.LdJSONFormatter()
//	//	w2 := &lumberjack.Logger{
//	//		Filename:   path.Join(repoPath, "logs", "ipfs.log"),
//	//		MaxSize:    10, // Megabytes
//	//		MaxBackups: 3,
//	//		MaxAge:     30, // Days
//	//	}
//	//	ipfslogging.Output(w2)()
//	//}
//	//
//	//var level logging.Level
//	//switch strings.ToLower(x.LogLevel) {
//	//case "debug":
//	//	level = logging.DEBUG
//	//case "info":
//	//	level = logging.INFO
//	//case "notice":
//	//	level = logging.NOTICE
//	//case "warning":
//	//	level = logging.WARNING
//	//case "error":
//	//	level = logging.ERROR
//	//case "critical":
//	//	level = logging.CRITICAL
//	//default:
//	//	level = logging.DEBUG
//	//}
//	//logging.SetLevel(level, "")
//
//	err = core.CheckAndSetUlimit()
//	if err != nil {
//		return err
//	}
//
//	// If the database cannot be decrypted, exit
//	if sqliteDB.Config().IsEncrypted() {
//		sqliteDB.Close()
//		fmt.Print("Database is encrypted, enter your password: ")
//		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
//		fmt.Println("")
//		pw := string(bytePassword)
//		sqliteDB, err = InitializeRepo(repoPath, pw, "", time.Now())
//		if err != nil && err != repo.ErrRepoExists {
//			return err
//		}
//		if sqliteDB.Config().IsEncrypted() {
//			log.Error("Invalid password")
//			os.Exit(3)
//		}
//	}
//
//	// acquire the repo lock _before_ constructing a node. we need to make
//	// sure we are permitted to access the resources (datastore, etc.)
//	// IPFS node setup
//	r, err := fsrepo.Open(repoPath)
//	switch err {
//	default:
//		log.Error(err)
//		return err
//	case fsrepo.ErrNeedMigration:
//		fmt.Println("Please get fs-repo-migrations from https://dist.ipfs.io")
//		log.Error(err)
//		return err
//	case nil:
//		break
//	}
//
//	// Start assembling node config
//	ncfg := &ipfscore.BuildCfg{
//		Repo:   r,
//		Permanent: true,
//		Online: true,
//		ExtraOpts: map[string]bool{
//			"pubsub": false,
//			"ipnsps": false,
//			"mplex":  true,
//		},
//		Routing: ipfscore.DHTOption,
//	}
//
//	cctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	nd, err := ipfscore.NewNode(cctx, ncfg)
//	if err != nil {
//		log.Error("error from node construction: ", err)
//		return err
//	}
//	nd.SetLocal(false)
//
//	printSwarmAddrs(nd)
//
//	//defer func() {
//	//	// We wait for the node to close first, as the node has children
//	//	// that it will wait for before closing, such as the API server.
//	//	nd.Close()
//	//
//	//	select {
//	//	case <-cctx.Done():
//	//		log.Info("Gracefully shut down daemon")
//	//	default:
//	//	}
//	//}()
//
//	ctx := commands.Context{}
//	ctx.Online = true
//	ctx.ConfigRoot = repoPath
//	ctx.LoadConfig = func(path string) (*config.Config, error) {
//		return fsrepo.ConfigAt(repoPath)
//	}
//	ctx.ConstructNode = func() (*ipfscore.IpfsNode, error) {
//		return nd, nil
//	}
//
//	// Textile node setup
//	core.Node = &core.TextileNode{
//		Context:   ctx,
//		IpfsNode:  nd,
//		RepoPath:  repoPath,
//		Datastore: sqliteDB,
//	}
//
//	// construct api endpoint - every time
//	apiErrc, err := serveHTTPApi(&core.Node.Context)
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//
//	//cfg, err := ctx.GetConfig()
//	//if err != nil {
//	//	log.Error(err)
//	//	return err
//	//}
//	//
//	//// construct http gateway - if it is set in the config
//	//var gwErrc <-chan error
//	//if len(cfg.Addresses.Gateway) > 0 {
//	//	var err error
//	//	gwErrc, err = serveHTTPGateway(&core.Node.Context)
//	//	if err != nil {
//	//		log.Error(err)
//	//		return err
//	//	}
//	//}
//
//	// initialize metrics collector
//	prometheus.MustRegister(&corehttp.IpfsNodeCollector{Node: nd})
//
//	fmt.Printf("Daemon is ready\n")
//	// collect long-running errors and block for shutdown
//	// TODO(cryptix): our fuse currently doesnt follow this pattern for graceful shutdown
//	for err := range merge(apiErrc) { //, gwErrc) {
//		if err != nil {
//			log.Error(err)
//			return err
//		}
//	}
//
//	return nil
//}
//
//func InitializeRepo(dataDir, password, mnemonic string, creationDate time.Time) (*db.SQLiteDatastore, error) {
//	// Database
//	sqliteDB, err := db.Create(dataDir, password)
//	if err != nil {
//		return sqliteDB, err
//	}
//
//	// Initialize the IPFS repo if it does not already exist
//	err = repo.DoInit(dataDir, 4096, password, mnemonic, creationDate, sqliteDB.Config().Init)
//	if err != nil {
//		return sqliteDB, err
//	}
//	return sqliteDB, nil
//}
//
////type DummyWriter struct{}
////
////func (d *DummyWriter) Write(p []byte) (n int, err error) {
////	return 0, nil
////}
////
////type DummyListener struct {
////	addr net.Addr
////}
////
////func (d *DummyListener) Addr() net.Addr {
////	return d.addr
////}
////
////func (d *DummyListener) Accept() (net.Conn, error) {
////	conn, _ := net.FileConn(nil)
////	return conn, nil
////}
////
////func (d *DummyListener) Close() error {
////	return nil
////}
//
//// serveHTTPApi collects options, creates listener, prints status message and starts serving requests
//func serveHTTPApi(cctx *commands.Context) (<-chan error, error) {
//	cfg, err := cctx.GetConfig()
//	if err != nil {
//		return nil, fmt.Errorf("serveHTTPApi: GetConfig() failed: %s", err)
//	}
//
//	apiAddr := cfg.Addresses.API
//	apiMaddr, err := ma.NewMultiaddr(apiAddr)
//	if err != nil {
//		return nil, fmt.Errorf("serveHTTPApi: invalid API address: %q (err: %s)", apiAddr, err)
//	}
//
//	apiLis, err := manet.Listen(apiMaddr)
//	if err != nil {
//		return nil, fmt.Errorf("serveHTTPApi: manet.Listen(%s) failed: %s", apiMaddr, err)
//	}
//	// we might have listened to /tcp/0 - lets see what we are listing on
//	apiMaddr = apiLis.Multiaddr()
//	fmt.Printf("API server listening on %s\n", apiMaddr)
//
//	// by default, we don't let you load arbitrary ipfs objects through the api,
//	// because this would open up the api to scripting vulnerabilities.
//	// only the webui objects are allowed.
//	// if you know what you're doing, go ahead and pass --unrestricted-api.
//	unrestricted := false
//	gatewayOpt := corehttp.GatewayOption(false, corehttp.WebUIPaths...)
//	if unrestricted {
//		gatewayOpt = corehttp.GatewayOption(true, "/ipfs", "/ipns")
//	}
//
//	var opts = []corehttp.ServeOption{
//		corehttp.MetricsCollectionOption("api"),
//		corehttp.CheckVersionOption(),
//		corehttp.CommandsOption(*cctx),
//		corehttp.WebUIOption,
//		gatewayOpt,
//		corehttp.VersionOption(),
//		defaultMux("/debug/vars"),
//		defaultMux("/debug/pprof/"),
//		corehttp.MetricsScrapingOption("/debug/metrics/prometheus"),
//		corehttp.LogOption(),
//	}
//
//	if len(cfg.Gateway.RootRedirect) > 0 {
//		opts = append(opts, corehttp.RedirectOption("", cfg.Gateway.RootRedirect))
//	}
//
//	node, err := cctx.ConstructNode()
//	if err != nil {
//		return nil, fmt.Errorf("serveHTTPApi: ConstructNode() failed: %s", err)
//	}
//
//	if err := node.Repo.SetAPIAddr(apiMaddr); err != nil {
//		return nil, fmt.Errorf("serveHTTPApi: SetAPIAddr() failed: %s", err)
//	}
//
//	errc := make(chan error)
//	go func() {
//		errc <- corehttp.Serve(node, apiLis.NetListener(), opts...)
//		close(errc)
//	}()
//	return errc, nil
//}
//
//// printSwarmAddrs prints the addresses of the host
//func printSwarmAddrs(node *ipfscore.IpfsNode) {
//	if !node.OnlineMode() {
//		fmt.Println("Swarm not listening, running in offline mode.")
//		return
//	}
//
//	var lisAddrs []string
//	ifaceAddrs, err := node.PeerHost.Network().InterfaceListenAddresses()
//	if err != nil {
//		log.Errorf("failed to read listening addresses: %s", err)
//	}
//	for _, addr := range ifaceAddrs {
//		lisAddrs = append(lisAddrs, addr.String())
//	}
//	sort.Sort(sort.StringSlice(lisAddrs))
//	for _, addr := range lisAddrs {
//		fmt.Printf("Swarm listening on %s\n", addr)
//	}
//
//	var addrs []string
//	for _, addr := range node.PeerHost.Addrs() {
//		addrs = append(addrs, addr.String())
//	}
//	sort.Sort(sort.StringSlice(addrs))
//	for _, addr := range addrs {
//		fmt.Printf("Swarm announcing %s\n", addr)
//	}
//
//}
//
//// serveHTTPGateway collects options, creates listener, prints status message and starts serving requests
//func serveHTTPGateway(cctx *commands.Context) (<-chan error, error) {
//	cfg, err := cctx.GetConfig()
//	if err != nil {
//		return nil, fmt.Errorf("serveHTTPGateway: GetConfig() failed: %s", err)
//	}
//
//	gatewayMaddr, err := ma.NewMultiaddr(cfg.Addresses.Gateway)
//	if err != nil {
//		return nil, fmt.Errorf("serveHTTPGateway: invalid gateway address: %q (err: %s)", cfg.Addresses.Gateway, err)
//	}
//
//	writable := cfg.Gateway.Writable
//
//	gwLis, err := manet.Listen(gatewayMaddr)
//	if err != nil {
//		return nil, fmt.Errorf("serveHTTPGateway: manet.Listen(%s) failed: %s", gatewayMaddr, err)
//	}
//	// we might have listened to /tcp/0 - lets see what we are listing on
//	gatewayMaddr = gwLis.Multiaddr()
//
//	if writable {
//		fmt.Printf("Gateway (writable) server listening on %s\n", gatewayMaddr)
//	} else {
//		fmt.Printf("Gateway (readonly) server listening on %s\n", gatewayMaddr)
//	}
//
//	var opts = []corehttp.ServeOption{
//		corehttp.MetricsCollectionOption("gateway"),
//		corehttp.CheckVersionOption(),
//		corehttp.CommandsROOption(*cctx),
//		corehttp.VersionOption(),
//		corehttp.IPNSHostnameOption(),
//		corehttp.GatewayOption(writable, "/ipfs", "/ipns"),
//	}
//
//	if len(cfg.Gateway.RootRedirect) > 0 {
//		opts = append(opts, corehttp.RedirectOption("", cfg.Gateway.RootRedirect))
//	}
//
//	node, err := cctx.ConstructNode()
//	if err != nil {
//		return nil, fmt.Errorf("serveHTTPGateway: ConstructNode() failed: %s", err)
//	}
//
//	errc := make(chan error)
//	go func() {
//		errc <- corehttp.Serve(node, gwLis.NetListener(), opts...)
//		close(errc)
//	}()
//	return errc, nil
//}
//
//// merge does fan-in of multiple read-only error channels
//// taken from http://blog.golang.org/pipelines
//func merge(cs ...<-chan error) <-chan error {
//	var wg sync.WaitGroup
//	out := make(chan error)
//
//	// Start an output goroutine for each input channel in cs.  output
//	// copies values from c to out until c is closed, then calls wg.Done.
//	output := func(c <-chan error) {
//		for n := range c {
//			out <- n
//		}
//		wg.Done()
//	}
//	for _, c := range cs {
//		if c != nil {
//			wg.Add(1)
//			go output(c)
//		}
//	}
//
//	// Start a goroutine to close out once all the output goroutines are
//	// done.  This must start after the wg.Add call.
//	go func() {
//		wg.Wait()
//		close(out)
//	}()
//	return out
//}
//
//func printSplashScreen(verbose bool) {
//	white := color.New(color.FgWhite)
//	white.Println("  __                   __  .__.__")
//	white.Println("_/  |_  ____ ___  ____/  |_|__|  |   ____")
//	white.Println("\\   __\\/ __ \\\\  \\/  /\\   __\\  |  | _/ __ \\")
//	white.Println(" |  | \\  ___/ >    <  |  | |  |  |_\\  ___/")
//	white.Println(" |__|  \\___  >__/\\_ \\ |__| |__|____/\\___  >")
//	white.Println("           \\/      \\/                   \\/")
//	white.DisableColor()
//	fmt.Println("")
//	fmt.Println("textile server v" + core.VERSION)
//	if !verbose {
//		fmt.Println("[Press Ctrl+C to exit]")
//	}
//}
