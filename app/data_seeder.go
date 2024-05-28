package app

// import (
// 	"assignment-go-rest-api/dto"

// 	"github.com/shopspring/decimal"
// )

// func UserData() []dto.AuthRegisterRequest {
// 	users := []dto.AuthRegisterRequest{
// 		{
// 			Email:    "user1@gmail.com",
// 			Username: "User1",
// 			Password: "password1",
// 		},
// 		{
// 			Email:    "user2@gmail.com",
// 			Username: "User2",
// 			Password: "password2",
// 		},
// 		{
// 			Email:    "user3@gmail.com",
// 			Username: "User3",
// 			Password: "password3",
// 		},
// 		{
// 			Email:    "user4@gmail.com",
// 			Username: "User4",
// 			Password: "password4",
// 		},
// 		{
// 			Email:    "user5@gmail.com",
// 			Username: "User5",
// 			Password: "password5",
// 		},
// 	}

// 	return users
// }

// func TransactionTopUpData() []*dto.TopUpCreateRequest {
// 	topUps := []*dto.TopUpCreateRequest{
// 		{
// 			Amount:       decimal.NewFromInt(1000000),
// 			SourceFundId: 2,
// 			UserId:       1,
// 		},
// 		{
// 			Amount:       decimal.NewFromInt(1500000),
// 			SourceFundId: 1,
// 			UserId:       2,
// 		},
// 		{
// 			Amount:       decimal.NewFromInt(2000000),
// 			SourceFundId: 3,
// 			UserId:       3,
// 		},
// 		{
// 			Amount:       decimal.NewFromInt(2500000),
// 			SourceFundId: 1,
// 			UserId:       4,
// 		},
// 		{
// 			Amount:       decimal.NewFromInt(3000000),
// 			SourceFundId: 2,
// 			UserId:       5,
// 		},
// 		{
// 			Amount:       decimal.NewFromInt(3500000),
// 			SourceFundId: 3,
// 			UserId:       1,
// 		},
// 		{
// 			Amount:       decimal.NewFromInt(4000000),
// 			SourceFundId: 3,
// 			UserId:       2,
// 		},
// 		{
// 			Amount:       decimal.NewFromInt(4500000),
// 			SourceFundId: 3,
// 			UserId:       3,
// 		},
// 	}

// 	return topUps
// }

// func TransactionTransferData() []*dto.TransferCreateRequest {
// 	transfers := []*dto.TransferCreateRequest{
// 		{
// 			SenderUserId:          1,
// 			Amount:                decimal.NewFromInt(50000),
// 			RecipientWalletNumber: "9990000000002",
// 			SourceFundId:          1,
// 		},
// 		{
// 			SenderUserId:          2,
// 			Amount:                decimal.NewFromInt(50000),
// 			RecipientWalletNumber: "9990000000001",
// 			SourceFundId:          1,
// 		},
// 		{
// 			SenderUserId:          3,
// 			Amount:                decimal.NewFromInt(2500000),
// 			RecipientWalletNumber: "9990000000004",
// 			SourceFundId:          1,
// 		},
// 		{
// 			SenderUserId:          4,
// 			Amount:                decimal.NewFromInt(2500000),
// 			RecipientWalletNumber: "9990000000003",
// 			SourceFundId:          1,
// 		},
// 		{
// 			SenderUserId:          5,
// 			Amount:                decimal.NewFromInt(250000),
// 			RecipientWalletNumber: "9990000000003",
// 			SourceFundId:          1,
// 		},
// 		{
// 			SenderUserId:          3,
// 			Amount:                decimal.NewFromInt(250000),
// 			RecipientWalletNumber: "9990000000005",
// 			SourceFundId:          1,
// 		},
// 		{
// 			SenderUserId:          4,
// 			Amount:                decimal.NewFromInt(300000),
// 			RecipientWalletNumber: "9990000000001",
// 			SourceFundId:          3,
// 		},
// 		{
// 			SenderUserId:          1,
// 			Amount:                decimal.NewFromInt(300000),
// 			RecipientWalletNumber: "9990000000004",
// 			SourceFundId:          3,
// 		},
// 		{
// 			SenderUserId:          2,
// 			Amount:                decimal.NewFromInt(450000),
// 			RecipientWalletNumber: "9990000000003",
// 			SourceFundId:          1,
// 		},
// 		{
// 			SenderUserId:          3,
// 			Amount:                decimal.NewFromInt(450000),
// 			RecipientWalletNumber: "9990000000002",
// 			SourceFundId:          1,
// 		},
// 		{
// 			SenderUserId:          4,
// 			Amount:                decimal.NewFromInt(200000),
// 			RecipientWalletNumber: "9990000000005",
// 			SourceFundId:          3,
// 		},
// 		{
// 			SenderUserId:          5,
// 			Amount:                decimal.NewFromInt(200000),
// 			RecipientWalletNumber: "9990000000004",
// 			SourceFundId:          3,
// 		},
// 	}

// 	return transfers
// }
