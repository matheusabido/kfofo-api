package controllers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/matheusabido/kfofo-api/db"
	"github.com/matheusabido/kfofo-api/utils"
	"github.com/matheusabido/kfofo-api/validator"
)

type StoreHomeDTO struct {
	UserId        int     `json:"user_id" validate:"required"`
	Address       string  `json:"address" validate:"required"`
	City          string  `json:"city" validate:"required"`
	Description   string  `json:"description" validate:"required"`
	CostDay       float64 `json:"cost_day" validate:"required,min=1"`
	CostWeek      float64 `json:"cost_week" validate:"omitempty,min=1"`
	CostMonth     float64 `json:"cost_month" validate:"omitempty,min=1"`
	RestrictionId int     `json:"restriction_id" validate:"required"`
	ShareTypeId   int     `json:"share_type_id" validate:"required"`
}

func PostHome(ctx *gin.Context) {
	var data StoreHomeDTO
	if !validator.BindAndValidate(ctx, &data) {
		return
	}

	user := utils.GetUser(ctx)

	if data.UserId != user.Id {
		ctx.JSON(403, gin.H{"error": "You can't create a home for another user."})
		return
	}

	var userId int
	var userName string
	var restrictionName string
	var restrictionDescription string
	var shareName string
	var shareDescription string
	row := db.Instance.QueryRow(context.Background(), "SELECT u.id as user_id, u.name as user_name, r.name as restriction_name, r.description as restriction_desc, s.name as share_name, s.description as share_desc FROM users u INNER JOIN restrictions r ON r.id = $2 INNER JOIN share_types s ON s.id = $3 WHERE u.id = $1", data.UserId, data.RestrictionId, data.ShareTypeId)
	err := row.Scan(&userId, &userName, &restrictionName, &restrictionDescription, &shareName, &shareDescription)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid IDs"})
		return
	}

	var homeId int
	err = db.Instance.QueryRow(context.Background(), "INSERT INTO homes (user_id, address, city, description, cost_day, cost_week, cost_month, restriction_id, share_type_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id", data.UserId, data.Address, data.City, data.Description, data.CostDay, data.CostWeek, data.CostMonth, data.RestrictionId, data.ShareTypeId).Scan(&homeId)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{
		"id":                      homeId,
		"user_id":                 user.Id,
		"user_name":               userName,
		"address":                 data.Address,
		"city":                    data.City,
		"description":             data.Description,
		"cost_day":                data.CostDay,
		"cost_week":               data.CostWeek,
		"cost_month":              data.CostMonth,
		"restriction_id":          data.RestrictionId,
		"restriction_name":        restrictionName,
		"restriction_description": restrictionDescription,
		"share_type_id":           data.ShareTypeId,
		"share_type_name":         shareName,
		"share_type_description":  shareDescription,
	})
}
