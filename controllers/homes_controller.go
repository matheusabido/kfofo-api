package controllers

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/matheusabido/kfofo-api/db"
	"github.com/matheusabido/kfofo-api/utils"
	"github.com/matheusabido/kfofo-api/validator"
)

func GetHomes(ctx *gin.Context) {
	queryUser, _ := ctx.GetQuery("user")

	pageStr := ctx.DefaultQuery("page", "1")
	pageSize := 20

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		ctx.JSON(400, gin.H{"error": "Invalid page parameter"})
		return
	}

	whereString := ""
	var whereValues []any
	var allValues []any
	index := 1
	if queryUser != "" {
		whereString = " WHERE user_id = $1"
		queryUserId, err := strconv.Atoi(queryUser)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid user"})
			return
		}
		whereValues = append(whereValues, queryUserId)
		allValues = append(allValues, queryUserId)
		index++
	}

	offset := (page - 1) * pageSize

	totalChan := make(chan int)
	homesChan := make(chan []gin.H)

	go func() {
		var total int
		err = db.Instance.QueryRow(context.Background(), "SELECT COUNT(*) as total FROM homes"+whereString, whereValues...).Scan(&total)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, gin.H{"error": "Internal server error"})
			totalChan <- -1
			return
		}
		totalChan <- total
	}()

	allValues = append(allValues, pageSize)
	allValues = append(allValues, offset)

	go func() {
		query := "SELECT h.id, h.picture_path, h.user_id, u.name as user_name, h.address, h.city, h.description, h.cost_day, h.cost_week, h.cost_month, h.restriction_id, r.name as restriction_name, r.description as restriction_description, h.share_type_id, s.name as share_name, s.description as share_description FROM homes h INNER JOIN users u ON h.user_id = u.id INNER JOIN restrictions r ON h.restriction_id = r.id INNER JOIN share_types s ON h.share_type_id = s.id" + whereString + " ORDER BY h.id DESC LIMIT $" + strconv.Itoa(index) + " OFFSET $" + strconv.Itoa(index+1)
		index += 2
		rows, err := db.Instance.Query(context.Background(), query, allValues...)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, gin.H{"error": "Internal server error"})
			homesChan <- nil
			return
		}
		defer rows.Close()

		homes := make([]gin.H, 0, pageSize)
		for rows.Next() {
			var homeId int
			var picturePath string
			var userId int
			var userName string
			var address string
			var city string
			var description string
			var costDay float64
			var costWeek float64
			var costMonth float64
			var restrictionId int
			var restrictionName string
			var restrictionDesc string
			var shareTypeId int
			var shareName string
			var shareDesc string

			err = rows.Scan(&homeId, &picturePath, &userId, &userName, &address, &city, &description,
				&costDay, &costWeek, &costMonth, &restrictionId, &restrictionName, &restrictionDesc,
				&shareTypeId, &shareName, &shareDesc)
			if err != nil {
				ctx.JSON(500, gin.H{"error": "Internal server error"})
				homesChan <- nil
				return
			}

			homeData := gin.H{
				"id":                      homeId,
				"picture_path":            picturePath,
				"user_id":                 userId,
				"user_name":               userName,
				"address":                 address,
				"city":                    city,
				"description":             description,
				"cost_day":                costDay,
				"cost_week":               costWeek,
				"cost_month":              costMonth,
				"restriction_id":          restrictionId,
				"restriction_name":        restrictionName,
				"restriction_description": restrictionDesc,
				"share_type_id":           shareTypeId,
				"share_type_name":         shareName,
				"share_type_description":  shareDesc,
			}
			homes = append(homes, homeData)
		}
		homesChan <- homes
	}()

	total := <-totalChan
	homes := <-homesChan
	if total == -1 || homes == nil {
		return
	}

	ctx.JSON(200, gin.H{
		"total":    total,
		"homes":    homes,
		"lastPage": math.Ceil(float64(total) / float64(pageSize)),
	})
}

