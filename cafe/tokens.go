package cafe

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/cafe/auth"
	"github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/net/service"
	"net/http"
	"strings"
	"time"
)

func (c *Cafe) verify(token *jwt.Token) (interface{}, error) {
	return []byte(c.TokenSecret), nil
}

func (c *Cafe) auth(g *gin.Context) {
	if g.Request.URL.Path == "/api/v0/users" || g.Request.URL.Path == "/api/v0/referrals" {
		return
	}

	// extract token string from request header
	var tokenString string
	parsed := strings.Split(g.Request.Header.Get("Authorization"), " ")
	if len(parsed) == 2 {
		tokenString = parsed[1]
	}

	// parse it
	token, pErr := jwt.Parse(tokenString, c.verify)

	// pull out claims
	claims, err := auth.ParseClaims(token.Claims)
	if err != nil {
		g.AbortWithError(403, err)
		return
	}

	// check valid
	if pErr != nil {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			// 401 indicates a retry is expected after a token refresh
			g.AbortWithError(401, auth.ErrInvalidClaims)
			return
		}
		g.AbortWithError(403, pErr)
		return
	}

	// check scope
	switch claims.Scope {
	case auth.Access:
		break
	case auth.Refresh:
		if g.Request.URL.Path != "/api/v0/tokens" {
			g.AbortWithError(403, auth.ErrInvalidClaims)
			return
		}
	default:
		g.AbortWithError(403, auth.ErrInvalidClaims)
		return
	}

	// verify extra fields
	if !claims.VerifyAudience(string(service.TextileProtocol), true) {
		g.AbortWithError(403, auth.ErrInvalidClaims)
		return
	}
}

func (c *Cafe) refreshToken(g *gin.Context) {
	var session models.Session
	if err := g.BindJSON(&session); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ensure bearer matches payload refresh token
	var tokenString string
	parsed := strings.Split(g.Request.Header.Get("Authorization"), " ")
	if len(parsed) == 2 {
		tokenString = parsed[1]
	}
	if session.RefreshToken != tokenString {
		g.AbortWithError(403, auth.ErrInvalidClaims)
		return
	}

	// ensure access and token are a valid pair
	access, _ := jwt.Parse(session.AccessToken, c.verify)
	if access == nil {
		g.AbortWithError(403, auth.ErrInvalidClaims)
		return
	}
	refresh, _ := jwt.Parse(session.RefreshToken, c.verify)
	if refresh == nil {
		g.AbortWithError(403, auth.ErrInvalidClaims)
		return
	}
	accessClaims, err := auth.ParseClaims(access.Claims)
	if err != nil {
		g.AbortWithError(403, err)
		return
	}
	refreshClaims, err := auth.ParseClaims(refresh.Claims)
	if err != nil {
		g.AbortWithError(403, err)
		return
	}
	if refreshClaims.Id[1:] != accessClaims.Id {
		g.AbortWithError(403, auth.ErrInvalidClaims)
		return
	}

	// get a new session
	refreshed, err := auth.NewSession(session.SubjectId, c.TokenSecret, c.Ipfs().Identity.Pretty(), service.TextileProtocol, month)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ship it
	g.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Session: refreshed,
	})
}
