package dto

import (
	"assignment-go-rest-api/entity"
	"time"

	"github.com/shopspring/decimal"
)

type TopUpCreateRequest struct {
	Amount       decimal.Decimal `json:"amount" binding:"required"`
	SourceFundId uint            `json:"source_of_fund_id" binding:"required"`
}

type TopUpResponse struct {
	Id           uint            `json:"id"`
	Amount       decimal.Decimal `json:"amount"`
	SourceOfFund string          `json:"source_of_fund"`
	Description  string          `json:"description"`
	CreatedAt    time.Time       `json:"created_at"`
}

type TransferCreateRequest struct {
	Amount                decimal.Decimal `json:"amount" binding:"required"`
	RecipientWalletNumber string          `json:"recipient_wallet_number" binding:"required,gte=0,lte=13"`
	SourceFundId          uint            `json:"source_of_fund_id" binding:"required"`
	Description           string          `json:"description" binding:"lte=35"`
}

type TransferResponse struct {
	Id              uint            `json:"id"`
	SenderWallet    WalletResponse  `json:"sender_wallet"`
	RecipientWallet WalletResponse  `json:"recipient_wallet"`
	Amount          decimal.Decimal `json:"amount"`
	SourceFund      string          `json:"source_of_fund"`
	Description     string          `json:"description"`
	TransactionType string          `json:"transaction_type"`
	CreatedAt       time.Time       `json:"created_at"`
}

type TransactionListResponse struct {
	Id              uint            `json:"id"`
	Amount          decimal.Decimal `json:"amount"`
	Description     string          `json:"description"`
	CreatedAt       time.Time       `json:"created_at"`
	SenderWallet    WalletResponse  `json:"sender_wallet"`
	RecipientWallet WalletResponse  `json:"recipient_wallet"`
	FundName        string          `json:"source_of_fund"`
	TransactionType string          `json:"transaction_type"`
}

type TransactionListResponses struct {
	Pagination       PaginationResponse        `json:"pagination_info"`
	TransactionsList []TransactionListResponse `json:"transactions"`
}

func ToTopUpResponse(ts entity.Transaction) *TopUpResponse {
	return &TopUpResponse{
		Id:           ts.Id,
		Amount:       ts.Amount,
		SourceOfFund: ts.SourceOfFund.FundName,
		Description:  ts.Description,
		CreatedAt:    ts.CreatedAt,
	}
}

func ToTransferResponse(ts entity.Transaction) *TransferResponse {
	return &TransferResponse{
		Id:              ts.Id,
		SenderWallet:    *ToWalletResponse(ts.SenderWallet),
		RecipientWallet: *ToWalletResponse(ts.RecipientWallet),
		Amount:          ts.Amount,
		SourceFund:      ts.SourceOfFund.FundName,
		Description:     ts.Description,
		TransactionType: ts.TransactionType.Name,
		CreatedAt:       ts.CreatedAt,
	}
}

func ToTransactionListResponse(tc entity.Transaction) *TransactionListResponse {
	return &TransactionListResponse{
		Id:          tc.Id,
		Amount:      tc.Amount,
		Description: tc.Description,
		CreatedAt:   tc.CreatedAt,
		SenderWallet: WalletResponse{
			UserName:     tc.SenderWallet.User.Username,
			WalletNumber: tc.SenderWallet.WalletNumber,
		},
		RecipientWallet: WalletResponse{
			UserName:     tc.RecipientWallet.User.Username,
			WalletNumber: tc.RecipientWallet.WalletNumber,
		},
		FundName:        tc.SourceOfFund.FundName,
		TransactionType: tc.TransactionType.Name,
	}
}

func ToTransactionListResponses(tcs []entity.Transaction, pagination entity.PaginationInfo) *TransactionListResponses {
	responses := []TransactionListResponse{}

	for _, tc := range tcs {
		responses = append(responses, *ToTransactionListResponse(tc))
	}

	return &TransactionListResponses{Pagination: *ToPaginationResponse(pagination), TransactionsList: responses}
}
