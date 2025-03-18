package controllers

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/matheusabido/kfofo-api/db"
	"github.com/matheusabido/kfofo-api/middleware"
	"github.com/matheusabido/kfofo-api/validator"
	"golang.org/x/crypto/bcrypt"
)

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func PostLogin(ctx *gin.Context) {
	var data LoginDTO

	if !validator.BindAndValidate(ctx, &data) {
		return
	}

	var id int
	var email string
	var password string
	err := db.Instance.QueryRow(context.Background(), "SELECT id, email, password FROM users WHERE email = $1", data.Email).Scan(&id, &email, &password)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid credentials."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(data.Password))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid credentials."})
		return
	}

	claims := middleware.JWTClaims{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	sign := []byte(os.Getenv("JWT_SIGN"))
	tokenString, err := token.SignedString(sign)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{
		"id":    id,
		"token": tokenString,
	})
}
