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
	"github.com/textileio/go-textile/cluster"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/keypair"
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

	secret, err := cluster.NewClusterSecret()
	if err != nil {
		t.Fatal(err)
	}

	err = core.InitRepo(core.InitConfig{
		Account:       accnt1,
		RepoPath:      vars.repoPath1,
		ApiAddr:       fmt.Sprintf("127.0.0.1:%s", core.GetRandomPort()),
		SwarmPorts:    swarmPort1,
		ClusterSecret: secret,
		Debug:         true,
	})
	if err != nil {
		t.Fatalf("init node1 failed: %s", err)
	}
	err = core.InitRepo(core.InitConfig{
		Account:       accnt2,
		RepoPath:      vars.repoPath2,
		ApiAddr:       fmt.Sprintf("127.0.0.1:%s", core.GetRandomPort()),
		SwarmPorts:    swarmPort2,
		ClusterSecret: secret,
		Debug:         true,
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

func TestNewTextileCluster(t *testing.T) {
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

	<-vars.node1.OnlineCh()
	<-vars.node2.OnlineCh()

	cid, err := pinTestData(vars.node1.Ipfs())
	if err != nil {
		t.Fatal(err)
	}
	vars.cid = *cid
}

func TestTextileClusterSync(t *testing.T) {
	ctx, cancel := context.WithTimeout(vars.node1.Ipfs().Context(), time.Minute)
	defer cancel()

	info, err := vars.node1.Cluster().SyncAll(ctx)
	if err != nil {
		t.Fatalf("sync all failed: %s", err)
	}

	var foundCid bool
	for _, i := range info {
		fmt.Println(i.String())
		if i.Cid.Equals(vars.cid) {
			foundCid = true
		}
	}

	if !foundCid {
		t.Fatalf("failed to find cid in cluster: %s", vars.cid.String())
	}
}

func TestTextileCluster_Teardown(t *testing.T) {
	vars.node1 = nil
	vars.node2 = nil
}

func pinTestData(node *icore.IpfsNode) (*icid.Cid, error) {
	f, err := os.Open("../mill/testdata/image.jpeg")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ipfs.AddData(node, f, true, false)
}
