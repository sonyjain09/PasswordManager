package middleware

import (
	"net/http"
	"strings"

	"schedvault/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func init() {
	util.InitEnv()
}

func AuthMiddleware() gin.HandlerFunc {
	var jwt_secret = os.Getenv("JWT_SECRET")

	return func(c *gin.Context) {
		// Retrieve the token from the Authorization header
		token_string := c.GetHeader("Authorization")
		if token_string == "" { // no token
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// validate token
		token_string = strings.TrimPrefix(token_string, "Bearer ")
		token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwt_secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// pass control to next handler
		c.Next()
	}
}
