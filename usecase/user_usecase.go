package usecase

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/repository"
	"assignment-go-rest-api/utils"
	"context"
	"database/sql"
)

type UserUsecaseOpts struct {
	Db               *sql.DB
	UserRepository   repository.UserRepository
	WalletRepository repository.WalletRepository
}

type UserUsecase interface {
	GetUser(ctx context.Context, request *dto.UserRequestParam) (*entity.User, error)
}

type UserUsecaseImpl struct {
	db               *sql.DB
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
}

func NewUserUsecase(userUOpts *UserUsecaseOpts) UserUsecase {
	return &UserUsecaseImpl{
		db:               userUOpts.Db,
		userRepository:   userUOpts.UserRepository,
		walletRepository: userUOpts.WalletRepository,
	}
}

func (u *UserUsecaseImpl) GetUser(ctx context.Context, request *dto.UserRequestParam) (*entity.User, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	user, err := u.userRepository.GetById(ctx, tx, uint(request.UserId))
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if user.Id == 0 {
		return nil, apperror.ErrUserNotFound()
	}
	return user, nil
}
