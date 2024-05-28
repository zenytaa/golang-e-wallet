package handler

import (
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/usecase"
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

	tc := entity.Transaction{
		User: entity.User{Id: payload.SenderUserId},
		SenderWallet: entity.Wallet{
			User: entity.User{Id: payload.SenderUserId},
		},
		RecipientWallet: entity.Wallet{
			User:         entity.User{},
			WalletNumber: payload.RecipientWalletNumber,
		},
		Amount:         payload.Amount,
		SourceOfFundId: payload.SourceFundId,
		Description:    payload.Description,
	}

	transaferResponse, err := h.TransactionUsecase.TransferWithTransactor(ctx, tc)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Message: constant.ResponseMsgTransferSucces,
		Data:    dto.ToTransferResponse(*transaferResponse),
	})
}
