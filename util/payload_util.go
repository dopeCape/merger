package util

import (
	"os"
	"strconv"
)

func PayloadSizeChecker(payload string) (bool, error) {
	size := len(payload)
	maxSize := os.Getenv("MAX_PAYLOAD")
	maxSizeInt, err := strconv.Atoi(maxSize)
	if err != nil {
		return false, err
	}
	if size > (maxSizeInt * 1024) {
		return false, nil
	}
	return true, nil

}
