package model

import "bank/validator"

type SignUp struct {
	Name string `json:"name"`
	Surname string `json:"surname"`
	Phone string `json:"phone"`
	Password string `json:"password"`
	Confirm string `json:"confirm"`
	validator.Validator
}

type SignIn struct {
	Phone string `json:"phone"`
	Password string `json:"password"`
	validator.Validator
}

type User struct {
	ID int 
	IsAdmin bool
}

type UserInfo struct {
	Name string `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
	Phone string `json:"phone,omitempty"`
	ID int `json:"account_id,omitempty"`
}