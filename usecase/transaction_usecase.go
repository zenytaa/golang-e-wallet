package usecase

import (
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/repository"
	"context"
)

type TransactionUsecaseOpts struct {
	TransactionRepo repository.TransactionRepository
	SourceFundRepo  repository.SourceFundRepository
	WalletRepo      repository.WalletRepository
}

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, tc entity.Transaction) (uint, error)
}

type TransactionUsecaseImpl struct {
	TransactionRepository repository.TransactionRepository
	SourceFundRepository  repository.SourceFundRepository
	WalletRepository      repository.WalletRepository
}

func NewTransactionUsecase(transUOpts *TransactionUsecaseOpts) TransactionUsecase {
	return &TransactionUsecaseImpl{
		TransactionRepository: transUOpts.TransactionRepo,
		SourceFundRepository:  transUOpts.SourceFundRepo,
		WalletRepository:      transUOpts.WalletRepo,
	}
}

func (u *TransactionUsecaseImpl) CreateTransaction(ctx context.Context, tc entity.Transaction) (uint, error) {
	_, err := u.SourceFundRepository.GetById(ctx, tc.SourceOfFundId)
	if err != nil {
		return 0, err
	}
	panic("unimplement")
}
