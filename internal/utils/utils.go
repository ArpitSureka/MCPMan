package utils

import (
	"fmt"
)


func CheckError(err error) {
	if err != nil {
		fmt.Printf("Error, Failed to run: %v\n", err)
	}
}