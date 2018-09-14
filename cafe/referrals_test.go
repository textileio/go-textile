package cafe

import (
	"github.com/textileio/textile-go/cafe/models"
	util "github.com/textileio/textile-go/util/testing"
	"testing"
)

func TestReferrals_CreateReferral(t *testing.T) {
	num := 10
	res, err := util.CreateReferral(util.CafeReferralKey, num, 2, "test")
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		t.Errorf("got bad status: %d", res.StatusCode)
		return
	}
	resp := &models.ReferralResponse{}
	if err := util.UnmarshalJSON(res.Body, resp); err != nil {
		t.Error(err)
		return
	}
	if len(resp.RefCodes) != num {
		t.Error("got bad ref codes")
		return
	}
	res2, err := util.CreateReferral("canihaz?", 1, 1, "test")
	if err != nil {
		t.Error(err)
		return
	}
	defer res2.Body.Close()
	if res2.StatusCode != 403 {
		t.Errorf("got bad status: %d", res.StatusCode)
		return
	}
}

func TestReferrals_ListReferrals(t *testing.T) {
	res, err := util.ListReferrals(util.CafeReferralKey)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Errorf("got bad status: %d", res.StatusCode)
		return
	}
	resp := &models.ReferralResponse{}
	if err := util.UnmarshalJSON(res.Body, resp); err != nil {
		t.Error(err)
		return
	}
	if len(resp.RefCodes) == 0 {
		t.Error("got bad ref codes")
		return
	}
	res2, err := util.ListReferrals("canihaz?")
	if err != nil {
		t.Error(err)
		return
	}
	defer res2.Body.Close()
	if res2.StatusCode != 403 {
		t.Errorf("got bad status: %d", res.StatusCode)
		return
	}
}
