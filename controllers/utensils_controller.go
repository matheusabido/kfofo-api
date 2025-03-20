package controllers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/matheusabido/kfofo-api/db"
	"github.com/matheusabido/kfofo-api/utils"
	"github.com/matheusabido/kfofo-api/validator"
)

func GetUtensils(ctx *gin.Context) {
	homeQuery, exists := ctx.GetQuery("home_id")
	var rows pgx.Rows
	var err error

	if exists && homeQuery != "" {
		homeId, err := strconv.Atoi(homeQuery)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		query := `
			SELECT u.id, u.name 
			FROM utensils u
			INNER JOIN home_utensils_pivot hup ON hup.utensil_id = u.id 
			WHERE hup.home_id = $1
		`
		rows, err = db.Instance.Query(context.Background(), query, homeId)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
	} else {
		query := "SELECT id, name FROM utensils"
		rows, err = db.Instance.Query(context.Background(), query)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
	}
	defer rows.Close()

	utensils := []gin.H{}
	for rows.Next() {
		var id int
		var name string
		if err = rows.Scan(&id, &name); err != nil {
			ctx.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		utensils = append(utensils, gin.H{"id": id, "name": name})
	}

	ctx.JSON(200, gin.H{"data": utensils})
}

type UpdateUtensilsDTO struct {
	HomeId     int   `json:"home_id" validate:"required"`
	UtensilIds []int `json:"utensil_ids" validate:"required,dive,gt=0"`
}

func UpdateUtensils(ctx *gin.Context) {
	var data UpdateUtensilsDTO
	if !validator.BindAndValidate(ctx, &data) {
		return
	}

	user := utils.GetUser(ctx)
	var ownerId int
	err := db.Instance.QueryRow(context.Background(), "SELECT user_id FROM homes WHERE id = $1", data.HomeId).Scan(&ownerId)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Home not found"})
		return
	}
	if user.Id != ownerId {
		ctx.JSON(403, gin.H{"error": "You can't edit this home"})
		return
	}

	tx, err := db.Instance.Begin(context.Background())
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), "DELETE FROM home_utensils_pivot WHERE home_id = $1", data.HomeId)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	if len(data.UtensilIds) > 0 {
		values := make([]interface{}, 0, len(data.UtensilIds)*2)
		placeholders := make([]string, 0, len(data.UtensilIds))
		for i, utensilId := range data.UtensilIds {
			values = append(values, data.HomeId, utensilId)
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		}

		query := "INSERT INTO home_utensils_pivot (home_id, utensil_id) VALUES " +
			strings.Join(placeholders, ", ") +
			" ON CONFLICT (home_id, utensil_id) DO NOTHING"

		_, err = tx.Exec(context.Background(), query, values...)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
	}

	if err = tx.Commit(context.Background()); err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Ok."})
}
