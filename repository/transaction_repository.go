package repository

import (
	"assignment-go-rest-api/dto"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/utils"
	"context"
	"database/sql"

	"github.com/shopspring/decimal"
)

type TransactionRepository interface {
	GetAllByUserId(ctx context.Context, tx *sql.Tx, userId uint) ([]entity.Transaction, error)
	GetAllOrderBy(ctx context.Context, tx *sql.Tx, userId uint, query *dto.ListTransactionQuery) ([]entity.Transaction, error)
	Save(ctx context.Context, tx *sql.Tx, tc *entity.Transaction) (*entity.Transaction, error)
	Count(ctx context.Context, tx *sql.Tx, userId int) (int, error)
	GetAllOrderBySearch(ctx context.Context, tx *sql.Tx, userId uint, query *dto.ListTransactionQuery) ([]entity.Transaction, error)
	GetAllBySearch(ctx context.Context, tx *sql.Tx, userId uint, query *dto.ListTransactionQuery) ([]entity.Transaction, error)
	Income(ctx context.Context, tx *sql.Tx, userId uint) (decimal.Decimal, error)
	Expense(ctx context.Context, tx *sql.Tx, userId uint) (decimal.Decimal, error)
}

type TransactionRepositoryImpl struct {
}

func NewTransactionRepository() TransactionRepository {
	return &TransactionRepositoryImpl{}
}

func (r *TransactionRepositoryImpl) Income(ctx context.Context, tx *sql.Tx, userId uint) (decimal.Decimal, error) {

	SQL := `
	SELECT 
	COALESCE(SUM(amount), 0)
	FROM transactions
	WHERE deleted_at IS NULL
	AND sender_wallet_id = $1 AND sender_wallet_id = recipient_wallet_id
	`

	var count decimal.Decimal
	row := tx.QueryRowContext(ctx, SQL, userId)
	err := row.Scan(&count)
	utils.IfErrorLogPrint(err)
	return count, err
}

func (r *TransactionRepositoryImpl) Expense(ctx context.Context, tx *sql.Tx, userId uint) (decimal.Decimal, error) {

	SQL := `
	SELECT 
	COALESCE(SUM(amount), 0)
	FROM transactions
	WHERE deleted_at IS NULL
	AND sender_wallet_id = $1 AND sender_wallet_id != recipient_wallet_id
	`

	var count decimal.Decimal
	row := tx.QueryRowContext(ctx, SQL, userId)
	err := row.Scan(&count)
	utils.IfErrorLogPrint(err)
	return count, err
}

func (r *TransactionRepositoryImpl) GetAllByUserId(ctx context.Context, tx *sql.Tx, userId uint) ([]entity.Transaction, error) {
	tcs := []entity.Transaction{}

	SQL := `
	SELECT 
	id, sender_wallet_id, recipient_wallet_id, recipient_username,
	amount, source_of_fund_id, description, 
	created_at, updated_at, deleted_at
	FROM transactions
	WHERE deleted_at IS NULL
	AND sender_wallet_id = $1
	OR recipient_wallet_id = $1
	ORDER BY created_at DESC
	;`

	rows, err := tx.QueryContext(ctx, SQL, userId)
	utils.IfErrorLogPrint(err)
	defer rows.Close()

	for rows.Next() {
		tc := entity.Transaction{}
		err := rows.Scan(
			&tc.Id, &tc.SenderWalletId, &tc.RecipientWalletId, &tc.RecipientUsername,
			&tc.Amount, &tc.SourceOfFundId, &tc.Description,
			&tc.CreatedAt, &tc.UpdatedAt, &tc.DeletedAt,
		)
		utils.IfErrorLogPrint(err)
		tcs = append(tcs, tc)
	}
	return tcs, err
}

func (r *TransactionRepositoryImpl) Count(ctx context.Context, tx *sql.Tx, userId int) (int, error) {
	panic("unimplemented")
}

func (r *TransactionRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, tc *entity.Transaction) (*entity.Transaction, error) {
	SQL := `
	INSERT 	INTO transactions 
		(sender_wallet_id, recipient_wallet_id, recipient_username ,amount, source_of_fund_id, description, 
		created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id, created_at, updated_at
	;`

	err := tx.QueryRowContext(ctx, SQL,
		tc.SenderWalletId, tc.RecipientWalletId, tc.RecipientUsername,
		tc.Amount, tc.SourceOfFundId, tc.Description,
	).Scan(&tc.Id, &tc.CreatedAt, &tc.UpdatedAt)
	utils.IfErrorLogPrint(err)

	return tc, err
}

