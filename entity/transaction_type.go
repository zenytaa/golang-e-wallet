package entity

import "time"

type TransactionTypes struct {
	Id        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
