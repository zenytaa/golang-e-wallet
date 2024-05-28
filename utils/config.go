package utils

import (
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DbHost         string
	DbUser         string
	DbPassword     string
	DbName         string
	DbPort         string
	DbSSL          string
	DbTimeZone     string
	SecretKey      string
	HashCost       int
	Issuer         string
	AccessTokenExp int
}

func ConfigInit() (Config, error) {
	env, err := godotenv.Read()
	if err != nil {
		return Config{}, err
	}

	hashCost, err := strconv.Atoi(env["HASH_COST"])
	if err != nil {
		return Config{}, err
	}

	accTokenExp, err := strconv.Atoi(env["ACCESS_TOKEN_EXP"])
	if err != nil {
		return Config{}, err
	}

	return Config{
		DbHost:         env["DB_HOST"],
		DbUser:         env["DB_USER"],
		DbPassword:     env["DB_PASSWORD"],
		DbName:         env["DB_DBNAME"],
		DbPort:         env["DB_PORT"],
		DbSSL:          env["DB_SSLMODE"],
		DbTimeZone:     env["DB_TIMEZONE"],
		SecretKey:      env["SECRET_KEY"],
		HashCost:       hashCost,
		Issuer:         env["ISS"],
		AccessTokenExp: accTokenExp,
	}, nil
}
