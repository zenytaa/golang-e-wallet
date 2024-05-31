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
	containsFe := fe.Param()
	if fe.Param() == `!#$%&'()*+,-./:"\;<=>?@[]^_{|}~` {
		containsFe = "special character"
	}
	if fe.Param() == "abcdefghijklmnopqrstuvwxyz" {
		containsFe = "lowercase"
	}
	if fe.Param() == "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		containsFe = "uppercase"
	}
	if fe.Param() == "1234567890" {
		containsFe = "number"
	}
	if fe.Param() == ` []{|}"\%~#=&<>?/.` {
		containsFe = fe.Param()[1:]
	}
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "use valid email address"
	case "max":
		return "maximum value for this field is " + fe.Param() + " characters"
	case "min":
		return "minimum value for this field is " + fe.Param() + " characters"
	case "lte":
		return "maximum value for this field is " + fe.Param()
	case "gte":
		return "minimum value for this field is " + fe.Param()
	case "len=0|url":
		return "this field just receive empty string or valid url"
	case "len=0|e164":
		return "this field just receive empty string or e.164 phone number format (ex: +1123456789)"
	case "e164":
		return "use e.164 phone number format (ex: +1123456789)"
	case "datetime":
		return "use valid date format: yyyy-mm-dd (ex: 2026-01-02)"
	case "excludes":
		return "this field cannot contain space"
	case "containsany":
		return "this field must contain at least 1 " + containsFe
	case "excludesall":
		return "this field cannot contain space and any of this characters " + containsFe
	case "lowercase":
		return "this field cannot contain uppercase character"
	case "alphanum":
		return "this field just receive alphanumeric"
	default:
		return "unknown error"
	}

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
