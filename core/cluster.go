package core

import (
	"context"
	"fmt"
	"time"

	ds "github.com/ipfs/go-datastore"
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

type cfgs struct {
	clusterCfg          *ipfscluster.Config
	crdtCfg             *crdt.Config
	maptrackerCfg       *maptracker.Config
	statelessTrackerCfg *stateless.Config
	pubsubmonCfg        *pubsubmon.Config
	diskInfCfg          *disk.Config
	numpinInfCfg        *numpin.Config
}

func makeConfigs() (*config.Manager, *cfgs) {
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
	return cfg, &cfgs{
		clusterCfg,
		crdtCfg,
		maptrackerCfg,
		statelessCfg,
		pubsubmonCfg,
		diskInfCfg,
		numpinInfCfg,
	}
}

func makeAndLoadConfigs(repoPath string) (*config.Manager, *cfgs, error) {
	cfgMgr, cfgs := makeConfigs()
	err := cfgMgr.LoadJSONFileAndEnv(repoPath)
	if err != nil {
		return nil, nil, err
	}
	return cfgMgr, cfgs, nil
}

// createCluster creates all the necessary things to produce the cluster object
func (t *Textile) createCluster(ctx context.Context, bootstraps []ma.Multiaddr) (*ipfscluster.Cluster, error) {
	cfgMgr, cfgs, err := makeAndLoadConfigs(t.repoPath)
	if err != nil {
		return nil, err
	}
	defer cfgMgr.Shutdown()

	cfgs.clusterCfg.LeaveOnShutdown = true

	tracker, err := setupPinTracker(
		"map",
		t.node.PeerHost,
		cfgs.maptrackerCfg,
		cfgs.statelessTrackerCfg,
		cfgs.clusterCfg.Peername,
	)
	if err != nil {
		return nil, err
	}

	informer, alloc, err := setupAllocation(
		"disk-freespace",
		cfgs.diskInfCfg,
		cfgs.numpinInfCfg,
	)
	if err != nil {
		return nil, err
	}

	ipfscluster.ReadyTimeout = raft.DefaultWaitForLeaderTimeout + 5*time.Second

	cons, err := setupConsensus(
		t.node.PeerHost,
		t.node.DHT,
		t.node.PubSub,
		cfgs.crdtCfg,
		t.node.Repo.Datastore(),
	)
	if err != nil {
		return nil, err
	}

	var peersF func(context.Context) ([]peer.ID, error)
	mon, err := pubsubmon.New(ctx, cfgs.pubsubmonCfg, t.node.PubSub, peersF)
	if err != nil {
		return nil, err
	}

	cluster, err := ipfscluster.NewCluster(
		ctx,
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
		return nil, err
	}

	// noop if no bootstraps
	// if bootstrapping fails, consensus will never be ready
	// and timeout. So this can happen in background and we
	// avoid worrying about error handling here (since Cluster
	// will realize).
	go bootstrap(ctx, cluster, bootstraps)

	return cluster, nil
}

// bootstrap will bootstrap this peer to one of the bootstrap addresses
// if there are any.
func bootstrap(ctx context.Context, cluster *ipfscluster.Cluster, bootstraps []ma.Multiaddr) {
	for _, bstrap := range bootstraps {
		log.Infof("Bootstrapping to %s", bstrap)
		err := cluster.Join(ctx, bstrap)
		if err != nil {
			log.Errorf("bootstrap to %s failed: %s", bstrap, err)
		}
	}
}

func setupAllocation(
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

func setupPinTracker(
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

func setupConsensus(
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
