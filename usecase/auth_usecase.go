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

type AuthUsecaseOpts struct {
	UserRepository    repository.UserRepository
	WalletRepository  repository.WalletRepository
	Transactor        repository.Transactor
	AuthTokenProvider utils.AuthTokenProvider
}

type AuthUsecase interface {
	Register(ctx context.Context, request dto.AuthRegisterRequest) error
	RegisterWithInTransactor(ctx context.Context, request dto.AuthRegisterRequest) error
	Login(ctx context.Context, request dto.AuthLoginRequest) (*string, error)
}

type AuthUsecaseImpl struct {
	userRepository    repository.UserRepository
	walletRepository  repository.WalletRepository
	transactor        repository.Transactor
	authTokenProvider utils.AuthTokenProvider
}

func NewAuthUsecase(authuOpts *AuthUsecaseOpts) AuthUsecase {
	return &AuthUsecaseImpl{
		userRepository:    authuOpts.UserRepository,
		walletRepository:  authuOpts.WalletRepository,
		transactor:        authuOpts.Transactor,
		authTokenProvider: authuOpts.AuthTokenProvider,
	}
}

func (u *AuthUsecaseImpl) Register(ctx context.Context, request dto.AuthRegisterRequest) error {
	hashPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return err
	}

	createUser := entity.User{
		Email:    request.Email,
		Username: request.Username,
		Password: string(hashPassword),
	}

	newUser, err := u.userRepository.Save(ctx, &createUser)
	if err != nil {
		return err
	}

	createWallet := entity.Wallet{
		User:         *newUser,
		WalletNumber: utils.GenerateWalletNumber(newUser.Id),
		Balance:      decimal.New(0, 0),
	}

	_, err = u.walletRepository.Save(ctx, &createWallet)
	if err != nil {
		return err
	}

	return nil
}

func (u *AuthUsecaseImpl) RegisterWithInTransactor(ctx context.Context, request dto.AuthRegisterRequest) error {
	_, err := u.transactor.WithinTransactor(ctx, func(ctx context.Context) (interface{}, error) {
		err := u.Register(ctx, request)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *AuthUsecaseImpl) Login(ctx context.Context, request dto.AuthLoginRequest) (*string, error) {
	if len(request.Email) == 0 {
		return nil, apperror.ErrEmailRequired()
	}
	if len(request.Password) == 0 {
		return nil, apperror.ErrPasswordRequired()
	}

	user, err := u.userRepository.GetByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if user.Id == 0 {
		return nil, apperror.ErrIncorrectCredentials()
	}

	ok, err := utils.CheckPassword(request.Password, []byte(user.Password))
	if err != nil {
		return nil, apperror.ErrIncorrectCredentials()
	}

	dataTokenMap := make(map[string]interface{})
	dataTokenMap["id"] = user.Id

	if ok {
		tokenString, err := u.authTokenProvider.TokenCreateAndSign(dataTokenMap)
		if err != nil {
			return nil, apperror.ErrUnauthorized()
		}
		return &tokenString, nil
	}

	return nil, apperror.ErrInternalServer()
}
