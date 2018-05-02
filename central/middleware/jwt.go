package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/v1/users" || c.Request.URL.Path == "/api/v1/referrals" {
			return
		}
		_, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			c.AbortWithError(401, err)
		}
	}
}
