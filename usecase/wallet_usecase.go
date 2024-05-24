package usecase

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/repository"
	"assignment-go-rest-api/utils"
	"context"
	"database/sql"

	"github.com/shopspring/decimal"
)

type WalletUsecaseOpts struct {
	Db               *sql.DB
	UserRepository   repository.UserRepository
	WalletRepository repository.WalletRepository
}

type WalletUsecase interface {
	CreateWallet(ctx context.Context, request *dto.WalletRequest) (*entity.Wallet, error)
	GetWalletByUserId(ctx context.Context, request *dto.WalletRequest) (*entity.Wallet, error)
}

type WalletUsecaseImpl struct {
	db               *sql.DB
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
}

func NewWalletUsecase(walletUOpts *WalletUsecaseOpts) WalletUsecase {
	return &WalletUsecaseImpl{
		db:               walletUOpts.Db,
		userRepository:   walletUOpts.UserRepository,
		walletRepository: walletUOpts.WalletRepository,
	}
}

func (u *WalletUsecaseImpl) CreateWallet(ctx context.Context, request *dto.WalletRequest) (*entity.Wallet, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	user, err := u.userRepository.GetById(ctx, tx, request.UserId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if user.Id == 0 {
		return nil, apperror.ErrUserNotFound()
	}

	wallet, err := u.walletRepository.GetByUserId(ctx, tx, user.Id)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if wallet.Id == 0 {
		return nil, apperror.ErrWalletAlreadyCreated()
	}

	wallet.UserId = user.Id
	wallet.Balance = decimal.New(0, 0)
	wallet.WalletNumber = utils.GenerateWalletNumber(user.Id)

	newWallet, err := u.walletRepository.Save(ctx, tx, wallet)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	return newWallet, nil
}

func (u *WalletUsecaseImpl) GetWalletByUserId(ctx context.Context, request *dto.WalletRequest) (*entity.Wallet, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	wallet, err := u.walletRepository.GetByUserId(ctx, tx, request.UserId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if wallet.Id == 0 {
		return nil, apperror.ErrWalletNotFound()
	}

	return wallet, nil
}
