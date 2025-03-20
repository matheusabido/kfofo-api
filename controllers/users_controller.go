package controllers

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/matheusabido/kfofo-api/db"
	"github.com/matheusabido/kfofo-api/middleware"
	"github.com/matheusabido/kfofo-api/utils"
	"github.com/matheusabido/kfofo-api/validator"
	"golang.org/x/crypto/bcrypt"
)

type StoreUserDTO struct {
	Name      string `json:"name" validate:"required,min=5"`
	Email     string `json:"email" validate:"required,email"`
	BirthDate string `json:"birth_date" validate:"required,datetime=2006-01-02"`
	Password  string `json:"password" validate:"required,min=8"`
}

func GetUser(ctx *gin.Context) {
	idRaw := ctx.Param("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	user := utils.GetUser(ctx)
	if user == nil {
		return
	}

	if id != user.Id {
		ctx.JSON(403, gin.H{"error": "You can't see this user's details."})
		return
	}

	ctx.JSON(200, gin.H{
		"id":         id,
		"name":       user.Name,
		"email":      user.Email,
		"birth_date": user.BirthDate.Format("2006-01-02"),
	})
}

func PostUser(ctx *gin.Context) {
	var data StoreUserDTO

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

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{
			"error": "Internal server error",
		})
		return
	}

	var id int
	err = db.Instance.QueryRow(context.Background(), "INSERT INTO users (name, email, birth_date, password) VALUES ($1, $2, $3, $4) RETURNING id", data.Name, data.Email, data.BirthDate, string(encryptedPassword)).Scan(&id)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{
			"error": "Internal server error",
		})
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
		"name":       data.Name,
		"email":      data.Email,
		"birth_date": data.BirthDate,
		"token":      tokenString,
	})
}

type UpdateUserDTO struct {
	Name        string `json:"name" validate:"omitempty,min=5"`
	BirthDate   string `json:"birth_date" validate:"omitempty,datetime=2006-01-02"`
	NewPassword string `json:"new_password" validate:"omitempty,min=8"`
	Password    string `json:"password" validate:"required,min=8"`
}

func PutUser(ctx *gin.Context) {
	var data UpdateUserDTO
	idValue := ctx.Param("id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	if !validator.BindAndValidate(ctx, &data) {
		return
	}

	user := utils.GetUser(ctx)
	if user.Id != id {
		ctx.JSON(403, gin.H{"error": "You can't edit this user's info."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid credentials."})
		return
	}

	index := 1
	var updates []string
	var values []any
	if len(data.Name) > 0 {
		updates = append(updates, "name = $"+strconv.Itoa(index))
		values = append(values, data.Name)
		index++
	}
	if len(data.BirthDate) > 0 {
		updates = append(updates, "birth_date = $"+strconv.Itoa(index))
		values = append(values, data.BirthDate)
		index++
	}
	if len(data.NewPassword) > 0 {
		updates = append(updates, "password = $"+strconv.Itoa(index))
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		values = append(values, string(encryptedPassword))
		index++
	}

	if len(updates) == 0 {
		ctx.JSON(200, gin.H{"message": "Data updated."})
		return
	}

	var builder strings.Builder
	builder.WriteString("UPDATE users SET ")
	builder.WriteString(strings.Join(updates, ", "))
	builder.WriteString(" WHERE id = $" + strconv.Itoa(index))
	values = append(values, id)

	_, err = db.Instance.Exec(context.Background(), builder.String(), values...)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Data updated."})
}

func DeleteUser(ctx *gin.Context) {
	idValue := ctx.Param("id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	claims := utils.GetClaims(ctx)
	if claims.Id != id {
		ctx.JSON(403, gin.H{"error": "You can't delete this user's account."})
		return
	}

	_, err = db.Instance.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"message": "User deleted."})
}
