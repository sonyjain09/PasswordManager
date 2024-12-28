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
        // get token from Authorization header
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
            c.Abort()
            return
        }

        // Remove "Bearer " prefix
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        // Validate token

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                fmt.Printf("Unexpected signing method: %v\n", token.Header["alg"])
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(jwt_secret), nil
        })

        if token != nil && token.Header != nil {
            fmt.Printf("Token algorithm: %v\n", token.Header["alg"])
        } else {
            fmt.Println("Token parsing failed")
        }

        if err != nil || !token.Valid {
            fmt.Printf("Token validation error: %v\n", err)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Pass control to next handler
        c.Next()
    }
}

