package main

import "github.com/google/uuid"

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUser() User {
	return User{
		Email:    GetEmail(),
		Username: uuid.New().String(),
		Password: uuid.New().String(),
	}
}

func (u User) MarshalJSON() ([]byte, error) {
	return []byte(`{"email":"` + u.Email + `","username":"` + u.Username + `","password":"` + u.Password + `"}`), nil
}
