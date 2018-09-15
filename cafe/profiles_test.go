package cafe

import (
	"crypto/rand"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/util"
	tutil "github.com/textileio/textile-go/util/testing"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"testing"
)

var profileKey libp2pc.PrivKey
var challengeRequest = map[string]string{
	"pk": "sneakypk",
}
var challengeResponse *models.ChallengeResponse
var profileRefCode string
var profileRegistration = map[string]interface{}{
	"challenge": map[string]string{
		"pk":        "sneakypk",
		"value":     "invalid",
		"nonce":     "invalid",
		"signature": "invalid",
	},
	"ref_code": "canihaz?",
}

func TestProfiles_Setup(t *testing.T) {
	// create a referral for the test
	res, err := tutil.CreateReferral(tutil.CafeReferralKey, 1, 1, "test")
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
	if err := tutil.UnmarshalJSON(res.Body, resp); err != nil {
		t.Error(err)
		return
	}
	if len(resp.RefCodes) > 0 {
		profileRefCode = resp.RefCodes[0]
	} else {
		t.Error("got bad ref codes")
	}
}

func TestProfiles_Challenge(t *testing.T) {
	res, err := tutil.ProfileChallenge(challengeRequest)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 400 {
		t.Errorf("bad status from profile challenge with bad pk: %d", res.StatusCode)
		return
	}

	// make a key pair
	sk, pk, err := libp2pc.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Error(err)
	}
	profileKey = sk
	pks, err := util.EncodeKey(pk)
	if err != nil {
		t.Error(err)
	}
	challengeRequest["pk"] = pks
	res2, err := tutil.ProfileChallenge(challengeRequest)
	if err != nil {
		t.Error(err)
		return
	}
	defer res2.Body.Close()
	if res2.StatusCode != 201 {
		t.Errorf("bad status from profile challenge: %d", res2.StatusCode)
		return
	}
	resp2 := &models.ChallengeResponse{}
	if err := tutil.UnmarshalJSON(res2.Body, resp2); err != nil {
		t.Error(err)
		return
	}
	if resp2.Value == nil {
		t.Error("get challenge did not return a value")
		return
	}
	challengeResponse = resp2
}

func TestProfiles_Register(t *testing.T) {
	res, err := tutil.RegisterProfile(profileRegistration)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 404 {
		t.Errorf("bad status from registration with bad ref code: %d", res.StatusCode)
		return
	}
	profileRegistration["ref_code"] = profileRefCode
	res2, err := tutil.RegisterProfile(profileRegistration)
	if err != nil {
		t.Error(err)
		return
	}
	defer res2.Body.Close()
	if res2.StatusCode != 403 {
		t.Errorf("bad status from registration with bad nonce: %d", res2.StatusCode)
		return
	}
	var snonce string
	if challengeResponse.Value != nil {
		snonce = *challengeResponse.Value
	}
	profileRegistration["challenge"] = map[string]string{
		"pk":        "sneakypk",
		"value":     snonce,
		"nonce":     "invalid",
		"signature": "invalid",
	}
	res3, err := tutil.RegisterProfile(profileRegistration)
	if err != nil {
		t.Error(err)
		return
	}
	defer res3.Body.Close()
	if res3.StatusCode != 403 {
		t.Errorf("bad status from registration with bad pk: %d", res3.StatusCode)
		return
	}
	cnonce := ksuid.New().String()
	sigb, err := profileKey.Sign([]byte(snonce + cnonce))
	if err != nil {
		t.Error(err)
		return
	}
	profileRegistration["challenge"] = map[string]string{
		"pk":        challengeRequest["pk"],
		"value":     snonce,
		"nonce":     cnonce,
		"signature": libp2pc.ConfigEncodeKey(sigb),
	}
	res4, err := tutil.RegisterProfile(profileRegistration)
	if err != nil {
		t.Error(err)
		return
	}
	defer res4.Body.Close()
	if res4.StatusCode != 201 {
		t.Errorf("bad status from good registration: %d", res4.StatusCode)
		return
	}
	resp4 := &models.SessionResponse{}
	if err := tutil.UnmarshalJSON(res4.Body, resp4); err != nil {
		t.Error(err)
		return
	}
	if resp4.Session == nil {
		t.Error("registration response missing session")
		return
	}
	res5, err := tutil.RegisterProfile(profileRegistration)
	if err != nil {
		t.Error(err)
		return
	}
	defer res5.Body.Close()
	if res5.StatusCode != 404 {
		t.Errorf("bad status from registration with already used ref code: %d", res5.StatusCode)
		return
	}
}

func TestProfiles_Login(t *testing.T) {
	res, err := tutil.ProfileChallenge(challengeRequest)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		t.Errorf("bad status from profile challenge: %d", res.StatusCode)
		return
	}
	resp := &models.ChallengeResponse{}
	if err := tutil.UnmarshalJSON(res.Body, resp); err != nil {
		t.Error(err)
		return
	}
	if resp.Value == nil {
		t.Error("get challenge did not return a value")
		return
	}
	res2, err := tutil.LoginProfile(map[string]string{
		"pk":        "sneakypk",
		"value":     "invalid",
		"nonce":     "invalid",
		"signature": "invalid",
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer res2.Body.Close()
	if res2.StatusCode != 404 {
		t.Errorf("bad status from login with bad pk: %d", res2.StatusCode)
		return
	}
	res3, err := tutil.LoginProfile(map[string]string{
		"pk":        challengeRequest["pk"],
		"value":     "invalid",
		"nonce":     "invalid",
		"signature": "invalid",
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer res3.Body.Close()
	if res3.StatusCode != 403 {
		t.Errorf("bad status from loign with bad nonce: %d", res3.StatusCode)
		return
	}
	var snonce string
	if resp.Value != nil {
		snonce = *resp.Value
	}
	cnonce := ksuid.New().String()
	badsigb, err := profileKey.Sign([]byte(ksuid.New().String() + cnonce))
	if err != nil {
		t.Error(err)
		return
	}
	res4, err := tutil.LoginProfile(map[string]string{
		"pk":        challengeRequest["pk"],
		"value":     snonce,
		"nonce":     cnonce,
		"signature": libp2pc.ConfigEncodeKey(badsigb),
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer res4.Body.Close()
	if res4.StatusCode != 403 {
		t.Errorf("bad status from login with bad sig: %d", res4.StatusCode)
		return
	}
	sigb, err := profileKey.Sign([]byte(snonce + cnonce))
	if err != nil {
		t.Error(err)
		return
	}
	signed := map[string]string{
		"pk":        challengeRequest["pk"],
		"value":     snonce,
		"nonce":     cnonce,
		"signature": libp2pc.ConfigEncodeKey(sigb),
	}
	res5, err := tutil.LoginProfile(signed)
	if err != nil {
		t.Error(err)
		return
	}
	defer res5.Body.Close()
	if res5.StatusCode != 200 {
		t.Errorf("bad status from good login: %d", res5.StatusCode)
		return
	}
	resp5 := &models.SessionResponse{}
	if err := tutil.UnmarshalJSON(res5.Body, resp5); err != nil {
		t.Error(err)
		return
	}
	if resp5.Session == nil {
		t.Error("login response missing session")
		return
	}
	res6, err := tutil.LoginProfile(signed)
	if err != nil {
		t.Error(err)
		return
	}
	defer res6.Body.Close()
	if res6.StatusCode != 403 {
		t.Errorf("bad status from login with already used nonce: %d", res6.StatusCode)
		return
	}
}
