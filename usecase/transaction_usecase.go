package usecase

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/repository"
	"assignment-go-rest-api/utils"
	"context"
	"database/sql"

	"github.com/shopspring/decimal"
)

type TransactionUsecase interface {
	ListTransaction(ctx context.Context, userId uint, query *dto.ListTransactionQuery) ([]dto.ListTransactionResponse, error)
	TopUp(ctx context.Context, request *dto.TopUpCreateRequest) (*dto.TopUpResponse, error)
	Transfer(ctx context.Context, request *dto.TransferCreateRequest) (*dto.TransferResponse, error)
	Income(ctx context.Context, userId uint) (*decimal.Decimal, error)
	Expense(ctx context.Context, userId uint) (*decimal.Decimal, error)
}

type TransactionUsecaseImpl struct {
	db                    *sql.DB
	transactionRepository repository.TransactionRepository
	walletRepository      repository.WalletRepository
	sourceFundRepository  repository.SourceFundRepository
	userRepository        repository.UserRepository
}

func NewTransactionUsecase(
	db *sql.DB,
	transactionRepository repository.TransactionRepository,
	walletRepository repository.WalletRepository,
	sourceFundRepository repository.SourceFundRepository,
	userRepository repository.UserRepository,
) TransactionUsecase {
	return &TransactionUsecaseImpl{
		db:                    db,
		transactionRepository: transactionRepository,
		walletRepository:      walletRepository,
		sourceFundRepository:  sourceFundRepository,
		userRepository:        userRepository,
	}
}

