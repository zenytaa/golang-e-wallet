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

type UserRepoOpts struct {
	Db *sql.DB
}

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) (*entity.User, error)
	GetById(ctx context.Context, userId uint) (*entity.User, error)
	GetByEmail(ctx context.Context, userEmail string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(userROpts *UserRepoOpts) UserRepository {
	return &UserRepositoryImpl{db: userROpts.Db}
}

func (r *UserRepositoryImpl) Save(ctx context.Context, user *entity.User) (*entity.User, error) {
	newUser := entity.User{}
	var err error

	values := []interface{}{}
	values = append(values, user.Email)
	values = append(values, user.Username)
	values = append(values, user.Password)

	SQL := `
	INSERT INTO users (email, username ,password, created_at, updated_at) VALUES
	($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at
	;`

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, SQL, values...).Scan(&newUser.Id, &newUser.CreatedAt, &newUser.UpdatedAt)
	} else {
		err = r.db.QueryRowContext(ctx, SQL, values...).Scan(&newUser.Id, &newUser.CreatedAt, &newUser.UpdatedAt)
	}

	if err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constant.ViolatesUniqueConstraintPgErrCode {
			return nil, apperror.ErrEmailAlreadyRegistered()
		}
		return nil, err
	}

	return &newUser, err
}

func (r *UserRepositoryImpl) GetById(ctx context.Context, userId uint) (*entity.User, error) {
	user := entity.User{}
	var err error

	SQL := `
	SELECT id, email, username, password, created_at, updated_at, deleted_at
	FROM users
	WHERE deleted_at IS NULL
	AND id = $1
	;`

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, SQL, userId).Scan(
			&user.Id, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		)
	} else {
		err = r.db.QueryRowContext(ctx, SQL, userId).Scan(
			&user.Id, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ErrUserNotFound()
		}
		return nil, err
	}

	return &user, err
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, userEmail string) (*entity.User, error) {
	user := entity.User{}
	var err error

	SQL := `
	SELECT id, email, username, password, created_at, updated_at, deleted_at
	FROM users
	WHERE deleted_at IS NULL
	AND email = $1 
	;`

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, SQL, userEmail).Scan(
			&user.Id, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		)
	} else {
		err = r.db.QueryRowContext(ctx, SQL, userEmail).Scan(
			&user.Id, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ErrUserNotFound()
		}
		return nil, err
	}

	return &user, err
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	var err error
	var stmt *sql.Stmt

	SQL := `
	UPDATE users
	SET password = $1, updated_at = NOW()
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

	res, err := stmt.ExecContext(ctx, user.Password, user.Id)
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

	return user, err
}
