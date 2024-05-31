package middleware

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(config utils.Config) func(*gin.Context) {
	return func(ctx *gin.Context) {
		authorized, data, err := utils.NewJwtProvider(config).IsAuthorized(ctx)
		if !authorized && err != nil && data == nil {
			if err.Error() == "token expired" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, apperror.ErrTokenExpired())
				return
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, apperror.ErrUnauthorized())
			return
		}
		ctx.Set("data", data)
		ctx.Next()
	}
}
