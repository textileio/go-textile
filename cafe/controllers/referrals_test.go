package controllers_test

import (
	util "github.com/textileio/textile-go/util/testing"
	"testing"
)

func TestReferrals_CreateReferral(t *testing.T) {
	num := 10
	stat, res, err := util.CreateReferral(util.RefKey, num, 2, "TestReferrals_CreateReferral")
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 201 {
		t.Errorf("got bad status: %d", stat)
		return
	}
	if len(res.RefCodes) != num {
		t.Error("got bad ref codes")
		return
	}
}

func TestReferrals_ListReferrals(t *testing.T) {
	stat, res, err := util.ListReferrals(util.RefKey)
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 200 {
		t.Errorf("got bad status: %d", stat)
		return
	}
	if len(res.RefCodes) == 0 {
		t.Error("got bad ref codes")
		return
	}
}
