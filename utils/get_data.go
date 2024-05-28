package utils

import (
	"assignment-go-rest-api/apperror"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetQuery(c *gin.Context, queryKey string) (string, bool) {
	return c.GetQuery(queryKey)
}

func GetDataFromContext(ctx *gin.Context) (*uint, error) {
	data, isExist := ctx.Get("data")
	if !isExist {
		return nil, &apperror.CustomError{
			Code:    http.StatusNotFound,
			Message: "context not found",
		}
	}

	dataUint := data.(*uint)

	return dataUint, nil
}
