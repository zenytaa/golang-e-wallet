package handler

import (
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/middleware"
	"assignment-go-rest-api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandlerOpts struct {
	AuthUsecase usecase.AuthUsecase
}

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authHOpts *AuthHandlerOpts) *AuthHandler {
	return &AuthHandler{
		authUsecase: authHOpts.AuthUsecase,
	}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
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

func (h *AuthHandler) Login(ctx *gin.Context) {
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
