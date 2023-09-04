package data

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type NullInt64 struct {
	sql.NullInt64
}

type NullString struct {
	sql.NullString
}

type NullTime struct {
	mysql.NullTime
}

type NullBool struct {
	sql.NullBool
}
