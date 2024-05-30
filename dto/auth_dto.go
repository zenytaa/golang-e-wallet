package dto

import (
	"assignment-go-rest-api/entity"
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
	Token string `json:"token"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"new_password" binding:"required,gte=8"`
}

func ToAuthRegisterResponse(user entity.User, wallet entity.Wallet) *AuthRegisterResponse {
	return &AuthRegisterResponse{
		Id:            user.Id,
		Email:         user.Email,
		Username:      user.Username,
		WalletNumbber: wallet.WalletNumber,
	}
}

func ToForgotPasswordResponse(pwdReset entity.PasswordReset) *ForgotPasswordResponse {
	return &ForgotPasswordResponse{
		Token: pwdReset.Token,
	}
}
