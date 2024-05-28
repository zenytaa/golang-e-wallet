package usecase

import (
	"assignment-go-rest-api/apperror"
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
	_, err := u.SourceFundRepository.GetById(ctx, tc.SourceOfFundId)
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

	if senderWallet.Balance.LessThan(tc.Amount) || senderWallet.Balance.LessThan(decimal.NewFromInt(0)) {
		return nil, apperror.ErrInsufficientBalance()
	}

	senderTcId, err := u.TransactionRepository.CreateOne(ctx, entity.Transaction{
		User:            entity.User{Id: senderWallet.User.Id},
		SenderWallet:    *senderWallet,
		RecipientWallet: *recipientWallet,
		Amount:          senderWallet.Balance.Sub(tc.Amount),
		SourceOfFundId:  tc.SourceOfFundId,
		Description:     tc.Description,
	})

	if err != nil {
		return nil, err
	}

	_, err = u.TransactionRepository.CreateOne(ctx, entity.Transaction{
		User:            entity.User{Id: recipientWallet.User.Id},
		SenderWallet:    *senderWallet,
		RecipientWallet: *recipientWallet,
		Amount:          recipientWallet.Balance.Add(tc.Amount),
		SourceOfFundId:  tc.SourceOfFundId,
		Description:     tc.Description,
	})

	if err != nil {
		return nil, err
	}

	return senderTcId, nil
}

func (u *TransactionUsecaseImpl) TransferWithTransactor(ctx context.Context, tc entity.Transaction) (*entity.Transaction, error) {
	var newTc *entity.Transaction
	var err error

	_, err = u.Transactor.WithinTransactor(ctx, func(ctx context.Context) (interface{}, error) {
		newTc, err = u.Transfer(ctx, tc)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return newTc, nil
}
