package validator

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func SetupValidator() {
	validate = validator.New()
}

func Bind(ctx *gin.Context, data any) bool {
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return false
	}
	return true
}

func Validate(ctx *gin.Context, data interface{}) bool {
	if err := validate.Struct(data); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return false
	}
	return true
}

func BindAndValidate(ctx *gin.Context, data interface{}) bool {
	return Bind(ctx, &data) && Validate(ctx, data)
}
