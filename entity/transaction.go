package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	Id              uint
	Amount          decimal.Decimal
	Description     string
	SenderWallet    Wallet
	RecipientWallet Wallet
	SourceOfFund    SourceOfFund
	TransactionType TransactionTypes
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
