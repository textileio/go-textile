package core

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	ds "github.com/ipfs/go-datastore"
	util "github.com/ipfs/go-ipfs-util"
	ipfscluster "github.com/ipfs/ipfs-cluster"
	"github.com/ipfs/ipfs-cluster/allocator/ascendalloc"
	"github.com/ipfs/ipfs-cluster/allocator/descendalloc"
	"github.com/ipfs/ipfs-cluster/config"
	"github.com/ipfs/ipfs-cluster/consensus/crdt"
	"github.com/ipfs/ipfs-cluster/consensus/raft"
	"github.com/ipfs/ipfs-cluster/informer/disk"
	"github.com/ipfs/ipfs-cluster/informer/numpin"
	"github.com/ipfs/ipfs-cluster/monitor/pubsubmon"
	"github.com/ipfs/ipfs-cluster/pintracker/maptracker"
	"github.com/ipfs/ipfs-cluster/pintracker/stateless"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/textileio/go-textile/ipfs"
)

type clusterCfgs struct {
	clusterCfg          *ipfscluster.Config
	crdtCfg             *crdt.Config
	maptrackerCfg       *maptracker.Config
	statelessTrackerCfg *stateless.Config
	pubsubmonCfg        *pubsubmon.Config
	diskInfCfg          *disk.Config
	numpinInfCfg        *numpin.Config
}

func makeClusterConfigs() (*config.Manager, *clusterCfgs) {
	cfg := config.NewManager()
	clusterCfg := &ipfscluster.Config{}
	crdtCfg := &crdt.Config{}
	maptrackerCfg := &maptracker.Config{}
	statelessCfg := &stateless.Config{}
	pubsubmonCfg := &pubsubmon.Config{}
	diskInfCfg := &disk.Config{}
	numpinInfCfg := &numpin.Config{}
	cfg.RegisterComponent(config.Cluster, clusterCfg)
	cfg.RegisterComponent(config.Consensus, crdtCfg)
	cfg.RegisterComponent(config.PinTracker, maptrackerCfg)
	cfg.RegisterComponent(config.PinTracker, statelessCfg)
	cfg.RegisterComponent(config.Monitor, pubsubmonCfg)
	cfg.RegisterComponent(config.Informer, diskInfCfg)
	cfg.RegisterComponent(config.Informer, numpinInfCfg)
	return cfg, &clusterCfgs{
		clusterCfg,
		crdtCfg,
		maptrackerCfg,
		statelessCfg,
		pubsubmonCfg,
		diskInfCfg,
		numpinInfCfg,
	}
}

func clusterConfigPath(repoPath string) string {
	return filepath.Join(repoPath, "service.json")
}

func makeAndLoadClusterConfigs(repoPath string) (*config.Manager, *clusterCfgs, error) {
	cfgMgr, cfgs := makeClusterConfigs()
	err := cfgMgr.LoadJSONFromFile(clusterConfigPath(repoPath))
	if err != nil {
		return nil, nil, err
	}
	return cfgMgr, cfgs, nil
}

func parseClusterBootstraps(addrs []string) ([]ma.Multiaddr, error) {
	var parsed []ma.Multiaddr
	for _, a := range addrs {
		p, err := ma.NewMultiaddr(a)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, p)
	}
	return parsed, nil
}

func initCluster(repoPath, secret string) error {
	decoded, err := ipfscluster.DecodeClusterSecret(secret)
	if err != nil {
		return err
	}

	cfgMgr, cfgs := makeClusterConfigs()
	err = cfgMgr.Default()
	if err != nil {
		return err
	}
	cfgs.clusterCfg.Secret = decoded

	return cfgMgr.SaveJSON(clusterConfigPath(repoPath))
}

func (t *Textile) clusterExists() bool {
	return util.FileExists(clusterConfigPath(t.repoPath))
}

