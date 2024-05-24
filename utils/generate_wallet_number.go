package utils

import "fmt"

func GenerateWalletNumber(userId uint) string {
	if userId < 10 {
		return fmt.Sprintf("999000000000%d", userId)
	}
	if userId < 100 {
		return fmt.Sprintf("99900000000%d", userId)
	} else {
		return fmt.Sprintf("9990000000%d", userId)
	}
}
