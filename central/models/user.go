package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type User struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Username   string        `bson:"username" json:"username"`
	Password   string        `bson:"password" json:"password"`
	Created    time.Time     `bson:"created" json:"created"`
	LastSeen   time.Time     `bson:"last_seen" json:"last_seen"`
	Identities []Identity    `bson:"identities" json:"identities"`
}

type Registration struct {
	Username string    `json:"username" binding:"required"`
	Password string    `json:"password" binding:"required"`
	Identity *Identity `json:"identity" binding:"required"`
}

type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type IdentityType string

const (
	PhoneNumber  IdentityType = "phone_number"
	EmailAddress IdentityType = "email_address"
)

type Identity struct {
	Type     IdentityType `bson:"type" json:"type" binding:"required"`
	Value    string       `bson:"value" json:"value" binding:"required"`
	Verified bool         `bson:"verified" json:"verified"`
}
