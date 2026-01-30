package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Protect(signature []byte) gin.HandlerFunc {

	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(auth, "Bearer ")

		_, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return signature, nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
