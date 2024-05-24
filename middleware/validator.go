package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorValidation struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetErrorValidationMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return strings.ToLower(fe.Param()) + "should be field"
	case "email":
		return "wrong email format"
	case "gte":
		return strings.ToLower(fe.Field()) + " should be more than " + fe.Param()
	case "lte":
		return strings.ToLower(fe.Field()) + " should be less than " + fe.Param()
	case "min":
		return strings.ToLower(fe.Field()) + " should be " + fe.Param() + " character"
	}
	return "unknown error"

}

func Validation(ctx *gin.Context, err error) {
	var validationError validator.ValidationErrors
	if errors.As(err, &validationError) {
		out := make([]ErrorValidation, len(validationError))
		for i, fe := range validationError {
			out[i] = ErrorValidation{strings.ToLower(fe.Field()), GetErrorValidationMsg(fe)}
		}
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": out,
			},
		)
	}
}
