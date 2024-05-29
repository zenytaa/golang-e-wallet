package usecase

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/repository"
	"context"

	"github.com/shopspring/decimal"
)

type TransactionUsecaseOpts struct {
	TransactionRepo repository.TransactionRepository
	SourceFundRepo  repository.SourceFundRepository
	WalletRepo      repository.WalletRepository
	Transactor      repository.Transactor
}

type TransactionUsecase interface {
	Transfer(ctx context.Context, tc entity.Transaction) (*entity.Transaction, error)
	TransferWithTransactor(ctx context.Context, tc entity.Transaction) (*entity.Transaction, error)
	TopUp(ctx context.Context, tc entity.Transaction) (*entity.Transaction, error)
	TopUpWithTransactor(ctx context.Context, tc entity.Transaction) (*entity.Transaction, error)
}

type TransactionUsecaseImpl struct {
	TransactionRepository repository.TransactionRepository
	SourceFundRepository  repository.SourceFundRepository
	WalletRepository      repository.WalletRepository
	Transactor            repository.Transactor
}

func NewTransactionUsecase(transUOpts *TransactionUsecaseOpts) TransactionUsecase {
	return &TransactionUsecaseImpl{
		TransactionRepository: transUOpts.TransactionRepo,
		SourceFundRepository:  transUOpts.SourceFundRepo,
		WalletRepository:      transUOpts.WalletRepo,
		Transactor:            transUOpts.Transactor,
	}
}

func (u *TransactionUsecaseImpl) Transfer(ctx context.Context, tc entity.Transaction) (*entity.Transaction, error) {
	if tc.Amount.IntPart() < int64(constant.MinTransfer) || tc.Amount.IntPart() > int64(constant.MaxTransfer) {
		return nil, apperror.ErrLimitTransfer()
	}

	sourceFund, err := u.SourceFundRepository.GetById(ctx, tc.SourceOfFund.Id)
	if err != nil {
		return nil, err
	}

	recipientWallet, err := u.WalletRepository.GetByWalletNumber(ctx, tc.RecipientWallet.WalletNumber)
	if err != nil {
		return nil, err
	}

	senderWallet, err := u.WalletRepository.GetByUserId(ctx, tc.SenderWallet.User.Id)
	if err != nil {
		return nil, err
	}

	if tc.RecipientWallet.WalletNumber == senderWallet.WalletNumber {
		return nil, apperror.ErrBadRequest()
	}

	if senderWallet.Balance.LessThan(tc.Amount) || senderWallet.Balance.LessThan(decimal.NewFromInt(0)) {
		return nil, apperror.ErrInsufficientBalance()
	}

	senderTc, err := u.TransactionRepository.CreateOne(ctx, entity.Transaction{
		SenderWallet:    *senderWallet,
		RecipientWallet: *recipientWallet,
		Amount:          tc.Amount,
		SourceOfFund:    tc.SourceOfFund,
		Description:     tc.Description,
	})

	if err != nil {
		return nil, err
	}

	_, err = u.TransactionRepository.CreateOne(ctx, entity.Transaction{
		SenderWallet:    *senderWallet,
		RecipientWallet: *recipientWallet,
		Amount:          tc.Amount,
		SourceOfFund:    tc.SourceOfFund,
		Description:     tc.Description,
	})
	if err != nil {
		return nil, err
	}

	// update wallet sender
	senderWallet.Balance = senderWallet.Balance.Sub(tc.Amount)
	_, err = u.WalletRepository.Update(ctx, senderWallet)
	if err != nil {
		return nil, err
	}

	// update wallet recipient
	recipientWallet.Balance = recipientWallet.Balance.Add(tc.Amount)
	_, err = u.WalletRepository.Update(ctx, recipientWallet)
	if err != nil {
		return nil, err
	}

	tcResponse := entity.Transaction{
		Id:              senderTc.Id,
		SenderWallet:    *senderWallet,
		RecipientWallet: *recipientWallet,
		Amount:          tc.Amount,
		SourceOfFund:    *sourceFund,
		Description:     tc.Description,
	}

	tcResponse.SenderWallet.Balance = decimal.NewFromInt(0)
	tcResponse.RecipientWallet.Balance = decimal.NewFromInt(0)

	return &tcResponse, nil
}

func (u *TransactionUsecaseImpl) TransferWithTransactor(ctx context.Context, tc entity.Transaction) (*entity.Transaction, error) {
	var newTc *entity.Transaction
	var err error

	_, err = u.Transactor.WithinTransactor(ctx, func(ctx context.Context) (interface{}, error) {
		res, err := u.Transfer(ctx, tc)
		if err != nil {
			return nil, err
		}
		newTc = res
		return res, nil
	})
	if err != nil {
		return nil, err
	}

	return newTc, nil
}

func (u *TransactionUsecaseImpl) TopUp(ctx context.Context, tc entity.Transaction) (*entity.Transaction, error) {
	if tc.Amount.IntPart() < int64(constant.MinTopUp) || tc.Amount.IntPart() > int64(constant.MaxTopUp) {
		return nil, apperror.ErrLimitTopUp()
	}

	sourceFund, err := u.SourceFundRepository.GetById(ctx, tc.SourceOfFund.Id)
	if err != nil {
		return nil, err
	}

	wallet, err := u.WalletRepository.GetByUserId(ctx, tc.SenderWallet.User.Id)
	if err != nil {
		return nil, err
	}

	res, err := u.TransactionRepository.CreateOne(ctx, entity.Transaction{
		Id:              0,
		SenderWallet:    *wallet,
		RecipientWallet: *wallet,
		Amount:          tc.Amount,
		SourceOfFund:    *sourceFund,
		Description:     `top up from ` + sourceFund.FundName,
	})
	if err != nil {
		return nil, err
	}

	wallet.Balance = wallet.Balance.Add(tc.Amount)
	_, err = u.WalletRepository.Update(ctx, wallet)
	if err != nil {
		return nil, err
	}

	tcResponse := entity.Transaction{
		Id:              res.Id,
		SenderWallet:    *wallet,
		RecipientWallet: *wallet,
		Amount:          tc.Amount,
		SourceOfFund:    *sourceFund,
		Description:     res.Description,
	}

	return &tcResponse, nil
}

func (u *TransactionUsecaseImpl) TopUpWithTransactor(ctx context.Context, tc entity.Transaction) (*entity.Transaction, error) {
	var newTc *entity.Transaction
	var err error

	_, err = u.Transactor.WithinTransactor(ctx, func(ctx context.Context) (interface{}, error) {
		res, err := u.TopUp(ctx, tc)
		if err != nil {
			return nil, err
		}
		newTc = res
		return res, nil
	})
	if err != nil {
		return nil, err
	}

	return newTc, nil
}
