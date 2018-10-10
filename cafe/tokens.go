package cafe

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/cafe/auth"
	"github.com/textileio/textile-go/cafe/models"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var errForbidden = "forbidden"
var forbiddenResponse = models.SessionResponse{
	Error: &errForbidden,
}
var ErrUnauthorized = "unauthorized"
var unauthorizedResponse = models.SessionResponse{
	Error: &ErrUnauthorized,
}

func (c *Cafe) verify(token *jwt.Token) (interface{}, error) {
	return []byte(c.TokenSecret), nil
}

func (c *Cafe) authSession(g *gin.Context) {
	// extract token string from request header
	var tokenString string
	parsed := strings.Split(g.Request.Header.Get("Authorization"), " ")
	if len(parsed) == 2 {
		tokenString = parsed[1]
	}

	// parse it
	token, pErr := jwt.Parse(tokenString, c.verify)
	if token == nil {
		g.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedResponse)
		return
	}

	// pull out claims
	claims, err := auth.ParseClaims(token.Claims)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusForbidden, forbiddenResponse)
		return
	}

	// check valid
	if pErr != nil {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			// 401 indicates a retry is expected after a token refresh
			g.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedResponse)
			return
		}
		g.AbortWithStatusJSON(http.StatusForbidden, forbiddenResponse)
		return
	}

	// check scope
	tokenRoute := strings.Contains(g.Request.URL.Path, "/tokens")
	switch claims.Scope {
	case auth.Access:
		if tokenRoute {
			g.AbortWithStatusJSON(http.StatusForbidden, forbiddenResponse)
			return
		}
	case auth.Refresh:
		if !tokenRoute {
			g.AbortWithStatusJSON(http.StatusForbidden, forbiddenResponse)
			return
		}
	default:
		g.AbortWithStatusJSON(http.StatusForbidden, forbiddenResponse)
		return
	}

	// verify extra fields
	if !claims.VerifyAudience(string("FIXME"), true) {
		g.AbortWithStatusJSON(http.StatusForbidden, forbiddenResponse)
		return
	}
}

func (c *Cafe) refreshSession(g *gin.Context) {
	body, err := ioutil.ReadAll(g.Request.Body)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accessToken := string(body)

	// ensure bearer matches payload refresh token
	var refreshToken string
	parsed := strings.Split(g.Request.Header.Get("Authorization"), " ")
	if len(parsed) == 2 {
		refreshToken = parsed[1]
	}

	// ensure access and token are a valid pair
	access, _ := jwt.Parse(accessToken, c.verify)
	if access == nil {
		g.JSON(http.StatusForbidden, forbiddenResponse)
		return
	}
	refresh, _ := jwt.Parse(refreshToken, c.verify)
	if refresh == nil {
		g.JSON(http.StatusForbidden, forbiddenResponse)
		return
	}
	accessClaims, err := auth.ParseClaims(access.Claims)
	if err != nil {
		g.JSON(http.StatusForbidden, forbiddenResponse)
		return
	}
	refreshClaims, err := auth.ParseClaims(refresh.Claims)
	if err != nil {
		g.JSON(http.StatusForbidden, forbiddenResponse)
		return
	}
	if refreshClaims.Id[1:] != accessClaims.Id {
		g.JSON(http.StatusForbidden, forbiddenResponse)
		return
	}
	if refreshClaims.Subject != accessClaims.Subject {
		g.JSON(http.StatusForbidden, forbiddenResponse)
		return
	}

	// get a new session
	refreshed, err := auth.NewSession(accessClaims.Subject, c.TokenSecret, c.Ipfs().Identity.Pretty(), "FIXME", oneMonth)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ship it
	g.JSON(http.StatusOK, models.SessionResponse{
		Session: refreshed,
	})
}
