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

type UserUsecase interface {
	GetUser(ctx context.Context, request *dto.UserRequestParam) (*entity.User, error)
	// CreateUser(ctx context.Context, request *dto.UserRequest) (*entity.User, error)
}

type UserUsecaseImpl struct {
	db                    *sql.DB
	userRepository        repository.UserRepository
	walletRepository      repository.WalletRepository
	transactionRepository repository.TransactionRepository
}

func NewUserUsecase(
	db *sql.DB,
	userRepository repository.UserRepository,
	walletRepository repository.WalletRepository,
	transactionRepository repository.TransactionRepository,
) UserUsecase {
	return &UserUsecaseImpl{
		db:                    db,
		userRepository:        userRepository,
		walletRepository:      walletRepository,
		transactionRepository: transactionRepository,
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
