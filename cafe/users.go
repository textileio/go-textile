package cafe

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/nbutton23/zxcvbn-go"
	"github.com/textileio/textile-go/cafe/auth"
	"github.com/textileio/textile-go/cafe/dao"
	"github.com/textileio/textile-go/cafe/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"time"
)

var usernameRx = regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9._]+[a-zA-Z0-9_]$`)
var emailRx = regexp.MustCompile(`^[^@^\s]+@[^@^\s]+$`)
var numbersOnlyRx = regexp.MustCompile(`[^+^0-9]+`)

const month = time.Hour * 24 * 7 * 4

func (c *Cafe) signUp(g *gin.Context) {
	var reg models.Registration
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

	// test username
	// username validation based on twitter:
	// - only contain letters, number, period, or underscore
	// - must not start or end with period
	// - max 30 characters (twitter is 15, instagram is 30)
	valid := usernameRx.Match([]byte(reg.Username))
	if !valid || len(reg.Username) > 30 {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}

	// test email address
	if reg.Identity.Type == models.EmailAddress {
		// not trying to be too strict here, just:
		// - make sure there's at least one "@"
		valid = emailRx.Match([]byte(reg.Identity.Value))
		if !valid {
			g.JSON(http.StatusBadRequest, gin.H{"error": "invalid email address"})
			return
		}
	}

	// clean phone number
	if reg.Identity.Type == models.PhoneNumber {
		// no way gonna try and validate phone numbers, just:
		// - remove everything but numbers and "+"
		// - make sure its not zero-length
		cleaned := numbersOnlyRx.ReplaceAllString(reg.Identity.Value, "")
		if len(cleaned) == 0 {
			g.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone number"})
			return
		}
		reg.Identity.Value = cleaned
	}

	// check password strength
	match := zxcvbn.PasswordStrength(reg.Password, []string{reg.Identity.Value})
	if match.Score < 1 {
		msg := fmt.Sprintf("weak password - crackable in %s", match.CrackTimeDisplay)
		g.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	// hash password
	password, err := hashAndSalt(reg.Password)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create a new user
	now := time.Now()
	reg.Identity.Verified = false
	user := models.User{
		ID:         bson.NewObjectId(),
		Username:   reg.Username,
		Password:   password,
		Created:    now,
		LastSeen:   now,
		Identities: []models.Identity{*reg.Identity},
	}
	if err := dao.Dao.InsertUser(user); err != nil {
		g.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// get a session
	session, err := auth.NewSession(user.ID.Hex(), c.TokenSecret, c.Ipfs().Identity.Pretty(), month)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// lastly, mark the code as used
	ref.Remaining = ref.Remaining - 1
	if err := dao.Dao.UpdateReferral(ref); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ship it
	g.JSON(http.StatusCreated, models.Response{
		Status:  http.StatusCreated,
		Session: session,
	})
}

func (c *Cafe) signIn(g *gin.Context) {
	var creds models.Credentials
	if err := g.BindJSON(&creds); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// lookup username
	user, err := dao.Dao.FindUserByUsername(creds.Username)
	if err != nil {
		g.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// check password
	if !checkPassword(user.Password, creds.Password) {
		g.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// get a session
	session, err := auth.NewSession(user.ID.Hex(), c.TokenSecret, c.Ipfs().Identity.Pretty(), month)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ship it
	g.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Session: session,
	})
}

func hashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPassword(hashed string, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
