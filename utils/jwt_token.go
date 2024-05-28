package utils

import (
	"assignment-go-rest-api/apperror"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthTokenProvider interface {
	TokenCreateAndSign(datas map[string]interface{}) (string, error)
	ParseAndVerify(signed string) (jwt.MapClaims, error)
	IsAuthorized(ctx *gin.Context) (bool, *uint, error)
	ExtractTokenFromHeader(ctx *gin.Context) (string, error)
}

type JwtProvider struct {
	config Config
}

func NewJwtProvider(config Config) AuthTokenProvider {
	return &JwtProvider{config: config}
}

func (j *JwtProvider) TokenCreateAndSign(datas map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": datas,
		"iss":  j.config.Issuer,
		"exp":  time.Now().Add(time.Duration(j.config.AccessTokenExp) * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (j *JwtProvider) ParseAndVerify(signed string) (jwt.MapClaims, error) {
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

func (j *JwtProvider) IsAuthorized(ctx *gin.Context) (bool, *uint, error) {
	token, err := j.ExtractTokenFromHeader(ctx)
	if err != nil {
		return false, nil, err
	}

	claims, err := j.ParseAndVerify(token)
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

func (j *JwtProvider) ExtractTokenFromHeader(ctx *gin.Context) (string, error) {
	bearerToken := ctx.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1], nil
	}
	return "", errors.New("token not found")
}
