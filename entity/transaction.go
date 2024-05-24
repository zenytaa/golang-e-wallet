package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	Id                uint
	UserId            uint
	SenderWalletId    uint
	RecipientWalletId uint
	RecipientUsername string
	Amount            decimal.Decimal
	SourceOfFundId    uint
	Description       string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}
