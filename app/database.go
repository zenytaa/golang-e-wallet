package app

import (
	"assignment-go-rest-api/utils"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

func ConnectDB(config utils.Config) (*sql.DB, error) {

	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.DbHost,
		config.DbUser,
		config.DbPassword,
		config.DbName,
		config.DbPort,
		config.DbSSL,
		config.DbTimeZone,
	)
	fmt.Println(dataSourceName)
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
