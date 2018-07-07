package wallet_test

import (
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"os"
	"testing"

	"github.com/segmentio/ksuid"
	cmodels "github.com/textileio/textile-go/central/models"
	util "github.com/textileio/textile-go/util/testing"
	. "github.com/textileio/textile-go/wallet"
)

var repo = "testdata/.textile"

var wallet *Wallet
var addedId string

var centralReg = &cmodels.Registration{
	Username: ksuid.New().String(),
	Password: ksuid.New().String(),
	Identity: &cmodels.Identity{
		Type:  cmodels.EmailAddress,
		Value: ksuid.New().String() + "@textile.io",
	},
	Referral: "",
}

func TestNewWallet(t *testing.T) {
	os.RemoveAll(repo)
	config := Config{
		RepoPath:   repo,
		CentralAPI: util.CentralApiURL,
	}
	var err error
	wallet, _, err = NewWallet(config)
	if err != nil {
		t.Errorf("create wallet failed: %s", err)
	}
}

func TestWallet_StartWallet(t *testing.T) {
	online, err := wallet.Start()
	if err != nil {
		t.Errorf("start wallet failed: %s", err)
	}
	<-online
}

func TestWallet_Started(t *testing.T) {
	if !wallet.Started() {
		t.Errorf("should report started")
	}
}

func TestWallet_Online(t *testing.T) {
	if !wallet.Online() {
		t.Errorf("should report online")
	}
}

func TestWallet_GetGatewayAddress(t *testing.T) {
	// TODO
}

func TestWallet_GetRepoPath(t *testing.T) {
	// TODO
}

func TestWallet_SignUp(t *testing.T) {
	_, ref, err := util.CreateReferral(util.RefKey, 1, 1, "TestWallet_SignUp")
	if err != nil {
		t.Errorf("create referral for signup failed: %s", err)
		return
	}
	if len(ref.RefCodes) == 0 {
		t.Error("create referral for signup got no codes")
		return
	}
	centralReg.Referral = ref.RefCodes[0]

	err = wallet.SignUp(centralReg)
	if err != nil {
		t.Errorf("signup failed: %s", err)
		return
	}
}

func TestWallet_SignIn(t *testing.T) {
	creds := &cmodels.Credentials{
		Username: centralReg.Username,
		Password: centralReg.Password,
	}
	err := wallet.SignIn(creds)
	if err != nil {
		t.Errorf("signin failed: %s", err)
		return
	}
}

func TestWallet_IsSignedIn(t *testing.T) {
	// TODO
}

func TestWallet_GetUsername(t *testing.T) {
	// TODO
}

func TestWallet_GetId(t *testing.T) {
	// TODO
}

func TestWallet_GetPrivKey(t *testing.T) {
	// TODO
}

func TestWallet_GetPubKey(t *testing.T) {
	// TODO
}

func TestWallet_GetPubKeyString(t *testing.T) {
	key, err := wallet.GetPubKeyString()
	if err != nil {
		t.Errorf("get pub key failed: %s", err)
		return
	}
	if key == "" {
		t.Error("pub key empty")
		return
	}
}

func TestWallet_GetAccessToken(t *testing.T) {
	// TODO
}

func TestWallet_GetCentralAPI(t *testing.T) {
	// TODO
}

func TestWallet_GetCentralUserAPI(t *testing.T) {
	// TODO
}

func TestWallet_Threads(t *testing.T) {
	// TODO
}

func TestWallet_GetThread(t *testing.T) {
	// TODO
}

func TestWallet_GetThreadByName(t *testing.T) {
	// TODO
}

func TestWallet_AddThread(t *testing.T) {
	sk, _, err := libp2pc.GenerateKeyPair(libp2pc.Ed25519, 0)
	if err != nil {
		t.Error(err)
	}
	thrd, err := wallet.AddThread("test", sk)
	if err != nil {
		t.Errorf("add thread failed: %s", err)
		return
	}
	if thrd == nil {
		t.Error("add thread didn't return thread")
	}
}

func TestWallet_AddThreadWithMnemonic(t *testing.T) {
	// TODO
}

func TestWallet_RemoveThread(t *testing.T) {
	// TODO
}

func TestWallet_PublishThreads(t *testing.T) {
	// TODO
}

func TestWallet_Devices(t *testing.T) {
	// TODO
}

func TestWallet_AddDevice(t *testing.T) {
	// TODO
}

func TestWallet_RemoveDevice(t *testing.T) {
	// TODO
}

func TestWallet_AddPhoto(t *testing.T) {
	added, err := wallet.AddPhoto("testdata/image.jpg")
	if err != nil {
		t.Errorf("add photo failed: %s", err)
		return
	}
	if len(added.Id) == 0 {
		t.Errorf("add photo got bad id")
	}
	addedId = added.Id
}

func TestWallet_GetBlock(t *testing.T) {
	// TODO
}

func TestWallet_GetBlockByTarget(t *testing.T) {
	// TODO
}

func TestWallet_GetDataAtPath(t *testing.T) {
	// TODO
}

func TestWallet_ConnectPeer(t *testing.T) {
	// TODO
}

func TestWallet_PingPeer(t *testing.T) {
	// TODO
}

func TestWallet_Peers(t *testing.T) {
	// TODO
}

func TestWallet_SignOut(t *testing.T) {
	err := wallet.SignOut()
	if err != nil {
		t.Errorf("signout failed: %s", err)
		return
	}
}

func TestWallet_Stop(t *testing.T) {
	err := wallet.Stop()
	if err != nil {
		t.Errorf("stop wallet failed: %s", err)
	}
}

func TestWallet_StartedAgain(t *testing.T) {
	if wallet.Started() {
		t.Errorf("should report stopped")
	}
}

func TestWallet_OnlineAgain(t *testing.T) {
	if wallet.Online() {
		t.Errorf("should report offline")
	}
}

// test signin in stopped state, should re-connect to db
func TestWallet_SignInAgain(t *testing.T) {
	creds := &cmodels.Credentials{
		Username: centralReg.Username,
		Password: centralReg.Password,
	}
	err := wallet.SignIn(creds)
	if err != nil {
		t.Errorf("signin failed: %s", err)
		return
	}
}

func Test_Teardown(t *testing.T) {
	os.RemoveAll(wallet.GetRepoPath())
}
