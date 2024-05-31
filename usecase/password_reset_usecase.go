package usecase

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/repository"
	"assignment-go-rest-api/utils"
	"context"
	"errors"
	"time"
)

type PasswordResetUsecaseOpts struct {
	PasswordResetRepository repository.PasswordResetRepository
	UserRepository          repository.UserRepository
	Transactor              repository.Transactor
	AuthTokenProvider       utils.AuthTokenProvider
}

type PasswordResetUsecase interface {
	ForgotPassword(ctx context.Context, email string) (*entity.PasswordReset, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type PasswordResetUsecaseImpl struct {
	passwordResetRepository repository.PasswordResetRepository
	userRepository          repository.UserRepository
	transactor              repository.Transactor
	authTokenProvider       utils.AuthTokenProvider
}

func NewPasswordResetUsecase(passUOpts *PasswordResetUsecaseOpts) PasswordResetUsecase {
	return &PasswordResetUsecaseImpl{
		passwordResetRepository: passUOpts.PasswordResetRepository,
		userRepository:          passUOpts.UserRepository,
		transactor:              passUOpts.Transactor,
		authTokenProvider:       passUOpts.AuthTokenProvider,
	}
}

func (u *PasswordResetUsecaseImpl) ForgotPassword(ctx context.Context, email string) (*entity.PasswordReset, error) {
	config, err := utils.ConfigInit()
	if err != nil {
		return nil, err
	}

	user, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	passwordReset, err := u.passwordResetRepository.GetByUserId(ctx, user.Id)
	if err != nil {
		return nil, err
	}
	if passwordReset != nil {
		u.passwordResetRepository.SoftDelete(ctx, passwordReset.Id)
	}

	dataTokenMap := make(map[string]interface{})
	dataTokenMap["id"] = user.Id

	token, _ := u.authTokenProvider.TokenCreateAndSign(dataTokenMap)

	passwordReset, err = u.passwordResetRepository.Save(ctx, &entity.PasswordReset{
		UserId:    user.Id,
		Token:     token,
		ExpiredAt: time.Now().Add(time.Duration(config.ResetTokenExp) * time.Minute),
	})
	if err != nil {
		return nil, err
	}

	return passwordReset, nil
}

func (u *PasswordResetUsecaseImpl) ResetPassword(ctx context.Context, token, newPassword string) error {

	passwordReset, err := u.passwordResetRepository.GetByToken(ctx, token)
	if err != nil {
		return err
	}

	if time.Now().Unix() > passwordReset.ExpiredAt.Unix() {
		return apperror.ErrTokenExpired()
	}

	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user, err := u.userRepository.GetById(ctx, passwordReset.UserId)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)

	_, err = u.userRepository.Update(ctx, user)
	if err != nil {
		return err
	}

	err = u.passwordResetRepository.SoftDelete(ctx, passwordReset.Id)
	if err != nil {
		return err
	}

	return nil
}
