package service

import (
	"math/rand"
	"regexp"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var alphanumericRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateShortCode(length int) string {
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func ValidateCustomCode(code string) bool {
	length := len(code)
	if length < 3 || length > 20 {
		return false
	}
	return alphanumericRegex.MatchString(code)
}
