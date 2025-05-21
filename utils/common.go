package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

func NewSlug(input string) string {
	slug := regexp.MustCompile("[^a-zA-Z0-9+]").ReplaceAllString(input, "-")
	slug = strings.ToLower(strings.Trim(slug, "-"))
	return fmt.Sprintf("%s-%s", slug, RandomString(7))
}

func RandomString(length int) string {
	char := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)

	for i := range length {
		b[i] = char[rand.Intn(len(char))]
	}

	return string(b)
}
