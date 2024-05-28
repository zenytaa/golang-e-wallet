package dto

import (
	"assignment-go-rest-api/entity"

	"github.com/shopspring/decimal"
)

type WalletRequest struct {
	UserId uint `json:"user_id" binding:"required"`
}
type WalletResponse struct {
	UserName     UserResponse    `json:"user"`
	WalletNumber string          `json:"wallet_number"`
	Balance      decimal.Decimal `json:"balance"`
}

func ToWalletResponse(wallet entity.Wallet) *WalletResponse {
	return &WalletResponse{
		UserName:     *ToUserResponse(wallet.User),
		WalletNumber: wallet.WalletNumber,
		Balance:      wallet.Balance,
	}
}
