package utils

import "log"

func CheckError(error error) {
	if error != nil {
		log.Fatal(error)
	}
}
