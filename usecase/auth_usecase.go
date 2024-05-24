package usecase

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/repository"
	"assignment-go-rest-api/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type AuthUsecase interface {
	Register(ctx context.Context, request dto.AuthRegisterRequest) (*dto.AuthRegisterResponse, error)
	Login(ctx context.Context, request dto.AuthLoginRequest) (*string, error)
	ForgotPassword(ctx context.Context, request dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error)
}

type AuthUsecaseImpl struct {
	db                      *sql.DB
	userRepository          repository.UserRepository
	walletRepository        repository.WalletRepository
	passwordResetRepository repository.PasswordResetRepository
}

func NewAuthUsecase(db *sql.DB, userRepository repository.UserRepository, walletRepository repository.WalletRepository, passwordResetRepository repository.PasswordResetRepository) AuthUsecase {
	return &AuthUsecaseImpl{
		db:                      db,
		userRepository:          userRepository,
		walletRepository:        walletRepository,
		passwordResetRepository: passwordResetRepository,
	}
}

func (u *AuthUsecaseImpl) Register(ctx context.Context, request dto.AuthRegisterRequest) (*dto.AuthRegisterResponse, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	user, err := u.userRepository.GetByEmail(ctx, tx, request.Email)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if user.Id != 0 {
		return nil, apperror.ErrEmailAlreadyRegistered()
	}

	hashCostString := os.Getenv("HASH_COST")
	hashCost, _ := strconv.Atoi(hashCostString)
	hashPassword, err := utils.HashPassword(request.Password, hashCost)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	createUser := entity.User{
		Email:    request.Email,
		Username: request.Username,
		Password: string(hashPassword),
	}

	newUser, err := u.userRepository.Save(ctx, tx, &createUser)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	wallet, err := u.walletRepository.GetByUserId(ctx, tx, newUser.Id)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if wallet.Id != 0 {
		return nil, apperror.ErrWalletAlreadyCreated()
	}

	createWallet := entity.Wallet{
		UserId:       newUser.Id,
		WalletNumber: utils.GenerateWalletNumber(newUser.Id),
		Balance:      decimal.New(0, 0),
	}

	newWallet, err := u.walletRepository.Save(ctx, tx, &createWallet)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	return dto.ToAuthRegisterResponse(*newUser, *newWallet), nil
}

func (u *AuthUsecaseImpl) Login(ctx context.Context, request dto.AuthLoginRequest) (*string, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	if len(request.Email) == 0 {
		return nil, apperror.ErrEmailRequired()
	}
	if len(request.Password) == 0 {
		return nil, apperror.ErrPasswordRequired()
	}

	user, err := u.userRepository.GetByEmail(ctx, tx, request.Email)
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
		tokenString, err := utils.TokenCreateAndSign(user.Id, os.Getenv("SECRET_KEY"))
		if err != nil {
			return nil, apperror.ErrUnauthorized()
		}
		return &tokenString, nil
	}

	return nil, apperror.ErrInternalServer()
}

func (u *AuthUsecaseImpl) ForgotPassword(ctx context.Context, request dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	user, err := u.userRepository.GetByEmail(ctx, tx, request.Email)
	fmt.Println(user)
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
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

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

	user, err := u.userRepository.GetById(ctx, tx, passwordReset.UserId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if user.Id == 0 {
		return nil, apperror.ErrUserNotFound()
	}
	user.Password = string(passwordHash)

	_, err = u.userRepository.Update(ctx, tx, user)
	if err != nil {
		return nil, err
	}

	_, err = u.passwordResetRepository.SoftDelete(ctx, passwordReset)
	if err != nil {
		return nil, err
	}

	return dto.ToResetPasswordResponse(*user, *passwordReset), nil
}