// Expense implements TransactionUsecase.
func (u *TransactionUsecaseImpl) Expense(ctx context.Context, userId uint) (*decimal.Decimal, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	expense, err := u.transactionRepository.Expense(ctx, tx, userId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	return &expense, nil
}

// Income implements TransactionUsecase.
func (u *TransactionUsecaseImpl) Income(ctx context.Context, userId uint) (*decimal.Decimal, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	income, err := u.transactionRepository.Income(ctx, tx, userId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	return &income, nil
}

func (u *TransactionUsecaseImpl) ListTransaction(ctx context.Context, userId uint, query *dto.ListTransactionQuery) ([]dto.ListTransactionResponse, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	// fmt.Println("!query.SearchExist ; ", !query.SearchExist)
	// fmt.Println("!query.SortByExist ; ", !query.SortByExist)
	// fmt.Println("!query.SortExist ; ", !query.SortExist)

	if query.SearchExist && query.SortByExist {
		query = dto.ToListTransactionQuery(query)
		trs, err := u.transactionRepository.GetAllOrderBySearch(ctx, tx, userId, query)
		if err != nil {
			return nil, apperror.ErrInternalServer()
		}
		if len(query.SortBy) == 0 {
			return nil, apperror.ErrBadRequest()
		}
		if len(trs) == 0 {
			return nil, apperror.ErrTransactionNotFound()
		}
		return dto.ToListTransactionResponses(trs), nil
	}

	if query.SearchExist {
		query = dto.ToListTransactionQuery(query)
		trs, err := u.transactionRepository.GetAllBySearch(ctx, tx, userId, query)
		if err != nil {
			return nil, apperror.ErrInternalServer()
		}
		if len(query.Search) == 0 {
			return nil, apperror.ErrBadRequest()
		}
		if len(trs) == 0 {
			return nil, apperror.ErrTransactionNotFound()
		}
		return dto.ToListTransactionResponses(trs), nil
	}

	if query.SortByExist {
		query = dto.ToListTransactionQuery(query)
		trs, err := u.transactionRepository.GetAllOrderBy(ctx, tx, userId, query)
		if err != nil {
			return nil, apperror.ErrInternalServer()
		}
		if len(query.SortBy) == 0 {
			return nil, apperror.ErrBadRequest()
		}
		if len(trs) == 0 {
			return nil, apperror.ErrTransactionNotFound()
		}
		return dto.ToListTransactionResponses(trs), nil
	}

	trs, err := u.transactionRepository.GetAllByUserId(ctx, tx, userId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	recipientWallets := []entity.Wallet{}
	for _, t := range trs {
		recipientWallet, err := u.walletRepository.GetByWalletId(ctx, tx, t.RecipientWalletId)
		if err != nil {
			return nil, apperror.ErrInternalServer()
		}
		if recipientWallet.Id == 0 {
			return nil, apperror.ErrWalletNotFound()
		}
		recipientWallets = append(recipientWallets, *recipientWallet)
	}

	recipientUserDatas := []entity.User{}
	for _, w := range recipientWallets {
		recipientUserData, err := u.userRepository.GetById(ctx, tx, w.UserId)
		if err != nil {
			return nil, apperror.ErrInternalServer()
		}
		if recipientUserData.Id == 0 {
			return nil, apperror.ErrWalletNotFound()
		}
		recipientUserDatas = append(recipientUserDatas, *recipientUserData)
	}

	return dto.ToListTransactionResponses(trs), nil

}

func (u *TransactionUsecaseImpl) TopUp(ctx context.Context, request *dto.TopUpCreateRequest) (*dto.TopUpResponse, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	if request.Amount.LessThanOrEqual(decimal.NewFromInt(int64(constant.MinTopUp))) {
		return nil, apperror.ErrMiniumTopUp()
	}
	if request.Amount.GreaterThan(decimal.NewFromInt(int64(constant.MaxTopUp))) {
		return nil, apperror.ErrMaximumTopUp()
	}

	sourceFund, err := u.sourceFundRepository.GetById(ctx, request.SourceFundId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if sourceFund.Id == 0 {
		return nil, apperror.ErrSourceFundNotFound()
	}

	wallet, err := u.walletRepository.GetByUserId(ctx, tx, request.UserId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if wallet.Id == 0 {
		return nil, apperror.ErrWalletNotFound()
	}

	transaction := &entity.Transaction{
		SenderWalletId:    request.UserId,
		RecipientWalletId: request.UserId,
		RecipientUsername: "-",
		Amount:            request.Amount,
		SourceOfFundId:    sourceFund.Id,
		Description:       "Top Up from " + sourceFund.FundName,
	}

	newTransaction, err := u.transactionRepository.Save(ctx, tx, transaction)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	wallet.Balance = decimal.Sum(wallet.Balance, request.Amount)

	_, err = u.walletRepository.Update(ctx, tx, wallet)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	return dto.ToTopUpResponse(*newTransaction), nil
}

func (u *TransactionUsecaseImpl) Transfer(ctx context.Context, request *dto.TransferCreateRequest) (*dto.TransferResponse, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	if len(request.Description) <= 0 {
		request.Description = "-"
	}

	if request.Amount.LessThanOrEqual(decimal.NewFromInt(int64(constant.MinTransfer))) {
		return nil, apperror.ErrMiniumTransfer()
	}
	if request.Amount.GreaterThan(decimal.NewFromInt(int64(constant.MaxTransfer))) {
		return nil, apperror.ErrMaximumTopUp()
	}

	sourceFund, err := u.sourceFundRepository.GetById(ctx, request.SourceFundId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if sourceFund.Id == 0 {
		return nil, apperror.ErrSourceFundNotFound()
	}

	senderWallet, err := u.walletRepository.GetByUserId(ctx, tx, request.SenderUserId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if senderWallet.Id == 0 {
		return nil, apperror.ErrWalletNotFound()
	}
	if senderWallet.Balance.LessThan(request.Amount) {
		return nil, apperror.ErrInsufficientBalance()
	}

	recipientWallet, err := u.walletRepository.GetByWalletNumber(ctx, tx, request.RecipientWalletNumber)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if recipientWallet.Id == 0 {
		return nil, apperror.ErrWalletRecipientNotFound()
	}

	if senderWallet.Id == recipientWallet.Id {
		return nil, apperror.ErrCantTransferToOwnWallet()
	}

	senderUserData, err := u.userRepository.GetById(ctx, tx, senderWallet.UserId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if senderUserData.Id == 0 {
		return nil, apperror.ErrUserNotFound()
	}

	transaction := &entity.Transaction{
		SenderWalletId:    recipientWallet.Id,
		RecipientWalletId: senderWallet.Id,
		RecipientUsername: senderUserData.Username,
		Amount:            request.Amount,
		SourceOfFundId:    sourceFund.Id,
		Description:       request.Description,
	}

	sender := &entity.User{
		Email:    senderUserData.Email,
		Username: senderUserData.Username,
	}

	transaction, err = u.transactionRepository.Save(ctx, tx, transaction)
	if err != nil {
		return dto.ToTransferResponse(*transaction, *sender), apperror.ErrInternalServer()
	}

	// ============================================================================================

	recipientUserData, err := u.userRepository.GetById(ctx, tx, recipientWallet.UserId)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}
	if recipientUserData.Id == 0 {
		return nil, apperror.ErrUserNotFound()
	}

	transaction = &entity.Transaction{
		SenderWalletId:    senderWallet.Id,
		RecipientWalletId: recipientWallet.Id,
		RecipientUsername: recipientUserData.Username,
		Amount:            request.Amount,
		SourceOfFundId:    sourceFund.Id,
		Description:       request.Description,
	}

	recipient := &entity.User{
		Email:    recipientUserData.Email,
		Username: recipientUserData.Username,
	}

	transaction, err = u.transactionRepository.Save(ctx, tx, transaction)
	if err != nil {
		return dto.ToTransferResponse(*transaction, *recipient), apperror.ErrInternalServer()
	}

	senderWallet.Balance = senderWallet.Balance.Sub(request.Amount)
	_, err = u.walletRepository.Update(ctx, tx, senderWallet)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	recipientWallet.Balance = recipientWallet.Balance.Add(request.Amount)
	_, err = u.walletRepository.Update(ctx, tx, recipientWallet)
	if err != nil {
		return nil, apperror.ErrInternalServer()
	}

	return dto.ToTransferResponse(*transaction, *recipient), nil
}
