package controllers

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/matheusabido/kfofo-api/utils"
	"github.com/oracle/oci-go-sdk/v49/objectstorage"
)

type UploadImage struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func PostImage(ctx *gin.Context) {
	var form UploadImage

	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	file, err := form.File.Open()
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	defer file.Close()

	if !utils.IsImage(file) {
		ctx.JSON(400, gin.H{"error": "Imagem inválida."})
		return
	}

	namespace := os.Getenv("BUCKET_NAMESPACE")
	bucket := os.Getenv("BUCKET_NAME")
	objectName := "randomgeneratedname.png"

	request := objectstorage.PutObjectRequest{
		NamespaceName: &namespace,
		BucketName:    &bucket,
		ObjectName:    &objectName,
		PutObjectBody: file,
		ContentLength: &form.File.Size,
	}

	_, err = utils.GetClient().PutObject(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(200, gin.H{"message": "File stored"})
}
