package cafe

import (
	"fmt"
	"github.com/segmentio/ksuid"
	util "github.com/textileio/textile-go/util/testing"
	"testing"
)

var refCode string
var registration = map[string]interface{}{
	"username": ksuid.New().String(),
	"password": ksuid.New().String(),
	"identity": map[string]string{
		"type":  "email_address",
		"value": fmt.Sprintf("%s@textile.io", ksuid.New().String()),
	},
	"ref_code": "canihaz?",
}
var credentials = map[string]interface{}{
	"username": registration["username"],
	"password": registration["password"],
}

func TestUsers_Setup(t *testing.T) {
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
		refCode = ref.RefCodes[0]
	} else {
		t.Error("got bad ref codes")
	}
}

func TestUsers_SignUp(t *testing.T) {
	stat, _, err := util.SignUp(registration)
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 404 {
		t.Errorf("bad status from sign up with bad ref code: %d", stat)
		return
	}

	registration["ref_code"] = refCode
	stat2, res2, err := util.SignUp(registration)
	if err != nil {
		t.Error(err)
		return
	}
	if stat2 != 201 {
		t.Errorf("got bad status: %d", stat)
		return
	}
	if res2.Session == nil {
		t.Error("got bad session")
		return
	}

	stat3, _, err := util.SignUp(registration)
	if err != nil {
		t.Error(err)
		return
	}
	if stat3 != 404 {
		t.Errorf("bad status from sign up with already used ref code: %d", stat)
		return
	}
}

func TestUsers_SignIn(t *testing.T) {
	stat, res, err := util.SignIn(credentials)
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 200 {
		t.Errorf("got bad status: %d", stat)
		return
	}
	if res.Session == nil {
		t.Error("got bad session")
		return
	}
	credentials["password"] = "doh!"
	stat1, _, err := util.SignIn(credentials)
	if err != nil {
		t.Error(err)
		return
	}
	if stat1 != 403 {
		t.Errorf("got bad status: %d", stat1)
		return
	}
	credentials["username"] = "bart"
	stat2, _, err := util.SignIn(credentials)
	if err != nil {
		t.Error(err)
		return
	}
	if stat2 != 404 {
		t.Errorf("got bad status: %d", stat2)
		return
	}
}
