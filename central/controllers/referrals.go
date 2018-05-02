package controllers

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"

	"github.com/textileio/textile-go/central/dao"
	"github.com/textileio/textile-go/central/models"
)

func CreateReferral(c *gin.Context) {
	// cheap way to lock down this endpoint
	if os.Getenv("REF_KEY") != c.GetHeader("X-Referral-Key") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// how many should we make?
	count := 1
	params := c.Request.URL.Query()
	if params["count"] != nil && len(params["count"]) > 0 {
		tmp, err := strconv.ParseInt(params["count"][0], 10, 64)
		if err == nil {
			count = int(tmp)
		}
	}

	// hodl 'em
	refs := make([]string, count)
	for i := range refs {
		code, err := createReferral()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		refs[i] = code
	}

	// ship it
	c.JSON(http.StatusCreated, gin.H{
		"status":    http.StatusCreated,
		"ref_codes": refs,
	})
}

func ListReferrals(c *gin.Context) {
	// cheap way to lock down this endpoint
	if os.Getenv("REF_KEY") != c.GetHeader("X-Referral-Key") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// get 'em
	refs, err := dao.Dao.ListUnusedReferrals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	codes := make([]string, len(refs))
	for i, r := range refs {
		codes[i] = r.Code
	}

	// ship it
	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"ref_codes": codes,
	})
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func createReferral() (string, error) {
	code := randString(5)
	ref := models.Referral{
		ID:      bson.NewObjectId(),
		Code:    code,
		Created: time.Now(),
	}
	if err := dao.Dao.InsertReferral(ref); err != nil {
		return "", err
	}
	return code, nil
}
