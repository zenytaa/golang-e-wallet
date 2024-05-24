package entity

import "time"

type PasswordReset struct {
	Id        uint
	UserId    uint
	Token     string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
