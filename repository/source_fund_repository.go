package repository

import (
	"assignment-go-rest-api/apperror"
	"assignment-go-rest-api/entity"
	"context"
	"database/sql"
)

type SourceFundRepository interface {
	GetById(ctx context.Context, sourceId uint) (*entity.SourceOfFund, error)
}

type SourceFundRepositoryImpl struct {
	db *sql.DB
}

func NewSourceFundRepository(db *sql.DB) SourceFundRepository {
	return &SourceFundRepositoryImpl{
		db: db,
	}
}

func (r *SourceFundRepositoryImpl) GetById(ctx context.Context, sourceId uint) (*entity.SourceOfFund, error) {
	sourceFund := entity.SourceOfFund{}

	var err error

	SQL := `
	SELECT id, fund_name FROM source_of_funds WHERE id = $1
	;`

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, SQL, sourceId).Scan(&sourceFund.Id, &sourceFund.FundName)
	} else {
		err = r.db.QueryRowContext(ctx, SQL, sourceId).Scan(&sourceFund.Id, &sourceFund.FundName)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ErrSourceFundNotFound()
		}
		return nil, err
	}

	return &sourceFund, err
}
