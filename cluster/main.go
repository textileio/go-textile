package cluster

import (
	"context"
	"fmt"
	"path/filepath"

	ds "github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log"
	ipfscluster "github.com/ipfs/ipfs-cluster"
	"github.com/ipfs/ipfs-cluster/allocator/ascendalloc"
	"github.com/ipfs/ipfs-cluster/allocator/descendalloc"
	"github.com/ipfs/ipfs-cluster/config"
	"github.com/ipfs/ipfs-cluster/consensus/crdt"
	"github.com/ipfs/ipfs-cluster/informer/disk"
	"github.com/ipfs/ipfs-cluster/informer/numpin"
	"github.com/ipfs/ipfs-cluster/monitor/pubsubmon"
	"github.com/ipfs/ipfs-cluster/observations"
	"github.com/ipfs/ipfs-cluster/pintracker/maptracker"
	"github.com/ipfs/ipfs-cluster/pintracker/stateless"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	ma "github.com/multiformats/go-multiaddr"
)

var log = logging.Logger("tex-cluster")

func InitCluster(repoPath, listenAddr string) error {
	cfgMgr, cfgs := makeClusterConfigs()
	err := cfgMgr.Default()
	if err != nil {
		return err
	}

	if listenAddr != "" {
		addr, err := ma.NewMultiaddr(listenAddr)
		if err != nil {
			return err
		}
		cfgs.ClusterCfg.ListenAddr = addr
	}

	return cfgMgr.SaveJSON(ConfigPath(repoPath))
}

func ConfigPath(repoPath string) string {
	return filepath.Join(repoPath, "service.json")
}

func MakeAndLoadConfigs(repoPath string) (*config.Manager, *cfgs, error) {
	cfgMgr, cfgs := makeClusterConfigs()
	err := cfgMgr.LoadJSONFromFile(ConfigPath(repoPath))
	if err != nil {
		return nil, nil, err
	}
	return cfgMgr, cfgs, nil
}

func ParseBootstraps(addrs []string) ([]ma.Multiaddr, error) {
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

func Bootstrap(ctx context.Context, cluster *ipfscluster.Cluster, cons ipfscluster.Consensus, bootstraps []ma.Multiaddr) {
	for _, bstrap := range bootstraps {
		log.Infof("Bootstrapping to %s", bstrap)
		err := cluster.Join(ctx, bstrap)
		if err != nil {
			log.Errorf("bootstrap to %s failed: %s", bstrap, err)
		} else {
			for _, p := range cluster.Peers(ctx) {
				err = cons.Trust(ctx, p.ID)
				if err != nil {
					log.Errorf("failed to trust %s: %s", p.ID.Pretty(), err)
				}
			}
		}
	}
}

func SetupAllocation(
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

func SetupPinTracker(
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

func SetupConsensus(
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

type cfgs struct {
	ClusterCfg          *ipfscluster.Config
	CrdtCfg             *crdt.Config
	MaptrackerCfg       *maptracker.Config
	StatelessTrackerCfg *stateless.Config
	PubsubmonCfg        *pubsubmon.Config
	DiskInfCfg          *disk.Config
	NumpinInfCfg        *numpin.Config
	TracingCfg          *observations.TracingConfig
}

func makeClusterConfigs() (*config.Manager, *cfgs) {
	cfg := config.NewManager()
	clusterCfg := &ipfscluster.Config{}
	crdtCfg := &crdt.Config{}
	maptrackerCfg := &maptracker.Config{}
	statelessCfg := &stateless.Config{}
	pubsubmonCfg := &pubsubmon.Config{}
	diskInfCfg := &disk.Config{}
	numpinInfCfg := &numpin.Config{}
	tracingCfg := &observations.TracingConfig{}
	cfg.RegisterComponent(config.Cluster, clusterCfg)
	cfg.RegisterComponent(config.Consensus, crdtCfg)
	cfg.RegisterComponent(config.PinTracker, maptrackerCfg)
	cfg.RegisterComponent(config.PinTracker, statelessCfg)
	cfg.RegisterComponent(config.Monitor, pubsubmonCfg)
	cfg.RegisterComponent(config.Informer, diskInfCfg)
	cfg.RegisterComponent(config.Informer, numpinInfCfg)
	cfg.RegisterComponent(config.Observations, tracingCfg)
	return cfg, &cfgs{
		clusterCfg,
		crdtCfg,
		maptrackerCfg,
		statelessCfg,
		pubsubmonCfg,
		diskInfCfg,
		numpinInfCfg,
		tracingCfg,
	}
}
