package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	Id              uint
	User            User
	SenderWallet    Wallet
	RecipientWallet Wallet
	Amount          decimal.Decimal
	SourceOfFundId  uint
	Description     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
