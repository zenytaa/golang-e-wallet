package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	Id              uint
	SenderWallet    Wallet
	RecipientWallet Wallet
	Amount          decimal.Decimal
	SourceOfFund    SourceOfFund
	Description     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
