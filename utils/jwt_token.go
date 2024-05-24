package utils

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/constant"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TokenCreateAndSign(userId uint, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"my-data": userId,
		"iss":     os.Getenv("ISS"),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		// and other claims, alternatively, you may want to explore on how to create custome claims
	})

	signed, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func VerifyExtractTokenClaim(ctx *gin.Context) (int, error) {
	tokenString := ExtractTokenFromHeader(ctx)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	}, jwt.WithIssuer(os.Getenv("ISS")),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		err = apperror.NewCustomError(http.StatusUnauthorized, constant.ResponseMsgErrorUnauthorized)
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["my-data"]), 10, 32)
		if err != nil {
			err = apperror.NewCustomError(http.StatusUnauthorized, constant.ResponseMsgErrorUnauthorized)
			return 0, err
		}
		return int(uid), nil
	}
	return 0, nil
}

func ExtractTokenFromHeader(ctx *gin.Context) string {
	tokenString := ctx.Query("token")
	if tokenString != "" {
		return tokenString
	}

	bearerToken := ctx.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
