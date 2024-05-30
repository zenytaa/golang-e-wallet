package repository

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/utils"
	"context"
	"database/sql"
)

type PasswordResetRepoOpts struct {
	Db *sql.DB
}

type PasswordResetRepository interface {
	Save(ctx context.Context, pwdReset *entity.PasswordReset) (*entity.PasswordReset, error)
	GetByUserId(ctx context.Context, userId uint) (*entity.PasswordReset, error)
	GetByToken(ctx context.Context, token string) (*entity.PasswordReset, error)
	SoftDelete(ctx context.Context, id uint) error
}

type PasswordResetRepositoryImpl struct {
	Db *sql.DB
}

func NewPasswordResetRepository(passROpts *PasswordResetRepoOpts) PasswordResetRepository {
	return &PasswordResetRepositoryImpl{
		Db: passROpts.Db,
	}
}

func (r *PasswordResetRepositoryImpl) Save(ctx context.Context, pwdReset *entity.PasswordReset) (*entity.PasswordReset, error) {
	SQL := `
	INSERT INTO password_resets
		(user_id, token, expired_at, created_at, updated_at)
	VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, token, created_at, updated_at
	;`

	err := r.Db.QueryRowContext(ctx, SQL,
		pwdReset.UserId, pwdReset.Token,
		pwdReset.ExpiredAt,
	).Scan(&pwdReset.Id, &pwdReset.Token, &pwdReset.CreatedAt, &pwdReset.UpdatedAt)

	return pwdReset, err
}

func (r *PasswordResetRepositoryImpl) GetByUserId(ctx context.Context, userId uint) (*entity.PasswordReset, error) {
	SQL := `
	SELECT 
		id, user_id, token, expired_at, created_at, updated_at, deleted_at
	FROM password_resets
	WHERE deleted_at IS NULL
	AND user_id = $1
	;`

	rows, err := r.Db.QueryContext(ctx, SQL, userId)
	utils.IfErrorLogPrint(err)
	defer rows.Close()

	pwdReset := entity.PasswordReset{}
	if rows.Next() {
		err := rows.Scan(&pwdReset.Id, &pwdReset.UserId, &pwdReset.Token, &pwdReset.ExpiredAt, &pwdReset.CreatedAt, &pwdReset.UpdatedAt, &pwdReset.DeletedAt)
		utils.IfErrorLogPrint(err)
		return &pwdReset, nil
	}

	return &pwdReset, err
}

func (r *PasswordResetRepositoryImpl) GetByToken(ctx context.Context, token string) (*entity.PasswordReset, error) {
	SQL := `
	SELECT 
		id, user_id, token, expired_at
	FROM password_resets
	WHERE deleted_at IS NULL
	AND token = $1
	AND expired_at >= NOW()
	;`

	ps := entity.PasswordReset{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, SQL, token).Scan(&ps.Id, &ps.UserId, &ps.Token, &ps.ExpiredAt)
	} else {
		err = r.Db.QueryRowContext(ctx, SQL, token).Scan(&ps.Id, &ps.UserId, &ps.Token, &ps.ExpiredAt)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ErrResetTokenNotFound()
		}
		return nil, err
	}

	return &ps, err
}

func (r *PasswordResetRepositoryImpl) SoftDelete(ctx context.Context, id uint) error {
	SQL := `
	UPDATE password_resets
	SET deleted_at = NOW()
	WHERE id = $1
	;`

	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, SQL)
	} else {
		stmt, err = r.Db.PrepareContext(ctx, SQL)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperror.ErrResetTokenNotFound()
	}

	return nil
}
