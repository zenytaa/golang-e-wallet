package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) ([]byte, error) {
	var config Config
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), config.HashCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func CheckPassword(pwd string, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(pwd))
	if err != nil {
		return false, err
	}
	return true, nil
}
