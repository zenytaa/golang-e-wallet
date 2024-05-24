package entity

import "time"

type User struct {
	Id        uint
	Email     string
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
