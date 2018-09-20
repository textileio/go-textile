package wallet_test

import (
	"fmt"
	util "github.com/textileio/textile-go/util/testing"
	. "github.com/textileio/textile-go/wallet"
	"github.com/tyler-smith/go-bip32"
	"log"
	"os"
	"testing"
)

var repo = "testdata/.textile"

var wallet *Wallet

func TestNewWallet(t *testing.T) {
	os.RemoveAll(repo)
	config := Config{
		RepoPath: repo,
		CafeAddr: util.CafeAddr,
	}
	var err error
	wallet, err = NewWallet(config)
	if err != nil {
		t.Errorf("create wallet failed: %s", err)
	}

	seed, err := bip32.NewSeed()
	if err != nil {
		log.Fatalln("Error generating seed:", err)
	}

	// Create master private key from seed
	computerVoiceMasterKey, _ := bip32.NewMasterKey(seed)

	// Map departments to keys
	// There is a very small chance a given child index is invalid
	// If so your real program should handle this by skipping the index
	departmentKeys := map[string]*bip32.Key{}
	departmentKeys["Sales"], _ = computerVoiceMasterKey.NewChildKey(0)
	departmentKeys["Marketing"], _ = computerVoiceMasterKey.NewChildKey(1)
	departmentKeys["Engineering"], _ = computerVoiceMasterKey.NewChildKey(2)
	departmentKeys["Customer Support"], _ = computerVoiceMasterKey.NewChildKey(3)

	// Create public keys for record keeping, auditors, payroll, etc
	departmentAuditKeys := map[string]*bip32.Key{}
	departmentAuditKeys["Sales"] = departmentKeys["Sales"].PublicKey()
	departmentAuditKeys["Marketing"] = departmentKeys["Marketing"].PublicKey()
	departmentAuditKeys["Engineering"] = departmentKeys["Engineering"].PublicKey()
	departmentAuditKeys["Customer Support"] = departmentKeys["Customer Support"].PublicKey()

	// Print public keys
	for department, pubKey := range departmentAuditKeys {
		k, err := pubKey.NewChildKey(0)
		if err != nil {
			panic(err)
		}
		fmt.Println(department, k.B58Serialize())
	}
}

//func TestWallet_StartWallet(t *testing.T) {
//	if err := wallet.Start(); err != nil {
//		t.Errorf("start wallet failed: %s", err)
//	}
//	<-wallet.Online()
//}
//
//func TestWallet_Started(t *testing.T) {
//	if !wallet.Started() {
//		t.Errorf("should report started")
//	}
//}
//
//func TestWallet_IsOnline(t *testing.T) {
//	if !wallet.IsOnline() {
//		t.Errorf("should report online")
//	}
//}
//
//func TestWallet_CafeRegister(t *testing.T) {
//	res, err := util.CreateReferral(util.CafeReferralKey, 1, 1, "test")
//	if err != nil {
//		t.Errorf("create referral for registration failed: %s", err)
//		return
//	}
//	defer res.Body.Close()
//	resp := &models.ReferralResponse{}
//	if err := util.UnmarshalJSON(res.Body, resp); err != nil {
//		t.Error(err)
//		return
//	}
//	if len(resp.RefCodes) == 0 {
//		t.Error("create referral for registration got no codes")
//		return
//	}
//
//	if err := wallet.CafeRegister(resp.RefCodes[0]); err != nil {
//		t.Errorf("register failed: %s", err)
//		return
//	}
//}
//
//func TestWallet_CafeLogin(t *testing.T) {
//	if err := wallet.CafeLogin(); err != nil {
//		t.Errorf("login failed: %s", err)
//		return
//	}
//}
//
//func TestWallet_AddThread(t *testing.T) {
//	sk, _, err := libp2pc.GenerateEd25519Key(rand.Reader)
//	if err != nil {
//		t.Error(err)
//	}
//	thrd, err := wallet.AddThread("test", sk, true)
//	if err != nil {
//		t.Errorf("add thread failed: %s", err)
//		return
//	}
//	if thrd == nil {
//		t.Error("add thread didn't return thread")
//	}
//}
//
//func TestWallet_AddPhoto(t *testing.T) {
//	added, err := wallet.AddPhoto("../util/testdata/image.jpg")
//	if err != nil {
//		t.Errorf("add photo failed: %s", err)
//		return
//	}
//	if len(added.Id) == 0 {
//		t.Errorf("add photo got bad id")
//	}
//	// test adding an image w/o the orientation tag
//	added2, err := wallet.AddPhoto("../util/testdata/image-no-orientation.jpg")
//	if err != nil {
//		t.Errorf("add photo w/o orientation tag failed: %s", err)
//		return
//	}
//	if len(added2.Id) == 0 {
//		t.Errorf("add photo w/o orientation tag got bad id")
//	}
//}
//
//func TestWallet_CafeLogout(t *testing.T) {
//	err := wallet.CafeLogout()
//	if err != nil {
//		t.Errorf("logout failed: %s", err)
//		return
//	}
//}
//
//func TestWallet_Stop(t *testing.T) {
//	err := wallet.Stop()
//	if err != nil {
//		t.Errorf("stop wallet failed: %s", err)
//	}
//}
//
//func TestWallet_StartedAgain(t *testing.T) {
//	if wallet.Started() {
//		t.Errorf("should report stopped")
//	}
//}
//
//func TestWallet_OnlineAgain(t *testing.T) {
//	if wallet.IsOnline() {
//		t.Errorf("should report offline")
//	}
//}
//
//// test cafe login in stopped state, should re-connect to db
//func TestWallet_LoginAgain(t *testing.T) {
//	if err := wallet.CafeLogin(); err != nil {
//		t.Errorf("login from stopped failed: %s", err)
//		return
//	}
//}
//
func Test_Teardown(t *testing.T) {
	os.RemoveAll(wallet.GetRepoPath())
}
