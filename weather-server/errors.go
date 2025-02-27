package main

import (
	"database/sql"
	"errors"
)

func sqlError(err error, value interface{}) {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Info("Zero Rows Found ", value)
		} else {
			logger.Error("Scan: ", err)
		}
	}
}

func errorHandler(err error, message string) {
	if err != nil {
		logger.Error(message, err)
	}
}
