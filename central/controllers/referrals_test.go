package controllers_test

import (
	"testing"

	"github.com/textileio/textile-go/test"
)

func TestReferrals_CreateReferral(t *testing.T) {
	num := 10
	stat, res, err := test.CreateReferral(test.RefKey, num)
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
	stat, res, err := test.ListReferrals(test.RefKey)
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
