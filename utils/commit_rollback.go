package utils

import "database/sql"

func CommitOrRollback(tx *sql.Tx) {
	err := recover()
	if err != nil {
		errRollback := tx.Rollback()
		IfErrorLogPrint(errRollback)
		return
	} else {
		errorCommit := tx.Commit()
		IfErrorLogPrint(errorCommit)
	}
}
