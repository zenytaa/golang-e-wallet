package repository

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/constant"
	"assignment-go-rest-api/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx"
)

type TransactionRepoOpts struct {
	Db *sql.DB
}

type TransactionRepository interface {
	CreateOne(ctx context.Context, tc entity.Transaction) (*entity.Transaction, error)
	GetAllByUser(ctx context.Context, senderId uint64, params entity.TransactionParams) ([]entity.Transaction, int, error)
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
	values = append(values, tc.SourceOfFund.Id)
	values = append(values, tc.Description)
	values = append(values, tc.TransactionType.Id)

	SQL := `
		INSERT INTO transactions
		(sender_wallet_id, recipient_wallet_id, amount, source_of_fund_id, description, transaction_type_id)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, description;
	`

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, SQL, values...).Scan(&newTc.Id, &newTc.Description)
	} else {
		err = r.Db.QueryRowContext(ctx, SQL, values...).Scan(&newTc.Id, &newTc.Description)
	}

	if err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constant.ViolatesUniqueConstraintPgErrCode {
			return nil, apperror.ErrBadRequest()
		}
	}

	return &newTc, nil
}

func (r *TransactionRepositoryImpl) GetAllByUser(ctx context.Context, senderId uint64, params entity.TransactionParams) ([]entity.Transaction, int, error) {
	tcs := []entity.Transaction{}

	var totalRows int

	sqlCountTotalRows := `SELECT COUNT(*) OVER() AS total_rows`
	sqlGetAllColl := `
		, t.id, t.amount, t.description, t.created_at, 
		sw.id, sw.wallet_number, ss.id, ss.username, 
		rw.id, rw.wallet_number, rs.id, rs.username,
		sf.id, sf.fund_name,
		tt.id, tt.name
	`
	sqlGetAllCommand := `
		FROM transactions t
		JOIN wallets sw ON t.sender_wallet_id = sw.id
		JOIN users ss ON sw.user_id = ss.id
		JOIN wallets rw ON t.recipient_wallet_id = rw.id
		JOIN users rs ON rw.user_id = rs.id
		JOIN source_of_funds sf ON t.source_of_fund_id = sf.id
		JOIN transaction_types tt ON t.transaction_type_id = tt.id
		WHERE t.deleted_at IS NULL AND ss.id = $1 or rs.id = $1
	`

	var sb strings.Builder
	sb.WriteString(sqlCountTotalRows)
	sb.WriteString(sqlGetAllColl)
	sb.WriteString(sqlGetAllCommand)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(sqlCountTotalRows)
	sbTotalRows.WriteString(sqlGetAllCommand)

	values := []interface{}{}
	valuesCountTotal := []interface{}{}

	numberArguments := 2
	values = append(values, senderId)
	valuesCountTotal = append(valuesCountTotal, senderId)

	if params.Keyword != "" {
		sb.WriteString(`AND t.description ILIKE '%`)
		sb.WriteString(params.Keyword)
		sb.WriteString(`%' `)

		sbTotalRows.WriteString(`AND t.description ILIKE '%`)
		sbTotalRows.WriteString(params.Keyword)
		sbTotalRows.WriteString(`%' `)
	}

	var sortBy string
	switch params.SortBy {
	case "amount":
		sortBy = `t.amount `
	default:
		sortBy = `t.created_at `
	}
	sb.WriteString(fmt.Sprintf(`ORDER BY %s `, sortBy))
	sbTotalRows.WriteString(fmt.Sprintf(`ORDER BY %s `, sortBy))

	if params.Sort == "" {
		params.Sort = `DESC `
	}
	sb.WriteString(fmt.Sprintf(`%s `, params.Sort))
	sbTotalRows.WriteString(fmt.Sprintf(`%s `, params.Sort))

	if params.Limit != 0 {
		sb.WriteString(`LIMIT `)
		sb.WriteString(fmt.Sprintf(`$%d `, numberArguments))
		values = append(values, params.Limit)
		numberArguments++
	}

	if params.Page != 0 {
		sb.WriteString(`OFFSET `)
		sb.WriteString(fmt.Sprintf(`$%d`, numberArguments))
		values = append(values, params.Limit*(params.Page-1))
		numberArguments++
	}

	rows, err := r.Db.QueryContext(ctx, sb.String(), values...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		tc := entity.Transaction{}
		err := rows.Scan(&totalRows,
			&tc.Id, &tc.Amount, &tc.Description, &tc.CreatedAt,
			&tc.SenderWallet.Id, &tc.SenderWallet.WalletNumber,
			&tc.SenderWallet.User.Id, &tc.SenderWallet.User.Username,
			&tc.RecipientWallet.Id, &tc.RecipientWallet.WalletNumber, &tc.RecipientWallet.User.Id, &tc.RecipientWallet.User.Username,
			&tc.SourceOfFund.Id, &tc.SourceOfFund.FundName,
			&tc.TransactionType.Id, &tc.TransactionType.Name,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, apperror.ErrTransactionNotFound()
			}
			return nil, 0, err
		}
		tcs = append(tcs, tc)
	}

	if totalRows == 0 {
		err := r.Db.QueryRowContext(ctx, sbTotalRows.String(), valuesCountTotal...).Scan(&totalRows)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, apperror.ErrTransactionNotFound()
			}
			return nil, 0, err
		}
	}

	return tcs, totalRows, nil
}
