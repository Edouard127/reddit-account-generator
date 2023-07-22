package main

import "math/rand"

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewRandomUser() User {
	return User{
		Email:    GetEmail(),
		Username: fillWithNumbers(generateId(rand.Intn(12) + 5)),
		Password: fillWithNumbers(generateId(rand.Intn(34))),
	}
}
