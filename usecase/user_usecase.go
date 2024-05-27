package usecase

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/repository"
	"context"
)

type UserUsecaseOpts struct {
	UserRepository   repository.UserRepository
	WalletRepository repository.WalletRepository
}

type UserUsecase interface {
	GetUser(ctx context.Context, request *dto.UserRequestParam) (*entity.User, error)
}

type UserUsecaseImpl struct {
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
}

func NewUserUsecase(userUOpts *UserUsecaseOpts) UserUsecase {
	return &UserUsecaseImpl{
		userRepository:   userUOpts.UserRepository,
		walletRepository: userUOpts.WalletRepository,
	}
}

func (u *UserUsecaseImpl) GetUser(ctx context.Context, request *dto.UserRequestParam) (*entity.User, error) {
	user, err := u.userRepository.GetById(ctx, uint(request.UserId))
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if user.Id == 0 {
		return nil, apperror.ErrUserNotFound()
	}
	return user, nil
}
