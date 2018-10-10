package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Profile struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Address  string        `bson:"address" json:"address"`
	Created  time.Time     `bson:"created" json:"created"`
	LastSeen time.Time     `bson:"last_seen" json:"last_seen"`
}

type ChallengeRequest struct {
	Address string `json:"address" binding:"required"`
}

type ChallengeResponse struct {
	Value *string `json:"value,omitempty"`
	Error *string `json:"error,omitempty"`
}

type SignedChallenge struct {
	Address   string `json:"address" binding:"required"`
	Value     string `json:"value" binding:"required"`
	Nonce     string `json:"nonce" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type ProfileRegistration struct {
	Challenge SignedChallenge `json:"challenge" binding:"required"`
	Referral  string          `json:"ref_code" binding:"required"`
}

// SESSION

type Session struct {
	AccessToken      string `json:"access_token"`
	ExpiresAt        int64  `json:"expires_at"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresAt int64  `json:"refresh_expires_at"`
	SubjectId        string `json:"subject_id"`
	TokenType        string `json:"token_type"`
}

type SessionResponse struct {
	Session *Session `json:"session,omitempty"`
	Error   *string  `json:"error,omitempty"`
}

// REFERRALS

type Referral struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Code      string        `bson:"code" json:"code"`
	Created   time.Time     `bson:"created" json:"created"`
	Remaining int           `bson:"remaining" json:"remaining"`
	Requester string        `bson:"requester" json:"requester"`
}

type ReferralRequest struct {
	Key         string
	Count       int
	Limit       int
	RequestedBy string
}

type ReferralResponse struct {
	RefCodes []string `json:"ref_codes,omitempty"`
	Error    *string  `json:"error,omitempty"`
}

// NONCES

type Nonce struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Address string        `bson:"address" json:"address"`
	Value   string        `bson:"value" json:"value"`
	Created time.Time     `bson:"created" json:"created"`
}

// PINS

type PinResponse struct {
	Id    *string `json:"id,omitempty"`
	Error *string `json:"error,omitempty"`
}
