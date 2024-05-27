package repository

import (
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/utils"
	"context"
	"database/sql"
)

type WalletRepository interface {
	Save(ctx context.Context, tx *sql.Tx, wallet *entity.Wallet) (*entity.Wallet, error)
	GetByUserId(ctx context.Context, tx *sql.Tx, userId uint) (*entity.Wallet, error)
	GetByWalletNumber(ctx context.Context, tx *sql.Tx, walletNumber string) (*entity.Wallet, error)
	Update(ctx context.Context, tx *sql.Tx, wallet *entity.Wallet) (*entity.Wallet, error)
	GetByWalletId(ctx context.Context, tx *sql.Tx, walletId uint) (*entity.Wallet, error)
}

type WalletRepositoryImpl struct {
}

func NewWalletRepository() WalletRepository {
	return &WalletRepositoryImpl{}
}

func (r *WalletRepositoryImpl) GetByWalletId(ctx context.Context, tx *sql.Tx, walletId uint) (*entity.Wallet, error) {
	SQL := `
	SELECT id, user_id, wallet_number, balance, created_at, updated_at, deleted_at
	FROM wallets
	WHERE deleted_at IS NULL
	AND id = $1
	;`

	rows, err := tx.QueryContext(ctx, SQL, walletId)
	utils.IfErrorLogPrint(err)
	defer rows.Close()

	wallet := entity.Wallet{}
	if rows.Next() {
		err := rows.Scan(&wallet.Id, &wallet.UserId, &wallet.WalletNumber, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
		utils.IfErrorLogPrint(err)
		return &wallet, nil
	}

	return &wallet, err
}

func (r *WalletRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, wallet *entity.Wallet) (*entity.Wallet, error) {
	SQL := `
	INSERT INTO wallets (user_id, wallet_number, balance, created_at, updated_at)
	VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at
	;`

	err := tx.QueryRowContext(ctx, SQL, wallet.UserId, wallet.WalletNumber, wallet.Balance).
		Scan(&wallet.Id, &wallet.CreatedAt, &wallet.UpdatedAt)
	utils.IfErrorLogPrint(err)

	return wallet, err
}

func (r *WalletRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, wallet *entity.Wallet) (*entity.Wallet, error) {
	SQL := `
	UPDATE wallets
	SET balance = $1, updated_at = NOW()
	WHERE id = $2
	;`

	_, err := tx.ExecContext(ctx, SQL, wallet.Balance, wallet.Id)
	utils.IfErrorLogPrint(err)

	return wallet, err
}

func (r *WalletRepositoryImpl) GetByUserId(ctx context.Context, tx *sql.Tx, userId uint) (*entity.Wallet, error) {
	SQL := `
	SELECT id, user_id, wallet_number, balance, created_at, updated_at, deleted_at
	FROM wallets
	WHERE deleted_at IS NULL
	AND user_id = $1
	;`

	rows, err := tx.QueryContext(ctx, SQL, userId)
	utils.IfErrorLogPrint(err)
	defer rows.Close()

	wallet := entity.Wallet{}
	if rows.Next() {
		err := rows.Scan(&wallet.Id, &wallet.UserId, &wallet.WalletNumber, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
		utils.IfErrorLogPrint(err)
		return &wallet, nil
	}

	return &wallet, err
}

func (r *WalletRepositoryImpl) GetByWalletNumber(ctx context.Context, tx *sql.Tx, walletNumber string) (*entity.Wallet, error) {
	SQL := `
	SELECT id, user_id, wallet_number, balance, created_at, updated_at, deleted_at
	FROM wallets
	WHERE deleted_at IS NULL
	AND wallet_number = $1
	;`

	rows, err := tx.QueryContext(ctx, SQL, walletNumber)
	utils.IfErrorLogPrint(err)
	defer rows.Close()

	wallet := entity.Wallet{}
	if rows.Next() {
		err := rows.Scan(&wallet.Id, &wallet.UserId, &wallet.WalletNumber, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
		utils.IfErrorLogPrint(err)
		return &wallet, nil
	}

	return &wallet, err
}
