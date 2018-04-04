package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type User struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Created    time.Time     `bson:"created" json:"created"`
	LastSeen   time.Time     `bson:"last_seen" json:"last_seen"`
	Identities []Identity    `bson:"identities" json:"identities"`
}
