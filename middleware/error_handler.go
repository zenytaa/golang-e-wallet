package middleware

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/constant"
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return errorHandler(gin.ErrorTypeAny)
}

func errorHandler(errType gin.ErrorType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		detectedErrors := ctx.Errors.ByType(errType)

		log.Print("App error")
		if len(detectedErrors) > 0 {
			err := detectedErrors[0].Err
			var parsedError *apperror.CustomError
			switch err.(type) {
			case *apperror.CustomError:
				parsedError = err.(*apperror.CustomError)
			default:
				parsedError = &apperror.CustomError{
					Message: constant.ResponseMsgErrorInternalServer,
				}
			}
			ctx.IndentedJSON(parsedError.Code, parsedError)
			ctx.Abort()
			return
		}
	}
}
