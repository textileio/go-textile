package controllers

import (
	"time"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/textileio/textile-go/central/models"
	"github.com/textileio/textile-go/central/dao"
	"github.com/segmentio/ksuid"
)

func SignUp(c *gin.Context) {
	var reg models.Registration
	if err := c.BindJSON(&reg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: check password strength
	password, err := hashAndSalt(reg.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// create a token
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims = jwt.MapClaims{
		"Id":  ksuid.New().String(),
		"exp": time.Now().Add(time.Hour*24*7).Unix(),
	}
	signed, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// ship it
	c.JSON(http.StatusCreated, models.Response{
		Status:     http.StatusCreated,
		ResourceID: user.ID.Hex(),
		Token:      signed,
	})
}

func SignIn(c *gin.Context) {
	var creds models.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// lookup username
	user, err := dao.Dao.FindUserByUsername(creds.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// check password
	if !checkPassword(user.Password, creds.Password) {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}

	// create a token
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims = jwt.MapClaims{
		"Id":  ksuid.New().String(),
		"exp": time.Now().Add(time.Hour*24*7).Unix(),
	}
	signed, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// ship it
	c.JSON(http.StatusOK, models.Response{
		Status:     http.StatusOK,
		ResourceID: user.ID.Hex(),
		Token:      signed,
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
