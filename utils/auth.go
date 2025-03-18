package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/matheusabido/kfofo-api/middleware"
)

func GetClaims(ctx *gin.Context) *middleware.JWTClaims {
	claimsValue, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return nil
	}

	claims, ok := claimsValue.(*middleware.JWTClaims)
	if !ok {
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return nil
	}

	return claims
}

func GetUser(ctx *gin.Context) *middleware.User {
	userValue, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return nil
	}

	user, ok := userValue.(*middleware.User)
	if !ok {
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return nil
	}

	return user
}
