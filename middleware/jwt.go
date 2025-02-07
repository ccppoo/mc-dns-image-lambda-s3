package middleware

import (
	// "log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey = "asdjanskdansdkajsndklasdnlaskndlasndjklasdasd"
)

type MyCustomClaims struct {
	AccessType string `json:"typ"`
	UserID     string `json:"uid"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie
		// log.Printf("\n========== Auth Middleware Start =============\n")
		tokenString, err := c.Cookie("token")
		// log.Printf("tokenString : %s", tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// Parse JWT token
		// log.Println("start parsing token")
		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		}, jwt.WithLeeway(2*time.Second))

		if err != nil || !token.Valid {
			// log.Printf("tooken is not valid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			c.Set("claims", claims)
			// log.Printf("claims : ")
			// log.Printf("claims.AccessType : %s", claims.AccessType)
			// log.Printf("claims.Audience : %s", claims.Audience)
			// log.Printf("claims.ExpiresAt : %s", claims.ExpiresAt)
			// log.Printf("claims.ID : %s", claims.ID)
			// log.Printf("claims.IssuedAt : %s", claims.IssuedAt)
			// log.Printf("claims.ExpiresAt : %s", claims.ExpiresAt)
			// log.Printf("claims.Subject : %s", claims.Subject)
			// log.Printf("claims.UserID : %s", claims.UserID)

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		// log.Printf("\n========== Auth Middleware End =============\n")

		c.Next() // Proceed if authorized
	}
}
