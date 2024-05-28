package usecase

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/repository"
	"assignment-go-rest-api/utils"
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type AuthUsecaseOpts struct {
	UserRepository          repository.UserRepository
	WalletRepository        repository.WalletRepository
	PasswordResetRepository repository.PasswordResetRepository
	Transactor              repository.Transactor
}

type AuthUsecase interface {
	Register(ctx context.Context, request dto.AuthRegisterRequest) error
	RegisterWithInTransactor(ctx context.Context, request dto.AuthRegisterRequest) error
	Login(ctx context.Context, request dto.AuthLoginRequest) (*string, error)
	ForgotPassword(ctx context.Context, request dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error)
}

type AuthUsecaseImpl struct {
	userRepository          repository.UserRepository
	walletRepository        repository.WalletRepository
	passwordResetRepository repository.PasswordResetRepository
	transactor              repository.Transactor
}

func NewAuthUsecase(authuOpts *AuthUsecaseOpts) AuthUsecase {
	return &AuthUsecaseImpl{
		userRepository:          authuOpts.UserRepository,
		walletRepository:        authuOpts.WalletRepository,
		passwordResetRepository: authuOpts.PasswordResetRepository,
		transactor:              authuOpts.Transactor,
	}
}

func (u *AuthUsecaseImpl) Register(ctx context.Context, request dto.AuthRegisterRequest) error {
	hashCostString := os.Getenv("HASH_COST")
	hashCost, _ := strconv.Atoi(hashCostString)
	hashPassword, err := utils.HashPassword(request.Password, hashCost)
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
	if ok {
		tokenString, err := utils.TokenCreateAndSign()
		if err != nil {
			return nil, apperror.ErrUnauthorized()
		}
		return &tokenString, nil
	}

	return nil, apperror.ErrInternalServer()
}

func (u *AuthUsecaseImpl) ForgotPassword(ctx context.Context, request dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error) {
	user, err := u.userRepository.GetByEmail(ctx, request.Email)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if user.Id == 0 {
		return nil, apperror.ErrUserNotFound()
	}

	passwordReset, err := u.passwordResetRepository.GetByUserId(ctx, user.Id)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	passwordReset.UserId = user.Id
	token, _ := utils.TokenCreateAndSign(user.Id, os.Getenv("SECRET_KEY"))
	passwordReset.Token = token
	passwordReset.ExpiredAt = time.Now().Add(1 * time.Hour)

	passwordReset, err = u.passwordResetRepository.Save(ctx, passwordReset)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	return dto.ToForgotPasswordResponse(*user, *passwordReset), nil
}

func (u *AuthUsecaseImpl) ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error) {
	passwordReset, err := u.passwordResetRepository.GetByToken(ctx, request.Token)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	if passwordReset.UserId == 0 {
		return nil, apperror.ErrResetTokenNotFound()
	}
	if request.Password != request.ConfirmPassword {
		return nil, apperror.ErrPasswordNotMatch()
	}

	hashCostInt, _ := strconv.Atoi(os.Getenv("HASH_COST"))
	passwordHash, err := utils.HashPassword(request.Password, hashCostInt)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user, err := u.userRepository.GetById(ctx, passwordReset.UserId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if user.Id == 0 {
		return nil, apperror.ErrUserNotFound()
	}
	user.Password = string(passwordHash)

	_, err = u.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	_, err = u.passwordResetRepository.SoftDelete(ctx, passwordReset)
	if err != nil {
		return nil, err
	}

	return dto.ToResetPasswordResponse(*user, *passwordReset), nil
}
