package core

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/textileio/go-textile/repo/config"

	"github.com/textileio/go-textile/keypair"
)

var clusterVars = struct {
	repoPath1 string
	repoPath2 string
	node1     *Textile
	node2     *Textile
}{
	repoPath1: "testdata/.cluster1",
	repoPath2: "testdata/.cluster2",
}

func TestInitCluster(t *testing.T) {
	_ = os.RemoveAll(clusterVars.repoPath1)
	_ = os.RemoveAll(clusterVars.repoPath2)

	accnt1 := keypair.Random()
	accnt2 := keypair.Random()

	swarmPort1 := GetRandomPort()
	swarmPort2 := GetRandomPort()

	secret, err := NewClusterSecret()
	if err != nil {
		t.Fatal(err)
	}

	err = InitRepo(InitConfig{
		Account:       accnt1,
		RepoPath:      clusterVars.repoPath1,
		ApiAddr:       fmt.Sprintf("127.0.0.1:%s", GetRandomPort()),
		SwarmPorts:    swarmPort1,
		ClusterSecret: secret,
		Debug:         true,
	})
	if err != nil {
		t.Fatalf("init node1 failed: %s", err)
	}
	err = InitRepo(InitConfig{
		Account:       accnt2,
		RepoPath:      clusterVars.repoPath2,
		ApiAddr:       fmt.Sprintf("127.0.0.1:%s", GetRandomPort()),
		SwarmPorts:    swarmPort2,
		ClusterSecret: secret,
		Debug:         true,
	})
	if err != nil {
		t.Fatalf("init node2 failed: %s", err)
	}

	// update bootstraps
	addr1, err := getPeerAddress(clusterVars.repoPath1, swarmPort1)
	if err != nil {
		t.Fatal(err)
	}
	addr2, err := getPeerAddress(clusterVars.repoPath2, swarmPort2)
	if err != nil {
		t.Fatal(err)
	}
	err = updateClusterBootstraps(clusterVars.repoPath1, []string{addr2})
	if err != nil {
		t.Fatal(err)
	}
	err = updateClusterBootstraps(clusterVars.repoPath2, []string{addr1})
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
	clusterVars.node1, err = NewTextile(RunConfig{
		RepoPath: clusterVars.repoPath1,
		Debug:    true,
	})
	if err != nil {
		t.Fatalf("create node1 failed: %s", err)
	}
	clusterVars.node2, err = NewTextile(RunConfig{
		RepoPath: clusterVars.repoPath2,
		Debug:    true,
	})
	if err != nil {
		t.Fatalf("create node2 failed: %s", err)
	}

	<-clusterVars.node1.OnlineCh()
	<-clusterVars.node2.OnlineCh()

	err = addTestData(clusterVars.node1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTextileClusterSync(t *testing.T) {
	ctx, cancel := context.WithTimeout(clusterVars.node1.node.Context(), time.Minute)
	defer cancel()

	info, err := clusterVars.node1.cluster.SyncAll(ctx)
	if err != nil {
		t.Fatalf("sync failed: %s", err)
	}

	for _, i := range info {
		fmt.Println(i.String())
	}
}

func TestTextileCluster_Teardown(t *testing.T) {
	clusterVars.node1 = nil
	clusterVars.node2 = nil
}
