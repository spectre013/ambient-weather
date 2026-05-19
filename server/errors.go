package main

import (
	"database/sql"
	"errors"
)

// sqlError logs SQL errors with a method name and query/context. Kept
// for back-compat with handler code; new code should prefer
// logger.WithError(...).WithField("query", q).Error(...).
func sqlError(method string, err error, value interface{}) {
	if err == nil {
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		logger.WithField("method", method).WithField("context", value).Error("Zero rows found")
	} else {
		logger.WithField("method", method).WithField("context", value).WithError(err).Error("scan failed")
	}
}

// errorHandler removed -- unused after inlining at call sites.
// Use logger.WithError(err).Error(message) directly instead.
