package drivers

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MySQLConnection struct {
	Conn *sqlx.DB
}

func (m *MySQLConnection) GetDB() *sqlx.DB {
	return m.Conn
}

func NewMySQLConnection(host string, port int, username string, password string, dbName string) (Connection, error) {

	// Build the connection string from the config
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		username,
		password,
		host,
		port,
		dbName,
	)

	// Open a connection to the MySQL database
	conn, err := sqlx.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Test the connection to ensure it's working
	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return &MySQLConnection{Conn: conn}, nil
}
