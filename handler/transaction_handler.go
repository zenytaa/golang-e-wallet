package handler

import (
	"assignment-go-rest-api/usecase"

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

}
