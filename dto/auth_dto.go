package dto

import (
	"assignment-go-rest-api/entity"
	"time"

	"github.com/shopspring/decimal"
)

type AuthRegisterResponse struct {
	Id            uint   `json:"id"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	WalletNumbber string `json:"wallet_number"`
}

type AuthLoginResponse struct {
	Message string `json:"message" binding:"required"`
	Token   string `json:"token" binding:"required"`
}

type AuthRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type AuthLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	Password        string `json:"new_password" binding:"required,gte=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,gte=8"`
}

type ResetPasswordResponse struct {
	Id     uint   `json:"id"`
	UserId uint   `json:"user_id"`
	Email  string `json:"email"`
}

type ListTransactionResponse struct {
	Id              uint            `json:"id"`
	SenderWallet    WalletResponse  `json:"sender_wallet"`
	RecipientWallet WalletResponse  `json:"recipient_wallet"`
	Amount          decimal.Decimal `json:"amount"`
	SourceFund      string          `json:"source_of_fund"`
	Description     string          `json:"description"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func ToAuthRegisterResponse(user entity.User, wallet entity.Wallet) *AuthRegisterResponse {
	return &AuthRegisterResponse{
		Id:            user.Id,
		Email:         user.Email,
		Username:      user.Username,
		WalletNumbber: wallet.WalletNumber,
	}
}

func ToForgotPasswordResponse(user entity.User, pwdReset entity.PasswordReset) *ForgotPasswordResponse {
	return &ForgotPasswordResponse{
		Email: user.Email,
		Token: pwdReset.Token,
	}
}

func ToResetPasswordResponse(user entity.User, pwdReset entity.PasswordReset) *ResetPasswordResponse {
	return &ResetPasswordResponse{
		Id:     pwdReset.Id,
		UserId: user.Id,
		Email:  user.Email,
	}
}
