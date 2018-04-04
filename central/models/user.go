package models

import (
	"time"
	"github.com/globalsign/mgo/bson"
)

type User struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
	Created time.Time `bson:"created" json:"created"`
	LastSeen time.Time `bson:"last_seen" json:"last_seen"`
	Identities []Identity `bson:"identities" json:"identities"`
}
