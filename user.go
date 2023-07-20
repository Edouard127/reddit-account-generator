package main

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewRandomUser() User {
	return User{
		Email:    GetEmail(),
		Username: generateId(18),
		Password: generateId(18),
	}
}
