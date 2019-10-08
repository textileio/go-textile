package mobile

import (
	"os"

	"github.com/textileio/go-textile/wallet"

	"github.com/gogo/protobuf/proto"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/pb"
)

func createAndStartPeer(conf InitConfig, wait bool, handler core.CafeOutboxHandler, messenger Messenger) (*Mobile, error) {
	w, err := wallet.WalletFromWordCount(12)
	if err != nil {
		return nil, err
	}

	res, err := WalletAccountAt(w.RecoveryPhrase, 0, "")
	if err != nil {
		return nil, err
	}

	accnt := new(pb.MobileWalletAccount)
	err = proto.Unmarshal(res, accnt)
	if err != nil {
		return nil, err
	}
	conf.Seed = accnt.Seed

	repoPath, err := conf.Repo()
	if err != nil {
		return nil, err
	}

	_ = os.RemoveAll(repoPath)

	err = InitRepo(&conf)
	if err != nil {
		return nil, err
	}

	node, err := NewTextile(&RunConfig{
		RepoPath:          repoPath,
		Debug:             conf.Debug,
		CafeOutboxHandler: handler,
	}, messenger)
	if err != nil {
		return nil, err
	}

	err = node.Start()
	if err != nil {
		return nil, err
	}

	if wait {
		<-node.onlineCh()
	}

	return node, nil
}

func addTestThread(node *Mobile, conf *pb.AddThreadConfig) (*pb.Thread, error) {
	mconf, err := proto.Marshal(conf)
	if err != nil {
		return nil, err
	}
	res, err := node.AddThread(mconf)
	if err != nil {
		return nil, err
	}
	thrd := new(pb.Thread)
	err = proto.Unmarshal(res, thrd)
	if err != nil {
		return nil, err
	}
	return thrd, nil
}
