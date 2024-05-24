package repository

import (
	"assignment-go-rest-api/entity"
	"assignment-go-rest-api/utils"
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
	SQL := `
	SELECT id, fund_name FROM source_of_funds WHERE id = $1
	;`

	rows, err := r.db.QueryContext(ctx, SQL, sourceId)
	utils.IfErrorLogPrint(err)
	defer rows.Close()

	sourceFund := &entity.SourceOfFund{}
	if rows.Next() {
		err := rows.Scan(&sourceFund.Id, &sourceFund.FundName)
		utils.IfErrorLogPrint(err)
		return sourceFund, nil
	}

	return sourceFund, err
}
