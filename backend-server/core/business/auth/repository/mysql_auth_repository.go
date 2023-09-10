package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLAuthRepository struct {
	DB *sqlx.DB
}

func NewSQLUserRepository(conn drivers.Connection) AuthRepository {
	return &SQLAuthRepository{DB: conn.GetDB()}
}
