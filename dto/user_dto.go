package dto

import (
	"assignment-go-rest-api/entity"

	"github.com/shopspring/decimal"
)

type UserRequestParam struct {
	UserId int `json:"id" binding:"required"`
}

type UserDetailResponse struct {
	UserId       uint                      `json:"user_id"`
	Email        string                    `json:"email"`
	Username     string                    `json:"username"`
	WalletId     uint                      `json:"wallet_id"`
	WalletNumber string                    `json:"wallet_number"`
	Balance      decimal.Decimal           `json:"balance"`
	Transactions []ListTransactionResponse `json:"transactions"`
	Income       decimal.Decimal           `json:"income"`
	Expense      decimal.Decimal           `json:"expense"`
}

type UserResponse struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
}

func ToUserResponse(user entity.User) *UserResponse {
	return &UserResponse{UserId: user.Id, Username: user.Username}
}

func ToUserDetailResponse(user entity.User, wallet entity.Wallet) *UserDetailResponse {
	return &UserDetailResponse{
		UserId:       user.Id,
		Email:        user.Email,
		Username:     user.Username,
		WalletId:     wallet.Id,
		WalletNumber: wallet.WalletNumber,
		Balance:      wallet.Balance,
	}
}
