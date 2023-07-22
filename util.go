package main

import (
	"math/rand"
)

func generateId(size int) string {
	var str string
	for i := 0; i < size; i++ {
		str += string(rune(rand.Intn(26) + 97))
	}
	return str
}

func fillWithNumbers(str string) string {
	for i := 0; i < len(str); i += 2 {
		str = str[:i] + string(rune(rand.Intn(10)+48)) + str[i:]
	}
	return str
}
