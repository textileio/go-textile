package cafe

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/cafe/models"
	util "github.com/textileio/textile-go/util/testing"
	"os"
	"testing"
)

var pRefCode string
var pRegistration = map[string]interface{}{
	"username": ksuid.New().String(),
	"password": ksuid.New().String(),
	"identity": map[string]string{
		"type":  "email_address",
		"value": fmt.Sprintf("%s@textile.io", ksuid.New().String()),
	},
	"ref_code": "canihaz?",
}
var pSession *models.Session
var blockHash = "QmbQ4K3vXNJ3DjCNdG2urCXs7BuHqWQG1iSjZ8fbnF8NMs"
var photoHash = "QmSUnsZi9rGvPZLWy2v5N7fNxUWVNnA5nmppoM96FbLqLp"

func TestPin_Setup(t *testing.T) {
	// create a referral for the test
	stat, ref, err := util.CreateReferral(util.CafeReferralKey, 1, 1, "test")
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 201 {
		t.Errorf("could not create referral, bad status: %d", stat)
		return
	}
	if len(ref.RefCodes) > 0 {
		pRefCode = ref.RefCodes[0]
	} else {
		t.Error("got bad ref codes")
		return
	}
	pRegistration["ref_code"] = pRefCode
	stat, res, err := util.SignUp(pRegistration)
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 201 {
		t.Errorf("got bad status: %d", stat)
		return
	}
	pSession = res.Session
}

func TestPin_Pin(t *testing.T) {
	block, err := os.Open("testdata/" + blockHash)
	if err != nil {
		t.Error(err)
		return
	}
	defer block.Close()
	stat, res, err := util.Pin(block, pSession.AccessToken, "application/octet-stream")
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 201 {
		t.Errorf("got bad status: %d", stat)
		return
	}
	if res.Id == nil {
		t.Error("response should contain id")
		return
	}
	if *res.Id != blockHash {
		t.Errorf("hashes do not match: %s, %s", *res.Id, blockHash)
	}
}

func TestPin_PinArchive(t *testing.T) {
	archive, err := os.Open("testdata/" + photoHash + ".tar.gz")
	if err != nil {
		t.Error(err)
		return
	}
	defer archive.Close()
	stat, res, err := util.Pin(archive, pSession.AccessToken, "application/gzip")
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 201 {
		t.Errorf("got bad status: %d", stat)
		return
	}
	if res.Id == nil {
		t.Error("response should contain id")
		return
	}
	if *res.Id != photoHash {
		t.Errorf("hashes do not match: %s, %s", *res.Id, photoHash)
	}
}
