package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
	Created time.Time `bson:"created" json:"created"`
	LastSeen time.Time `bson:"last_seen" json:"last_seen"`
	Identities []Identity `bson:"identities" json:"identities"`
}
