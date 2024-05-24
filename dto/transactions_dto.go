package dto

import (
	"assignment-go-rest-api/entity"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type TopUpCreateRequest struct {
	Amount       decimal.Decimal `json:"amount" binding:"required"`
	SourceFundId uint            `json:"source_of_fund_id" binding:"required"`
	UserId       uint
}

type TopUpResponse struct {
	Id                uint            `json:"id" binding:"required"`
	SenderWalletId    uint            `json:"sender_wallet_id" binding:"required"`
	RecipientWalletId uint            `json:"recipient_wallet_id" binding:"required"`
	Amount            decimal.Decimal `json:"amount" binding:"required"`
	SourceOfFundId    uint            `json:"source_of_fund_id" binding:"required"`
	Description       string          `json:"description" binding:"required,len=35"`
	CreatedAt         time.Time       `json:"created_at" binding:"required"`
}

type TransferCreateRequest struct {
	SenderUserId          uint
	Amount                decimal.Decimal `json:"amount" binding:"required"`
	RecipientWalletNumber string          `json:"recipient_wallet_number" binding:"required,gte=0,lte=13"`
	SourceFundId          uint            `json:"source_of_fund_id" binding:"required"`
	Description           string          `json:"description" binding:"lte=35"`
}

type TransferResponse struct {
	Id                uint            `json:"id" binding:"required"`
	SenderWalletId    uint            `json:"sender_wallet_id" binding:"required"`
	RecipientWalletId uint            `json:"recipient_wallet_id" binding:"required"`
	RecipientUsername string          `json:"recipient_username" binding:"required"`
	Amount            decimal.Decimal `json:"amount" binding:"required,gte=50000,lte=10000000"`
	SourceFundId      uint            `json:"source_of_fund_id" binding:"required"`
	Description       string          `json:"description" binding:"len=35"`
	CreatedAt         time.Time       `json:"created_at" binding:"required"`
	UpdatedAt         time.Time       `json:"updated_at" binding:"required"`
}

type ListTransactionQuery struct {
	Sort        string `form:"sort"`
	SortExist   bool
	Search      string `form:"search"`
	SearchExist bool
	SortBy      string `form:"sort_by"`
	SortByExist bool
	Limit       int `form:"limit"`
	Page        int `form:"page"`
}

func ToTopUpResponse(ts entity.Transaction) *TopUpResponse {
	return &TopUpResponse{
		Id:                ts.Id,
		SenderWalletId:    ts.SenderWalletId,
		RecipientWalletId: ts.RecipientWalletId,
		Amount:            ts.Amount,
		SourceOfFundId:    ts.SourceOfFundId,
		Description:       ts.Description,
		CreatedAt:         ts.CreatedAt,
	}
}

func ToTransferResponse(ts entity.Transaction, user entity.User) *TransferResponse {
	return &TransferResponse{
		Id:                ts.Id,
		SenderWalletId:    ts.SenderWalletId,
		RecipientWalletId: ts.RecipientWalletId,
		RecipientUsername: user.Username,
		Amount:            ts.Amount,
		SourceFundId:      ts.SourceOfFundId,
		Description:       ts.Description,
		CreatedAt:         ts.CreatedAt,
		UpdatedAt:         ts.UpdatedAt,
	}
}

func ToListTransactionResponse(ts entity.Transaction) *ListTransactionResponse {
	return &ListTransactionResponse{
		Id:                ts.Id,
		SenderWalletId:    ts.SenderWalletId,
		RecipientWalletId: ts.RecipientWalletId,
		RecipientUsername: ts.RecipientUsername,
		Amount:            ts.Amount,
		SourceFundId:      ts.SourceOfFundId,
		Description:       ts.Description,
		CreatedAt:         ts.CreatedAt,
		UpdatedAt:         ts.UpdatedAt,
	}
}

func ToListTransactionResponses(trs []entity.Transaction) []ListTransactionResponse {
	var listResponses []ListTransactionResponse
	for _, t := range trs {
		listResponses = append(listResponses, *ToListTransactionResponse(t))
	}
	return listResponses
}

func ToListTransactionQuery(query *ListTransactionQuery) *ListTransactionQuery {
	if query.Limit == 0 {
		query.Limit = 8
	}
	if query.Page == 0 {
		query.Page = 1
	}

	query.SortBy = strings.ToLower(query.SortBy)
	if query.SortBy == "date" {
		query.SortBy = "updated_at"
	}
	if query.SortBy == "amount" {
		query.SortBy = "updated_at"
	}
	if query.SortBy == "to" {
		query.SortBy = "recipient_wallet_id"
	}

	query.Sort = strings.ToUpper(query.Sort)
	if query.Sort != "ASC" {
		query.Sort = "DESC"
	}

	return query
}
