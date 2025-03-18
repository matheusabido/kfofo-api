package middleware

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(ctx *gin.Context) {
	authorization := ctx.GetHeader("Authorization")

	if len(authorization) < 7 || authorization[:7] != "Bearer " {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		ctx.Abort()
		return
	}

	tokenString := authorization[7:]
	claims := &JWTClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			ctx.Abort()
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		secret := []byte(os.Getenv("JWT_SIGN"))
		return secret, nil
	})

	if err != nil {
		ctx.Abort()
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	ctx.Set("claims", claims)
	ctx.Next()
}
