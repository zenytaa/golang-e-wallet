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

type TransactionParams struct {
	SortBy  string
	Sort    string
	Limit   int
	Page    int
	Keyword string
}
