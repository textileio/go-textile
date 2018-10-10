package dao

import (
	"crypto/tls"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/textileio/textile-go/cafe/models"
	"log"
	"net"
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
	userCollection     = "users"
	profileCollection  = "profiles"
	referralCollection = "referrals"
	nonceCollection    = "nonce"
)

var indexes = map[string][]mgo.Index{
	profileCollection: {
		{
			Key:        []string{"address"},
			Unique:     true,
			DropDups:   true,
			Background: true,
		},
	},
	referralCollection: {
		{
			Key:        []string{"code"},
			Unique:     true,
			DropDups:   true,
			Background: true,
		},
		{
			Key:        []string{"user_id"},
			Background: true,
		},
	},
	nonceCollection: {
		{
			Key:        []string{"value"},
			Unique:     true,
			DropDups:   true,
			Background: true,
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

// REFERRALS

// Find a referral by code
func (m *DAO) FindReferralByCode(code string) (models.Referral, error) {
	var ref models.Referral
	err := db.C(referralCollection).Find(bson.M{"code": code}).One(&ref)
	return ref, err
}

// List referrals
func (m *DAO) ListUnusedReferrals() ([]models.Referral, error) {
	var refs []models.Referral
	err := db.C(referralCollection).Find(bson.M{"remaining": bson.M{"$gt": 0}}).All(&refs)
	return refs, err
}

// Insert a new referral
func (m *DAO) InsertReferral(ref models.Referral) error {
	err := db.C(referralCollection).Insert(&ref)
	return err
}

// Delete an existing referral
func (m *DAO) DeleteReferral(ref models.Referral) error {
	err := db.C(referralCollection).Remove(&ref)
	return err
}

// Update an existing referral
func (m *DAO) UpdateReferral(ref models.Referral) error {
	err := db.C(referralCollection).UpdateId(ref.ID, &ref)
	return err
}

// PROFILES

// Find a profile by id
func (m *DAO) FindProfileById(id string) (models.Profile, error) {
	var profile models.Profile
	err := db.C(profileCollection).FindId(bson.ObjectIdHex(id)).One(&profile)
	return profile, err
}

// Find a profile by public key
func (m *DAO) FindProfileByAddress(address string) (models.Profile, error) {
	var profile models.Profile
	err := db.C(profileCollection).Find(bson.M{"address": address}).One(&profile)
	return profile, err
}

// Insert a new profile
func (m *DAO) InsertProfile(profile models.Profile) error {
	err := db.C(profileCollection).Insert(&profile)
	return err
}

// Delete an existing profile
func (m *DAO) DeleteProfile(profile models.Profile) error {
	err := db.C(profileCollection).Remove(&profile)
	return err
}

// Update an existing profile
func (m *DAO) UpdateProfile(profile models.Profile) error {
	err := db.C(profileCollection).UpdateId(profile.ID, &profile)
	return err
}

// NONCES

// Find a nonce value
func (m *DAO) FindNonce(value string) (models.Nonce, error) {
	var nonce models.Nonce
	err := db.C(nonceCollection).Find(bson.M{"value": value}).One(&nonce)
	return nonce, err
}

// Insert a new nonce
func (m *DAO) InsertNonce(nonce models.Nonce) error {
	err := db.C(nonceCollection).Insert(&nonce)
	return err
}

// Delete an existing nonce
func (m *DAO) DeleteNonce(nonce models.Nonce) error {
	err := db.C(nonceCollection).Remove(&nonce)
	return err
}
