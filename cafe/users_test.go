package cafe

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/cafe/models"
	util "github.com/textileio/textile-go/util/testing"
	"testing"
)

var userRefCode string
var userRegistration = map[string]interface{}{
	"username": ksuid.New().String(),
	"password": ksuid.New().String(),
	"identity": map[string]string{
		"type":  "email_address",
		"value": fmt.Sprintf("%s@textile.io", ksuid.New().String()),
	},
	"ref_code": "canihaz?",
}
var credentials = map[string]interface{}{
	"username": userRegistration["username"],
	"password": userRegistration["password"],
}

func TestUsers_Setup(t *testing.T) {
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
		userRefCode = resp.RefCodes[0]
	} else {
		t.Error("got bad ref codes")
	}
}

func TestUsers_SignUp(t *testing.T) {
	res, err := util.SignUpUser(userRegistration)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 404 {
		t.Errorf("bad status from sign up with bad ref code: %d", res.StatusCode)
		return
	}
	userRegistration["ref_code"] = userRefCode
	res2, err := util.SignUpUser(userRegistration)
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
	if resp2.Session == nil {
		t.Error("signup response missing session")
		return
	}
	res3, err := util.SignUpUser(userRegistration)
	if err != nil {
		t.Error(err)
		return
	}
	defer res3.Body.Close()
	if res3.StatusCode != 404 {
		t.Errorf("bad status from sign up with already used ref code: %d", res3.StatusCode)
		return
	}
}

func TestUsers_SignIn(t *testing.T) {
	res, err := util.SignInUser(credentials)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Errorf("got bad status: %d", res.StatusCode)
		return
	}
	resp := &models.SessionResponse{}
	if err := util.UnmarshalJSON(res.Body, resp); err != nil {
		t.Error(err)
		return
	}
	if resp.Session == nil {
		t.Error("got bad session")
		return
	}
	credentials["password"] = "doh!"
	res2, err := util.SignInUser(credentials)
	if err != nil {
		t.Error(err)
		return
	}
	defer res2.Body.Close()
	if res2.StatusCode != 403 {
		t.Errorf("got bad status: %d", res2.StatusCode)
		return
	}
	credentials["username"] = "bart"
	res3, err := util.SignInUser(credentials)
	if err != nil {
		t.Error(err)
		return
	}
	defer res3.Body.Close()
	if res3.StatusCode != 404 {
		t.Errorf("got bad status: %d", res3.StatusCode)
		return
	}
}
