package main

import (
	"math/rand"
)

func generateId(size int) string {
	var password string
	for i := 0; i < size; i++ {
		password += string(rune(rand.Intn(26) + 97))
	}
	return password
}
