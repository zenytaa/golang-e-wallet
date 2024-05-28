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

type WalletRepoOpts struct {
	Db *sql.DB
}

type WalletRepository interface {
	Save(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error)
	GetByUserId(ctx context.Context, userId uint) (*entity.Wallet, error)
	GetByWalletNumber(ctx context.Context, walletNumber string) (*entity.Wallet, error)
	Update(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error)
	GetByWalletId(ctx context.Context, walletId uint) (*entity.Wallet, error)
}

type WalletRepositoryImpl struct {
	db *sql.DB
}

func NewWalletRepository(walletROpts *WalletRepoOpts) WalletRepository {
	return &WalletRepositoryImpl{
		db: walletROpts.Db,
	}
}

func (r *WalletRepositoryImpl) GetByWalletId(ctx context.Context, walletId uint) (*entity.Wallet, error) {
	wallet := entity.Wallet{}
	var err error

	SQL := `
	SELECT w.id, w.user_id, u.username, w.wallet_number, w.balance, w.created_at, w.updated_at, w.deleted_at
	FROM wallets w
	JOIN users u ON w.user_id = u.id
	WHERE deleted_at IS NULL
	AND id = $1
	;`

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, SQL, walletId).Scan(&wallet.Id, &wallet.User.Id, &wallet.User.Username, &wallet.WalletNumber, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
	} else {
		err = r.db.QueryRowContext(ctx, SQL, walletId).Scan(&wallet.Id, &wallet.User.Id, &wallet.User.Username, &wallet.WalletNumber, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ErrWalletNotFound()
		}
		return nil, err
	}

	return &wallet, err
}

func (r *WalletRepositoryImpl) Save(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error) {
	var err error
	newWallet := entity.Wallet{}
	values := []interface{}{}
	values = append(values, wallet.User.Id)
	values = append(values, wallet.WalletNumber)
	values = append(values, wallet.Balance)

	SQL := `
	INSERT INTO wallets (user_id, wallet_number, balance, created_at, updated_at)
	VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at
	;`

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, SQL, values...).Scan(&newWallet.Id, &newWallet.CreatedAt, &newWallet.UpdatedAt)
	} else {
		err = r.db.QueryRowContext(ctx, SQL, values...).Scan(&newWallet.Id, &newWallet.CreatedAt, &newWallet.UpdatedAt)
	}

	if err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constant.ViolatesUniqueConstraintPgErrCode {
			return nil, apperror.ErrBadRequest()
		}
	}

	return &newWallet, err
}

func (r *WalletRepositoryImpl) Update(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error) {
	var err error
	var stmt *sql.Stmt

	SQL := `
	UPDATE wallets
	SET balance = $1, updated_at = NOW()
	WHERE id = $2
	;`

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, SQL)
	} else {
		stmt, err = r.db.PrepareContext(ctx, SQL)
	}

	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(ctx, wallet.Balance, wallet.Id)
	if err != nil {
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, apperror.ErrWalletNotFound()
	}

	return wallet, err
}

func (r *WalletRepositoryImpl) GetByUserId(ctx context.Context, userId uint) (*entity.Wallet, error) {
	wallet := entity.Wallet{}
	var err error

	SQL := `
	SELECT w.id, w.user_id, u.username, w.wallet_number, w.balance, w.created_at, w.updated_at, w.deleted_at
	FROM wallets w
	JOIN users u ON w.user_id = u.id
	WHERE deleted_at IS NULL
	AND user_id = $1
	;`

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, SQL, userId).Scan(&wallet.Id, &wallet.User.Id, &wallet.User.Username, &wallet.WalletNumber, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
	} else {
		err = r.db.QueryRowContext(ctx, SQL, userId).Scan(&wallet.Id, &wallet.User.Id, &wallet.User.Username, &wallet.WalletNumber, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ErrWalletNotFound()
		}
		return nil, err
	}

	return &wallet, err
}

func (r *WalletRepositoryImpl) GetByWalletNumber(ctx context.Context, walletNumber string) (*entity.Wallet, error) {
	wallet := entity.Wallet{}
	var err error

	SQL := `
	SELECT w.id, w.user_id, u.username, w.wallet_number, w.balance, w.created_at, w.updated_at, w.deleted_at
	FROM wallets w
	JOIN users u ON w.user_id = u.id
	WHERE deleted_at IS NULL
	AND wallet_number = $1
	;`

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, SQL, walletNumber).Scan(&wallet.Id, &wallet.User.Id, &wallet.User.Username, &wallet.WalletNumber, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
	} else {
		err = r.db.QueryRowContext(ctx, SQL, walletNumber).Scan(&wallet.Id, &wallet.User.Id, &wallet.User.Username, &wallet.WalletNumber, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ErrWalletNotFound()
		}
		return nil, err
	}

	return &wallet, err
}
