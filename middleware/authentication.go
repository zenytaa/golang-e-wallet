package middleware

import (
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/usecase"
	"assignment-go-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(userUsecase usecase.UserUsecase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, err := utils.VerifyExtractTokenClaim(ctx)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		param := &dto.UserRequestParam{
			UserId: uid,
		}
		user, err := userUsecase.GetUser(ctx, param)

		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": constant.ResponseMsgErrorUnauthorized,
			})
			return
		}
		ctx.Set("user", user)
		ctx.Next()
	}
}
