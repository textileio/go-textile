package cafe

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/cafe/auth"
	"github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/net/service"
	"testing"
	"time"
)

var session *models.Session
var claims *auth.TextileClaims

func TestTokens_Setup(t *testing.T) {
	// create a referral for the test
	var code string
	res, err := createReferral(cafeReferralKey, 1, 1, "test")
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
	if err := unmarshalJSON(res.Body, resp); err != nil {
		t.Error(err)
		return
	}
	if len(resp.RefCodes) > 0 {
		code = resp.RefCodes[0]
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
	res2, err := signUpUser(reg)
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
	if err := unmarshalJSON(res2.Body, resp2); err != nil {
		t.Error(err)
		return
	}
	if resp2.Session == nil {
		t.Error("got bad session")
		return
	}
	session = resp2.Session
	token, err := jwt.Parse(session.AccessToken, verify)
	if err != nil {
		t.Error(err)
	}
	claims, err = auth.ParseClaims(token.Claims)
	if err != nil {
		t.Error(err)
	}
}

func TestTokens_Refresh(t *testing.T) {
	res, err := refreshSession(session.AccessToken, session.RefreshToken)
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
	if err := unmarshalJSON(res.Body, resp); err != nil {
		t.Error(err)
		return
	}
	if resp.Session == nil {
		t.Error("got bad session")
		return
	}
}

func TestTokens_RefreshBadSignature(t *testing.T) {
	session, err := auth.NewSession("abc", "bad", claims.Issuer, service.ThreadProtocol, time.Hour)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := refreshSession(session.AccessToken, session.RefreshToken)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 403 {
		t.Errorf("got bad status: %d", res.StatusCode)
	}
}

func TestTokens_RefreshBadAudience(t *testing.T) {
	session, err := auth.NewSession("abc", cafeTokenSecret, claims.Issuer, "trust_us", time.Hour)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := refreshSession(session.AccessToken, session.RefreshToken)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 403 {
		t.Errorf("got bad status: %d", res.StatusCode)
	}
}

func TestTokens_RefreshWrongToken(t *testing.T) {
	session, err := auth.NewSession("abc", cafeTokenSecret, claims.Issuer, service.ThreadProtocol, time.Hour)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := refreshSession(session.AccessToken, session.AccessToken)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 403 {
		t.Errorf("got bad status: %d", res.StatusCode)
	}
}

func TestTokens_RefreshExpired(t *testing.T) {
	session, err := auth.NewSession("abc", cafeTokenSecret, claims.Issuer, service.ThreadProtocol, 0)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Second)
	res, err := refreshSession(session.AccessToken, session.RefreshToken)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 401 {
		t.Errorf("got bad status: %d", res.StatusCode)
	}
}
