package dao

import (
	"log"

	. "github.com/textileio/textile-go/central/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	userCollection = "users"
)

// Establish a connection to database
func (m *DAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find a user by id
func (m *DAO) FindById(id string) (User, error) {
	var user User
	err := db.C(userCollection).FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

// Insert a user into database
func (m *DAO) Insert(user User) error {
	err := db.C(userCollection).Insert(&user)
	return err
}

// Delete an existing user
func (m *DAO) Delete(user User) error {
	err := db.C(userCollection).Remove(&user)
	return err
}

// Update an existing user
func (m *DAO) Update(user User) error {
	err := db.C(userCollection).UpdateId(user.ID, &user)
	return err
}
