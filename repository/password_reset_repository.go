package repository

import (
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/utils"
	"context"
	"database/sql"
)

type PasswordResetRepository interface {
	Save(ctx context.Context, pwdReset *entity.PasswordReset) (*entity.PasswordReset, error)
	GetByUserId(ctx context.Context, userId uint) (*entity.PasswordReset, error)
	GetByToken(ctx context.Context, pwdToken string) (*entity.PasswordReset, error)
	SoftDelete(ctx context.Context, pwdReset *entity.PasswordReset) (*entity.PasswordReset, error)
}

type PasswordResetRepositoryImpl struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) PasswordResetRepository {
	return &PasswordResetRepositoryImpl{
		db: db,
	}
}

func (r *PasswordResetRepositoryImpl) Save(ctx context.Context, pwdReset *entity.PasswordReset) (*entity.PasswordReset, error) {
	SQL := `
	INSERT INTO password_resets
		(user_id, token, expired_at, created_at, updated_at)
	VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at
	;`

	err := r.db.QueryRowContext(ctx, SQL,
		pwdReset.UserId, pwdReset.Token,
		pwdReset.ExpiredAt,
	).Scan(&pwdReset.Id, &pwdReset.CreatedAt, &pwdReset.UpdatedAt)

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

	rows, err := r.db.QueryContext(ctx, SQL, userId)
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

func (r *PasswordResetRepositoryImpl) GetByToken(ctx context.Context, pwdToken string) (*entity.PasswordReset, error) {
	SQL := `
	SELECT 
		id, user_id, token, expired_at, created_at, updated_at, deleted_at
	FROM password_resets
	WHERE deleted_at IS NULL
	AND token = $1
	AND expired_at >= NOW()
	;`

	rows, err := r.db.QueryContext(ctx, SQL, pwdToken)
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

func (r *PasswordResetRepositoryImpl) SoftDelete(ctx context.Context, pwdReset *entity.PasswordReset) (*entity.PasswordReset, error) {
	SQL := `
	UPDATE password_resets
	SET deleted_at = NOW()
	WHERE user_id = $1
	;`

	_, err := r.db.QueryContext(ctx, SQL, pwdReset.UserId)
	utils.IfErrorLogPrint(err)

	return pwdReset, err
}
