package main

import (
	"database/sql"
	"errors"
	"fmt"
)

func sqlError(method string, err error, value interface{}) {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error(fmt.Sprintf("%s - Zero Rows Found ", method), value)
		} else {
			logger.Error(fmt.Sprintf("%s - Scan: ", method), err)
		}
	}
}

func errorHandler(err error, message string) {
	if err != nil {
		logger.Error(message, err)
	}
}
