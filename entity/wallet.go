package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	Id           uint
	UserId       uint
	WalletNumber string
	Balance      decimal.Decimal
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
