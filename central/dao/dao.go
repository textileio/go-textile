package dao

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/textileio/textile-go/central/models"
)

type DAO struct {
	Hosts    string
	Name     string
	User     string
	Password string
	TLS      bool
}

var Dao *DAO

var db *mgo.Database

const (
	userCollection = "users"
)

var indexes = map[string][]mgo.Index{
	userCollection: {
		{
			Key:        []string{"username"},
			Unique:     true,
			DropDups:   true,
			Background: true,
		},
		{
			Key:        []string{"identities.value", "identities.type"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     true,
		},
	},
}

func (m *DAO) Index() {
	for cn, list := range indexes {
		for _, index := range list {
			if err := db.C(cn).EnsureIndex(index); err != nil {
				log.Fatal(err)
			}
		}
	}
}

// Establish a connection to database
func (m *DAO) Connect() {
	creds := fmt.Sprintf("%s:%s@", m.User, m.Password)
	if len(creds) == 2 {
		creds = ""
	}
	uri := fmt.Sprintf("mongodb://%s%s/%s", creds, m.Hosts, m.Name)
	dialInfo, err := mgo.ParseURL(uri)
	if err != nil {
		log.Fatal(err)
	}
	if m.TLS {
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			tlsConfig := &tls.Config{
				InsecureSkipVerify: true,
			}
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			if err != nil {
				log.Fatal(err)
			}
			return conn, err
		}
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Name)
}

// Find a user by id
func (m *DAO) FindUserById(id string) (models.User, error) {
	var user models.User
	err := db.C(userCollection).FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

// Find a user by username
func (m *DAO) FindUserByUsername(un string) (models.User, error) {
	var user models.User
	err := db.C(userCollection).Find(bson.M{"username": un}).One(&user)
	return user, err
}

// Find a user by email
func (m *DAO) FindUserByIdentity(id models.Identity) (models.User, error) {
	var user models.User
	err := db.C(userCollection).Find(bson.M{
		"identities": bson.M{"$elemMatch": bson.M{"type": id.Type, "value": id.Value}},
	}).One(&user)
	return user, err
}

// Insert a user into database
func (m *DAO) InsertUser(user models.User) error {
	err := db.C(userCollection).Insert(&user)
	return err
}

// Delete an existing user
func (m *DAO) DeleteUser(user models.User) error {
	err := db.C(userCollection).Remove(&user)
	return err
}

// Update an existing user
func (m *DAO) UpdateUser(user models.User) error {
	err := db.C(userCollection).UpdateId(user.ID, &user)
	return err
}
