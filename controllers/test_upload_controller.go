package controllers

import (
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/matheusabido/kfofo-api/utils"
)

type UploadHomePicture struct {
	HomeId int
	File   *multipart.FileHeader `form:"file" binding:"required"`
}

func PostHomePicture(ctx *gin.Context) {
	var form UploadHomePicture

	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := utils.UploadFile(form.File, "")
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"message": "File stored"})
}
