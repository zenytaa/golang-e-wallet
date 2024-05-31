package dto

import (
	"assignment-go-rest-api/entity"
)

type WalletRequest struct {
	UserId uint `json:"user_id" binding:"required"`
}
type WalletResponse struct {
	UserName     string `json:"username"`
	WalletNumber string `json:"wallet_number"`
}

func ToWalletResponse(wallet entity.Wallet) *WalletResponse {
	return &WalletResponse{
		UserName:     wallet.User.Username,
		WalletNumber: wallet.WalletNumber,
	}
}
