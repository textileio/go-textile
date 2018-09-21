package dao_test

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/cafe/dao"
	"github.com/textileio/textile-go/cafe/models"
	"os"
	"testing"
	"time"
)

var d = dao.DAO{}

var now = time.Now()
var unusedRefCnt int

var ref = models.Referral{
	ID:        bson.NewObjectId(),
	Code:      ksuid.New().String(),
	Created:   now,
	Remaining: 1,
}

var user = models.User{
	ID:       bson.NewObjectId(),
	Username: ksuid.New().String(),
	Password: ksuid.New().String(),
	Created:  now,
	LastSeen: now,
	Identities: []models.UserIdentity{
		{
			Type:     models.EmailAddress,
			Value:    fmt.Sprintf("%s@textile.io", ksuid.New().String()),
			Verified: false,
		},
	},
}

var profile = models.Profile{
	ID:       bson.NewObjectId(),
	Pk:       ksuid.New().String(),
	Created:  now,
	LastSeen: now,
}

var nonce = models.Nonce{
	ID:      bson.NewObjectId(),
	Pk:      ksuid.New().String(),
	Value:   ksuid.New().String(),
	Created: now,
}

func TestDao_Connect(t *testing.T) {
	d.Hosts = os.Getenv("CAFE_DB_HOSTS")
	d.Name = os.Getenv("CAFE_DB_NAME")
	d.Connect()
}

func TestDao_Index(t *testing.T) {
	d.Index()
}

func TestDAO_InsertUser(t *testing.T) {
	if err := d.InsertUser(user); err != nil {
		t.Errorf("insert user failed: %s", err)
		return
	}
}

func TestDAO_InsertUserAgain(t *testing.T) {
	var user2 = models.User{
		ID:       bson.NewObjectId(),
		Username: user.Username,
		Password: ksuid.New().String(),
		Created:  now,
		LastSeen: now,
		Identities: []models.UserIdentity{
			{
				Type:     models.EmailAddress,
				Value:    fmt.Sprintf("%s@textile.io", ksuid.New().String()),
				Verified: false,
			},
		},
	}
	if err := d.InsertUser(user2); err == nil {
		t.Error("username should be unique")
		return
	}
	var user3 = models.User{
		ID:         bson.NewObjectId(),
		Username:   ksuid.New().String(),
		Password:   ksuid.New().String(),
		Created:    now,
		LastSeen:   now,
		Identities: user.Identities,
	}
	if err := d.InsertUser(user3); err == nil {
		t.Error("identity should be unique")
		return
	}
}

func TestDAO_FindUserById(t *testing.T) {
	loaded, err := d.FindUserById(user.ID.Hex())
	if err != nil {
		t.Errorf("find user by id failed: %s", err)
		return
	}
	if loaded.Username != user.Username {
		t.Error("username mismatch")
	}
}

func TestDAO_FindUserByUsername(t *testing.T) {
	if _, err := d.FindUserByUsername(user.Username); err != nil {
		t.Errorf("find user by username failed: %s", err)
		return
	}
}

func TestDAO_FindUserByIdentity(t *testing.T) {
	if _, err := d.FindUserByIdentity(user.Identities[0]); err != nil {
		t.Errorf("find user by identity failed: %s", err)
		return
	}
}

func TestDAO_UpdateUser(t *testing.T) {
	un := ksuid.New().String()
	user.Username = un
	if err := d.UpdateUser(user); err != nil {
		t.Errorf("update user failed: %s", err)
		return
	}
	loaded, err := d.FindUserById(user.ID.Hex())
	if err != nil {
		t.Errorf("find user again by id failed: %s", err)
		return
	}
	if loaded.Username != user.Username {
		t.Error("username mismatch")
	}
}

func TestDAO_DeleteUser(t *testing.T) {
	if err := d.DeleteUser(user); err != nil {
		t.Errorf("delete user failed: %s", err)
		return
	}
	if _, err := d.FindUserById(user.ID.Hex()); err == nil {
		t.Error("user deleted, but found")
	}
}

func TestDAO_InsertProfile(t *testing.T) {
	if err := d.InsertProfile(profile); err != nil {
		t.Errorf("insert profile failed: %s", err)
		return
	}
}

func TestDAO_InsertProfileAgain(t *testing.T) {
	var profile2 = models.Profile{
		ID:       bson.NewObjectId(),
		Pk:       profile.Pk,
		Created:  now,
		LastSeen: now,
	}
	if err := d.InsertProfile(profile2); err == nil {
		t.Error("pk should be unique")
		return
	}
}

