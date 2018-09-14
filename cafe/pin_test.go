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
	res, err := util.CreateReferral(util.CafeReferralKey, 1, 1, "test")
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		t.Errorf("could not create referral, bad status: %d", res.StatusCode)
		return
	}
	resp := &models.ReferralResponse{}
	if err := util.UnmarshalJSON(res.Body, resp); err != nil {
		t.Error(err)
		return
	}
	if len(resp.RefCodes) > 0 {
		pRefCode = resp.RefCodes[0]
	} else {
		t.Error("got bad ref codes")
		return
	}
	pRegistration["ref_code"] = pRefCode
	res2, err := util.SignUpUser(pRegistration)
	if err != nil {
		t.Error(err)
		return
	}
	defer res2.Body.Close()
	if res2.StatusCode != 201 {
		t.Errorf("got bad status: %d", res2.StatusCode)
		return
	}
	resp2 := &models.SessionResponse{}
	if err := util.UnmarshalJSON(res2.Body, resp2); err != nil {
		t.Error(err)
		return
	}
	pSession = resp2.Session
}

func TestPin_Pin(t *testing.T) {
	block, err := os.Open("testdata/" + blockHash)
	if err != nil {
		t.Error(err)
		return
	}
	defer block.Close()
	res, err := util.Pin(block, pSession.AccessToken, "application/octet-stream")
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		t.Errorf("got bad status: %d", res.StatusCode)
		return
	}
	resp := &models.PinResponse{}
	if err := util.UnmarshalJSON(res.Body, resp); err != nil {
		t.Error(err)
		return
	}
	if resp.Id == nil {
		t.Error("response should contain id")
		return
	}
	if *resp.Id != blockHash {
		t.Errorf("hashes do not match: %s, %s", *resp.Id, blockHash)
	}
}

func TestPin_PinArchive(t *testing.T) {
	archive, err := os.Open("testdata/" + photoHash + ".tar.gz")
	if err != nil {
		t.Error(err)
		return
	}
	defer archive.Close()
	res, err := util.Pin(archive, pSession.AccessToken, "application/gzip")
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		t.Errorf("got bad status: %d", res.StatusCode)
		return
	}
	resp := &models.PinResponse{}
	if err := util.UnmarshalJSON(res.Body, resp); err != nil {
		t.Error(err)
		return
	}
	if resp.Id == nil {
		t.Error("response should contain id")
		return
	}
	if *resp.Id != photoHash {
		t.Errorf("hashes do not match: %s, %s", *resp.Id, photoHash)
	}
}
