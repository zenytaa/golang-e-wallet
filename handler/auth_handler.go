package handler

import (
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/middleware"
	"assignment-go-rest-api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandlerImpl struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandlerImpl) Register(ctx *gin.Context) {
	authRequest := dto.AuthRegisterRequest{}

	err := ctx.ShouldBindJSON(&authRequest)
	if err != nil {
		middleware.Validation(ctx, err)
		return
	}

	err = h.authUsecase.RegisterWithInTransactor(ctx, authRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: constant.ResponseMsgRegisterSuccess,
	})
}

func (h *AuthHandlerImpl) Login(ctx *gin.Context) {
	authRequest := dto.AuthLoginRequest{}

	err := ctx.ShouldBindJSON(&authRequest)
	if err != nil {
		middleware.Validation(ctx, err)
		return
	}

	token, err := h.authUsecase.Login(ctx, authRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: constant.ResponseMsgLoginSuccess,
		Data: gin.H{
			"token": token,
		},
	})
}

func (h *AuthHandlerImpl) ForgotPassword(ctx *gin.Context) {
	forgotPasswordRequest := dto.ForgotPasswordRequest{}

	err := ctx.ShouldBindJSON(&forgotPasswordRequest)
	if err != nil {
		middleware.Validation(ctx, err)
		return
	}

	forgotPassword, err := h.authUsecase.ForgotPassword(ctx, forgotPasswordRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: constant.ResponseMsgForgotPasswordSuccess,
		Data:    forgotPassword,
	})
}

func (h *AuthHandlerImpl) ResetPassword(ctx *gin.Context) {
	resetPasswordRequest := dto.ResetPasswordRequest{}

	err := ctx.ShouldBindJSON(&resetPasswordRequest)
	if err != nil {
		middleware.Validation(ctx, err)
		return
	}

	resetPassword, err := h.authUsecase.ResetPassword(ctx, resetPasswordRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: constant.ResponseMsgResetPasswordSuccess,
		Data:    resetPassword,
	})
}
