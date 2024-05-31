package handler

import (
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandlerImpl struct {
	walletUsecase usecase.WalletUsecase
}

func NewUserHandler(walletUsecase usecase.WalletUsecase) *UserHandlerImpl {
	return &UserHandlerImpl{
		walletUsecase: walletUsecase,
	}
}

func (h *UserHandlerImpl) GetProfile(ctx *gin.Context) {
	user := ctx.MustGet("user").(*entity.User)
	walletRequest := dto.WalletRequest{}

	walletRequest.UserId = user.Id
	wallet, err := h.walletUsecase.GetWalletByUserId(ctx, &walletRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	userDetailResponse := dto.ToUserDetailResponse(*user, *wallet)

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: constant.ResponseMsgGetProfileSuccsess,
		Data:    userDetailResponse,
	})
}
