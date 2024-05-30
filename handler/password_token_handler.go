package handler

import (
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/middleware"
	"assignment-go-rest-api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PasswordResetHandlerOpts struct {
	PasswordResetUsecase usecase.PasswordResetUsecase
}

type PasswordResetHandler struct {
	passwordResetUsecase usecase.PasswordResetUsecase
}

func NewPasswordResetHandler(resetHOpts *PasswordResetHandlerOpts) *PasswordResetHandler {
	return &PasswordResetHandler{resetHOpts.PasswordResetUsecase}
}

func (h *PasswordResetHandler) ForgotPassword(ctx *gin.Context) {
	payload := dto.ForgotPasswordRequest{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	forgotPassword, err := h.passwordResetUsecase.ForgotPassword(ctx, payload.Email)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: constant.ResponseMsgForgotPasswordSuccess,
		Data:    dto.ToForgotPasswordResponse(*forgotPassword),
	})
}

func (h *PasswordResetHandler) PasswordReset(ctx *gin.Context) {
	payload := dto.ResetPasswordRequest{}

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		middleware.Validation(ctx, err)
		return
	}

	err = h.passwordResetUsecase.ResetPassword(ctx, payload.Token, payload.Password)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: constant.ResponseMsgResetPasswordSuccess,
		Data:    nil,
	})
}
