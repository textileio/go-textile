package dao_test

import (
	"testing"
	"time"
	"fmt"
	"os"

	"github.com/globalsign/mgo/bson"
	"github.com/segmentio/ksuid"
	"github.com/joho/godotenv"

	"github.com/textileio/textile-go/central/dao"
	"github.com/textileio/textile-go/central/models"
)

var d = dao.DAO{}

var reg = &models.Registration{
	Identity: &models.Identity{
		Type: models.EmailAddress,
		Value: fmt.Sprintf("%s@textile.io", ksuid.New().String()),
		Verified: false,
	},
}
var now = time.Now()
var user = models.User{
	ID:         bson.NewObjectId(),
	Username:   ksuid.New().String(),
	Password:   ksuid.New().String(),
	Created:    now,
	LastSeen:   now,
	Identities: []models.Identity{*reg.Identity},
}

func TestDao_Connect(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	d.Hostname = os.Getenv("HOSTNAME")
	d.DatabaseName = os.Getenv("DATABASENAME")
	d.Connect()
}

func TestDAO_Insert(t *testing.T) {
	if err := d.InsertUser(user); err != nil {
		t.Errorf("insert user failed: %s", err)
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
		t.Errorf("find user by usernamem failed: %s", err)
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
