package dao_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/joho/godotenv"
	"github.com/segmentio/ksuid"

	"github.com/textileio/textile-go/central/dao"
	"github.com/textileio/textile-go/central/models"
)

var d = dao.DAO{}

var now = time.Now()
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
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	d.Hostname = os.Getenv("HOSTNAME")
	d.DatabaseName = os.Getenv("DATABASE")
	d.Connect()
}

func TestDao_Index(t *testing.T) {
	d.Index()
}

func TestDAO_Insert(t *testing.T) {
	if err := d.InsertUser(user); err != nil {
		t.Errorf("insert user failed: %s", err)
		return
	}
}

func TestDAO_InsertAgain(t *testing.T) {
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

func TestDAO_FindById(t *testing.T) {
	loaded, err := d.FindUserById(user.ID.Hex())
	if err != nil {
		t.Errorf("find user by id failed: %s", err)
		return
	}
	if loaded.Username != user.Username {
		t.Error("username mismatch")
	}
}

func TestDAO_FindByUsername(t *testing.T) {
	_, err := d.FindUserByUsername(user.Username)
	if err != nil {
		t.Errorf("find user by username failed: %s", err)
		return
	}
}

func TestDAO_FindByIdentity(t *testing.T) {
	_, err := d.FindUserByIdentity(user.Identities[0])
	if err != nil {
		t.Errorf("find user by identity failed: %s", err)
		return
	}
}

func TestDAO_Update(t *testing.T) {
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

func TestDAO_Delete(t *testing.T) {
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
