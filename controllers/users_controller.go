package controllers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/matheusabido/kfofo-api/db"
	"github.com/matheusabido/kfofo-api/validator"
	"golang.org/x/crypto/bcrypt"
)

type StoreDTO struct {
	Name      string `json:"name" validate:"required,min=5"`
	Email     string `json:"email" validate:"required,email"`
	BirthDate string `json:"birth_date" validate:"required,datetime=2006-01-02"`
	Password  string `json:"password" validate:"required,min=8"`
}

func PostUser(ctx *gin.Context) {
	var data StoreDTO

	if !validator.BindAndValidate(ctx, &data) {
		return
	}

	var exists bool
	err := db.Instance.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", data.Email).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if exists {
		ctx.JSON(400, gin.H{"error": "This e-mail is already registered."})
		return
	}

	var id int
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{
			"error": "Internal server error",
		})
		return
	}

	err = db.Instance.QueryRow(context.Background(), "INSERT INTO users (name, email, birth_date, password) VALUES ($1, $2, $3, $4) RETURNING id", data.Name, data.Email, data.BirthDate, string(encryptedPassword)).Scan(&id)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"id":         id,
		"name":       data.Name,
		"email":      data.Email,
		"birth_date": data.BirthDate,
	})
}
