package core

import (
	"context"
	"time"

	"github.com/ipfs/ipfs-cluster/observations"

	util "github.com/ipfs/go-ipfs-util"
	ipfscluster "github.com/ipfs/ipfs-cluster"
	capi "github.com/ipfs/ipfs-cluster/api"
	"github.com/ipfs/ipfs-cluster/consensus/raft"
	"github.com/ipfs/ipfs-cluster/monitor/pubsubmon"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/textileio/go-textile/cluster"
	"github.com/textileio/go-textile/ipfs"
)

func (t *Textile) clusterExists() bool {
	return util.FileExists(cluster.ConfigPath(t.repoPath))
}

// startCluster creates all the necessary things to produce the cluster object
func (t *Textile) startCluster() error {
	cfgMgr, cfgs, err := cluster.MakeAndLoadConfigs(t.repoPath)
	if err != nil {
		return err
	}
	defer cfgMgr.Shutdown()

	cfgs.ClusterCfg.LeaveOnShutdown = true

	tracker, err := cluster.SetupPinTracker(
		"map",
		t.node.PeerHost,
		cfgs.MaptrackerCfg,
		cfgs.StatelessTrackerCfg,
		cfgs.ClusterCfg.Peername,
	)
	if err != nil {
		return err
	}

	informer, alloc, err := cluster.SetupAllocation(
		"disk-freespace",
		cfgs.DiskInfCfg,
		cfgs.NumpinInfCfg,
	)
	if err != nil {
		return err
	}

	ipfscluster.ReadyTimeout = raft.DefaultWaitForLeaderTimeout + 5*time.Second

	cons, err := cluster.SetupConsensus(
		t.node.PeerHost,
		t.node.DHT,
		t.node.PubSub,
		cfgs.CrdtCfg,
		t.node.Repo.Datastore(),
	)
	if err != nil {
		return err
	}

	tracer, err := observations.SetupTracing(cfgs.TracingCfg)
	if err != nil {
		return err
	}

	var peersF func(context.Context) ([]peer.ID, error)
	mon, err := pubsubmon.New(t.node.Context(), cfgs.PubsubmonCfg, t.node.PubSub, peersF)
	if err != nil {
		return err
	}

	connector, err := ipfs.NewClusterConnector(t.node, func(ctx context.Context) []*capi.ID {
		return t.cluster.Peers(ctx)
	})
	if err != nil {
		return err
	}

	t.cluster, err = ipfscluster.NewCluster(
		t.node.Context(),
		t.node.PeerHost,
		t.node.DHT,
		cfgs.ClusterCfg,
		t.node.Repo.Datastore(),
		cons,
		nil,
		connector,
		tracker,
		mon,
		alloc,
		informer,
		tracer,
	)
	if err != nil {
		return err
	}

	bootstraps, err := cluster.ParseBootstraps(t.config.Cluster.Bootstraps)
	if err != nil {
		return err
	}

	// noop if no bootstraps
	// if bootstrapping fails, consensus will never be ready
	// and timeout. So this can happen in background and we
	// avoid worrying about error handling here (since Cluster
	// will realize).
	go cluster.Bootstrap(t.node.Context(), t.cluster, bootstraps)

	return nil
}
