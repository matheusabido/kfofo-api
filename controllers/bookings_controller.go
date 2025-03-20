package controllers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/matheusabido/kfofo-api/db"
	"github.com/matheusabido/kfofo-api/utils"
	"github.com/matheusabido/kfofo-api/validator"
)

func GetBookings(ctx *gin.Context) {
	user := utils.GetUser(ctx)
	query := `
		SELECT b.id, b.home_id, b.from_date, b.to_date, b.cost_per_cycle, h.address, h.city, h.picture_path
		FROM bookings b
		INNER JOIN homes h ON h.id = b.home_id
		WHERE b.user_id = $1
		ORDER BY id DESC
	`
	rows, err := db.Instance.Query(context.Background(), query, user.Id)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	bookings := make([]gin.H, 0)
	for rows.Next() {
		var address, city, picturePath string
		var id, homeId int
		var fromDate, toDate time.Time
		var costPerCycle float64

		err = rows.Scan(&id, &homeId, &fromDate, &toDate, &costPerCycle, &address, &city, &picturePath)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		booking := gin.H{
			"id":             id,
			"home_id":        homeId,
			"from_date":      fromDate.Format("2006-01-02"),
			"to_date":        toDate.Format("2006-01-02"),
			"cost_per_cycle": costPerCycle,
			"address":        address,
			"city":           city,
			"picture_path":   picturePath,
		}
		bookings = append(bookings, booking)
	}

	ctx.JSON(200, bookings)
}

type StoreBookingDTO struct {
	HomeId   int    `json:"home_id" validate:"required"`
	FromDate string `json:"from_date" validate:"required"`
	ToDate   string `json:"to_date" validate:"required"`
}

func PostBooking(ctx *gin.Context) {
	var data StoreBookingDTO
	if !validator.BindAndValidate(ctx, &data) {
		return
	}

	user := utils.GetUser(ctx)

	fromDate, err := time.Parse("2006-01-02", data.FromDate)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid fromDate."})
		return
	}
	toDate, err := time.Parse("2006-01-02", data.ToDate)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid toDate."})
		return
	}

	var homeId int
	var costDay float64
	err = db.Instance.QueryRow(context.Background(), "SELECT id, cost_day FROM homes WHERE id = $1", data.HomeId).Scan(&homeId, &costDay)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(404, gin.H{"error": "Home not found"})
		return
	}

	var bookingId int
	insertQuery := `
		INSERT INTO bookings (user_id, home_id, from_date, to_date, payment_type, cost_per_cycle)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err = db.Instance.QueryRow(context.Background(), insertQuery, user.Id, data.HomeId, fromDate, toDate, 0, costDay).Scan(&bookingId)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{
		"id":             bookingId,
		"user_id":        user.Id,
		"home_id":        data.HomeId,
		"from_date":      fromDate.Format("2006-01-02"),
		"to_date":        toDate.Format("2006-01-02"),
		"cost_per_cycle": costDay,
	})
}

func DeleteBooking(ctx *gin.Context) {
	bookingIdStr := ctx.Param("id")
	bookingId, err := strconv.Atoi(bookingIdStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid booking id"})
		return
	}

	user := utils.GetUser(ctx)
	var ownerId int
	err = db.Instance.QueryRow(context.Background(), "SELECT user_id FROM bookings WHERE id = $1", bookingId).Scan(&ownerId)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Booking not found"})
		return
	}

	if ownerId != user.Id {
		ctx.JSON(403, gin.H{"error": "You can't delete this booking."})
		return
	}

	_, err = db.Instance.Exec(context.Background(), "DELETE FROM bookings WHERE id = $1", bookingId)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Booking deleted."})
}
