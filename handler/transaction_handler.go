package handler

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/usecase"
	"assignment-go-rest-api/utils"
	"net/http"
	"strconv"

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
		Amount:          payload.Amount,
		SourceOfFund:    entity.SourceOfFund{Id: payload.SourceFundId},
		Description:     payload.Description,
		TransactionType: entity.TransactionTypes{Id: constant.TransferTypeId, Name: constant.TransferTypeName},
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
		Amount:          payload.Amount,
		SourceOfFund:    entity.SourceOfFund{Id: payload.SourceFundId},
		TransactionType: entity.TransactionTypes{Id: constant.TopUpTypeId},
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

func (h *TransactionHandler) GetListTransaction(ctx *gin.Context) {
	dataUserId, err := utils.GetDataFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	sortBy := ctx.Query("sortBy")
	sort := ctx.Query("sort")
	keyword := ctx.Query("keyword")

	limit := constant.DefaultLimit
	page := constant.DefaultPage

	limitStr := ctx.Query("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			ctx.Error(apperror.ErrInvalidIntegerInput())
			return
		}
		if limit == 0 {
			ctx.Error(apperror.ErrInvalidZeroLimitInput())
		}
	}

	pageStr := ctx.Query("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			ctx.Error(apperror.ErrInvalidIntegerInput())
			return
		}
		if page == 0 {
			ctx.Error(apperror.ErrInvalidZeroLimitInput())
		}
	}

	params := entity.TransactionParams{
		SortBy:  sortBy,
		Sort:    sort,
		Limit:   limit,
		Page:    page,
		Keyword: keyword,
	}

	tcs, pagination, err := h.TransactionUsecase.GetListTransaction(ctx, *dataUserId, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: constant.ResponseMsgGetListTransactionSuccess,
		Data:    dto.ToTransactionListResponses(tcs, *pagination),
	})

}