func TestDAO_FindProfileById(t *testing.T) {
	loaded, err := d.FindProfileById(profile.ID.Hex())
	if err != nil {
		t.Errorf("find profile by id failed: %s", err)
		return
	}
	if loaded.Pk != profile.Pk {
		t.Error("pk mismatch")
	}
}

func TestDAO_FindProfileByPk(t *testing.T) {
	if _, err := d.FindProfileByPk(profile.Pk); err != nil {
		t.Errorf("find profile by pk failed: %s", err)
		return
	}
}

func TestDAO_UpdateProfile(t *testing.T) {
	lastSeen := time.Now().Add(time.Minute)
	profile.LastSeen = lastSeen
	if err := d.UpdateProfile(profile); err != nil {
		t.Errorf("update profile failed: %s", err)
		return
	}
	loaded, err := d.FindProfileById(profile.ID.Hex())
	if err != nil {
		t.Errorf("find profile again by id failed: %s", err)
		return
	}
	if loaded.LastSeen.Unix() != profile.LastSeen.Unix() {
		t.Error("last seen mismatch")
	}
}

func TestDAO_DeleteProfile(t *testing.T) {
	if err := d.DeleteProfile(profile); err != nil {
		t.Errorf("delete profile failed: %s", err)
		return
	}
	if _, err := d.FindProfileById(profile.ID.Hex()); err == nil {
		t.Error("profile deleted, but found")
	}
}

func TestDAO_InsertReferral(t *testing.T) {
	if err := d.InsertReferral(ref); err != nil {
		t.Errorf("insert ref failed: %s", err)
		return
	}
}

func TestDAO_FindReferralByCode(t *testing.T) {
	if _, err := d.FindReferralByCode(ref.Code); err != nil {
		t.Errorf("find ref by code failed: %s", err)
		return
	}
}

func TestDAO_ListUnusedReferrals(t *testing.T) {
	refs, err := d.ListUnusedReferrals()
	if err != nil {
		t.Errorf("list unused refs failed: %s", err)
		return
	}
	unusedRefCnt = len(refs)
}

func TestDAO_UpdateReferral(t *testing.T) {
	ref.Remaining = 0
	err := d.UpdateReferral(ref)
	if err != nil {
		t.Errorf("update ref failed: %s", err)
		return
	}
	loaded, err := d.FindReferralByCode(ref.Code)
	if err != nil {
		t.Errorf("find ref again by code failed: %s", err)
		return
	}
	if loaded.Remaining != 0 {
		t.Error("remaining count mismatch")
	}
}

func TestDAO_ListUnusedReferralsAgain(t *testing.T) {
	refs, err := d.ListUnusedReferrals()
	if err != nil {
		t.Errorf("list unused refs failed: %s", err)
		return
	}
	if len(refs) != unusedRefCnt-1 {
		t.Error("incorrect number of unused refs")
	}
}

func TestDAO_DeleteReferral(t *testing.T) {
	if err := d.DeleteReferral(ref); err != nil {
		t.Errorf("delete ref failed: %s", err)
		return
	}
	if _, err := d.FindReferralByCode(ref.Code); err == nil {
		t.Error("ref deleted, but found")
	}
}

func TestDAO_InsertNonce(t *testing.T) {
	if err := d.InsertNonce(nonce); err != nil {
		t.Errorf("insert nonce failed: %s", err)
		return
	}
}

func TestDAO_InsertNonceAgain(t *testing.T) {
	var nonce2 = models.Nonce{
		ID:      bson.NewObjectId(),
		Pk:      ksuid.New().String(),
		Value:   nonce.Value,
		Created: now,
	}
	if err := d.InsertNonce(nonce2); err == nil {
		t.Error("nonce value should be unique")
		return
	}
}

func TestDAO_FindNonce(t *testing.T) {
	loaded, err := d.FindNonce(nonce.Value)
	if err != nil {
		t.Errorf("find nonce by value failed: %s", err)
		return
	}
	if loaded.Value != nonce.Value {
		t.Error("value mismatch")
	}
}

func TestDAO_DeleteNonce(t *testing.T) {
	if err := d.DeleteNonce(nonce); err != nil {
		t.Errorf("delete nonce failed: %s", err)
		return
	}
	if _, err := d.FindNonce(nonce.Value); err == nil {
		t.Error("nonce deleted, but found")
	}
}
