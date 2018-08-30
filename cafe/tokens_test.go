package cafe

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/cafe/auth"
	"github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/net/service"
	util "github.com/textileio/textile-go/util/testing"
	"testing"
	"time"
)

var session *models.Session
var claims *auth.TextileClaims

func TestTokens_Setup(t *testing.T) {
	// create a referral for the test
	var code string
	ref, err := util.CreateReferral(util.CafeReferralKey, 1, 1, "test")
	if err != nil {
		t.Error(err)
		return
	}
	if ref.Status != 201 {
		t.Errorf("could not create referral, bad status: %d", ref.Status)
		return
	}
	if len(ref.RefCodes) > 0 {
		code = ref.RefCodes[0]
	} else {
		t.Error("got bad ref codes")
	}
	var reg = map[string]interface{}{
		"username": ksuid.New().String(),
		"password": ksuid.New().String(),
		"identity": map[string]string{
			"type":  "email_address",
			"value": fmt.Sprintf("%s@textile.io", ksuid.New().String()),
		},
		"ref_code": code,
	}
	stat, res, err := util.SignUp(reg)
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 201 {
		t.Errorf("got bad status: %d", stat)
		return
	}
	if res.Session == nil {
		t.Error("got bad session")
		return
	}
	session = res.Session
	token, _ := jwt.Parse(session.AccessToken, util.Verify)
	claims, err = auth.ParseClaims(token.Claims)
	if err != nil {
		t.Error(err)
	}
}

func TestTokens_Refresh(t *testing.T) {
	stat, res, err := util.Refresh(session)
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
}

func TestTokens_RefreshBadSignature(t *testing.T) {
	session, err := auth.NewSession("abc", "bad", claims.Issuer, service.TextileProtocol, time.Hour)
	if err != nil {
		t.Error(err)
		return
	}
	stat, _, err := util.Refresh(session)
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 403 {
		t.Errorf("got bad status: %d", stat)
	}
}

func TestTokens_RefreshBadAudience(t *testing.T) {
	session, err := auth.NewSession("abc", util.CafeTokenSecret, claims.Issuer, "trust_us", time.Hour)
	if err != nil {
		t.Error(err)
		return
	}
	stat, _, err := util.Refresh(session)
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 403 {
		t.Errorf("got bad status: %d", stat)
	}
}

func TestTokens_RefreshExpired(t *testing.T) {
	session, err := auth.NewSession("abc", util.CafeTokenSecret, claims.Issuer, service.TextileProtocol, 0)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Second)
	stat, _, err := util.Refresh(session)
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 401 {
		t.Errorf("got bad status: %d", stat)
	}
}
