package repository

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/entity"
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx"
)

type TransactionRepoOpts struct {
	Db *sql.DB
}

type TransactionRepository interface {
	CreateOne(ctx context.Context, tc entity.Transaction) (*entity.Transaction, error)
}

type TransactionRepositoryImpl struct {
	Db *sql.DB
}

func NewTransactionRepository(trOpts *TransactionRepoOpts) TransactionRepository {
	return &TransactionRepositoryImpl{Db: trOpts.Db}
}

func (r *TransactionRepositoryImpl) CreateOne(ctx context.Context, tc entity.Transaction) (*entity.Transaction, error) {
	var err error
	newTc := entity.Transaction{}
	values := []interface{}{}
	values = append(values, tc.SenderWallet.Id)
	values = append(values, tc.RecipientWallet.Id)
	values = append(values, tc.Amount)
	values = append(values, tc.SourceOfFundId)
	values = append(values, tc.Description)

	SQL := `
		INSERT INTO transactions
		(sender_wallet_id, recipient_wallet_id, amount, source_of_fund_id, description)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, sender_wallet_id, recipient_wallet_id, amount, source_of_fund_id, description;
	`

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, SQL, values...).Scan(
			&newTc.Id, &newTc.SenderWallet.Id, &newTc.RecipientWallet.Id, &newTc.Amount, &newTc.Description,
		)
	} else {
		err = r.Db.QueryRowContext(ctx, SQL, values...).Scan(
			&newTc.Id, &newTc.SenderWallet.Id, &newTc.RecipientWallet.Id, &newTc.Amount, &newTc.Description,
		)
	}

	if err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constant.ViolatesUniqueConstraintPgErrCode {
			return nil, apperror.ErrBadRequest()
		}
	}

	return &newTc, nil
}
