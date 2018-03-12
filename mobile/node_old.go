package mobile
//
//import (
//	"context"
//	"os"
//	"path/filepath"
//
//	//manet "gx/ipfs/QmRK2LxanhK2gZq6k6R7vk5ZoYZk8ULSSTB7FzDsMUX6CB/go-multiaddr-net"
//	//ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
//	//
//	//"github.com/textileio/textile-go/api"
//	//"github.com/textileio/textile-go/repo"
//	//"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core/corehttp"
//
//	//"errors"
//	"fmt"
//	//routing "gx/ipfs/QmTiWLZ6Fo5j4KcTVutZJ5KWRRJrbxzmxA4td8NfEdrPh7/go-libp2p-routing"
//	//dht "gx/ipfs/QmUCS9EnqNq1kCnJds2eLDypBiS21aSiCf1MVzSUVB9TGA/go-libp2p-kad-dht"
//	//dhtutil "gx/ipfs/QmUCS9EnqNq1kCnJds2eLDypBiS21aSiCf1MVzSUVB9TGA/go-libp2p-kad-dht/util"
//	//proto "gx/ipfs/QmZ4Qi3GaRbjcx28Sme5eMH7RQjGkt8wHxt2a65oLaeFEV/gogo-protobuf/proto"
//	//p2phost "gx/ipfs/QmNmJZL7FQySMtE2BQuLMuZg2EB2CLEunJJUSVSc9YnnbV/go-libp2p-host"
//	//recpb "gx/ipfs/QmUpttFinNDmNPgFwKN8sZK6BUtBmA68Y4KdSBDXa8t9sJ/go-libp2p-record/pb"
//	//"io/ioutil"
//	//"net/http"
//	//"path"
//	"time"
//
//	"github.com/op/go-logging"
//	"github.com/textileio/textile-go/core"
//	"github.com/textileio/textile-go/repo"
//	//"github.com/textileio/textile-go/ipfs"
//	"github.com/textileio/textile-go/repo/db"
//	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/commands"
//	ipfscore "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/core"
//	//"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/namesys"
//	//namepb "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/namesys/pb"
//	//ipath "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/path"
//	//ipfsrepo "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo"
//	ipfsconfig "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/config"
//	"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo"
//	lockfile "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/repo/fsrepo/lock"
//	//"gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/thirdparty/ds-help"
//	//"gx/ipfs/QmdQTPWduSeyveSxeCAte33M592isSW5Z979g81aJphrgn/go-ipfs-ds-help"
//	//"gx/ipfs/QmVSep2WwKcXxMonPASsAJ3nZVjfVMKgMcaSigxKnUWpJv/go-libp2p-kad-dht"
//)
//
//type Node struct {
//	node       *core.TextileNode
//	config     NodeConfig
//	cancel     context.CancelFunc
//	ipfsConfig *ipfscore.BuildCfg
//	//apiConfig  *repo.APIConfig
//}
//type Mobile struct{}
//
//type Start struct {
//	Password   string `short:"p" long:"password" description:"the encryption password if the database is encrypted"`
//	LogLevel   string `short:"l" long:"loglevel" description:"set the logging level [debug, info, notice, warning, error, critical]" defaut:"debug"`
//	NoLogFiles bool   `short:"f" long:"nologfiles" description:"save logs on disk"`
//	DataDir    string `short:"d" long:"datadir" description:"specify the data directory to be used"`
//	Verbose    bool   `short:"v" long:"verbose" description:"print openbazaar logs to stdout"`
//}
//
//func (x *Start) Execute(args []string) error {
//	node := NewTextile(x.DataDir, x.Password, "")
//	node.Start()
//
//	//hash, err := ipfs.AddFile(node.node.Context, path.Join(node.node.RepoPath, "root", "files", "random.png"))
//	//if err != nil {
//	//	return err
//	//}
//	//fmt.Println(hash)
//
//	//data, err := ipfs.Cat(node.node.Context, "QmVTxHB89nbrS2qs3NTB7o3NR8mqMrA3BXM6kvoDJMKh99", 60 * time.Second)
//	//mhash, err := multihash.FromHexString("QmVTxHB89nbrS2qs3NTB7o3NR8mqMrA3BXM6kvoDJMKh99")
//	//if err != nil {
//	//	return err
//	//}
//	//ci := cid.NewCidV1(cid.Raw, mhash)
//	//
//	//node.node.IpfsNode.Pinning.PinWithMode(ci, pin.Recursive)
//	//node.node.IpfsNode.Pinning.Pin()
//
//	return nil
//}
//
//func NewTextile(repoPath string, password string, mnemonic string) *Node {
//
//	nodeconfig := NodeConfig{
//		RepoPath: repoPath,
//	}
//
//	var m Mobile
//	node, err := m.NewNode(nodeconfig, password, mnemonic)
//	if err != nil {
//		fmt.Println(err)
//	}
//	return node
//}
//
//func (m *Mobile) NewNode(config NodeConfig, password string, mnemonic string) (*Node, error) {
//
//	repoLockFile := filepath.Join(config.RepoPath, lockfile.LockFile)
//	os.Remove(repoLockFile)
//
//	// Logging
//	backendStdout := logging.NewLogBackend(os.Stdout, "", 0)
//	logger = logging.NewBackendFormatter(backendStdout, stdoutLogFormat)
//	logging.SetBackend(logger)
//
//	sqliteDB, err := initializeRepo(config.RepoPath, password, mnemonic, time.Now())
//	if err != nil && err != repo.ErrRepoExists {
//		return nil, err
//	}
//
//	// Get creation date. Ignore the error and use a default timestamp.
//	//creationDate, _ := sqliteDB.Config().GetCreationDate()
//
//	// Load config
//	//configFile, err := ioutil.ReadFile(path.Join(config.RepoPath, "config"))
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	//apiConfig, err := repo.GetAPIConfig(configFile)
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	//resolverConfig, err := repo.GetResolverConfig(configFile)
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	// Create user-agent file
//	//userAgentBytes := []byte(core.USERAGENT)
//	//ioutil.WriteFile(path.Join(config.RepoPath, "root", "user_agent"), userAgentBytes, os.ModePerm)
//
//	// IPFS node setup
//	r, err := fsrepo.Open(config.RepoPath)
//	if err != nil {
//		return nil, err
//	}
//
//	//cfg, err := r.Config()
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//identityKey, err := sqliteDB.Config().GetIdentityKey()
//	//if err != nil {
//	//	return nil, err
//	//}
//	//identity, err := ipfs.IdentityFromKey(identityKey)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//cfg.Identity = identity
//	//cfg.Swarm.DisableNatPortMap = true
//
//	// Setup testnet
//	//if config.Testnet {
//	//	testnetBootstrapAddrs, err := repo.GetTestnetBootstrapAddrs(configFile)
//	//	if err != nil {
//	//		return nil, err
//	//	}
//	//	cfg.Bootstrap = testnetBootstrapAddrs
//	//	//dht.ProtocolDHT = "/openbazaar/kad/testnet/1.0.0"
//	//	//bitswap.ProtocolBitswap = "/openbazaar/bitswap/testnet/1.1.0"
//	//	//service.ProtocolOpenBazaar = "/openbazaar/app/testnet/1.0.0"
//	//	//
//	//	//dataSharing.PushTo = []string{}
//	//}
//
//	ncfg := &ipfscore.BuildCfg{
//		Repo:    r,
//		Online:  true,
//		//Routing: DHTClientOption,
//	}
//
//	// Set IPNS query size
//	//querySize := cfg.Ipns.QuerySize
//	//if querySize <= 20 && querySize > 0 {
//	//	dhtutil.QuerySize = int(querySize)
//	//} else {
//	//	dhtutil.QuerySize = 16
//	//}
//	//namesys.UsePersistentCache = cfg.Ipns.UsePersistentCache
//
//	// Wallet
//	//mn, err := sqliteDB.Config().GetMnemonic()
//	//if err != nil {
//	//	return nil, err
//	//}
//	//var params chaincfg.Params
//	//if config.Testnet {
//	//	params = chaincfg.TestNet3Params
//	//} else {
//	//	params = chaincfg.MainNetParams
//	//}
//
//	//var tp net.Addr
//
//	//spvwalletConfig := &bitcoincash.Config{
//	//	Mnemonic:             mn,
//	//	Params:               &params,
//	//	MaxFee:               uint64(walletCfg.MaxFee),
//	//	LowFee:               uint64(walletCfg.LowFeeDefault),
//	//	MediumFee:            uint64(walletCfg.MediumFeeDefault),
//	//	HighFee:              uint64(walletCfg.HighFeeDefault),
//	//	FeeAPI:               *feeApi,
//	//	RepoPath:             config.RepoPath,
//	//	CreationDate:         creationDate,
//	//	DB:                   sqliteDB,
//	//	UserAgent:            "OpenBazaar",
//	//	TrustedPeer:          tp,
//	//	Logger:               logger,
//	//	ExchangeRateProvider: exchangeRates,
//	//}
//	//core.PublishLock.Lock()
//
//	// Set up the ban manager
//	//settings, err := sqliteDB.Settings().Get()
//	//if err != nil && err != db.SettingsNotSetError {
//	//	return nil, err
//	//}
//	//var blockedNodes []peer.ID
//	//if settings.BlockedNodes != nil {
//	//	for _, pid := range *settings.BlockedNodes {
//	//		id, err := peer.IDB58Decode(pid)
//	//		if err != nil {
//	//			continue
//	//		}
//	//		blockedNodes = append(blockedNodes, id)
//	//	}
//	//}
//	//bm := obnet.NewBanManager(blockedNodes)
//
//	// Create namesys resolvers
//	//resolvers := []obns.Resolver{
//	//	bstk.NewBlockStackClient(resolverConfig.Id, nil),
//	//	obns.NewDNSResolver(),
//	//}
//	//ns, err := obns.NewNameSystem(resolvers)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//// Push nodes
//	//var pushNodes []peer.ID
//	//for _, pnd := range dataSharing.PushTo {
//	//	p, err := peer.IDB58Decode(pnd)
//	//	if err != nil {
//	//		return nil, err
//	//	}
//	//	pushNodes = append(pushNodes, p)
//	//}
//
//	// OpenBazaar node setup
//	core.Node = &core.TextileNode{
//		RepoPath:  config.RepoPath,
//		Datastore: sqliteDB,
//		//Wallet:        wallet,
//		//NameSystem:    ns,
//		//ExchangeRates: exchangeRates,
//		//UserAgent:     core.USERAGENT,
//		//PushNodes:     pushNodes,
//		//BanManager:    bm,
//	}
//
//	//if len(cfg.Addresses.Gateway) <= 0 {
//	//	return nil, errors.New("No gateway addresses configured")
//	//}
//
//	//return &Node{config: config, node: core.Node, ipfsConfig: ncfg, apiConfig: apiConfig}, nil
//	return &Node{config: config, node: core.Node, ipfsConfig: ncfg}, nil
//}
//
//func (n *Node) startIPFSNode(repoPath string, config *ipfscore.BuildCfg) (*ipfscore.IpfsNode, commands.Context, error) {
//	cctx, cancel := context.WithCancel(context.Background())
//	n.cancel = cancel
//
//	ctx := commands.Context{}
//	nd, err := ipfscore.NewNode(cctx, config)
//	if err != nil {
//		return nil, ctx, err
//	}
//
//	ctx.Online = true
//	ctx.ConfigRoot = repoPath
//	ctx.LoadConfig = func(path string) (*ipfsconfig.Config, error) {
//		return fsrepo.ConfigAt(repoPath)
//	}
//	ctx.ConstructNode = func() (*ipfscore.IpfsNode, error) {
//		return nd, nil
//	}
//	return nd, ctx, nil
//}
//
//func (n *Node) Start() error {
//	fmt.Println("Starting IPFS Node")
//	fmt.Println("Repository: ", n.config.RepoPath)
//	nd, ctx, err := n.startIPFSNode(n.config.RepoPath, n.ipfsConfig)
//	if err != nil {
//		return err
//	}
//
//	n.node.Context = ctx
//	n.node.IpfsNode = nd
//
//	// Get current directory root hash
//	//fmt.Println("Getting IPNS keys")
//	//_, ipnskey := namesys.IpnsKeysForID(nd.Identity)
//	//ival, hasherr := nd.Repo.Datastore().Get(dshelp.NewKeyFromBinary([]byte(ipnskey)))
//	//if hasherr != nil {
//	//	return hasherr
//	//}
//	//val := ival.([]byte)
//	//dhtrec := new(recpb.Record)
//	//proto.Unmarshal(val, dhtrec)
//	//e := new(namepb.IpnsEntry)
//	//proto.Unmarshal(dhtrec.GetValue(), e)
//	//n.node.RootHash = ipath.Path(e.Value).String()
//	//
//	//fmt.Println("Reading config file...")
//	//configFile, err := ioutil.ReadFile(path.Join(n.node.RepoPath, "config"))
//	//if err != nil {
//	//	return err
//	//}
//	//republishInterval, err := repo.GetRepublishInterval(configFile)
//	//if err != nil {
//	//	return err
//	//}
//	//
//	//// Offline messaging storage
//	//n.node.MessageStorage = selfhosted.NewSelfHostedStorage(n.node.RepoPath, ctx, n.node.PushNodes, n.node.SendStore)
//
//	// Start gateway
//	// Create authentication cookie
//	//var authCookie http.Cookie
//	//authCookie.Name = "OpenBazaar_Auth_Cookie"
//	//
//	//if n.config.AuthenticationToken != "" {
//	//	authCookie.Value = n.config.AuthenticationToken
//	//	n.apiConfig.Authenticated = true
//	//}
//	//fmt.Println("Starting HTTP Gateway...")
//	//gateway, err := newHTTPGateway(core.Node, authCookie, *n.apiConfig)
//	//if err != nil {
//	//	return err
//	//}
//	//err = gateway.Serve()
//	//if err != nil {
//	//	return err
//	//}
//	//go gateway.Serve()
//	//fmt.Println("Gateway serving...")
//
//	//go func() {
//	//	<-dht.DefaultBootstrapConfig.DoneChan
//	//	fmt.Println("Creating new node service...")
//	//	n.node.Service = service.New(n.node, n.node.Context, n.node.Datastore)
//	//	mrCfg := ret.MRConfig{
//	//		Db:        n.node.Datastore,
//	//		Ctx:       n.node.Context,
//	//		IPFSNode:  nd,
//	//		BanManger: n.node.BanManager,
//	//		Service:   core.Node.Service,
//	//		PrefixLen: 14,
//	//		PushNodes: core.Node.PushNodes,
//	//		Dialer:    nil,
//	//		SendAck:   core.Node.SendOfflineAck,
//	//		SendError: core.Node.SendError,
//	//	}
//	//	MR := ret.NewMessageRetriever(mrCfg)
//	//	go MR.Run()
//	//	n.node.MessageRetriever = MR
//	//	PR := rep.NewPointerRepublisher(n.node.IpfsNode, n.node.Datastore, n.node.PushNodes, n.node.IsModerator)
//	//	go PR.Run()
//	//	n.node.PointerRepublisher = PR
//	//	MR.Wait()
//	//	TL := lis.NewTransactionListener(n.node.Datastore, n.node.Broadcast, n.node.Wallet)
//	//	WL := lis.NewWalletListener(n.node.Datastore, n.node.Broadcast)
//	//	n.node.Wallet.AddTransactionListener(TL.OnTransactionReceived)
//	//	n.node.Wallet.AddTransactionListener(WL.OnTransactionReceived)
//	//	su := bitcoin.NewStatusUpdater(n.node.Wallet, n.node.Broadcast, n.node.IpfsNode.Context())
//	//	fmt.Println("Starting Status Updater...")
//	//	go su.Start()
//	//	fmt.Println("Starting wallet...")
//	//	go n.node.Wallet.Start()
//	//
//	//	core.PublishLock.Unlock()
//	//	core.Node.UpdateFollow()
//	//	if !core.InitalPublishComplete {
//	//		core.Node.SeedNode()
//	//	}
//	//	fmt.Println("Initial publish complete")
//	//	core.Node.SetUpRepublisher(republishInterval)
//	//}()
//
//	return nil
//}
//
//func (n *Node) Stop() error {
//	core.Node.Datastore.Close()
//	repoLockFile := filepath.Join(core.Node.RepoPath, lockfile.LockFile)
//	os.Remove(repoLockFile)
//	core.Node.IpfsNode.Close()
//	return nil
//}
//
//func initializeRepo(dataDir, password, mnemonic string, creationDate time.Time) (*db.SQLiteDatastore, error) {
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
//// Collects options, creates listener, prints status message and starts serving requests
////func newHTTPGateway(node *core.TextileNode, authCookie http.Cookie, config repo.APIConfig) (*api.Gateway, error) {
////	// Get API configuration
////	fmt.Println("Getting config...")
////	cfg, err := node.Context.GetConfig()
////
////	if err != nil {
////		return nil, err
////	}
////
////	// Create a network listener
////	fmt.Println("Getting gateway multi-address...")
////	gatewayMaddr, err := ma.NewMultiaddr(cfg.Addresses.Gateway)
////	if err != nil {
////		return nil, fmt.Errorf("newHTTPGateway: invalid gateway address: %q (err: %s)", cfg.Addresses.Gateway, err)
////	}
////
////	gwLis, err := manet.Listen(gatewayMaddr)
////	if err != nil {
////		return nil, fmt.Errorf("newHTTPGateway: manet.Listen(%s) failed: %s", gatewayMaddr, err)
////	}
////
////	// We might have listened to /tcp/0 - let's see what we are listing on
////	gatewayMaddr = gwLis.Multiaddr()
////
////	// Setup an options slice
////	var opts = []corehttp.ServeOption{
////		corehttp.MetricsCollectionOption("gateway"),
////		corehttp.CommandsROOption(node.Context),
////		corehttp.VersionOption(),
////		corehttp.IPNSHostnameOption(),
////		corehttp.GatewayOption(cfg.Gateway.Writable, "/ipfs", "/ipns"),
////	}
////
////	if len(cfg.Gateway.RootRedirect) > 0 {
////		opts = append(opts, corehttp.RedirectOption("", cfg.Gateway.RootRedirect))
////	}
////
////	if err != nil {
////		return nil, fmt.Errorf("newHTTPGateway: ConstructNode() failed: %s", err)
////	}
////
////	return api.NewGateway(node, authCookie, gwLis.NetListener(), config, logger, opts...)
////}
//
////var DHTClientOption ipfscore.RoutingOption = constructClientDHTRouting
////
////func constructClientDHTRouting(ctx context.Context, host p2phost.Host, dstore ipfsrepo.Datastore) (routing.IpfsRouting, error) {
////	dhtRouting := dht.NewDHTClient(ctx, host, dstore)
////	dhtRouting.Validator[ipfscore.IpnsValidatorTag] = namesys.IpnsRecordValidator
////	dhtRouting.Selector[ipfscore.IpnsValidatorTag] = namesys.IpnsSelectorFunc
////	return dhtRouting, nil
////}
