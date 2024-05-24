package repository

import (
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/utils"
	"context"
	"database/sql"
)

type UserRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error)
	GetById(ctx context.Context, tx *sql.Tx, userId uint) (*entity.User, error)
	GetByEmail(ctx context.Context, tx *sql.Tx, userEmail string) (*entity.User, error)
	Update(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error)
}

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error) {
	SQL := `
	INSERT INTO users (email, username ,password, created_at, updated_at) VALUES
	($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at
	;`

	err := tx.QueryRowContext(ctx, SQL, user.Email, user.Username, user.Password).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)
	utils.IfErrorLogPrint(err)

	return user, err
}

func (r *UserRepositoryImpl) GetById(ctx context.Context, tx *sql.Tx, userId uint) (*entity.User, error) {
	SQL := `
	SELECT id, email, username, password, created_at, updated_at, deleted_at
	FROM users
	WHERE deleted_at IS NULL
	AND id = $1
	;`

	rows, err := tx.QueryContext(ctx, SQL, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := &entity.User{}
	if rows.Next() {
		err := rows.Scan(
			&user.Id, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		)
		utils.IfErrorLogPrint(err)
		return user, nil
	}

	return user, err
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, tx *sql.Tx, userEmail string) (*entity.User, error) {
	SQL := `
	SELECT id, email, username, password, created_at, updated_at, deleted_at
	FROM users
	WHERE deleted_at IS NULL
	AND email = $1 
	;`

	rows, err := tx.QueryContext(ctx, SQL, userEmail)
	utils.IfErrorLogPrint(err)
	defer rows.Close()

	user := &entity.User{}
	if rows.Next() {
		err := rows.Scan(
			&user.Id, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		)
		utils.IfErrorLogPrint(err)
		return user, nil
	}

	return user, err
}

func (r *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error) {
	SQL := `
	UPDATE users
	SET password = $1, updated_at = NOW()
	WHERE id = $2
	;`

	_, err := tx.ExecContext(ctx, SQL, user.Password, user.Id)
	utils.IfErrorLogPrint(err)

	return user, err
}
