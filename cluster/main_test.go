package cluster_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	icid "github.com/ipfs/go-cid"
	icore "github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo/config"
)

var vars = struct {
	repoPath1 string
	repoPath2 string

	node1 *core.Textile
	node2 *core.Textile

	cid icid.Cid
}{
	repoPath1: "testdata/.cluster1",
	repoPath2: "testdata/.cluster2",

	cid: icid.Undef,
}

func TestInitCluster(t *testing.T) {
	_ = os.RemoveAll(vars.repoPath1)
	_ = os.RemoveAll(vars.repoPath2)

	accnt1 := keypair.Random()
	accnt2 := keypair.Random()

	swarmPort1 := core.GetRandomPort()
	swarmPort2 := core.GetRandomPort()

	err := core.InitRepo(core.InitConfig{
		Account:              accnt1,
		RepoPath:             vars.repoPath1,
		ApiAddr:              fmt.Sprintf("127.0.0.1:%s", core.GetRandomPort()),
		SwarmPorts:           swarmPort1,
		Cluster:              true,
		ClusterBindMultiaddr: "/ip4/0.0.0.0/tcp/9096",
		Debug:                true,
	})
	if err != nil {
		t.Fatalf("init node1 failed: %s", err)
	}
	err = core.InitRepo(core.InitConfig{
		Account:              accnt2,
		RepoPath:             vars.repoPath2,
		ApiAddr:              fmt.Sprintf("127.0.0.1:%s", core.GetRandomPort()),
		SwarmPorts:           swarmPort2,
		Cluster:              true,
		ClusterBindMultiaddr: "/ip4/0.0.0.0/tcp/9097",
		Debug:                true,
	})
	if err != nil {
		t.Fatalf("init node2 failed: %s", err)
	}

	// update bootstraps
	addr1, err := getPeerAddress(vars.repoPath1, swarmPort1)
	if err != nil {
		t.Fatal(err)
	}
	addr2, err := getPeerAddress(vars.repoPath2, swarmPort2)
	if err != nil {
		t.Fatal(err)
	}
	err = updateClusterBootstraps(vars.repoPath1, []string{addr2})
	if err != nil {
		t.Fatal(err)
	}
	err = updateClusterBootstraps(vars.repoPath2, []string{addr1})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartCluster(t *testing.T) {
	var err error
	vars.node1, err = core.NewTextile(core.RunConfig{
		RepoPath: vars.repoPath1,
		Debug:    true,
	})
	if err != nil {
		t.Fatalf("create node1 failed: %s", err)
	}
	vars.node2, err = core.NewTextile(core.RunConfig{
		RepoPath: vars.repoPath2,
		Debug:    true,
	})
	if err != nil {
		t.Fatalf("create node2 failed: %s", err)
	}

	// set cluster logs to debug
	level := &pb.LogLevel{
		Systems: map[string]pb.LogLevel_Level{
			"cluster": pb.LogLevel_DEBUG,
		},
	}
	err = vars.node1.SetLogLevel(level)
	if err != nil {
		t.Fatal(err)
	}
	err = vars.node2.SetLogLevel(level)
	if err != nil {
		t.Fatal(err)
	}

	// start nodes
	err = vars.node1.Start()
	if err != nil {
		t.Fatalf("start node1 failed: %s", err)
	}
	<-vars.node1.OnlineCh()
	<-vars.node1.Cluster().Ready()

	// let node1 warm up
	timer := time.NewTimer(time.Second * 5)
	<-timer.C

	err = vars.node2.Start()
	if err != nil {
		t.Fatalf("start node2 failed: %s", err)
	}
	<-vars.node2.OnlineCh()
	<-vars.node2.Cluster().Ready()

	// let node2 warm up
	timer = time.NewTimer(time.Second * 5)
	<-timer.C

	// pin some data to node1
	cid, err := pinTestData(vars.node1.Ipfs())
	if err != nil {
		t.Fatal(err)
	}
	vars.cid = *cid
}

func TestTextileClusterPeers(t *testing.T) {
	ctx, cancel := context.WithTimeout(vars.node1.Ipfs().Context(), time.Minute)
	defer cancel()

	var ok bool
	for _, p := range vars.node1.Cluster().Peers(ctx) {
		if p.ID.Pretty() == vars.node2.Ipfs().Identity.Pretty() {
			ok = true
			break
		}
	}
	if !ok {
		t.Fatal("node2 not found in node1's peers")
	}
	ok = false
	for _, p := range vars.node2.Cluster().Peers(ctx) {
		if p.ID.Pretty() == vars.node1.Ipfs().Identity.Pretty() {
			ok = true
			break
		}
	}
	if !ok {
		t.Fatal("node1 not found in node2's peers")
	}
}

func TestTextileClusterSync(t *testing.T) {
	ctx, cancel := context.WithTimeout(vars.node1.Ipfs().Context(), time.Minute)
	defer cancel()

	_, err := vars.node1.Cluster().SyncAll(ctx)
	if err != nil {
		t.Fatalf("sync all failed: %s", err)
	}

	err = vars.node1.Cluster().StateSync(ctx)
	if err != nil {
		t.Fatalf("state sync failed: %s", err)
	}

	info, err := vars.node1.Cluster().Status(ctx, vars.cid)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info.String())
}

func TestTextileCluster_Stop(t *testing.T) {
	err := vars.node1.Stop()
	if err != nil {
		t.Fatalf("stop node1 failed: %s", err)
	}
	err = vars.node2.Stop()
	if err != nil {
		t.Fatalf("stop node2 failed: %s", err)
	}
}

func TestTextileCluster_Teardown(t *testing.T) {
	vars.node1 = nil
	vars.node2 = nil
}

func getPeerAddress(repoPath, swarmPort string) (string, error) {
	r, err := fsrepo.Open(repoPath)
	if err != nil {
		return "", err
	}
	defer r.Close()
	id, err := r.GetConfigKey("Identity.PeerID")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/ip4/127.0.0.1/tcp/%s/ipfs/%s", swarmPort, id), nil
}

func updateClusterBootstraps(repoPath string, bootstraps []string) error {
	conf, err := config.Read(repoPath)
	if err != nil {
		return err
	}
	conf.Cluster.Bootstraps = bootstraps
	return config.Write(repoPath, conf)
}

func pinTestData(node *icore.IpfsNode) (*icid.Cid, error) {
	f, err := os.Open("../mill/testdata/image.jpeg")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ipfs.AddData(node, f, true, false)
}
