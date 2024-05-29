package handler

import (
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/usecase"
	"assignment-go-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandlerOpts struct {
	TransactionUsecase usecase.TransactionUsecase
}

type TransactionHandler struct {
	TransactionUsecase usecase.TransactionUsecase
}

func NewTransactionHandler(transHOpts *TransactionHandlerOpts) *TransactionHandler {
	return &TransactionHandler{TransactionUsecase: transHOpts.TransactionUsecase}
}

func (h *TransactionHandler) Transfer(ctx *gin.Context) {
	var payload dto.TransferCreateRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	dataUserId, err := utils.GetDataFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	tc := entity.Transaction{
		SenderWallet: entity.Wallet{
			User: entity.User{Id: *dataUserId},
		},
		RecipientWallet: entity.Wallet{
			User:         entity.User{},
			WalletNumber: payload.RecipientWalletNumber,
		},
		Amount:       payload.Amount,
		SourceOfFund: entity.SourceOfFund{Id: payload.SourceFundId},
		Description:  payload.Description,
	}

	transaferResponse, err := h.TransactionUsecase.TransferWithTransactor(ctx, tc)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: constant.ResponseMsgTransferSucces,
		Data:    dto.ToTransferResponse(*transaferResponse),
	})
}

func (h *TransactionHandler) TopUp(ctx *gin.Context) {
	var payload dto.TopUpCreateRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	dataUserId, err := utils.GetDataFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	tc := entity.Transaction{
		SenderWallet: entity.Wallet{
			User: entity.User{Id: *dataUserId},
		},
		RecipientWallet: entity.Wallet{
			User: entity.User{Id: *dataUserId},
		},
		Amount:       payload.Amount,
		SourceOfFund: entity.SourceOfFund{Id: payload.SourceFundId},
	}

	topUpResponse, err := h.TransactionUsecase.TopUpWithTransactor(ctx, tc)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: constant.ResponseMsgTopUpSucces,
		Data:    dto.ToTopUpResponse(*topUpResponse),
	})
}