func GetHome(ctx *gin.Context) {
	idRaw := ctx.Param("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid home id"})
		return
	}

	query := `
		SELECT h.id, h.picture_path, h.user_id, u.name, h.address, h.city, h.description, h.cost_day, h.cost_week, h.cost_month,
			   h.restriction_id, r.name, r.description, r.icon,
			   h.share_type_id, s.name, s.description, s.icon
		FROM homes h
		INNER JOIN users u ON h.user_id = u.id
		INNER JOIN restrictions r ON h.restriction_id = r.id
		INNER JOIN share_types s ON h.share_type_id = s.id
		WHERE h.id = $1
	`
	row := db.Instance.QueryRow(context.Background(), query, id)

	var homeId int
	var picturePath string
	var userId int
	var userName string
	var address string
	var city string
	var description string
	var costDay float64
	var costWeek float64
	var costMonth float64
	var restrictionId int
	var restrictionName string
	var restrictionDesc string
	var restrictionIcon string
	var shareTypeId int
	var shareName string
	var shareDesc string
	var shareIcon string
	err = row.Scan(&homeId, &picturePath, &userId, &userName, &address, &city, &description, &costDay, &costWeek, &costMonth,
		&restrictionId, &restrictionName, &restrictionDesc, &restrictionIcon,
		&shareTypeId, &shareName, &shareDesc, &shareIcon)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(404, gin.H{"error": "Home not found"})
		return
	}

	ctx.JSON(200, gin.H{
		"id":                      homeId,
		"picture_path":            picturePath,
		"user_id":                 userId,
		"user_name":               userName,
		"address":                 address,
		"city":                    city,
		"description":             description,
		"cost_day":                costDay,
		"cost_week":               costWeek,
		"cost_month":              costMonth,
		"restriction_id":          restrictionId,
		"restriction_name":        restrictionName,
		"restriction_description": restrictionDesc,
		"restriction_icon":        restrictionIcon,
		"share_type_id":           shareTypeId,
		"share_type_name":         shareName,
		"share_type_description":  shareDesc,
		"share_type_icon":         shareIcon,
	})
}

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

type UpdateHomeDTO struct {
	Address       string  `json:"address" validate:"omitempty"`
	City          string  `json:"city" validate:"omitempty"`
	Description   string  `json:"description" validate:"omitempty"`
	CostDay       float64 `json:"cost_day" validate:"omitempty,min=1"`
	CostWeek      float64 `json:"cost_week" validate:"omitempty,min=1"`
	CostMonth     float64 `json:"cost_month" validate:"omitempty,min=1"`
	RestrictionId int     `json:"restriction_id" validate:"omitempty"`
	ShareTypeId   int     `json:"share_type_id" validate:"omitempty"`
}

func PutHome(ctx *gin.Context) {
	idRaw := ctx.Param("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid home id"})
		return
	}

	user := utils.GetUser(ctx)
	var ownerId int
	err = db.Instance.QueryRow(context.Background(), "SELECT user_id FROM homes WHERE id = $1", id).Scan(&ownerId)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Home not found"})
		return
	}

	if user.Id != ownerId {
		ctx.JSON(403, gin.H{"error": "You can't edit this home"})
		return
	}

	var data UpdateHomeDTO
	if !validator.BindAndValidate(ctx, &data) {
		return
	}

	index := 1
	var updates []string
	var values []any

	if data.Address != "" {
		updates = append(updates, "address = $"+strconv.Itoa(index))
		values = append(values, data.Address)
		index++
	}
	if data.City != "" {
		updates = append(updates, "city = $"+strconv.Itoa(index))
		values = append(values, data.City)
		index++
	}
	if data.Description != "" {
		updates = append(updates, "description = $"+strconv.Itoa(index))
		values = append(values, data.Description)
		index++
	}
	if data.CostDay != 0 {
		updates = append(updates, "cost_day = $"+strconv.Itoa(index))
		values = append(values, data.CostDay)
		index++
	}
	if data.CostWeek != 0 {
		updates = append(updates, "cost_week = $"+strconv.Itoa(index))
		values = append(values, data.CostWeek)
		index++
	}
	if data.CostMonth != 0 {
		updates = append(updates, "cost_month = $"+strconv.Itoa(index))
		values = append(values, data.CostMonth)
		index++
	}
	if data.RestrictionId != 0 {
		updates = append(updates, "restriction_id = $"+strconv.Itoa(index))
		values = append(values, data.RestrictionId)
		index++
	}
	if data.ShareTypeId != 0 {
		updates = append(updates, "share_type_id = $"+strconv.Itoa(index))
		values = append(values, data.ShareTypeId)
		index++
	}

	if len(updates) == 0 {
		ctx.JSON(200, gin.H{"message": "Home updated."})
		return
	}

	var builder strings.Builder
	builder.WriteString("UPDATE homes SET ")
	builder.WriteString(strings.Join(updates, ", "))
	builder.WriteString(" WHERE id = $" + strconv.Itoa(index))
	values = append(values, id)

	_, err = db.Instance.Exec(context.Background(), builder.String(), values...)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Home updated."})
}

func DeleteHome(ctx *gin.Context) {
	idRaw := ctx.Param("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid home id"})
		return
	}

	user := utils.GetUser(ctx)

	var ownerId int
	err = db.Instance.QueryRow(context.Background(), "SELECT user_id FROM homes WHERE id = $1", id).Scan(&ownerId)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Home not found"})
		return
	}

	if user.Id != ownerId {
		ctx.JSON(403, gin.H{"error": "You can't delete this home"})
		return
	}

	_, err = db.Instance.Exec(context.Background(), "DELETE FROM homes WHERE id = $1", id)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Home deleted."})
}
