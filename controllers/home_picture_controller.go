package controllers

import (
	"context"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/matheusabido/kfofo-api/db"
	"github.com/matheusabido/kfofo-api/utils"
)

type UploadHomePicture struct {
	HomeId int                   `form:"home_id" binding:"required"`
	File   *multipart.FileHeader `form:"file" binding:"required"`
}

func GetHomePicture(ctx *gin.Context) {
	path, _ := ctx.GetQuery("path")
	if path == "" {
		ctx.JSON(404, gin.H{"error": "Not found"})
		return
	}

	var returnType string
	if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
		returnType = "image/jpeg"
	} else {
		returnType = "image/png"
	}

	file, err := utils.GetFile(path)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Not found"})
		return
	}

	ctx.Data(200, returnType, file)
}

func PostHomePicture(ctx *gin.Context) {
	var form UploadHomePicture

	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	claims := utils.GetClaims(ctx)

	var userId int
	err := db.Instance.QueryRow(context.Background(), "SELECT user_id FROM homes WHERE id = $1", form.HomeId).Scan(&userId)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Home not found"})
		return
	}

	if userId != claims.Id {
		ctx.JSON(400, gin.H{"error": "You can't add a picture to this home"})
		return
	}

	err, path := utils.UploadFile(form.File, "homes/"+strconv.Itoa(form.HomeId))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	_, err = db.Instance.Exec(context.Background(), "UPDATE homes SET picture_path = $1 WHERE id = $2", path, form.HomeId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Picture updated"})
}
