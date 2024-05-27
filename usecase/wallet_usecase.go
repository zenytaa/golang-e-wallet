package usecase

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/repository"
	"assignment-go-rest-api/utils"
	"context"

	"github.com/shopspring/decimal"
)

type WalletUsecaseOpts struct {
	UserRepository   repository.UserRepository
	WalletRepository repository.WalletRepository
}

type WalletUsecase interface {
	CreateWallet(ctx context.Context, request *dto.WalletRequest) (*entity.Wallet, error)
	GetWalletByUserId(ctx context.Context, request *dto.WalletRequest) (*entity.Wallet, error)
}

type WalletUsecaseImpl struct {
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
}

func NewWalletUsecase(walletUOpts *WalletUsecaseOpts) WalletUsecase {
	return &WalletUsecaseImpl{
		userRepository:   walletUOpts.UserRepository,
		walletRepository: walletUOpts.WalletRepository,
	}
}

func (u *WalletUsecaseImpl) CreateWallet(ctx context.Context, request *dto.WalletRequest) (*entity.Wallet, error) {
	user, err := u.userRepository.GetById(ctx, request.UserId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if user.Id == 0 {
		return nil, apperror.ErrUserNotFound()
	}

	wallet, err := u.walletRepository.GetByUserId(ctx, user.Id)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if wallet.Id == 0 {
		return nil, apperror.ErrWalletAlreadyCreated()
	}

	wallet.UserId = user.Id
	wallet.Balance = decimal.New(0, 0)
	wallet.WalletNumber = utils.GenerateWalletNumber(user.Id)

	newWallet, err := u.walletRepository.Save(ctx, wallet)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	return newWallet, nil
}

func (u *WalletUsecaseImpl) GetWalletByUserId(ctx context.Context, request *dto.WalletRequest) (*entity.Wallet, error) {
	wallet, err := u.walletRepository.GetByUserId(ctx, request.UserId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if wallet.Id == 0 {
		return nil, apperror.ErrWalletNotFound()
	}

	return wallet, nil
}
