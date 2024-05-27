package repository

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/entity"
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type TransactionRepoOpts struct {
	Db *sql.DB
}

type TransactionRepository interface {
	CreateOne(ctx context.Context, tc entity.Transaction) (uint, error)
}

type TransactionRepositoryImpl struct {
	Db *sql.DB
}

func NewTransactionRepository(trOpts *TransactionRepoOpts) TransactionRepository {
	return &TransactionRepositoryImpl{Db: trOpts.Db}
}

func (r *TransactionRepositoryImpl) CreateOne(ctx context.Context, tc entity.Transaction) (uint, error) {
	var err error

	SQL := `
		INSERT INTO transactions
		(sender_wallet_id, recipient_wallet_id, amount, source_of_fund_id, description)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, SQL, tc).Scan(&tc.Id)
	} else {
		err = r.Db.QueryRowContext(ctx, SQL, tc).Scan(&tc.Id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constant.ViolatesUniqueConstraintPgErrCode {
			return 0, apperror.ErrBadRequest()
		}
	}

	return tc.Id, nil
}
