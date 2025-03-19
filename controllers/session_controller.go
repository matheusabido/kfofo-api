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

type loginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func PostLogin(ctx *gin.Context) {
	var data loginDTO

	if !validator.BindAndValidate(ctx, &data) {
		return
	}

	var id int
	var name string
	var email string
	var birthDate time.Time
	var password string
	err := db.Instance.QueryRow(context.Background(), "SELECT id, name, email, birth_date, password FROM users WHERE email = $1", data.Email).Scan(&id, &name, &email, &birthDate, &password)
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
		"id":         id,
		"name":       name,
		"email":      email,
		"birth_date": birthDate.Format("2006-01-02"),
		"token":      tokenString,
	})
}
