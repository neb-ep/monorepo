package entities

import "time"

type User struct {
	Username  string
	FirstName string
	LastName  string
	Password  string
	Email     string
	IsActive  bool
	CreatedAt time.Time
	DeletedAt time.Time
}

type UserCreds struct {
	UserId       int
	Username     string
	Email        string
	PasswordHash string
}

type AccessToken struct {
	Access          string
	Refresh         string
	AccessExpireIn  time.Time
	RefreshExpireIn time.Time
}