func (r *TransactionRepositoryImpl) GetAllOrderBy(ctx context.Context, tx *sql.Tx, userId uint, query *dto.ListTransactionQuery) ([]entity.Transaction, error) {
	SQL := `
	SELECT 
		id, sender_wallet_id, recipient_wallet_id, recipient_username,
		amount, source_of_fund_id, description, 
		created_at, updated_at, deleted_at
	FROM transactions
	WHERE deleted_at IS NULL
	AND sender_wallet_id = $1 
	OR recipient_wallet_id = $1
	ORDER BY ` + query.SortBy + " " + query.Sort + `;`

	// fmt.Println(query)

	// fmt.Println(SQL)
	rows, err := tx.QueryContext(ctx, SQL, userId)
	utils.IfErrorLogPrint(err)
	defer rows.Close()

	tcs := []entity.Transaction{}
	for rows.Next() {
		tc := entity.Transaction{}
		err := rows.Scan(
			&tc.Id, &tc.SenderWalletId, &tc.RecipientWalletId, &tc.RecipientUsername,
			&tc.Amount, &tc.SourceOfFundId, &tc.Description,
			&tc.CreatedAt, &tc.UpdatedAt, &tc.DeletedAt,
		)
		utils.IfErrorLogPrint(err)
		tcs = append(tcs, tc)
	}
	return tcs, err
}

func (r *TransactionRepositoryImpl) GetAllOrderBySearch(ctx context.Context, tx *sql.Tx, userId uint, query *dto.ListTransactionQuery) ([]entity.Transaction, error) {
	SQL := `
	SELECT 
		id, sender_wallet_id, recipient_wallet_id, recipient_username
		amount, source_of_fund_id, description, 
		created_at, updated_at, deleted_at
	FROM transactions
	WHERE deleted_at IS NULL
	AND sender_wallet_id = $1
	OR recipient_wallet_id = $1
	AND description ILIKE '%` + query.Search + `%'` + query.SortBy + " " + query.Sort + `;`

	// fmt.Println(query)

	// fmt.Println(SQL)
	rows, err := tx.QueryContext(ctx, SQL, userId)
	utils.IfErrorLogPrint(err)
	defer rows.Close()

	tcs := []entity.Transaction{}
	for rows.Next() {
		tc := entity.Transaction{}
		err := rows.Scan(
			&tc.Id, &tc.SenderWalletId, &tc.RecipientWalletId, &tc.RecipientUsername,
			&tc.Amount, &tc.SourceOfFundId, &tc.Description,
			&tc.CreatedAt, &tc.UpdatedAt, &tc.DeletedAt,
		)
		utils.IfErrorLogPrint(err)
		tcs = append(tcs, tc)
	}
	return tcs, err
}

func (r *TransactionRepositoryImpl) GetAllBySearch(ctx context.Context, tx *sql.Tx, userId uint, query *dto.ListTransactionQuery) ([]entity.Transaction, error) {
	SQL := `
	SELECT
		id, sender_wallet_id, recipient_wallet_id, recipient_username,
		amount, source_of_fund_id, description,
		created_at, updated_at, deleted_at
	FROM transactions
	WHERE deleted_at IS NULL
	AND sender_wallet_id = $1
	OR recipient_wallet_id = $1
	AND description ILIKE '%` + query.Search + `%';`

	// fmt.Println(query)

	// fmt.Println(SQL)
	rows, err := tx.QueryContext(ctx, SQL, userId)
	utils.IfErrorLogPrint(err)
	defer rows.Close()

	tcs := []entity.Transaction{}
	for rows.Next() {
		tc := entity.Transaction{}
		err := rows.Scan(
			&tc.Id, &tc.SenderWalletId, &tc.RecipientWalletId, &tc.RecipientUsername,
			&tc.Amount, &tc.SourceOfFundId, &tc.Description,
			&tc.CreatedAt, &tc.UpdatedAt, &tc.DeletedAt,
		)
		utils.IfErrorLogPrint(err)
		tcs = append(tcs, tc)
	}
	return tcs, err
}
