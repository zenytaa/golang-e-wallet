package utils

import (
	"assignment-go-rest-api/apperror"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TokenCreateAndSign(datas map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": datas,
		"iss":  os.Getenv("ISS"),
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func ParseAndVerify(signed string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(signed, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.SecretKey), nil
	}, jwt.WithIssuer(j.config.Issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if err.Error() == "token has invalid claims: token is expired" {
			return nil, errors.New("token expired")
		}
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, apperror.ErrUnauthorized()
}

func IsAuthorized(ctx *gin.Context) (bool, *uint, error) {
	token, err := ExtractTokenFromHeader(ctx)
	if err != nil {
		return false, nil, err
	}

	claims, err := ParseAndVerify(token)
	if err != nil {
		return false, nil, err
	}
	dataMap := claims["data"]
	data, _ := dataMap.(map[string]interface{})

	id := uint(data["id"].(float64))

	if id != 0 {
		return true, &id, nil
	}

	return false, nil, err
}

func ExtractTokenFromHeader(ctx *gin.Context) (string, error) {
	bearerToken := ctx.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1], nil
	}
	return "", errors.New("token not found")
}
