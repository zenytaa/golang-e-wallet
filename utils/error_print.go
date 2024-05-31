package utils

import "log"

func IfErrorLogPrint(err error) {
	if err != nil {
		log.Println("error query: ", err)
	}
	return
}
