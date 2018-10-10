package cafe

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/cafe/auth"
	"github.com/textileio/textile-go/cafe/dao"
	"github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/keypair"
	"net/http"
	"time"
)

func (c *Cafe) profileChallenge(g *gin.Context) {
	var req models.ChallengeRequest
	if err := g.BindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate public key
	accnt, err := keypair.Parse(req.Address)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := accnt.Sign([]byte{0x00}); err == nil {
		// we dont want to handle account seeds, just addresses
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// generate a new random nonce
	nonce := models.Nonce{
		ID:      bson.NewObjectId(),
		Address: req.Address,
		Value:   ksuid.New().String(),
		Created: time.Now(),
	}
	if err := dao.Dao.InsertNonce(nonce); err != nil {
		g.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// ship it
	g.JSON(http.StatusCreated, models.ChallengeResponse{
		Value: &nonce.Value,
	})
}

func (c *Cafe) registerProfile(g *gin.Context) {
	var reg models.ProfileRegistration
	if err := g.BindJSON(&reg); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// lookup the referral code
	ref, err := dao.Dao.FindReferralByCode(reg.Referral)
	if err != nil || ref.Remaining == 0 {
		g.JSON(http.StatusNotFound, gin.H{"error": "invalid or used referral code"})
		return
	}

	// lookup the nonce
	snonce, err := dao.Dao.FindNonce(reg.Challenge.Value)
	if err != nil {
		g.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if snonce.Address != reg.Challenge.Address {
		g.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// validate public key
	accnt, err := keypair.Parse(reg.Challenge.Address)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verify
	payload := []byte(reg.Challenge.Value + reg.Challenge.Nonce)
	sig, err := base64.StdEncoding.DecodeString(reg.Challenge.Signature)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := accnt.Verify(payload, sig); err != nil {
		g.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// create new
	now := time.Now()
	profile := models.Profile{
		ID:       bson.NewObjectId(),
		Address:  reg.Challenge.Address,
		Created:  now,
		LastSeen: now,
	}
	if err := dao.Dao.InsertProfile(profile); err != nil {
		g.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// get a session
	session, err := auth.NewSession(profile.ID.Hex(), c.TokenSecret, c.Ipfs().Identity.Pretty(), "FIXME", oneMonth)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// mark the code as used
	ref.Remaining = ref.Remaining - 1
	if err := dao.Dao.UpdateReferral(ref); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// delete the nonce
	if err := dao.Dao.DeleteNonce(snonce); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ship it
	g.JSON(http.StatusCreated, models.SessionResponse{
		Session: session,
	})
}

func (c *Cafe) loginProfile(g *gin.Context) {
	var cha models.SignedChallenge
	if err := g.BindJSON(&cha); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// lookup pk
	profile, err := dao.Dao.FindProfileByAddress(cha.Address)
	if err != nil {
		g.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// lookup the nonce
	snonce, err := dao.Dao.FindNonce(cha.Value)
	if err != nil {
		g.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if snonce.Address != cha.Address {
		g.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// validate public key
	accnt, err := keypair.Parse(profile.Address)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verify
	payload := []byte(cha.Value + cha.Nonce)
	sig, err := base64.StdEncoding.DecodeString(cha.Signature)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := accnt.Verify(payload, sig); err != nil {
		g.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// get a session
	session, err := auth.NewSession(profile.ID.Hex(), c.TokenSecret, c.Ipfs().Identity.Pretty(), "FIXME", oneMonth)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// delete the nonce
	if err := dao.Dao.DeleteNonce(snonce); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ship it
	g.JSON(http.StatusOK, models.SessionResponse{
		Session: session,
	})
}
