package drivers

import (
	"github.com/jmoiron/sqlx"
)

type Connection interface {
	GetDB() *sqlx.DB
}
