package dao_test

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/central/dao"
	"github.com/textileio/textile-go/central/models"
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
	Identities: []models.Identity{
		{
			Type:     models.EmailAddress,
			Value:    fmt.Sprintf("%s@textile.io", ksuid.New().String()),
			Verified: false,
		},
	},
}

func TestDao_Connect(t *testing.T) {
	d.Hosts = os.Getenv("DB_HOSTS")
	d.Name = os.Getenv("DB_NAME")
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
		Identities: []models.Identity{
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
	_, err := d.FindUserByUsername(user.Username)
	if err != nil {
		t.Errorf("find user by username failed: %s", err)
		return
	}
}

func TestDAO_FindUserByIdentity(t *testing.T) {
	_, err := d.FindUserByIdentity(user.Identities[0])
	if err != nil {
		t.Errorf("find user by identity failed: %s", err)
		return
	}
}

func TestDAO_UpdateUser(t *testing.T) {
	un := ksuid.New().String()
	user.Username = un
	err := d.UpdateUser(user)
	if err != nil {
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
	err := d.DeleteUser(user)
	if err != nil {
		t.Errorf("delete user failed: %s", err)
		return
	}
	_, err = d.FindUserById(user.ID.Hex())
	if err == nil {
		t.Error("user deleted, but found")
	}
}

func TestDAO_InsertReferral(t *testing.T) {
	if err := d.InsertReferral(ref); err != nil {
		t.Errorf("insert ref failed: %s", err)
		return
	}
}

func TestDAO_FindReferralByCode(t *testing.T) {
	_, err := d.FindReferralByCode(ref.Code)
	if err != nil {
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
	err := d.DeleteReferral(ref)
	if err != nil {
		t.Errorf("delete ref failed: %s", err)
		return
	}
	_, err = d.FindReferralByCode(ref.Code)
	if err == nil {
		t.Error("ref deleted, but found")
	}
}
