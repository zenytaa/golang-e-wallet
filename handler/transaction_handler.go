package handler

import (
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/middleware"
	"assignment-go-rest-api/usecase"
	"assignment-go-rest-api/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TransactionHandler interface {
	TopUp(ctx *gin.Context)
	Transfer(ctx *gin.Context)
	ListTransaction(ctx *gin.Context)
}

type TransactionHandlerImpl struct {
	transactionUsecase usecase.TransactionUsecase
	userUsecase        usecase.UserUsecase
	authUsecase        usecase.AuthUsecase
	walletUsecase      usecase.WalletUsecase
}

type TransactionHandlerConfig struct {
	TransactionUsecase usecase.TransactionUsecase
	UserUsecase        usecase.UserUsecase
	AuthUsecase        usecase.AuthUsecase
	WalletUsecase      usecase.WalletUsecase
}

func NewTransactionHandler(config *TransactionHandlerConfig) *TransactionHandlerImpl {
	return &TransactionHandlerImpl{
		transactionUsecase: config.TransactionUsecase,
		userUsecase:        config.UserUsecase,
		authUsecase:        config.AuthUsecase,
		walletUsecase:      config.WalletUsecase,
	}
}

func (h *TransactionHandlerImpl) ListTransaction(ctx *gin.Context) {
	query := &dto.ListTransactionQuery{}
	sortBy, sortByExist := utils.GetQuery(ctx, "sort_by")
	sort, sortExist := utils.GetQuery(ctx, "sort")
	search, searchExist := utils.GetQuery(ctx, "search")

	query.SortBy = sortBy
	query.SortByExist = sortByExist
	query.Sort = sort
	query.SortExist = sortExist
	query.Search = search
	query.SearchExist = searchExist

	query.SortBy = strings.ToLower(query.SortBy)
	query.Sort = strings.ToUpper(query.Sort)

	user := ctx.MustGet("user").(*entity.User)
	transactionResponses, err := h.transactionUsecase.ListTransaction(ctx, user.Id, query)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: constant.ResponseMsgGetListTransactionSuccess,
		Data:    transactionResponses,
	})
}

func (h *TransactionHandlerImpl) TopUp(ctx *gin.Context) {
	topUpRequest := dto.TopUpCreateRequest{}

	err := ctx.ShouldBindJSON(&topUpRequest)
	if err != nil {
		middleware.Validation(ctx, err)
		return
	}

	user := ctx.MustGet("user").(*entity.User)
	topUpRequest.UserId = user.Id

	topUpResponse, err := h.transactionUsecase.TopUp(ctx, &topUpRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: constant.ResponseMsgTopUpSucces,
		Data:    topUpResponse,
	})
}

func (h *TransactionHandlerImpl) Transfer(ctx *gin.Context) {
	transferRequest := dto.TransferCreateRequest{}
	fmt.Println(transferRequest)

	err := ctx.ShouldBindJSON(&transferRequest)
	if err != nil {
		middleware.Validation(ctx, err)
		return
	}

	user := ctx.MustGet("user").(*entity.User)
	transferRequest.SenderUserId = user.Id

	transferResponse, err := h.transactionUsecase.Transfer(ctx, &transferRequest)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.WebResponse{
		Code:    http.StatusOK,
		Message: constant.ResponseMsgTransferSucces,
		Data:    transferResponse,
	})
}