// startCluster creates all the necessary things to produce the cluster object
func (t *Textile) startCluster() error {
	cfgMgr, cfgs, err := makeAndLoadClusterConfigs(t.repoPath)
	if err != nil {
		return nil
	}
	defer cfgMgr.Shutdown()

	cfgs.clusterCfg.LeaveOnShutdown = true

	tracker, err := setupClusterPinTracker(
		"map",
		t.node.PeerHost,
		cfgs.maptrackerCfg,
		cfgs.statelessTrackerCfg,
		cfgs.clusterCfg.Peername,
	)
	if err != nil {
		return nil
	}

	informer, alloc, err := setupClusterAllocation(
		"disk-freespace",
		cfgs.diskInfCfg,
		cfgs.numpinInfCfg,
	)
	if err != nil {
		return nil
	}

	ipfscluster.ReadyTimeout = raft.DefaultWaitForLeaderTimeout + 5*time.Second

	cons, err := setupClusterConsensus(
		t.node.PeerHost,
		t.node.DHT,
		t.node.PubSub,
		cfgs.crdtCfg,
		t.node.Repo.Datastore(),
	)
	if err != nil {
		return nil
	}

	var peersF func(context.Context) ([]peer.ID, error)
	mon, err := pubsubmon.New(t.node.Context(), cfgs.pubsubmonCfg, t.node.PubSub, peersF)
	if err != nil {
		return nil
	}

	t.cluster, err = ipfscluster.NewCluster(
		t.node.Context(),
		t.node.PeerHost,
		t.node.DHT,
		cfgs.clusterCfg,
		t.node.Repo.Datastore(),
		cons,
		nil,
		ipfs.NewClusterConnector(t.node),
		tracker,
		mon,
		alloc,
		informer,
		nil,
	)
	if err != nil {
		return nil
	}

	bootstraps, err := parseClusterBootstraps(t.config.Cluster.Bootstraps)
	if err != nil {
		return nil
	}

	// noop if no bootstraps
	// if bootstrapping fails, consensus will never be ready
	// and timeout. So this can happen in background and we
	// avoid worrying about error handling here (since Cluster
	// will realize).
	go bootstrapCluster(t.node.Context(), t.cluster, bootstraps)

	return nil
}

// bootstrap will bootstrap this peer to one of the bootstrap addresses
// if there are any.
func bootstrapCluster(ctx context.Context, cluster *ipfscluster.Cluster, bootstraps []ma.Multiaddr) {
	for _, bstrap := range bootstraps {
		log.Infof("Bootstrapping to %s", bstrap)
		err := cluster.Join(ctx, bstrap)
		if err != nil {
			log.Errorf("bootstrap to %s failed: %s", bstrap, err)
		}
	}
}

func setupClusterAllocation(
	name string,
	diskInfCfg *disk.Config,
	numpinInfCfg *numpin.Config,
) (ipfscluster.Informer, ipfscluster.PinAllocator, error) {
	switch name {
	case "disk", "disk-freespace":
		informer, err := disk.NewInformer(diskInfCfg)
		if err != nil {
			return nil, nil, err
		}
		return informer, descendalloc.NewAllocator(), nil
	case "disk-reposize":
		informer, err := disk.NewInformer(diskInfCfg)
		if err != nil {
			return nil, nil, err
		}
		return informer, ascendalloc.NewAllocator(), nil
	case "numpin", "pincount":
		informer, err := numpin.NewInformer(numpinInfCfg)
		if err != nil {
			return nil, nil, err
		}
		return informer, ascendalloc.NewAllocator(), nil
	default:
		return nil, nil, fmt.Errorf("unknown allocation strategy")
	}
}

func setupClusterPinTracker(
	name string,
	h host.Host,
	mapCfg *maptracker.Config,
	statelessCfg *stateless.Config,
	peerName string,
) (ipfscluster.PinTracker, error) {
	switch name {
	case "map":
		ptrk := maptracker.NewMapPinTracker(mapCfg, h.ID(), peerName)
		log.Debug("map pintracker loaded")
		return ptrk, nil
	case "stateless":
		ptrk := stateless.New(statelessCfg, h.ID(), peerName)
		log.Debug("stateless pintracker loaded")
		return ptrk, nil
	default:
		return nil, fmt.Errorf("unknown pintracker type")
	}
}

func setupClusterConsensus(
	h host.Host,
	dht *dht.IpfsDHT,
	pubsub *pubsub.PubSub,
	crdtCfg *crdt.Config,
	store ds.Datastore,
) (ipfscluster.Consensus, error) {
	convrdt, err := crdt.New(h, dht, pubsub, crdtCfg, store)
	if err != nil {
		return nil, fmt.Errorf("error creating CRDT component: %s", err)
	}
	return convrdt, nil
}
