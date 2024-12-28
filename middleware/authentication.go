package middleware

import (
	"net/http"

	"schedvault/util"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"fmt"
)

func init() {
	util.InitEnv()
}

func AuthMiddleware() gin.HandlerFunc {
    var jwt_secret = os.Getenv("JWT_SECRET")

    return func(c *gin.Context) {

		// get token string 
		token_string := c.GetHeader("Authorization")
		if token_string == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}
		token_string = strings.TrimPrefix(token_string, "Bearer ")

		// get the jwt string and parse it
		token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwt_secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// get the user id from the token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		user_id, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		fmt.Printf("Token is valid, user_id: %v\n", user_id)

		// set the context user id
		c.Set("user_id", uint(user_id))

		// pass control to next middleware
		c.Next()
	}
}

